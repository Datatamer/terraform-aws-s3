package test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/testing"
)

// func uploadObjectE(awsRegion string, bucketName string, key string, encryption string, bodyString string) error {

// 	up := s3manager.NewUploader(sess)
// 	return uploadObjectWithUploaderE(awsRegion, bucketName, key, encryption, bodyString, up)
// }

func uploadObjectWithUploaderE(awsRegion string, bucketName string, key string, encryption string, bodyString string, up *s3manager.Uploader) error {
	var upInput *s3manager.UploadInput

	if encryption != "" {
		upInput = &s3manager.UploadInput{
			Bucket:               &bucketName,
			Key:                  &key,
			Body:                 strings.NewReader(bodyString),
			ServerSideEncryption: &encryption,
		}
	} else {
		upInput = &s3manager.UploadInput{
			Bucket: &bucketName,
			Key:    &key,
			Body:   strings.NewReader(bodyString),
		}
	}

	_, err := up.Upload(upInput)

	if err != nil {
		return err
	}

	return nil
}

func GetS3ObjectContentsWithSessionE(t testing.TestingT, awsRegion string, bucket string, key string, sess *session.Session) (string, error) {
	s3Client := s3.New(sess)

	res, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return "", err
	}

	contents := buf.String()
	logger.Log(t, fmt.Sprintf("Read contents from s3://%s/%s", bucket, key))

	return contents, nil
}
