package storage_client

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/RandySteven/CafeConnect/be/configs"
	"io"
	"log"
	"mime/multipart"
)

type (
	GoogleStorage interface {
		UploadFile(ctx context.Context, filePath, fileName string, fileRequest io.Reader, fileHeader *multipart.FileHeader, width, height uint) (resultPath string, err error)
		DeleteFile(ctx context.Context, objectFileName string) (err error)
	}

	googleStorage struct {
		c   *storage.Client
		bkt *storage.BucketHandle
	}
)

var _ GoogleStorage = &googleStorage{}

func NewGoogleStorage(config *configs.Config) (*googleStorage, error) {
	ctx := context.TODO()
	c, err := storage.NewClient(ctx)

	if err != nil {
		log.Fatalln(`error google storage : `, err)
		return nil, err
	}
	bkt := c.Bucket(config.Config.Storage.BucketName)
	return &googleStorage{
		c:   c,
		bkt: bkt,
	}, nil
}
