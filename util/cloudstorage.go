package util

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"

	"github.com/thoniwutr/-schedule-school-teachning-bsd13-backend/constant"
)

type CloudStorageManager interface {
	UploadFile(fileName string, data []byte) (string, error)
}

type cloudStorage struct {
	client     *storage.Client
	bucketName string
}

// NewCloudStorage defined to init cloud storage with bucket name
func NewCloudStorage(bucketName string) (*cloudStorage, error) {
	ctx := context.Background()
	// Create client as usual.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage new client: %w", err)
	}

	return &cloudStorage{client: client, bucketName: bucketName}, nil
}

// UploadFile upload file to cloud storage with specific path
func (cs *cloudStorage) UploadFile(fileName string, data []byte) (string, error) {
	ctx := context.Background()

	bucket := cs.client.Bucket(cs.bucketName)

	// Check if bucket exists
	if _, err := bucket.Attrs(ctx); err != nil {
		return "", fmt.Errorf("failed to connect to cloud storage: %w", err)
	}

	// Check if file already exists
	_, err := bucket.Object(fileName).Attrs(ctx)
	switch err {
	case storage.ErrObjectNotExist:
		// object does not already exists, continue execution
	case nil:
		// object already exists in cloud storage
		return "", constant.ErrConflict
	default:
		return "", fmt.Errorf("error checking file in cloud storage: %w", err)
	}

	// Upload an object with storage writer
	wc := bucket.Object(fileName).NewWriter(ctx)
	defer func() {
		_ = wc.Close()
	}()
	wc.ContentType = "application/zip"

	if _, err := wc.Write(data); err != nil {
		return "", fmt.Errorf("failed to write data to storage with error : %v ", err)
	}

	// Generate download URL with unlimited access
	path := fmt.Sprintf("https://storage.cloud.google.com/%v/%v", cs.bucketName, fileName)

	return path, nil

}
