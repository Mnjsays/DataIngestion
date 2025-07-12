package storage

import (
	"bytes"
	"dataIngestion/pkg/models"
	"dataIngestion/types"
	"dataIngestion/util"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"go.uber.org/zap"
	"io"
	"time"
)

func getAwsS3Session(app *types.App) (*s3.S3, error) {
	accessKey := app.Config.AwsS3.AccessKey
	secretKey := app.Config.AwsS3.SecretKey
	region := app.Config.AwsS3.Region

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)

	}
	obj := s3.New(sess)
	return obj, err
}
func AwsRead(fileName string, app *types.App) ([]byte, error) {
	ss, err := getAwsS3Session(app)
	if err != nil {
		app.Logger.Error("S3 Session error: %v", zap.Error(err))
		return nil, err
	}
	folder := app.Config.AwsS3.Folder
	key := fmt.Sprintf("%s/%s", folder, fileName)
	result, err := ss.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(app.Config.AwsS3.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer result.Body.Close()
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %v", err)
	}
	return data, nil
}
func AwsWrite(posts *models.Posts, app *types.App) error {
	jsonData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		app.Logger.Error("JSON marshal error: %v", zap.Error(err))
		return err
	}
	svc, err := getAwsS3Session(app)
	if err != nil {
		app.Logger.Error("S3 Session error: %v", zap.Error(err))
		return err
	}
	folder := app.Config.AwsS3.Folder
	savedAs := time.Now().UTC().Format("2025-01-02T15-04-05")
	savedAs = util.Sanitize(savedAs)

	key := fmt.Sprintf("%s/%s", folder, savedAs)
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
	app.Logger.Info("Data stored in s3 with file name ::", zap.String("file_name", savedAs))
	return nil
}
