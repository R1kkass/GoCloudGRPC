package db

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var SVC *s3.Client

func AwsConnect() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	fmt.Println(cfg.Region)
    if err != nil {
        log.Fatal(err)
    }

    // Создаем клиента для доступа к хранилищу S3
    SVC = s3.NewFromConfig(cfg)
}