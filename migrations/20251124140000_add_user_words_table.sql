-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS user_words (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    collection_id UUID NOT NULL REFERENCES word_collections(id) ON DELETE CASCADE,
    word VARCHAR(255) NOT NULL,
    translation VARCHAR(255) NOT NULL,
    example TEXT,
    next_review_date TIMESTAMP DEFAULT NOW(),
    review_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_words_collection ON user_words(collection_id);
CREATE INDEX idx_user_words_review_date ON user_words(next_review_date);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP INDEX IF EXISTS idx_user_words_review_date;
DROP INDEX IF EXISTS idx_user_words_collection;
DROP TABLE IF EXISTS user_words;

-- +goose StatementEnd

