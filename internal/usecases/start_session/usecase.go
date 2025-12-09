package start_session

import (
	"context"

	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/entity"
	"speech-processing-service/internal/errs"

	"go.uber.org/zap"
)

type SessionsCreator interface {
	CreateSession(ctx context.Context, sessionID string, topicID int) error
}

type QuestionsGetter interface {
	GetQuestionsByTopicID(ctx context.Context, topicID int) ([]storage.Question, error)
}

type Usecase struct {
	logger *zap.Logger

	sessionsCreator SessionsCreator
	questionsGetter QuestionsGetter
}

func New(
	logger *zap.Logger,
	sessionsCreator SessionsCreator,
	questionsGetter QuestionsGetter,
) Usecase {
	return Usecase{
		logger:          logger,
		sessionsCreator: sessionsCreator,
		questionsGetter: questionsGetter,
	}
}

func (u *Usecase) StartSession(ctx context.Context, sessionID string, topicID int) (entity.Session, error) {
	err := u.sessionsCreator.CreateSession(ctx, sessionID, topicID)
	if err != nil {
		return entity.Session{}, errs.Wrap("u.sessionsCreator.CreateSession", err)
	}

	questionsDB, err := u.questionsGetter.GetQuestionsByTopicID(ctx, topicID)
	if err != nil {
		return entity.Session{}, errs.Wrap("u.questionsGetter.GetQuestionsByTopicID", err)
	}

	if len(questionsDB) == 0 {
		return entity.Session{}, errs.New(errs.ErrNotFound, "questions not found")
	}

	questions := make([]entity.Question, len(questionsDB))
	for i, dbQuestion := range questionsDB {
		questions[i] = entity.Question{
			ID:   dbQuestion.ID,
			Text: dbQuestion.Question,
		}

	}

	return entity.Session{
		ID:        sessionID,
		TopicID:   topicID,
		Questions: questions,
	}, nil
}
