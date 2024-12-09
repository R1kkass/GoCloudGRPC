package controllers

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/R1kkass/GoCloudGRPC/helpers"
	"github.com/R1kkass/GoCloudGRPC/proto/keys"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type KeysServer struct {
	keys.UnimplementedKeysGreeterServer
}

func (s *KeysServer) UploadKeys(ctx context.Context, in *keys.FileUploadRequest) (*keys.FileUploadResponse, error) {
	user, err := helpers.GetUserFormMd(ctx)

	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Пользователь не найден")
	}
	var pathKeysFolder, _ = os.LookupEnv("PATH_KEYS")

	file, err := os.OpenFile(pathKeysFolder+strconv.Itoa(int(user.ID)), os.O_WRONLY, 0666)

	if err != nil {
		log.Println("UploadKeys err: ", err)
		return nil, status.Error(codes.Aborted, "Ошибка при передаче ключей")
	}
	defer file.Close()
	size := len(in.GetChunk())

	if size > 6*1024*1024 {
		fmt.Println(size)
		return nil, status.Error(codes.Aborted, "Файл слишком большой")
	}

	os.WriteFile(file.Name(), in.GetChunk(), 0644)
	return &keys.FileUploadResponse{
		Message: "Ключи успешно отправлены",
	}, nil
}

func (s *KeysServer) DownloadKeys(in *keys.Empty, responseStream keys.KeysGreeter_DownloadKeysServer) error {
	user, err := helpers.GetUserFormMd(responseStream.Context())
	fmt.Println("pathKeysFolder")

	if err != nil {
		fmt.Println(err)
		return err
	}
	var pathKeysFolder, _ = os.LookupEnv("PATH_KEYS")

	bufferSize := 64 * 1024
	file, err := os.Open(pathKeysFolder + strconv.Itoa(int(user.ID)))
	if err != nil {
		fmt.Println(err)
		return err
	}
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
		resp := &keys.FileDownloadResponse{
			Chunk: buff[:bytesRead],
		}
		fmt.Println(resp)
		err = responseStream.Send(resp)
		if err != nil {
			log.Println("error while sending chunk:", err)
			return err
		}
	}
	return nil
}
