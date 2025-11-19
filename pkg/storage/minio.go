package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"durich-be/pkg/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorage(cfg *config.MinioConfig) (*MinioStorage, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize minio client: %w", err)
	}

	storage := &MinioStorage{
		client:     minioClient,
		bucketName: cfg.BucketName,
	}

	if err := storage.ensureBucketExists(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return storage, nil
}

func (s *MinioStorage) ensureBucketExists(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *MinioStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, objectName string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(ctx, s.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to minio: %w", err)
	}

	return objectName, nil
}

func (s *MinioStorage) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file from minio: %w", err)
	}
	return nil
}

func (s *MinioStorage) GetFileURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

func (s *MinioStorage) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	object, err := s.client.GetObject(ctx, s.bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file from minio: %w", err)
	}
	return object, nil
}

func GenerateObjectName(userID, documentType, fileName string) string {
	timestamp := time.Now().Unix()
	ext := filepath.Ext(fileName)
	cleanFileName := strings.TrimSuffix(fileName, ext)

	return fmt.Sprintf("documents/%s/%s/%d_%s%s", userID, documentType, timestamp, cleanFileName, ext)
}

const (
	ContentTypePDF  = "application/pdf"
	ContentTypeJPEG = "image/jpeg"
	ContentTypePNG  = "image/png"
	ContentTypeDOC  = "application/msword"
	ContentTypeDOCX = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
)

func ValidateFileType(file *multipart.FileHeader, documentType string) error {
	contentType := file.Header.Get("Content-Type")

	allowedTypes := map[string][]string{
		"ktp":   {ContentTypePDF, ContentTypeJPEG, ContentTypePNG},
		"nib":   {ContentTypePDF, ContentTypeJPEG, ContentTypePNG},
		"npwp":  {ContentTypePDF, ContentTypeJPEG, ContentTypePNG},
		"other": {ContentTypePDF, ContentTypeJPEG, ContentTypePNG, ContentTypeDOC, ContentTypeDOCX},
	}

	types, exists := allowedTypes[documentType]
	if !exists {
		return fmt.Errorf("invalid document type")
	}

	for _, allowed := range types {
		if contentType == allowed {
			return nil
		}
	}

	return fmt.Errorf("invalid file type: allowed types for %s are %v", documentType, types)
}

func ValidateFileSize(file *multipart.FileHeader) error {
	const maxSize = 5 * 1024 * 1024 // 5MB
	if file.Size > maxSize {
		return fmt.Errorf("file size too large: maximum 5MB allowed")
	}
	return nil
}
