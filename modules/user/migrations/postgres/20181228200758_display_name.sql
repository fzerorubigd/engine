-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE aaa.users ADD COLUMN display_name VARCHAR;
UPDATE aaa.users SET display_name = email;
ALTER TABLE aaa.users ALTER COLUMN display_name SET NOT NULL;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE aaa.users DROP COLUMN display_name;