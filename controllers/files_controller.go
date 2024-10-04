package controllers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mypackages/consts"
	"mypackages/db"
	"mypackages/helpers"
	Model "mypackages/models"
	"mypackages/policy"
	"mypackages/proto/files"
	"os"
	"strconv"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)



func DownloadFile(in *files.FileDownloadRequest, responseStream files.FilesGreeter_DownloadFileServer) error {
	user, err := helpers.GetUserFormMd(responseStream.Context())

    if err != nil {
        fmt.Println(err)
        return err
    }

	var fileData *Model.File
	var result *gorm.DB;
	if in.GetFolderId() == 0 {
		result = db.DB.Model(&Model.File{}).Where(
			"id = ? AND folder_id is NULL AND user_id = ?",
			in.GetFileId(),
			user.ID).First(&fileData)
	} else {
		result = db.DB.Model(&Model.File{}).Where(
			"id = ? AND folder_id = ? AND user_id = ?",
			in.GetFileId(),
			in.GetFolderId(),
			user.ID).First(&fileData)
	}

	if result.RowsAffected == 0 && result.Error!=nil {
		return status.Error(codes.NotFound, "Файл не найден")
	}

	bufferSize := 256 * 1024
	var path string = getFilePath(fileData)

    file, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return err
    }
    
	stat, _ := os.Stat(path)
	var fileSize = 0
    defer file.Close()

    buff := make([]byte, bufferSize)
    for {
        bytesRead, err := file.Read(buff)
        if err != nil {
            if err != io.EOF {
                fmt.Println(err)
            }
            break
        }
		fileSize += bufferSize
        resp := &files.FileDownloadResponse{
            Chunk: buff[:bytesRead],
			FileName: fileData.FileName,
			Progress: float32(fileSize) / float32(stat.Size()) * 100,
        }
        err = responseStream.Send(resp)
        if err != nil {
            log.Println("error while sending chunk:", err)
            return err
        }
    }
    return nil
}

func UploadFile(stream files.FilesGreeter_UploadFileServer) error{
	user, err := helpers.GetUserFormMd(stream.Context())
	
	if err != nil {
		return stream.SendAndClose(&files.FileUploadResponse{Message: "Пользователь не найден"})
	}

	filesNameHash := uuid.New().String()

	var fileSize uint32
	req, err := stream.Recv()

	if err != nil || !policy.FolderPolicyID(req.GetFolderId(), user){
		return stream.SendAndClose(&files.FileUploadResponse{Message: "Не удалось загрузить файл"})
	}

    fileSize = 0

	result, file := createFile(req, user, filesNameHash)

	if result.RowsAffected == 0 || result.Error != nil {
		return status.Error(codes.Internal, "Не удалось создать файл")
	}
	path := getUploadPath(user, filesNameHash, req.GetFolderId())
	dst, _ := os.Create(path)
	
	err = writeInFile(req, dst, &fileSize, user, filesNameHash, file)

	if err != nil{
		return status.Error(codes.PermissionDenied, err.Error())
	}

    for {
        req, err := stream.Recv()
        if err == io.EOF {
            break
        }
        if err != nil {
            return status.Error(codes.Internal, err.Error())
        }

		err = writeInFile(req, dst, &fileSize, user, filesNameHash, file)
		
		if err != nil{
			return status.Error(codes.PermissionDenied, err.Error())
		}
		db.DB.Model(&Model.File{}).Where("id=?", file.ID).Update("size", fileSize)
    }
	
	defer func() {
		if stream.Context().Err() != nil {
			rollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		}

		if recover() != nil {
			rollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		}
	}()

	return stream.SendAndClose(&files.FileUploadResponse{Message: "Успешно загружено"})
}

func getFilePath(file *Model.File) string{
	var pathFileFolder, _ = os.LookupEnv("PATH_FILES")

	var path string
	if file.FolderID==0{
		path = pathFileFolder+strconv.Itoa(file.UserID)+"/"+file.FileNameHash;	
	} else{
		path = pathFileFolder+strconv.Itoa(file.UserID)+"/"+strconv.Itoa(file.FolderID)+"/"+file.FileNameHash;	
	}
	return path;
}

func getUploadPath(user *Model.User, filesNameHash string, folderId uint32) string {
	var pathFileFolder, _ = os.LookupEnv("PATH_FILES")
	var path string

	if folderId != 0 {
		path = pathFileFolder+strconv.Itoa(int(user.ID))+"/"+strconv.Itoa(int(folderId))+"/"+filesNameHash
	} else{
		path = pathFileFolder+strconv.Itoa(int(user.ID))+"/"+filesNameHash
	}
	return path
}

func createFile(req *files.FileUploadRequest, user *Model.User, filesNameHash string) (*gorm.DB, *Model.File) {
	var file *Model.File
	if req.GetFolderId()!=0 {
		file = &Model.File{
			FileName: req.GetFileName(),
			UserRelation: Model.UserRelation{
				UserID: int(user.ID),
			},
			FolderRelation: Model.FolderRelation{
				FolderID: int(req.GetFolderId()),
			},
			Size: int(len(req.GetChunk())),
			FileNameHash: filesNameHash,
			AccessId: consts.CLOSE,
		};
	} else{
		file = &Model.File{
			FileName:  req.GetFileName(),
			UserRelation: Model.UserRelation{
				UserID: int(user.ID),
			},
			Size:  int(len(req.GetChunk())),
			FileNameHash: filesNameHash,
			AccessId: consts.CLOSE,
		};
	}

	result := db.DB.Model(&Model.File{}).Create(&file)
	return result, file
}

func rollbackFile(user *Model.User, filesNameHash string, folderId uint32, fileId uint) {
	path := getUploadPath(user, filesNameHash, folderId)
	db.DB.Model(&Model.File{}).Where("id=?", fileId).Unscoped().Delete(&Model.File{})
	os.Remove(path)
}

func writeInFile(req *files.FileUploadRequest, dst *os.File, fileSize *uint32, user *Model.User, filesNameHash string, file *Model.File) error {
	chunk := req.GetChunk()
	fmt.Println(uint32(len(chunk)/12))
	dst.WriteAt(chunk, int64(*fileSize))
		
	*fileSize += uint32(len(chunk))

	if !policy.SpacePolicy(uint32(len(chunk))) {
		rollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		return errors.New("недостаточно места")
	}
	result := db.DB.Model(&Model.File{}).Where("id=?", file.ID).Update("size", fileSize)
	if result.Error != nil {
		rollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		return errors.New("недостаточно места")
	}
	return nil
}
