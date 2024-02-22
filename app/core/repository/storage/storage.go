package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/suraboy/upload-file-worksho/app/core/config"
	"log"
	"time"
)

type GcpStorageRepository interface {
	UploadFile(ctx context.Context, path string, fileName string, value []byte) error
	GetSignedURL(name string, expire time.Duration) (output string, err error)
}

// GCSStorage implements the StorageRepository interface.
type GCSStorage struct {
	AppConfig *config.AppConfig
	client    *storage.Client
}

type Config struct {
	AppConfig *config.AppConfig
	Storage   *storage.Client
}

// NewGCSStorage creates a new GCSStorage instance.
func NewGCSStorage(cfg Config) GcpStorageRepository {
	return &GCSStorage{
		AppConfig: cfg.AppConfig,
		client:    cfg.Storage,
	}
}

// UploadFile uploads a file to Google Cloud Storage.
func (s *GCSStorage) UploadFile(ctx context.Context, filepath string, fileName string, fileContent []byte) error {
	bucketName := s.AppConfig.Storage.Services[0].Bucket
	objectName := filepath + fileName

	// Create a new GCS object writer
	wc := s.client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer wc.Close()

	// Write the file content to the object writer
	if _, err := wc.Write(fileContent); err != nil {
		return fmt.Errorf("failed to write file content: %v", err)
	}

	// Close the object writer
	if err := wc.Close(); err != nil {
		return fmt.Errorf("failed to close object writer: %v", err)
	}

	log.Printf("File uploaded successfully to bucket: %s, object: %s", bucketName, objectName)
	
	return nil
}

func (s *GCSStorage) GetSignedURL(fileName string, expire time.Duration) (output string, err error) {
	bucketName := s.AppConfig.Storage.Services[0].Bucket
	name := fmt.Sprintf("branch/inbound/%s/", "20231222") + fileName

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(expire),
	}

	output, errSignedUrl := s.client.Bucket(bucketName).SignedURL(name, opts)
	if err != nil {
		err = fmt.Errorf("Bucket(%q).SignedURL: %v", bucketName, errSignedUrl)
	}

	return output, nil
}
