-- +goose Up
-- +goose StatementBegin
CREATE TABLE numbers (
		id SERIAL PRIMARY KEY,
		a NUMERIC NOT NULL,
		b NUMERIC NOT NULL,
		operation VARCHAR(10) NOT NULL,
		result NUMERIC NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS numbers;
-- +goose StatementEnd
