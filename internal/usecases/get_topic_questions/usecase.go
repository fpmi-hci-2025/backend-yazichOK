package get_topic_questions

import (
	"context"

	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"

	"go.uber.org/zap"
)

type QuestionsGetter interface {
	GetQuestionsByTopicID(ctx context.Context, topicID int) ([]storage.Question, error)
}

type Usecase struct {
	logger *zap.Logger

	questionsGetter QuestionsGetter
}

func New(
	logger *zap.Logger,
	questionsGetter QuestionsGetter,
) Usecase {
	return Usecase{
		logger:          logger,
		questionsGetter: questionsGetter,
	}
}

func (u *Usecase) GetTopicQuestions(ctx context.Context, topicID int) ([]entity.Question, error) {
	dbQuestions, err := u.questionsGetter.GetQuestionsByTopicID(ctx, topicID)
	if err != nil {
		return nil, errs.Wrap("u.questionsGetter.GetQuestionsByTopicID", err)
	}

	if len(dbQuestions) == 0 {
		return nil, errs.New(errs.ErrNotFound, "questions not found")
	}

	result := make([]entity.Question, len(dbQuestions))
	for i, dbQuestion := range dbQuestions {
		result[i] = entity.Question{
			ID:   dbQuestion.ID,
			Text: dbQuestion.Question,
		}
	}

	return result, nil
}
