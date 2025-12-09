package minio

import (
	"context"
	"mime/multipart"
	"time"

	"speech-processing-service/internal/config"
	"speech-processing-service/internal/errs"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	urlExpirationTime = time.Second * 24 * 60 * 60 * 7
)

type Minio struct {
	client *minio.Client

	imagesBucket  string
	answersBucket string
}

func New(cfg *config.Minio) (Minio, error) {
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return Minio{}, errs.New(errs.ErrInitialization, "minio:"+err.Error())
	}

	return Minio{
		client:        minioClient,
		imagesBucket:  cfg.ImagesBucket,
		answersBucket: cfg.AnswersBucket,
	}, nil
}

func (m *Minio) GenerateUrl(ctx context.Context, imagePath string, isAnswer bool) (string, error) {
	var bucketName string = m.imagesBucket

	if isAnswer {
		bucketName = m.answersBucket
	}

	presignedURL, err := m.client.PresignedGetObject(
		ctx,
		bucketName,
		imagePath,
		urlExpirationTime,
		nil)
	if err != nil {
		return "", errs.New(errs.ErrMinio, "m.client.PresignedGetObject:"+err.Error())
	}
	return presignedURL.String(), nil
}

func (m *Minio) UploadAnswer(
	ctx context.Context,
	filename string,
	file multipart.File,
	size int64) error {
	_, err := m.client.PutObject(
		ctx,
		m.answersBucket,
		filename,
		file,
		size,
		minio.PutObjectOptions{},
	)
	if err != nil {
		return errs.New(errs.ErrMinio, "m.client.PutObject:"+err.Error())
	}
	return nil
}

// UploadFile uploads a file to images bucket and returns the path
func (m *Minio) UploadFile(ctx context.Context, file *multipart.File, header *multipart.FileHeader, folder string) (string, error) {
	// Генерируем путь к файлу
	filename := folder + "/" + header.Filename

	_, err := m.client.PutObject(
		ctx,
		m.imagesBucket,
		filename,
		*file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		},
	)
	if err != nil {
		return "", errs.New(errs.ErrMinio, "m.client.PutObject: "+err.Error())
	}

	return filename, nil
}

// GenerateURL generates a presigned URL for a file in images bucket
func (m *Minio) GenerateURL(ctx context.Context, filename string) (string, error) {
	return m.GenerateUrl(ctx, filename, false)
}
