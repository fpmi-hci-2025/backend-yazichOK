-- +goose Up
-- +goose StatementBegin
CREATE TABLE answers (
    id SERIAL PRIMARY KEY,
    session_id UUID NOT NULL,
    question_id INT NOT NULL,
    minio_filename TEXT NOT NULL UNIQUE,
    FOREIGN KEY (question_id) REFERENCES questions(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS answers;
-- +goose StatementEnd
