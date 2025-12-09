package add_word_to_collection

import (
	"context"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"

	"github.com/google/uuid"
)

type StorageProvider interface {
	AddWordToCollection(ctx context.Context, collectionID, word, translation string, example *string) (storage.UserWord, error)
}

type UseCase struct {
	storage StorageProvider
}

func New(storage StorageProvider) UseCase {
	return UseCase{
		storage: storage,
	}
}

func (u *UseCase) AddWord(ctx context.Context, collectionID, word, translation string, example *string, userID int) (entity.UserWord, error) {
	// Валидация UUID коллекции
	if _, err := uuid.Parse(collectionID); err != nil {
		return entity.UserWord{}, errs.New(errs.ErrUseCaseExecution, "uuid.Parse: "+err.Error())
	}

	// TODO: Проверить, что коллекция принадлежит пользователю (добавить метод в storage)

	// Добавляем слово в коллекцию
	userWord, err := u.storage.AddWordToCollection(ctx, collectionID, word, translation, example)
	if err != nil {
		return entity.UserWord{}, errs.New(errs.ErrUseCaseExecution, "u.storage.AddWordToCollection: "+err.Error())
	}

	return entity.UserWord{
		ID:             userWord.ID,
		CollectionID:   userWord.CollectionID,
		Word:           userWord.Word,
		Translation:    userWord.Translation,
		Example:        userWord.Example,
		NextReviewDate: userWord.NextReviewDate,
		ReviewCount:    userWord.ReviewCount,
		CreatedAt:      userWord.CreatedAt,
		UpdatedAt:      userWord.UpdatedAt,
	}, nil
}
