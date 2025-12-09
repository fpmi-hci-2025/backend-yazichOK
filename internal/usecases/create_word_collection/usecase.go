package create_word_collection

import (
	"context"
	"mime/multipart"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"
)

type StorageProvider interface {
	CreateWordCollection(ctx context.Context, userID int, name, imagePath string) (storage.WordCollection, error)
}

type ImageUploader interface {
	UploadFile(ctx context.Context, file *multipart.File, header *multipart.FileHeader, folder string) (string, error)
}

type URLGetter interface {
	GenerateURL(ctx context.Context, filename string) (string, error)
}

type UseCase struct {
	storage       StorageProvider
	imageUploader ImageUploader
	urlGetter     URLGetter
}

func New(storage StorageProvider, imageUploader ImageUploader, urlGetter URLGetter) UseCase {
	return UseCase{
		storage:       storage,
		imageUploader: imageUploader,
		urlGetter:     urlGetter,
	}
}

func (u *UseCase) CreateCollection(
	ctx context.Context,
	userID int,
	name string,
	imageFile *multipart.File,
	imageHeader *multipart.FileHeader,
) (entity.WordCollection, error) {
	var imagePath string
	var imageURL string

	// Если есть изображение - загружаем в MinIO
	if imageFile != nil && imageHeader != nil {
		uploadedPath, err := u.imageUploader.UploadFile(ctx, imageFile, imageHeader, "collections")
		if err != nil {
			return entity.WordCollection{}, errs.New(errs.ErrUseCaseExecution, "u.imageUploader.UploadFile: "+err.Error())
		}
		imagePath = uploadedPath

		// Генерируем URL для доступа
		url, err := u.urlGetter.GenerateURL(ctx, imagePath)
		if err != nil {
			return entity.WordCollection{}, errs.New(errs.ErrUseCaseExecution, "u.urlGetter.GenerateURL: "+err.Error())
		}
		imageURL = url
	}

	// Создаем коллекцию в БД
	collection, err := u.storage.CreateWordCollection(ctx, userID, name, imagePath)
	if err != nil {
		return entity.WordCollection{}, errs.New(errs.ErrUseCaseExecution, "u.storage.CreateWordCollection: "+err.Error())
	}

	return entity.WordCollection{
		ID:                collection.ID.String(),
		UserID:            collection.UserID,
		Name:              collection.Name,
		ImageURL:          imageURL,
		TotalWordsCount:   collection.TotalWordsCount,
		LearnedWordsCount: collection.LearnedWordsCount,
		CurrentStreakDays: collection.CurrentStreakDays,
		LongestStreakDays: collection.LongestStreakDays,
		LastStudiedAt:     collection.LastStudiedAt,
		CreatedAt:         collection.CreatedAt,
		UpdatedAt:         collection.UpdatedAt,
	}, nil
}
