package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mypackages/helpers"
	"mypackages/proto/keys"
	"os"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type keysServer struct{
	keys.UnimplementedKeysGreeterServer	
}

func (s *keysServer) UploadKeys(ctx context.Context, in *keys.FileUploadRequest) (*keys.FileUploadResponse, error) {
	user, err := helpers.GetUserFormMd(ctx)
	
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "Пользователь не найден")
	}

	file, err := os.Create("keys/"+strconv.Itoa(int(user.ID)))

	if err != nil {
		log.Println("UploadKeys err: ", err)
		return nil, status.Error(codes.Aborted, "Ошибка при передаче ключей")
	}
    defer file.Close()
	size := len(in.GetChunk())

	if size > 6 * 1024 * 1024 {
		fmt.Println(size)
		return nil, status.Error(codes.Aborted, "Файл слишком большой")
	}

	os.WriteFile(file.Name(), in.GetChunk(), 0644)
	return &keys.FileUploadResponse{
		Message: "Ключи успешно отправлены",
	}, nil
}

func (s *keysServer) DownloadKeys(in *keys.Empty, responseStream keys.KeysGreeter_DownloadKeysServer) error {
	user, err := helpers.GetUserFormMd(responseStream.Context())

    if err != nil {
        fmt.Println(err)
        return err
    }
	var pathKeysFolder, _ = os.LookupEnv("PATH_KEYS")

	bufferSize := 64 *1024
    file, err := os.Open(pathKeysFolder+strconv.Itoa(int(user.ID)))
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
        err = responseStream.Send(resp)
        if err != nil {
            log.Println("error while sending chunk:", err)
            return err
        }
    }
    return nil
}