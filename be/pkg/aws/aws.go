package aws_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"mime/multipart"
)

type (
	AWS interface {
		ListBucket() (result *s3.ListBucketsOutput, err error)
		UploadImageFile(ctx context.Context, fileRequest io.Reader, filePath string, fileHeader *multipart.FileHeader, width, height uint) (resultLocation *string, err error)
		CreateBucket(name string) error
		UploadFileToS3(ctx context.Context, fileName, path string) (string, error)
	}

	awsClient struct {
		session *session.Session
		s3      *s3.S3
	}
)

var _ AWS = &awsClient{}

func NewAWS(configYml *configs.Config) (*awsClient, error) {
	if configYml == nil || &configYml.Config.AWS == nil {
		return nil, errors.New("AWS configuration is required")
	}

	awsCfg := configYml.Config.AWS
	if awsCfg.AccessKeyID == "" || awsCfg.SecretAccessKey == "" || awsCfg.Region == "" {
		return nil, errors.New("AWS access key ID, secret access key, and region are required")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsCfg.Region),
		Credentials: credentials.NewStaticCredentials(
			awsCfg.AccessKeyID,
			awsCfg.SecretAccessKey,
			"",
		),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &awsClient{
		session: sess,
		s3:      s3.New(sess),
	}, nil
}
