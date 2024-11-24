package chat_actions

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"errors"

	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func UploadFile(ctx context.Context,file []byte, resp *s3.CreateMultipartUploadOutput) {
	fmt.Println(resp)

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
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    aws.Int32(int32(partNumber)),
		
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

func DownloadChunk(ctx context.Context, resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) ([]byte, error) {
	partInput := &s3.GetObjectInput{
		Bucket:     resp.Bucket,
		Key:        resp.Key,
		PartNumber: aws.Int32(int32(partNumber)),
		Range:      aws.String("bytes=0-500"),
	}

	r, err := db.SVC.GetObject(ctx, partInput)

	if err != nil {
		return nil, errors.New("не удалось скачать файл")
	}

	body, err := io.ReadAll(r.Body)
	return body, nil
}
