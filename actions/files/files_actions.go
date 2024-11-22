package files_actions

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/R1kkass/GoCloudGRPC/consts"
	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/policy"
	"github.com/R1kkass/GoCloudGRPC/proto/files"
	"gorm.io/gorm"
)

func GetFilePath(file *Model.File) string {
	var pathFileFolder, _ = os.LookupEnv("PATH_FILES")

	var path string
	if file.FolderID == 0 {
		path = pathFileFolder + strconv.Itoa(int(file.UserID)) + "/" + file.FileNameHash
	} else {
		path = pathFileFolder + strconv.Itoa(int(file.UserID)) + "/" + strconv.Itoa(int(file.FolderID)) + "/" + file.FileNameHash
	}
	return path
}

func GetUploadPath(user *Model.User, filesNameHash string, folderId uint32) string {
	var pathFileFolder, _ = os.LookupEnv("PATH_FILES")
	var path string

	if folderId != 0 {
		path = pathFileFolder + strconv.Itoa(int(user.ID)) + "/" + strconv.Itoa(int(folderId)) + "/" + filesNameHash
	} else {
		path = pathFileFolder + strconv.Itoa(int(user.ID)) + "/" + filesNameHash
	}
	return path
}

func CreateFile(req *files.FileUploadRequest, user *Model.User, filesNameHash string) (*gorm.DB, *Model.File) {
	var file *Model.File
	if req.GetFolderId() != 0 {
		file = &Model.File{
			FileName: req.GetFileName(),
			UserRelation: Model.UserRelation{
				UserID: user.ID,
			},
			FolderRelation: Model.FolderRelation{
				FolderID: uint(req.GetFolderId()),
			},
			Size:         int(len(req.GetChunk())),
			FileNameHash: filesNameHash,
			AccessId:     consts.CLOSE,
		}
	} else {
		file = &Model.File{
			FileName: req.GetFileName(),
			UserRelation: Model.UserRelation{
				UserID: user.ID,
			},
			Size:         int(len(req.GetChunk())),
			FileNameHash: filesNameHash,
			AccessId:     consts.CLOSE,
		}
	}

	result := db.DB.Model(&Model.File{}).Create(&file)
	return result, file
}

func RollbackFile(user *Model.User, filesNameHash string, folderId uint32, fileId uint) {
	path := GetUploadPath(user, filesNameHash, folderId)
	db.DB.Model(&Model.File{}).Where("id=?", fileId).Unscoped().Delete(&Model.File{})
	os.Remove(path)
}

func WriteInFile(req *files.FileUploadRequest, dst *os.File, fileSize *uint32, user *Model.User, filesNameHash string, file *Model.File) error {
	chunk := req.GetChunk()
	fmt.Println(uint32(len(chunk) / 12))
	dst.WriteAt(chunk, int64(*fileSize))

	*fileSize += uint32(len(chunk))

	if !policy.SpacePolicy(uint32(len(chunk))) {
		RollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		return errors.New("недостаточно места")
	}
	result := db.DB.Model(&Model.File{}).Where("id=?", file.ID).Update("size", fileSize)
	if result.Error != nil {
	    RollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		return errors.New("недостаточно места")
	}
	return nil
}
