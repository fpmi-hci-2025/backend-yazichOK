-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS word_collections (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    image_path TEXT,
    
    -- Статистика (кэш для быстрого доступа)
    total_words_count INTEGER DEFAULT 0,
    learned_words_count INTEGER DEFAULT 0,
    current_streak_days INTEGER DEFAULT 0,
    longest_streak_days INTEGER DEFAULT 0,
    last_studied_at TIMESTAMP,
    
    -- AI рекомендации (кэш)
    ai_suggestions JSONB,
    ai_suggestions_generated_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_word_collections_user ON word_collections(user_id);
CREATE INDEX idx_word_collections_updated ON word_collections(updated_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_word_collections_updated;
DROP INDEX IF EXISTS idx_word_collections_user;
DROP TABLE IF EXISTS word_collections;
-- +goose StatementEnd

