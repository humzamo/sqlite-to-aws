package awsupload

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
)

const (
	awsRegion  = "eu-west-2"
	bucketName = "humza-mo-sqlite-to-aws"
	path       = "client-data"
	timeFormat = "2006-01-02-15:04:05.000"
)

func UploadToS3Bucket(fileName string, fileBytes []byte) error {
	// Create an AWS S3 configuration and client
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(awsRegion),
	)
	if err != nil {
		errors.Wrap(err, "error loading AWS configuration")
	}

	client := s3.NewFromConfig(cfg)

	uploadInput := &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(path + "/" + fileName),
		Body:        aws.ReadSeekCloser(bytes.NewReader(fileBytes)),
		ContentType: aws.String("application/json"),
	}

	// upload the JSON data to S3.
	_, err = client.PutObject(context.Background(), uploadInput)
	if err != nil {
		errors.Wrap(err, "error uploading JSON to S3")
	}

	return nil
}
