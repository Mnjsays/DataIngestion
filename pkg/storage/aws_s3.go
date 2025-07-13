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

type S3Storage struct {
	App *types.App
}

func (s *S3Storage) getAwsS3Session() (*s3.S3, error) {
	accessKey := s.App.Config.AwsS3.AccessKey
	secretKey := s.App.Config.AwsS3.SecretKey
	region := s.App.Config.AwsS3.Region

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %v", err)
	}
	return s3.New(sess), nil
}

func (s *S3Storage) ReadData(fileName string) ([]byte, error) {
	ss, err := s.getAwsS3Session()
	if err != nil {
		s.App.Logger.Error("S3 Session error: %v", zap.Error(err))
		return nil, err
	}
	folder := s.App.Config.AwsS3.Folder
	key := fmt.Sprintf("%s/%s", folder, fileName)
	result, err := ss.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.App.Config.AwsS3.BucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from S3: %v", err)
	}
	defer result.Body.Close()
	return io.ReadAll(result.Body)
}

func (s *S3Storage) StoreData(posts *models.Posts) error {
	jsonData, err := json.MarshalIndent(posts, "", "  ")
	if err != nil {
		s.App.Logger.Error("JSON marshal error: %v", zap.Error(err))
		return err
	}

	ss, err := s.getAwsS3Session()
	if err != nil {
		s.App.Logger.Error("S3 Session error: %v", zap.Error(err))
		return err
	}
	folder := s.App.Config.AwsS3.Folder
	savedAs := util.Sanitize(time.Now().UTC().Format("2025-01-02T15-04-05"))
	key := fmt.Sprintf("%s/%s", folder, savedAs)
	_, err = ss.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.App.Config.AwsS3.BucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(jsonData),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}
	s.App.Logger.Info("Data stored in S3 with file name:", zap.String("file_name", savedAs))
	return nil
}
