
-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied


CREATE SCHEMA aaa;

CREATE TABLE aaa.users (
	id uuid NOT NULL,
	email varchar NOT NULL,
	display_name varchar NOT NULL,
	password varchar NOT NULL,
	status smallint NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	last_login timestamp with time zone,
	change_pass_at timestamp with time zone,
	CONSTRAINT aaa_users_id_primary PRIMARY KEY (id),
	CONSTRAINT aaa_users_email_unique UNIQUE (email)
);

INSERT INTO aaa.users (id, email, display_name, password, status, created_at, updated_at) VALUES
	(uuid_generate_v4(), 'root@elbix.dev', 'root', '$2a$06$jDAy514SemGwCHhD..kfdedw/ibC3zyj.kqPtCHOoAwVYHC/RlDLa', 2, 'now', 'now');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP SCHEMA aaa CASCADE ;
