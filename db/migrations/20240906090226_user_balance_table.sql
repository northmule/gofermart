-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.user_balance (
	id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	user_id int8 NOT NULL,
	value numeric(11, 2) NOT NULL,
	updated_at timestamp DEFAULT now() NOT NULL,
	CONSTRAINT user_balance_pk PRIMARY KEY (id)
);


-- public.user_balance внешние включи

ALTER TABLE public.user_balance ADD CONSTRAINT user_balance_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_balance;
-- +goose StatementEnd
