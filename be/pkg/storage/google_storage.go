package storage_client

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/RandySteven/CafeConnect/be/configs"
	"github.com/RandySteven/CafeConnect/be/utils"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type (
	GoogleStorage interface {
		UploadFile(ctx context.Context, filePath string, fileRequest io.Reader, fileHeader *multipart.FileHeader, width, height uint) (resultPath string, err error)
	}

	googleStorage struct {
		c   *storage.Client
		bkt *storage.BucketHandle
	}
)

func (g *googleStorage) UploadFile(ctx context.Context, filePath string, fileRequest io.Reader, fileHeader *multipart.FileHeader, width, height uint) (resultPath string, err error) {
	err = os.MkdirAll("./temp-images", os.ModePerm)
	if err != nil {
		return "", err
	}

	if err = ctx.Err(); err != nil {
		return "", err
	}

	tempFile, err := os.CreateTemp("./temp-images", "upload-*.png")
	if err != nil {
		return "", fmt.Errorf("error during create temp file : %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	imageFile, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("open file header issue to : %w", err)
	}
	defer imageFile.Close()

	fileBytes, err := ioutil.ReadAll(imageFile)
	if err != nil {
		return "", err
	}
	if err = ctx.Err(); err != nil {
		return "", fmt.Errorf("operation canceled before writing to temp file: %w", err)
	}

	if _, err = tempFile.Write(fileBytes); err != nil {
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	if err = ctx.Err(); err != nil {
		return "", fmt.Errorf("operation canceled before resizing image: %w", err)
	}

	err = utils.ResizeImage(tempFile.Name(), tempFile.Name(), width, height)
	if err != nil {
		return "", err
	}

	fileExt := filepath.Ext(fileHeader.Filename)
	renamedImage := uuid.NewString() + fileExt
	if err = ctx.Err(); err != nil {
		return "", fmt.Errorf("operation canceled before opening resized file: %w", err)
	}

	file, err := os.Open(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err = ctx.Err(); err != nil {
		return "", fmt.Errorf("operation canceled before uploading to storage: %w", err)
	}

	fileWriter := g.bkt.Object(renamedImage).NewWriter(ctx)
	msg := `this is a test`
	if _, err = fileWriter.Write([]byte(msg)); err != nil {
		log.Fatal(err)
	}
	if err = fileWriter.Close(); err != nil {
		log.Fatalf("closing writer: %v", err)
	}

	return fileWriter.Name, nil
}

var _ GoogleStorage = &googleStorage{}

func NewGoogleStorage(config *configs.Config) (*googleStorage, error) {
	ctx := context.TODO()
	c, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	bkt := c.Bucket(``)
	return &googleStorage{
		c:   c,
		bkt: bkt,
	}, nil
}
