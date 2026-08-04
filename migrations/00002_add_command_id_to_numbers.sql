-- +goose Up
ALTER TABLE numbers
		ADD COLUMN command_id UUID;

UPDATE numbers
SET command_id = gen_random_uuid()
WHERE command_id IS NULL;

ALTER TABLE numbers
		ALTER COLUMN command_id SET NOT NULL,
		ADD CONSTRAINT numbers_command_id_unique UNIQUE (command_id);

-- +goose Down
ALTER TABLE numbers
		DROP CONSTRAINT IF EXISTS numbers_command_id_unique,
		DROP COLUMN IF EXISTS command_id;
