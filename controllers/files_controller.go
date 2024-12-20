package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	files_actions "github.com/R1kkass/GoCloudGRPC/actions/files"
	"github.com/R1kkass/GoCloudGRPC/db"
	"github.com/R1kkass/GoCloudGRPC/helpers"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/R1kkass/GoCloudGRPC/policy"
	"github.com/R1kkass/GoCloudGRPC/proto/files"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type FilesServer struct {
	*files.UnimplementedFilesGreeterServer
}

func (s *FilesServer) DownloadFile(in *files.FileDownloadRequest, responseStream files.FilesGreeter_DownloadFileServer) error {
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error DownloadFile: ", r)
		}
	}()
	
	user, err := helpers.GetUserFormMd(responseStream.Context())

	if err != nil {
		fmt.Println(err)
		return err
	}

	var fileData *Model.File
	var result *gorm.DB
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

	if result.RowsAffected == 0 && result.Error != nil {
		return status.Error(codes.NotFound, "Файл не найден")
	}

	bufferSize := 256 * 1024
	var path string = files_actions.GetFilePath(fileData)

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
			Chunk:    buff[:bytesRead],
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

func (s *FilesServer) UploadFile(stream files.FilesGreeter_UploadFileServer) error {
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error UploadFile: ", r)
		}
	}()
	
	user, err := helpers.GetUserFormMd(stream.Context())

	if err != nil {
		return status.Error(codes.Internal, "Пользователь не найден")
	}

	filesNameHash := uuid.New().String()

	var fileSize uint32 = 0
	req, err := stream.Recv()

	if err != nil || !policy.FolderPolicyID(req.GetFolderId(), user) {
		return status.Error(codes.Internal, "Не удалось загрузить файл")
	}

	result, file := files_actions.CreateFile(req, user, filesNameHash)

	if result.RowsAffected == 0 || result.Error != nil {
		return status.Error(codes.Internal, "Не удалось создать файл")
	}
	path := files_actions.GetUploadPath(user, filesNameHash, req.GetFolderId())
	dst, _ := os.Create(path)

	err = files_actions.WriteInFile(req, dst, &fileSize, user, filesNameHash, file)

	if err != nil {
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

		err = files_actions.WriteInFile(req, dst, &fileSize, user, filesNameHash, file)

		if err != nil {
			return status.Error(codes.PermissionDenied, err.Error())
		}
		db.DB.Model(&Model.File{}).Where("id=?", file.ID).Update("size", fileSize)
	}

	defer func() {
		if stream.Context().Err() != nil {
			files_actions.RollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		}

		if recover() != nil {
			files_actions.RollbackFile(user, filesNameHash, req.GetFolderId(), file.ID)
		}
	}()

	return stream.SendAndClose(&files.FileUploadResponse{Message: "Успешно загружено"})
}

func (s *FilesServer) FindFile(context context.Context, in *files.FindFileRequest) (*files.FindFileResponse, error) {
	
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Error FindFile: ", r)
		}
	}()
	
	user, err := helpers.GetUserFormMd(context)
	var file []*files.FileFind
	var folder []*files.FolderFind

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	qFile := db.DB.Model(&Model.File{}).Where("user_id = ?", user.ID)
	qFolder := db.DB.Model(&Model.Folder{}).Where("user_id = ?", user.ID)

	if !in.GetFindEveryWhere() {
		if in.GetFolderId() == 0 {
			qFile.Where("folder_id is NULL")
			qFolder.Where("folder_id is NULL")
		} else {
			qFile.Where("folder_id = ?", in.GetFolderId())
			qFolder.Where("folder_id = ?", in.GetFolderId())
		}
	}
	resultQFile := qFile.Where("file_name LIKE ?", "%"+in.GetSearch()+"%").Limit(10).Offset((int(in.GetPage()) - 1) * 10).Find(&file)
	resultQFolder := qFolder.Where("name_folder LIKE ?", "%"+in.GetSearch()+"%").Limit(10).Offset((int(in.GetPage()) - 1) * 10).Find(&folder)

	if resultQFile.Error != nil || resultQFolder.Error != nil {
		return nil, status.Error(codes.Unknown, "Не удалось найти файлы")
	}

	return &files.FindFileResponse{
		Files:   file,
		Folders: folder,
	}, nil
}
