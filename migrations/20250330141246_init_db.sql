-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS topics (
                        id SERIAL PRIMARY KEY,
                        title TEXT NOT NULL UNIQUE,
                        description TEXT NOT NULL,
                        image_path TEXT NOT NULL  -- path to image in minio
);

CREATE TABLE IF NOT EXISTS questions (
                           id SERIAL PRIMARY KEY,
                           topic_id INTEGER REFERENCES topics(id) ON DELETE CASCADE,
                           question_text TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS topics;
-- +goose StatementEnd
