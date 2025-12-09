package delete_word_collection

import (
	"context"
	"speech-processing-service/internal/errs"

	"github.com/google/uuid"
)

type StorageProvider interface {
	DeleteWordCollection(ctx context.Context, collectionID uuid.UUID, userID int) error
}

type UseCase struct {
	storage StorageProvider
}

func New(storage StorageProvider) UseCase {
	return UseCase{
		storage: storage,
	}
}

func (u *UseCase) DeleteCollection(ctx context.Context, collectionID string, userID int) error {
	id, err := uuid.Parse(collectionID)
	if err != nil {
		return errs.New(errs.ErrUseCaseExecution, "uuid.Parse: "+err.Error())
	}

	if err := u.storage.DeleteWordCollection(ctx, id, userID); err != nil {
		return errs.New(errs.ErrUseCaseExecution, "u.storage.DeleteWordCollection: "+err.Error())
	}

	// TODO: Удаление изображения из MinIO можно добавить позже
	// if collection.ImagePath != "" {
	//     u.imageDeleter.DeleteFile(ctx, collection.ImagePath)
	// }

	return nil
}
