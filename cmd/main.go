package main

import (
	_ "speech-processing-service/docs"
	"speech-processing-service/internal/app"
	"speech-processing-service/internal/config"
	"speech-processing-service/internal/drivers/apis/deepgram"
	"speech-processing-service/internal/drivers/apis/gemini"
	"speech-processing-service/internal/drivers/storage"
	"speech-processing-service/internal/drivers/tools/minio"
	"speech-processing-service/internal/errs"
	"speech-processing-service/internal/usecases/add_word_to_collection"
	"speech-processing-service/internal/usecases/attach_answer_to_session"
	"speech-processing-service/internal/usecases/create_word_collection"
	"speech-processing-service/internal/usecases/delete_word_collection"
	"speech-processing-service/internal/usecases/get_all_topics"
	"speech-processing-service/internal/usecases/get_article_by_id"
	"speech-processing-service/internal/usecases/get_articles"
	"speech-processing-service/internal/usecases/get_collection_detail"
	"speech-processing-service/internal/usecases/get_topic_questions"
	"speech-processing-service/internal/usecases/get_user_collections"
	"speech-processing-service/internal/usecases/session_completer"
	"speech-processing-service/internal/usecases/start_session"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
)

type drivers struct {
	storage  *storage.Storage
	minio    *minio.Minio
	deepgram *deepgram.Deepgram
	gemini   *gemini.Gemini
}

func newDrivers(cfg *config.Config) (drivers, error) {
	storage, err := storage.New(cfg.Postgres)
	if err != nil {
		return drivers{}, errs.Wrap("storage.New", err)
	}

	minio, err := minio.New(cfg.Minio)
	if err != nil {
		return drivers{}, errs.Wrap("minio.New", err)
	}

	deepgram := deepgram.New(cfg.Deepgram)

	gemini := gemini.New(cfg.Gemini)

	return drivers{
		storage:  &storage,
		minio:    &minio,
		deepgram: &deepgram,
		gemini:   &gemini,
	}, nil
}

type UseCases struct {
	allTopicsGetter       *get_all_topics.Usecase
	topicsQuestionsGetter *get_topic_questions.Usecase
	sessionStarter        *start_session.Usecase
	answerAttacher        *attach_answer_to_session.UseCase
	sessionCompleter      *session_completer.UseCase
	articlesGetter        *get_articles.UseCase
	articleByIDGetter     *get_article_by_id.UseCase
	createWordCollection  *create_word_collection.UseCase
	deleteWordCollection  *delete_word_collection.UseCase
	getUserCollections    *get_user_collections.UseCase
	getCollectionDetail   *get_collection_detail.UseCase
	addWordToCollection   *add_word_to_collection.UseCase
}

func newUseCases(logger *zap.Logger, drivers *drivers) UseCases {
	allTopicsGetter := get_all_topics.New(logger, drivers.storage, drivers.minio)
	topicsQuestionsGetter := get_topic_questions.New(logger, drivers.storage)
	sessionStarter := start_session.New(logger, drivers.storage, drivers.storage)
	answerAttacher := attach_answer_to_session.New(logger, drivers.minio, drivers.storage)
	sessionCompleter := session_completer.New(logger, drivers.storage, drivers.minio, drivers.deepgram, drivers.gemini)
	articlesGetter := get_articles.New(drivers.storage, drivers.minio)
	articleByIDGetter := get_article_by_id.New(drivers.storage, drivers.minio)
	createWordCollection := create_word_collection.New(drivers.storage, drivers.minio, drivers.minio)
	deleteWordCollection := delete_word_collection.New(drivers.storage)
	getUserCollections := get_user_collections.New(drivers.storage, drivers.minio)
	getCollectionDetail := get_collection_detail.New(drivers.storage, drivers.minio)
	addWordToCollection := add_word_to_collection.New(drivers.storage)

	return UseCases{
		allTopicsGetter:       &allTopicsGetter,
		topicsQuestionsGetter: &topicsQuestionsGetter,
		sessionStarter:        &sessionStarter,
		answerAttacher:        &answerAttacher,
		sessionCompleter:      &sessionCompleter,
		articlesGetter:        &articlesGetter,
		articleByIDGetter:     &articleByIDGetter,
		createWordCollection:  &createWordCollection,
		deleteWordCollection:  &deleteWordCollection,
		getUserCollections:    &getUserCollections,
		getCollectionDetail:   &getCollectionDetail,
		addWordToCollection:   &addWordToCollection,
	}
}

// @title Speech Processing Service API
// @version 1.0
// @description This is a sample server for speech processing.
func main() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.DisableStacktrace = true

	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	if err := godotenv.Load(); err != nil {
		logger.Warn(".env file not found, using environment variables from system", zap.Error(err))
	} else {
		logger.Info(".env file loaded successfully")
	}

	cfg := config.New()

	drivers, err := newDrivers(&cfg)
	if err != nil {
		logger.Error("new drivers", zap.Error(err))

		return
	}

	usecases := newUseCases(logger, &drivers)

	application := app.New(
		usecases.allTopicsGetter,
		usecases.topicsQuestionsGetter,
		usecases.sessionStarter,
		usecases.answerAttacher,
		usecases.sessionCompleter,
		usecases.articlesGetter,
		usecases.articleByIDGetter,
		usecases.createWordCollection,
		usecases.deleteWordCollection,
		usecases.getUserCollections,
		usecases.getCollectionDetail,
		usecases.addWordToCollection,
		&cfg,
		logger,
	)
	// TODO: fix this
	application.InitREST()
	application.Run()
}
