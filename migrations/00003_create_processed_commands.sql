-- +goose Up
CREATE TABLE processed_commands (
		command_id UUID PRIMARY KEY,
		command_type TEXT NOT NULL,
		payload_hash BYTEA NOT NULL,
		status TEXT NOT NULL CHECK (status IN ('completed', 'processing', 'failed')),
		response_payload JSONB,
		completed_at TIMESTAMPTZ,

		CHECK (
		(
			status = 'processing'
			AND response_payload IS NULL
			AND completed_at IS NULL
		)
		OR
		(
			status IN ('completed', 'failed')
			AND response_payload IS NOT NULL
			AND completed_at IS NOT NULL
		)
	)
);

-- +goose Down
DROP TABLE IF EXISTS processed_commands;
