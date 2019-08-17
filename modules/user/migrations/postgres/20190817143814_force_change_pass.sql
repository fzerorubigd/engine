-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied


ALTER TABLE aaa.users
    ADD COLUMN change_pass_at timestamp with time zone;

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE aaa.users
    DROP COLUMN change_pass_at;
