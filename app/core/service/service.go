package service

import (
	"context"
	"fmt"
	"github.com/suraboy/upload-file-worksho/app/core/config"
	"github.com/suraboy/upload-file-worksho/app/core/repository/mail"
	"github.com/suraboy/upload-file-worksho/app/core/repository/storage"
	"time"
)

type FileService interface {
	UploadFileService(ctx context.Context, fileName string, filePath []byte) (string, error)
	GetFileSignedURLService(fileName string) (string, error)
}

type Config struct {
	AppConfig *config.AppConfig
	Storage   storage.GcpStorageRepository
	Email     mail.SendgridMailRepository
}

type Service struct {
	appConfig *config.AppConfig
	storage   storage.GcpStorageRepository
	email     mail.SendgridMailRepository
}

// UploadFileService handles the file upload logic
func (s Service) UploadFileService(ctx context.Context, fileName string, fileContent []byte) (string, error) {
	path := fmt.Sprintf("upload/%s/", time.Now().Format("20060102"))
	err := s.storage.UploadFile(ctx, path, fileName, fileContent)

	if err != nil {
		return "", err
	}

	err = s.email.SendEmail("sirichai.jann@gmail.com", "sirichai.j@arise.tech")
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (s Service) GetFileSignedURLService(fileName string) (string, error) {
	return s.storage.GetSignedURL(fileName, 1*time.Hour)
}

func NewService(cfg Config) FileService {
	return &Service{
		storage:   cfg.Storage,
		email:     cfg.Email,
		appConfig: cfg.AppConfig,
	}
}
