package storage

import (
	"context"

	"speech-processing-service/internal/config"
	"speech-processing-service/internal/errs"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	errCodeViolation = "23503"
)

type Storage struct {
	db *sqlx.DB
}

func New(dbCfg *config.DB) (Storage, error) {
	connStr := dbCfg.GetConnStr()

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return Storage{}, errs.New(errs.ErrInitialization, "postgres:"+err.Error())
	}

	if err := db.Ping(); err != nil {
		return Storage{}, errs.New(errs.ErrInitialization, "postgres:"+err.Error())
	}

	// Run database migrations
	if err := goose.Up(db.DB, "migrations"); err != nil {
		return Storage{}, errs.New(errs.ErrInitialization, "migrations:"+err.Error())
	}

	return Storage{
		db: db,
	}, nil
}

func (s *Storage) GetAllTopics(ctx context.Context) ([]Topic, error) {
	var topics []Topic

	if err := s.db.SelectContext(
		ctx,
		&topics,
		"SELECT * FROM topics",
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext"+err.Error())
	}

	return topics, nil
}

func (s *Storage) GetQuestionsByTopicID(ctx context.Context, topicID int) ([]Question, error) {
	var questions []Question

	if err := s.db.SelectContext(
		ctx,
		&questions,
		"SELECT * FROM questions WHERE topic_id = $1",
		topicID,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext"+err.Error())
	}

	return questions, nil
}

func (s *Storage) CreateSession(ctx context.Context, sessionID string, topicID int) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO sessions (id, topic_id) VALUES ($1, $2)",
		sessionID,
		topicID,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == errCodeViolation {
			return errs.New(errs.ErrNotFound, "topic not found")
		}

		return errs.New(errs.ErrExecutionQuery, "s.db.ExecContext"+err.Error())
	}

	return nil
}

func (s *Storage) CreateAnswer(ctx context.Context, sessionID string, questionID int, filename string) error {
	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO answers (session_id, question_id, minio_filename) VALUES ($1, $2, $3)",
		sessionID,
		questionID,
		filename,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == errCodeViolation {
			return errs.New(errs.ErrNotFound, "session or question not found")
		}

		return errs.New(errs.ErrExecutionQuery, "s.db.ExecContext"+err.Error())
	}

	return nil
}

func (s *Storage) GetQuestionByID(ctx context.Context, id int) (Question, error) {
	var question Question

	if err := s.db.GetContext(
		ctx,
		&question,
		"SELECT * FROM questions WHERE id = $1",
		id,
	); err != nil {
		return Question{}, errs.New(errs.ErrExecutionQuery, "s.db.GetContext"+err.Error())
	}

	return question, nil
}

func (s *Storage) GetAnswerBySessionID(ctx context.Context, sessionID string) ([]Answer, error) {
	var answers []Answer
	if err := s.db.SelectContext(
		ctx,
		&answers,
		"SELECT * FROM answers WHERE session_id = $1",
		sessionID,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext"+err.Error())
	}

	return answers, nil
}

func (s *Storage) GetArticles(ctx context.Context, limit, offset int) ([]Article, error) {
	var articles []Article
	if err := s.db.SelectContext(
		ctx,
		&articles,
		"SELECT id, image_url, title, level, minutes_to_read FROM articles ORDER BY created_at DESC LIMIT $1 OFFSET $2",
		limit,
		offset,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext: "+err.Error())
	}

	return articles, nil
}

func (s *Storage) GetArticleByID(ctx context.Context, id int) (Article, error) {
	var article Article
	if err := s.db.GetContext(
		ctx,
		&article,
		"SELECT id, image_url, title, content, level, minutes_to_read FROM articles WHERE id = $1",
		id,
	); err != nil {
		return Article{}, errs.New(errs.ErrExecutionQuery, "s.db.GetContext: "+err.Error())
	}

	return article, nil
}

func (s *Storage) GetArticleVocabulary(ctx context.Context, articleID int) ([]ArticleVocabulary, error) {
	var vocabulary []ArticleVocabulary
	if err := s.db.SelectContext(
		ctx,
		&vocabulary,
		"SELECT id, article_id, word, part_of_speech, meaning FROM article_vocabulary WHERE article_id = $1",
		articleID,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext: "+err.Error())
	}

	return vocabulary, nil
}

func (s *Storage) GetArticleGrammarRules(ctx context.Context, articleID int) ([]ArticleGrammarRule, error) {
	var rules []ArticleGrammarRule
	if err := s.db.SelectContext(
		ctx,
		&rules,
		"SELECT id, article_id, name, example, note FROM article_grammar_rules WHERE article_id = $1",
		articleID,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext: "+err.Error())
	}

	return rules, nil
}

func (s *Storage) CreateWordCollection(ctx context.Context, userID int, name, imagePath string) (WordCollection, error) {
	var collection WordCollection
	if err := s.db.QueryRowxContext(
		ctx,
		`INSERT INTO word_collections (user_id, name, image_path) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, user_id, name, image_path, total_words_count, learned_words_count, 
		           current_streak_days, longest_streak_days, created_at, updated_at`,
		userID,
		name,
		imagePath,
	).StructScan(&collection); err != nil {
		return WordCollection{}, errs.New(errs.ErrExecutionQuery, "s.db.QueryRowxContext: "+err.Error())
	}

	return collection, nil
}

func (s *Storage) DeleteWordCollection(ctx context.Context, collectionID uuid.UUID, userID int) error {
	result, err := s.db.ExecContext(
		ctx,
		`DELETE FROM word_collections WHERE id = $1 AND user_id = $2`,
		collectionID,
		userID,
	)
	if err != nil {
		return errs.New(errs.ErrExecutionQuery, "s.db.ExecContext: "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errs.New(errs.ErrExecutionQuery, "result.RowsAffected: "+err.Error())
	}

	if rowsAffected == 0 {
		return errs.New(errs.ErrNotFound, "collection not found or access denied")
	}

	return nil
}

func (s *Storage) GetUserCollections(ctx context.Context, userID int) ([]WordCollection, error) {
	var collections []WordCollection
	if err := s.db.SelectContext(
		ctx,
		&collections,
		`SELECT id, user_id, name, image_path, total_words_count, learned_words_count,
		        current_streak_days, longest_streak_days, last_studied_at, created_at, updated_at
		 FROM word_collections 
		 WHERE user_id = $1
		 ORDER BY updated_at DESC`,
		userID,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext: "+err.Error())
	}

	return collections, nil
}

func (s *Storage) GetWordCollectionByID(ctx context.Context, collectionID string, userID int) (WordCollection, error) {
	var collection WordCollection
	if err := s.db.GetContext(
		ctx,
		&collection,
		`SELECT id, user_id, name, image_path, total_words_count, learned_words_count,
		        current_streak_days, longest_streak_days, last_studied_at, 
		        ai_suggestions, ai_suggestions_generated_at, created_at, updated_at
		 FROM word_collections 
		 WHERE id = $1 AND user_id = $2`,
		collectionID,
		userID,
	); err != nil {
		return WordCollection{}, errs.New(errs.ErrExecutionQuery, "s.db.GetContext: "+err.Error())
	}

	return collection, nil
}

func (s *Storage) GetUserWordsByCollectionID(ctx context.Context, collectionID string) ([]UserWord, error) {
	var words []UserWord
	if err := s.db.SelectContext(
		ctx,
		&words,
		`SELECT id, collection_id, word, translation, example, next_review_date,
		        review_count, created_at, updated_at
		 FROM user_words 
		 WHERE collection_id = $1
		 ORDER BY created_at DESC`,
		collectionID,
	); err != nil {
		return nil, errs.New(errs.ErrExecutionQuery, "s.db.SelectContext: "+err.Error())
	}

	return words, nil
}

func (s *Storage) AddWordToCollection(ctx context.Context, collectionID, word, translation string, example *string) (UserWord, error) {
	var userWord UserWord

	if err := s.db.GetContext(
		ctx,
		&userWord,
		`INSERT INTO user_words (collection_id, word, translation, example, next_review_date, review_count)
		 VALUES ($1, $2, $3, $4, NOW(), 0)
		 RETURNING id, collection_id, word, translation, example, next_review_date, review_count, created_at, updated_at`,
		collectionID,
		word,
		translation,
		example,
	); err != nil {
		return UserWord{}, errs.New(errs.ErrExecutionQuery, "s.db.GetContext: "+err.Error())
	}

	_, err := s.db.ExecContext(
		ctx,
		`UPDATE word_collections 
		 SET total_words_count = total_words_count + 1, updated_at = NOW()
		 WHERE id = $1`,
		collectionID,
	)
	if err != nil {
		return UserWord{}, errs.New(errs.ErrExecutionQuery, "s.db.ExecContext: "+err.Error())
	}

	return userWord, nil
}
