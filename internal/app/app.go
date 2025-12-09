package app

import (
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"speech-processing-service/internal/config"
	"speech-processing-service/internal/entity"

	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger"
)

type AllTopicsGetter interface {
	GetAllTopics(ctx context.Context) ([]entity.Topic, error)
}

type QuestionsGetter interface {
	GetTopicQuestions(ctx context.Context, topicID int) ([]entity.Question, error)
}

type SessionsCreator interface {
	StartSession(ctx context.Context, sessionID string, topicID int) (entity.Session, error)
}

type SessionCompleter interface {
	CompleteSession(ctx context.Context, sessionID string) (entity.AnalyzeTextResult, error)
}

type AnswerAttacher interface {
	AttachAnswerToSession(
		ctx context.Context,
		sessionID string,
		questionID int,
		file *multipart.File,
		header *multipart.FileHeader,
	) error
}

type ArticlesGetter interface {
	GetArticles(ctx context.Context, limit, offset int) ([]entity.ArticlePreview, error)
}

type ArticleByIDGetter interface {
	GetArticle(ctx context.Context, id int) (entity.Article, error)
}

type WordCollectionCreator interface {
	CreateCollection(ctx context.Context, userID int, name string, imageFile *multipart.File, imageHeader *multipart.FileHeader) (entity.WordCollection, error)
}

type WordCollectionDeleter interface {
	DeleteCollection(ctx context.Context, collectionID string, userID int) error
}

type UserCollectionsGetter interface {
	GetCollections(ctx context.Context, userID int) ([]entity.WordCollection, error)
}

type CollectionDetailGetter interface {
	GetCollectionDetail(ctx context.Context, collectionID string, userID int) (entity.WordCollectionDetail, error)
}

type WordAdder interface {
	AddWord(ctx context.Context, collectionID, word, translation string, example *string, userID int) (entity.UserWord, error)
}

type App struct {
	server *http.Server
	mux    *http.ServeMux

	topicsGetter          AllTopicsGetter
	questionsGetter       QuestionsGetter
	sessionsCreator       SessionsCreator
	answerAttacher        AnswerAttacher
	sessionCompleter      SessionCompleter
	getArticlesUC         ArticlesGetter
	getArticleByIDUC      ArticleByIDGetter
	createCollectionUC    WordCollectionCreator
	deleteCollectionUC    WordCollectionDeleter
	getUserCollectionsUC  UserCollectionsGetter
	getCollectionDetailUC CollectionDetailGetter
	addWordToCollectionUC WordAdder

	cfg    *config.Config
	logger *zap.Logger
}

func New(
	topicsGetter AllTopicsGetter,
	questionsGetter QuestionsGetter,
	sessionsCreator SessionsCreator,
	answerAttacher AnswerAttacher,
	completer SessionCompleter,
	getArticlesUC ArticlesGetter,
	getArticleByIDUC ArticleByIDGetter,
	createCollectionUC WordCollectionCreator,
	deleteCollectionUC WordCollectionDeleter,
	getUserCollectionsUC UserCollectionsGetter,
	getCollectionDetailUC CollectionDetailGetter,
	addWordToCollectionUC WordAdder,
	cfg *config.Config,
	logger *zap.Logger,
) App {
	mux := http.NewServeMux()

	// TODO: Add http.Server configuration here
	server := http.Server{
		Addr:    cfg.HTTPPort,
		Handler: mux,
	}

	return App{
		server:                &server,
		mux:                   mux,
		topicsGetter:          topicsGetter,
		questionsGetter:       questionsGetter,
		sessionsCreator:       sessionsCreator,
		answerAttacher:        answerAttacher,
		sessionCompleter:      completer,
		getArticlesUC:         getArticlesUC,
		getArticleByIDUC:      getArticleByIDUC,
		createCollectionUC:    createCollectionUC,
		deleteCollectionUC:    deleteCollectionUC,
		getUserCollectionsUC:  getUserCollectionsUC,
		getCollectionDetailUC: getCollectionDetailUC,
		addWordToCollectionUC: addWordToCollectionUC,
		cfg:                   cfg,
		logger:                logger,
	}
}

func (a *App) Run() {
	a.logger.Info("service starts", zap.String("addr", a.cfg.HTTPPort))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Fatal("listen and serve", zap.Error(err))
		}
	}()

	<-stop
	a.logger.Info("server shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Fatal("server shutdown", zap.Error(err))
	}

	wg.Wait()
	a.logger.Info("Сервер остановлен")
}

func (s *App) InitREST() {
	// Swagger UI - use JSON format (YAML has formatting issues)
	s.mux.HandleFunc("GET /swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	s.mux.HandleFunc("GET /topics", s.getAllTopics())
	s.mux.HandleFunc("GET /topics/{topicID}/questions", s.getTopicQuestions())

	s.mux.HandleFunc("POST /sessions", s.startSession())
	s.mux.HandleFunc("POST /sessions/{sessionID}/answer", s.attachAnswerToSession())
	s.mux.HandleFunc("POST /sessions/{sessionID}/complete", s.completeSession())

	s.mux.HandleFunc("GET /articles", s.getArticles())
	s.mux.HandleFunc("GET /articles/{id}", s.getArticleByID())

	s.mux.HandleFunc("GET /collections", s.getUserCollections())
	s.mux.HandleFunc("GET /collections/{id}", s.getCollectionDetail())
	s.mux.HandleFunc("POST /collections", s.createWordCollection())
	s.mux.HandleFunc("DELETE /collections/{id}", s.deleteWordCollection())
	s.mux.HandleFunc("POST /collections/{id}/words", s.addWordToCollection())
}
