-- +goose Up
-- +goose StatementBegin
CREATE TABLE numbers (
		id SERIAL PRIMARY KEY,
		number BIGINT NOT NULL,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS numbers;
-- +goose StatementEnd
