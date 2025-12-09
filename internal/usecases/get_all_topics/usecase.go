package get_all_topics

import (
	"context"

	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"

	"go.uber.org/zap"
)

type TopicsGetter interface {
	GetAllTopics(ctx context.Context) ([]storage.Topic, error)
}

type URLGetter interface {
	GenerateUrl(ctx context.Context, imagePath string, isAnswer bool) (string, error)
}

type Usecase struct {
	logger *zap.Logger

	topicsGetter TopicsGetter
	urlGetter    URLGetter
}

func New(
	logger *zap.Logger,
	topicsGetter TopicsGetter,
	urlGetter URLGetter,
) Usecase {
	return Usecase{
		logger:       logger,
		topicsGetter: topicsGetter,
		urlGetter:    urlGetter,
	}
}

func (u *Usecase) GetAllTopics(ctx context.Context) ([]entity.Topic, error) {
	dbTopics, err := u.topicsGetter.GetAllTopics(ctx)
	if err != nil {
		return nil, errs.Wrap("u.topicsGetter.GetAllTopics", err)
	}

	result := make([]entity.Topic, 0, len(dbTopics))
	for _, dbTopic := range dbTopics {
		photoURL, err := u.urlGetter.GenerateUrl(ctx, dbTopic.ImagePath, false)
		if err != nil {
			u.logger.Error("u.urlGetter.GenerateURl", zap.Error(err))

			continue
		}

		result = append(result, entity.Topic{
			ID:          dbTopic.ID,
			Title:       dbTopic.Title,
			Description: dbTopic.Description,
			PhotoURL:    photoURL,
			Questions:   nil,
		})
	}
	return result, nil
}
