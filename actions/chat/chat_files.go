package chat_actions

import (
	"bytes"
	"fmt"
	"io"

	"errors"

	"github.com/R1kkass/GoCloudGRPC/db"
	Model "github.com/R1kkass/GoCloudGRPC/models"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

func UploadFile(file []byte, resp *s3.CreateMultipartUploadOutput) {
	_, err := db.SVC.UploadPart(&s3.UploadPartInput{
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		UploadId:      resp.UploadId,
		Body:          bytes.NewReader(file),
		PartNumber:    aws.Int64(0),
		ContentLength: aws.Int64(int64(len(file))),
	})

	fmt.Println(err)
}

func CompleteMultipartUpload(resp *s3.CreateMultipartUploadOutput, completedParts []*s3.CompletedPart) (*s3.CompleteMultipartUploadOutput, error) {
	completeInput := &s3.CompleteMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	}
	return db.SVC.CompleteMultipartUpload(completeInput)
}

func UploadPart(resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) (*s3.CompletedPart, error) {
	tryNum := 1
	partInput := &s3.UploadPartInput{
		Body:          bytes.NewReader(fileBytes),
		Bucket:        resp.Bucket,
		Key:           resp.Key,
		PartNumber:    aws.Int64(int64(partNumber)),
		UploadId:      resp.UploadId,
		ContentLength: aws.Int64(int64(len(fileBytes))),
	}

	for tryNum <= 1 {
		uploadResult, err := db.SVC.UploadPart(partInput)
		if err != nil {
			if tryNum == 1 {
				if aerr, ok := err.(awserr.Error); ok {
					return nil, aerr
				}
				return nil, err
			}
			fmt.Printf("Retrying to upload part #%v\n", partNumber)
			tryNum++
		} else {
			fmt.Printf("Uploaded part #%v\n", partNumber)
			return &s3.CompletedPart{
				ETag:       uploadResult.ETag,
				PartNumber: aws.Int64(int64(partNumber)),
			}, nil
		}
	}
	return nil, nil
}

func AbortMultipartUpload(resp *s3.CreateMultipartUploadOutput) error {
	fmt.Println("Aborting multipart upload for UploadId#" + *resp.UploadId)
	abortInput := &s3.AbortMultipartUploadInput{
		Bucket:   resp.Bucket,
		Key:      resp.Key,
		UploadId: resp.UploadId,
	}
	_, err := db.SVC.AbortMultipartUpload(abortInput)
	return err
}

func Rollback(messageId uint) {
	db.DB.Unscoped().Where("id = ?", messageId).Delete(&Model.Message{})
	db.DB.Unscoped().Where("message_id = ?", messageId).Delete(&Model.ChatFile{})
}

func DownloadChunk(resp *s3.CreateMultipartUploadOutput, fileBytes []byte, partNumber int) ([]byte, error) {
	partInput := &s3.GetObjectInput{
		Bucket:     resp.Bucket,
		Key:        resp.Key,
		PartNumber: aws.Int64(int64(partNumber)),
		Range:      aws.String("bytes=0-500"),
	}

	r, err := db.SVC.GetObject(partInput)

	if err != nil {
		return nil, errors.New("не удалось скачать файл")
	}

	body, err := io.ReadAll(r.Body)
	return body, nil
}
