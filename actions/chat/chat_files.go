package chat_actions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"

	"errors"

	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func UploadFile(ctx context.Context, file []byte, resp *s3.CreateMultipartUploadOutput) {

	_, err := db.SVC.UploadPart(ctx, &s3.UploadPartInput{
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		UploadId:      resp.UploadId,
		Body:          bytes.NewReader(file),
		PartNumber:    aws.Int32(0),
		ContentLength: aws.Int64(int64(len(file))),
	})

	fmt.Println(err)
}

func CompleteMultipartUpload(ctx context.Context, resp *s3.CreateMultipartUploadOutput, completedParts []types.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return db.SVC.CompleteMultipartUpload(ctx, completeInput)
}

func UploadPart(ctx context.Context, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*types.CompletedPart, error) {
	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:       bytes.NewReader(fileBytes),
		Bucket:     resp.Bucket,
		Key:        resp.Key,
		PartNumber: aws.Int32(int32(partNumber)),

		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	}

	for tryNum <= 1 {
		uploadResult, err := db.SVC.UploadPart(ctx, partInput)
		if err != nil {
			if tryNum == 1 {
				return nil, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			return &types.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int32(int32(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func AbortMultipartUpload(ctx context.Context, resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := db.SVC.AbortMultipartUpload(ctx, abortInput)
	return err
}

func Rollback(messageId uint) {
	db.DB.Unscoped().Where("message_id = ?", messageId).Delete(&Model.ChatFile{})
	db.DB.Unscoped().Where("id = ?", messageId).Delete(&Model.Message{})
}

func DownloadChunk(ctx context.Context, key string, rangeInt int64, size int) ([]byte, error) {
	awsBucket, ok := os.LookupEnv("AWS_BUCKET")

	if !ok {
		return nil, errors.New("не удалось скачать файл")
	}

	rangeNum := int(rangeInt)
	rangeNumEnd := int(rangeInt) + 256*1024

	if rangeNumEnd > size {
		fmt.Println(rangeNumEnd, " ", size)
		rangeNumEnd = size
	}

	partInput := &s3.GetObjectInput{
		Bucket:     aws.String(awsBucket),
		Key:        aws.String("ChatFiles/"+key),
		PartNumber: aws.Int32(int32(rangeInt / (256*1024))),
		Range:      aws.String("bytes=" + strconv.Itoa(rangeNum) + "-" + strconv.Itoa(rangeNumEnd)),
	}
	r, err := db.SVC.GetObject(ctx, partInput)
	if err != nil {
		return nil, errors.New("не удалось скачать файл: " + err.Error())
	}
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.New("не удалось скачать файл: " + err.Error())
	}
	return bytes, nil

}

func GetFileSize(ctx context.Context, key string) (*int64, error) {
	awsBucket, ok := os.LookupEnv("AWS_BUCKET")

	if !ok {
		return nil, errors.New("не удалось скачать файл")
	}

	headInput := &s3.HeadObjectInput{
		Key:    aws.String("ChatFiles/"+key),
		Bucket: aws.String(awsBucket),
	}

	resp, err := db.SVC.HeadObject(ctx, headInput)

	if err != nil {
		return nil, err
	}

	return resp.ContentLength, nil
}

func CheckChatFile(user *Model.User,chatId uint32) error {
	var chatFile *Model.ChatFile
	r := db.DB.Model(&Model.ChatFile{}).Where("id = ?", chatId).First(&chatFile)
	
	if r.RowsAffected == 0 || r.Error != nil {
		return errors.New("файл чата не найден")
	}

	var chatUser *Model.ChatUser
	r = db.DB.Model(&Model.ChatUser{}).Where("chat_id = ? AND user_id = ?", chatFile.ChatID, user.ID).First(&chatUser)

	if r.RowsAffected == 0 || r.Error != nil {
		return errors.New("чат не найден")
	}

	return nil
}