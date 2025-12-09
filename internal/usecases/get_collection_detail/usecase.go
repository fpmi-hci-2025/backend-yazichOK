package get_collection_detail

import (
	"context"
	"encoding/json"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"

	"github.com/google/uuid"
)

type StorageProvider interface {
	GetWordCollectionByID(ctx context.Context, collectionID string, userID int) (storage.WordCollection, error)
	GetUserWordsByCollectionID(ctx context.Context, collectionID string) ([]storage.UserWord, error)
}

type URLGetter interface {
	GenerateUrl(ctx context.Context, objectPath string, isAnswer bool) (string, error)
}

type UseCase struct {
	storage   StorageProvider
	urlGetter URLGetter
}

func New(storage StorageProvider, urlGetter URLGetter) UseCase {
	return UseCase{
		storage:   storage,
		urlGetter: urlGetter,
	}
}

func (u *UseCase) GetCollectionDetail(ctx context.Context, collectionID string, userID int) (entity.WordCollectionDetail, error) {
	// Валидация UUID
	if _, err := uuid.Parse(collectionID); err != nil {
		return entity.WordCollectionDetail{}, errs.New(errs.ErrUseCaseExecution, "uuid.Parse: "+err.Error())
	}

	// Получение коллекции
	collection, err := u.storage.GetWordCollectionByID(ctx, collectionID, userID)
	if err != nil {
		return entity.WordCollectionDetail{}, errs.New(errs.ErrUseCaseExecution, "u.storage.GetWordCollectionByID: "+err.Error())
	}

	// Генерация URL для изображения
	var imageURL string
	if collection.ImagePath != "" {
		imageURL, err = u.urlGetter.GenerateUrl(ctx, collection.ImagePath, false)
		if err != nil {
			return entity.WordCollectionDetail{}, errs.New(errs.ErrUseCaseExecution, "u.urlGetter.GenerateUrl: "+err.Error())
		}
	}

	// Получение слов пользователя
	words, err := u.storage.GetUserWordsByCollectionID(ctx, collectionID)
	if err != nil {
		return entity.WordCollectionDetail{}, errs.New(errs.ErrUseCaseExecution, "u.storage.GetUserWordsByCollectionID: "+err.Error())
	}

	// Преобразование слов
	userWords := make([]entity.UserWord, 0, len(words))
	for _, word := range words {
		userWords = append(userWords, entity.UserWord{
			ID:             word.ID,
			CollectionID:   word.CollectionID,
			Word:           word.Word,
			Translation:    word.Translation,
			Example:        word.Example,
			NextReviewDate: word.NextReviewDate,
			ReviewCount:    word.ReviewCount,
			CreatedAt:      word.CreatedAt,
			UpdatedAt:      word.UpdatedAt,
		})
	}

	// Парсинг AI рекомендаций
	var aiSuggestions []entity.AISuggestion
	if collection.AISuggestions != nil && *collection.AISuggestions != "" {
		if err := json.Unmarshal([]byte(*collection.AISuggestions), &aiSuggestions); err != nil {
			// Если не удалось распарсить, просто оставляем пустой массив
			aiSuggestions = []entity.AISuggestion{}
		}
	}

	return entity.WordCollectionDetail{
		ID:                collection.ID.String(),
		Name:              collection.Name,
		ImageURL:          imageURL,
		TotalWordsCount:   collection.TotalWordsCount,
		LearnedWordsCount: collection.LearnedWordsCount,
		CurrentStreakDays: collection.CurrentStreakDays,
		LongestStreakDays: collection.LongestStreakDays,
		LastStudiedAt:     collection.LastStudiedAt,
		UserWords:         userWords,
		AISuggestions:     aiSuggestions,
	}, nil
}
