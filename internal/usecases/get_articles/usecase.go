package get_articles

import (
	"context"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"
)

type StorageProvider interface {
	GetArticles(ctx context.Context, limit, offset int) ([]storage.Article, error)
}

type URLGetter interface {
	GenerateUrl(ctx context.Context, imagePath string, isAnswer bool) (string, error)
}

type UseCase struct {
	storage   StorageProvider
	urlGetter URLGetter
}

func New(
	storage StorageProvider,
	urlGetter URLGetter,
) UseCase {
	return UseCase{
		storage:   storage,
		urlGetter: urlGetter,
	}
}

func (u *UseCase) GetArticles(ctx context.Context, limit, offset int) ([]entity.ArticlePreview, error) {
	articles, err := u.storage.GetArticles(ctx, limit, offset)
	if err != nil {
		return nil, errs.New(errs.ErrUseCaseExecution, "u.storage.GetArticles: "+err.Error())
	}

	result := make([]entity.ArticlePreview, 0, len(articles))
	for _, article := range articles {
		photoURL, err := u.urlGetter.GenerateUrl(ctx, article.ImageURL, false)
		if err != nil {
			return nil, errs.Wrap("u.urlGetter.GenerateURl", err)
		}

		result = append(result, entity.ArticlePreview{
			ID:            article.ID,
			ImageURL:      photoURL,
			Level:         article.Level,
			MinutesToRead: article.MinutesToRead,
			Title:         article.Title,
		})
	}

	return result, nil
}
