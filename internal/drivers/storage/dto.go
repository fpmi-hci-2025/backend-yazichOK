package storage

import "github.com/google/uuid"

type Topic struct {
	ID          int    `db:"id"`
	Title       string `db:"title"`
	Description string `db:"description"`
	ImagePath   string `db:"image_path"`
}

type Question struct {
	ID       int    `db:"id"`
	TopicID  int    `db:"topic_id"`
	Question string `db:"question_text"`
}

type Answer struct {
	ID         int    `db:"id"`
	QuestionID int    `db:"question_id"`
	SessionID  string `db:"session_id"`
	Filename   string `db:"minio_filename"`
}

type Article struct {
	ID            int    `db:"id"`
	ImageURL      string `db:"image_url"`
	Title         string `db:"title"`
	Content       string `db:"content"`
	Level         string `db:"level"`
	MinutesToRead int    `db:"minutes_to_read"`
}

type ArticleVocabulary struct {
	ID           int    `db:"id"`
	ArticleID    int    `db:"article_id"`
	Word         string `db:"word"`
	PartOfSpeech string `db:"part_of_speech"`
	Meaning      string `db:"meaning"`
}

type ArticleGrammarRule struct {
	ID        int    `db:"id"`
	ArticleID int    `db:"article_id"`
	Name      string `db:"name"`
	Example   string `db:"example"`
	Note      string `db:"note"`
}

type WordCollection struct {
	ID                       uuid.UUID `db:"id"`
	UserID                   int       `db:"user_id"`
	Name                     string    `db:"name"`
	ImagePath                string    `db:"image_path"`
	TotalWordsCount          int       `db:"total_words_count"`
	LearnedWordsCount        int       `db:"learned_words_count"`
	CurrentStreakDays        int       `db:"current_streak_days"`
	LongestStreakDays        int       `db:"longest_streak_days"`
	LastStudiedAt            *string   `db:"last_studied_at"`
	AISuggestions            *string   `db:"ai_suggestions"`
	AISuggestionsGeneratedAt *string   `db:"ai_suggestions_generated_at"`
	CreatedAt                string    `db:"created_at"`
	UpdatedAt                string    `db:"updated_at"`
}

type UserWord struct {
	ID             string  `db:"id"`
	CollectionID   string  `db:"collection_id"`
	Word           string  `db:"word"`
	Translation    string  `db:"translation"`
	Example        *string `db:"example"`
	NextReviewDate string  `db:"next_review_date"`
	ReviewCount    int     `db:"review_count"`
	CreatedAt      string  `db:"created_at"`
	UpdatedAt      string  `db:"updated_at"`
}
