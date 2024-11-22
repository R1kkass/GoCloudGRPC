package db

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var SVC *s3.S3

func AwsConnect() {

	awsRegion, _ := os.LookupEnv("AWS_DEFAULT_REGION")

	sess, _ := session.NewSession(&aws.Config{
		Region:   aws.String(awsRegion),
		Endpoint: aws.String("https://storage.yandexcloud.net/")},
	)

	SVC = s3.New(sess)

}