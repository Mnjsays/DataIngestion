package storage

import (
	"bytes"
	"dataIngestion/pkg/models"
	"dataIngestion/types"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"time"
)

func AwsStorage(posts *models.Posts, app *types.App) error {
	jsonData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		app.Logger.Error("JSON marshal error: %v", zap.Error(err))
		return err
	}
	accessKey := app.Config.AwsS3.AccessKey
	secretKey := app.Config.AwsS3.SecretKey
	region := app.Config.AwsS3.Region
	folder := app.Config.AwsS3.Folder
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return fmt.Errorf("failed to create AWS session: %v", err)

	}
	svc := s3.New(sess)
	timestamp := time.Now().UTC().Format("2025-01-02T15-04-05")
	key := fmt.Sprintf("%s/%s", folder, timestamp)
	params := &s3.PutObjectInput{
		Bucket:      aws.String(app.Config.AwsS3.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonData),
		ContentType: aws.String("application/json"),
	}
	_, err = svc.PutObject(params)
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	return nil

}
