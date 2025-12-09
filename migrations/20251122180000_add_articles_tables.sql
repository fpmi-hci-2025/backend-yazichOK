-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS articles (
    id SERIAL PRIMARY KEY,
    image_url TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    level TEXT NOT NULL,  -- A1, A2, B1, B2, C1, C2
    minutes_to_read INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS article_vocabulary (
    id SERIAL PRIMARY KEY,
    article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
    word TEXT NOT NULL,
    part_of_speech TEXT NOT NULL,
    meaning TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS article_grammar_rules (
    id SERIAL PRIMARY KEY,
    article_id INTEGER REFERENCES articles(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    example TEXT NOT NULL,
    note TEXT NOT NULL
);

CREATE INDEX idx_articles_level ON articles(level);
CREATE INDEX idx_vocabulary_article ON article_vocabulary(article_id);
CREATE INDEX idx_grammar_article ON article_grammar_rules(article_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_grammar_article;
DROP INDEX IF EXISTS idx_vocabulary_article;
DROP INDEX IF EXISTS idx_articles_level;
DROP TABLE IF EXISTS article_grammar_rules;
DROP TABLE IF EXISTS article_vocabulary;
DROP TABLE IF EXISTS articles;
-- +goose StatementEnd

