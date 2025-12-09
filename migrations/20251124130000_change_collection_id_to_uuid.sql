-- +goose Up
-- +goose StatementBegin

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE word_collections ADD COLUMN id_new UUID DEFAULT uuid_generate_v4();

UPDATE word_collections SET id_new = uuid_generate_v4() WHERE id_new IS NULL;

ALTER TABLE word_collections ALTER COLUMN id_new SET NOT NULL;

ALTER TABLE word_collections DROP CONSTRAINT word_collections_pkey;
ALTER TABLE word_collections DROP COLUMN id;

ALTER TABLE word_collections RENAME COLUMN id_new TO id;

ALTER TABLE word_collections ADD PRIMARY KEY (id);

DROP INDEX IF EXISTS idx_word_collections_user;
DROP INDEX IF EXISTS idx_word_collections_updated;

CREATE INDEX idx_word_collections_user ON word_collections(user_id);
CREATE INDEX idx_word_collections_updated ON word_collections(updated_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Возвращаем обратно на SERIAL (только для новой установки, данные потеряются)
ALTER TABLE word_collections DROP CONSTRAINT word_collections_pkey;
ALTER TABLE word_collections DROP COLUMN id;
ALTER TABLE word_collections ADD COLUMN id SERIAL PRIMARY KEY;

-- +goose StatementEnd

