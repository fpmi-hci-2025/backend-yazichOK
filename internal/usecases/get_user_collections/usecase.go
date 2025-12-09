package get_user_collections

import (
	"context"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"
)

type StorageProvider interface {
	GetUserCollections(ctx context.Context, userID int) ([]storage.WordCollection, error)
}

type URLGetter interface {
	GenerateURL(ctx context.Context, filename string) (string, error)
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

func (u *UseCase) GetCollections(ctx context.Context, userID int) ([]entity.WordCollection, error) {
	collections, err := u.storage.GetUserCollections(ctx, userID)
	if err != nil {
		return nil, errs.New(errs.ErrUseCaseExecution, "u.storage.GetUserCollections: "+err.Error())
	}

	result := make([]entity.WordCollection, 0, len(collections))
	for _, collection := range collections {
		var imageURL string
		if collection.ImagePath != "" {
			url, err := u.urlGetter.GenerateURL(ctx, collection.ImagePath)
			if err != nil {
				imageURL = ""
			} else {
				imageURL = url
			}
		}

		result = append(result, entity.WordCollection{
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
		})
	}

	return result, nil
}
