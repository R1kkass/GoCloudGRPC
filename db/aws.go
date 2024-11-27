package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/aws"
)

var SVC *s3.Client

func AwsConnect() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	cfg.BaseEndpoint = aws.String("https://storage.yandexcloud.net")
	// Создаем клиента для доступа к хранилищу S3
	SVC = s3.NewFromConfig(cfg)
}
