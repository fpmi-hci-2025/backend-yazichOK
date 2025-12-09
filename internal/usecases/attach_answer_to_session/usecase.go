package attach_answer_to_session

import (
	"context"
	"mime/multipart"

	"speech-processing-service/internal/errs"

	"go.uber.org/zap"
)

type AnswerUploader interface {
	UploadAnswer(
		ctx context.Context,
		filename string,
		file multipart.File,
		size int64,
	) error
}

type AnswerCreator interface {
	CreateAnswer(
		ctx context.Context,
		sessionID string,
		questionID int,
		filename string,
	) error
}

type UseCase struct {
	logger *zap.Logger

	answerUploader AnswerUploader
	answerCreator  AnswerCreator
}

func New(
	logger *zap.Logger,
	answerUploader AnswerUploader,
	creator AnswerCreator,
) UseCase {
	return UseCase{
		logger:         logger,
		answerUploader: answerUploader,
		answerCreator:  creator,
	}
}

func (u *UseCase) AttachAnswerToSession(
	ctx context.Context,
	sessionID string,
	questionID int,
	file *multipart.File,
	header *multipart.FileHeader,
) error {
	err := u.answerUploader.UploadAnswer(ctx, header.Filename, *file, header.Size)
	if err != nil {
		return errs.Wrap("u.answerUploader.UploadAnswer", err)
	}

	err = u.answerCreator.CreateAnswer(ctx, sessionID, questionID, header.Filename)
	if err != nil {
		return errs.Wrap("u.answerCreator.CreateAnswer", err)
	}

	return nil
}
