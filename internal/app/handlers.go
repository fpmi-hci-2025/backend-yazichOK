package app

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"speech-processing-service/internal/app/views"
	"speech-processing-service/internal/errs"

	"go.uber.org/zap"

	"github.com/google/uuid"
)

const (
	topicIDKey   = "topicID"
	sessionIDKey = "sessionID"

	answerKey     = "answer"
	questionIDKey = "questionID"
)

// getAllTopics godoc
// @Summary Get all topics
// @Description Get a list of all topics
// @Tags topics
// @Produce json
// @Success 200 {object} views.SuccessResponse{data=views.GetAllTopicsResponse}
// @Failure 500 {object} views.ErrorResponse{error=views.Error}
// @Router /topics [get]
func (s *App) getAllTopics() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		topics, err := s.topicsGetter.GetAllTopics(r.Context())
		if err != nil {
			s.logger.Error("handlers.getAllTopics", zap.Error(err))
		}

		views.Return(s.logger, w, r, views.NewGetAllToicsResponse(topics), err)
	}
}

// getAllTopics godoc
// @Summary Get questions by topic
// @Description Get a list of questions by topic
// @Tags topics
// @Produce json
// @Success 200 {object} views.SuccessResponse{data=views.GetTopicQuestionsResponse}
// @Failure 400 {object} views.ErrorResponse{error=views.Error}
// @Failure 404 {object} views.ErrorResponse{error=views.Error}
// @Failure 500 {object} views.ErrorResponse{error=views.Error}
// @Router /topics/{topicID}/questions [get]
func (s *App) getTopicQuestions() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		topicID := r.PathValue(topicIDKey)

		id, err := strconv.Atoi(topicID)
		if err != nil {
			s.logger.Error("handlers.getTopicQuestions", zap.Error(err))

			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeNumeric, fmt.Sprintf("topicID: %s", topicID)))

			return
		}

		questions, err := s.questionsGetter.GetTopicQuestions(r.Context(), id)
		if err != nil {
			s.logger.Error("handlers.getTopicQuestions", zap.Error(err))
		}

		views.Return(s.logger, w, r, views.NewGetTopicQuestionsResponse(questions), err)
	}
}

// startSession godoc
// @Summary Start session
// @Description Start a new session
// @Tags session
// @Produce json
// @Param session body views.StartSessionRequest true "Session data"
// @Success 200 {object} views.SuccessResponse{data=views.StartSessionResponse}
// @Failure 400 {object} views.ErrorResponse{error=views.Error}
// @Failure 404 {object} views.ErrorResponse{error=views.Error}
// @Failure 500 {object} views.ErrorResponse{error=views.Error}
// @Router /session [post]
func (s *App) startSession() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req views.StartSessionRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()

		if err := decoder.Decode(&req); err != nil {
			s.logger.Error("handlers.startSession", zap.Error(err))

			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, err.Error()))

			return
		}

		sessionID := uuid.New().String()
		session, err := s.sessionsCreator.StartSession(r.Context(), sessionID, req.TopicID)
		if err != nil {
			s.logger.Error("handlers.startSession", zap.Error(err))

			views.Return(s.logger, w, r, nil, err)

			return
		}

		views.Return(s.logger, w, r, views.NewStartSessionResponse(&session), nil)
	}
}

// attachAnswerToSession godoc
// @Summary Attach answer to session
// @Description Attach an answer to a session
// @Tags session
// @Accept multipart/form-data
// @Produce json
// @Param sessionID path string true "Session ID"
// @Param questionID formData int true "Question ID"
// @Param answer formData file true "Answer file"
// @Success 200 {object} views.SuccessResponse{data=string}
// @Failure 400 {object} views.ErrorResponse{error=views.Error}
// @Failure 404 {object} views.ErrorResponse{error=views.Error}
// @Failure 500 {object} views.ErrorResponse{error=views.Error}
// @Router /session/{sessionID}/answer [post]
func (s *App) attachAnswerToSession() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue(sessionIDKey)
		if err := uuid.Validate(sessionID); err != nil {
			s.logger.Error("handlers.attachAnswerToSession", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeUUID, fmt.Sprintf("sessionID: %s", sessionID)))
			return
		}

		questionIDString := r.FormValue(questionIDKey)
		questionID, err := strconv.Atoi(questionIDString)
		if err != nil {
			s.logger.Error(
				"handlers.attachAnswerToSession",
				zap.Error(fmt.Errorf("questionID: %s", questionIDString+err.Error())),
			)

			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeNumeric, fmt.Sprintf("questionID: %s", questionIDString)))
			return
		}

		file, header, err := r.FormFile(answerKey)
		if err != nil {
			s.logger.Error("handlers.attachAnswerToSession", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, err.Error()))
			return
		}
		defer file.Close()

		err = s.answerAttacher.AttachAnswerToSession(r.Context(), sessionID, questionID, &file, header)
		if err != nil {
			s.logger.Error("handlers.attachAnswerToSession", zap.Error(err))

			views.Return(s.logger, w, r, nil, err)

			return
		}

		views.Return(s.logger, w, r, nil, nil)
	}
}

// completeSession godoc
// @Summary Complete session
// @Description Complete a session
// @Tags session
// @Produce json
// @Param sessionID path string true "Session ID"
// @Success 200 {object} views.SuccessResponse{data=views.CompleteSessionResp}
// @Failure 400 {object} views.ErrorResponse{error=views.Error}
// @Router /session/{sessionID}/complete [post]
func (s *App) completeSession() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.PathValue(sessionIDKey)
		if err := uuid.Validate(sessionID); err != nil {
			s.logger.Error("handlers.attachAnswerToSession", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeUUID, fmt.Sprintf("sessionID: %s", sessionID)))
			return
		}

		result, err := s.sessionCompleter.CompleteSession(r.Context(), sessionID)
		if err != nil {
			s.logger.Error("handlers.completeSession", zap.Error(err))

			views.Return(s.logger, w, r, nil, err)

			return
		}

		views.Return(s.logger, w, r, views.NewCompleteSessionResp(&result), nil)
	}
}

// @Summary Get articles with pagination
// @Description Returns a list of article previews with pagination support
// @Tags articles
// @Accept json
// @Produce json
// @Param limit query int false "Number of articles to return" default(10)
// @Param offset query int false "Number of articles to skip" default(0)
// @Success 200 {object} views.SuccessResponse{data=views.ArticlesPreviewData}
// @Failure 400 {object} views.ErrorResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /articles [get]
func (s *App) getArticles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := 10
		offset := 0

		if limitParam := r.URL.Query().Get("limit"); limitParam != "" {
			if parsedLimit, err := strconv.Atoi(limitParam); err == nil && parsedLimit > 0 {
				limit = parsedLimit
			}
		}

		if offsetParam := r.URL.Query().Get("offset"); offsetParam != "" {
			if parsedOffset, err := strconv.Atoi(offsetParam); err == nil && parsedOffset >= 0 {
				offset = parsedOffset
			}
		}

		articles, err := s.getArticlesUC.GetArticles(r.Context(), limit, offset)
		if err != nil {
			s.logger.Error("handlers.getArticles", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, views.NewArticlesPreviewResponse(articles), nil)
	}
}

// @Summary Get article by ID
// @Description Returns full article details including vocabulary and grammar rules
// @Tags articles
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} views.ArticleData
// @Failure 400 {object} views.ErrorResponse
// @Failure 404 {object} views.ErrorResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /articles/{id} [get]
func (s *App) getArticleByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			s.logger.Error("handlers.getArticleByID: missing id parameter")
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeNumeric, "missing article id"))
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			s.logger.Error("handlers.getArticleByID: invalid id parameter", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeNumeric, fmt.Sprintf("article id: %s", idStr)))
			return
		}

		article, err := s.getArticleByIDUC.GetArticle(r.Context(), id)
		if err != nil {
			s.logger.Error("handlers.getArticleByID", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, views.NewArticleResponse(article), nil)
	}
}

// @Summary Create word collection
// @Description Create a new word collection with optional image
// @Tags collections
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Collection name"
// @Param image formData file false "Collection image"
// @Success 200 {object} views.CreateWordCollectionResponse
// @Failure 400 {object} views.ErrorResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /collections [post]
func (s *App) createWordCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get userID from auth context
		// For now, using hardcoded userID = 1
		userID := 1

		// Парсинг multipart form
		if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
			s.logger.Error("handlers.createWordCollection: parse multipart form", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, "invalid multipart form: "+err.Error()))
			return
		}

		// Получение имени коллекции
		name := r.FormValue("name")
		if name == "" {
			s.logger.Error("handlers.createWordCollection: missing name")
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, "name is required"))
			return
		}

		// Получение изображения (опционально)
		var imageFile *multipart.File
		var imageHeader *multipart.FileHeader
		file, header, err := r.FormFile("image")
		if err == nil {
			defer file.Close()
			imageFile = &file
			imageHeader = header
		}

		// Создание коллекции
		collection, err := s.createCollectionUC.CreateCollection(r.Context(), userID, name, imageFile, imageHeader)
		if err != nil {
			s.logger.Error("handlers.createWordCollection", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, views.NewCreateWordCollectionResponse(collection), nil)
	}
}

// @Summary Delete word collection
// @Description Delete a word collection by ID
// @Tags collections
// @Produce json
// @Param id path int true "Collection ID"
// @Success 200 {object} views.SuccessResponse
// @Failure 400 {object} views.ErrorResponse
// @Failure 404 {object} views.ErrorResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /collections/{id} [delete]
func (s *App) deleteWordCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get userID from auth context
		userID := 1

		// Получение ID коллекции (UUID)
		collectionID := r.PathValue("id")
		if collectionID == "" {
			s.logger.Error("handlers.deleteWordCollection: missing id parameter")
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, "missing collection id"))
			return
		}

		// Валидация UUID
		if err := uuid.Validate(collectionID); err != nil {
			s.logger.Error("handlers.deleteWordCollection: invalid uuid", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeUUID, fmt.Sprintf("collection id: %s", collectionID)))
			return
		}

		// Удаление коллекции
		if err := s.deleteCollectionUC.DeleteCollection(r.Context(), collectionID, userID); err != nil {
			s.logger.Error("handlers.deleteWordCollection", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, map[string]string{"message": "Collection deleted successfully"}, nil)
	}
}

// @Summary Get user collections
// @Description Get all word collections for the authenticated user
// @Tags collections
// @Produce json
// @Success 200 {object} views.GetUserCollectionsResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /collections [get]
func (s *App) getUserCollections() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get userID from auth context
		userID := 1

		collections, err := s.getUserCollectionsUC.GetCollections(r.Context(), userID)
		if err != nil {
			s.logger.Error("handlers.getUserCollections", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, views.NewGetUserCollectionsResponse(collections), nil)
	}
}

// @Summary Get collection detail
// @Description Get full information about a specific collection including user words and AI suggestions
// @Tags collections
// @Produce json
// @Param id path string true "Collection ID (UUID)"
// @Success 200 {object} views.WordCollectionDetailResponse
// @Failure 400 {object} views.ErrorResponse
// @Failure 404 {object} views.ErrorResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /collections/{id} [get]
func (s *App) getCollectionDetail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get userID from auth context
		userID := 1

		// Получение ID коллекции (UUID)
		collectionID := r.PathValue("id")
		if collectionID == "" {
			s.logger.Error("handlers.getCollectionDetail: missing id parameter")
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, "missing collection id"))
			return
		}

		// Валидация UUID
		if err := uuid.Validate(collectionID); err != nil {
			s.logger.Error("handlers.getCollectionDetail: invalid uuid", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeUUID, fmt.Sprintf("collection id: %s", collectionID)))
			return
		}

		// Получение детальной информации о коллекции
		detail, err := s.getCollectionDetailUC.GetCollectionDetail(r.Context(), collectionID, userID)
		if err != nil {
			s.logger.Error("handlers.getCollectionDetail", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, views.NewWordCollectionDetailResponse(detail), nil)
	}
}

// @Summary Add word to collection
// @Description Add a new word to a collection
// @Tags collections
// @Accept json
// @Produce json
// @Param id path string true "Collection ID (UUID)"
// @Param request body views.AddWordToCollectionRequest true "Word data"
// @Success 200 {object} views.AddWordToCollectionResponse
// @Failure 400 {object} views.ErrorResponse
// @Failure 404 {object} views.ErrorResponse
// @Failure 500 {object} views.ErrorResponse
// @Router /collections/{id}/words [post]
func (s *App) addWordToCollection() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Get userID from auth context
		userID := 1

		// Получение ID коллекции (UUID)
		collectionID := r.PathValue("id")
		if collectionID == "" {
			s.logger.Error("handlers.addWordToCollection: missing id parameter")
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, "missing collection id"))
			return
		}

		// Валидация UUID
		if err := uuid.Validate(collectionID); err != nil {
			s.logger.Error("handlers.addWordToCollection: invalid uuid", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrTypeMustBeUUID, fmt.Sprintf("collection id: %s", collectionID)))
			return
		}

		// Парсинг JSON body
		var req views.AddWordToCollectionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.logger.Error("handlers.addWordToCollection: failed to decode request", zap.Error(err))
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, err.Error()))
			return
		}

		// Валидация обязательных полей
		if req.Word == "" || req.Translation == "" {
			s.logger.Error("handlers.addWordToCollection: missing required fields")
			views.Return(s.logger, w, r, nil, errs.New(errs.ErrDecodingJSON, "word and translation are required"))
			return
		}

		// Добавление слова
		word, err := s.addWordToCollectionUC.AddWord(r.Context(), collectionID, req.Word, req.Translation, req.Example, userID)
		if err != nil {
			s.logger.Error("handlers.addWordToCollection", zap.Error(err))
			views.Return(s.logger, w, r, nil, err)
			return
		}

		views.Return(s.logger, w, r, views.NewAddWordToCollectionResponse(word), nil)
	}
}
