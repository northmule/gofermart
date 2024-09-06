-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users (
	id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	"name" varchar(100) NOT NULL,
	login varchar(100) NOT NULL,
	"password" varchar(200) NOT NULL,
	created_at timestamp DEFAULT now() NOT NULL,
	deleted_at timestamp NULL,
	"uuid" uuid NULL,
	CONSTRAINT users_pk PRIMARY KEY (id),
	CONSTRAINT users_uuid_unique UNIQUE (uuid)
);
CREATE UNIQUE INDEX users_login_idx ON public.users USING btree (login) WHERE (deleted_at IS NULL);
CREATE INDEX users_login_password_idx ON public.users USING btree (login, password);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
