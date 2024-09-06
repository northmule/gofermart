-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.orders (
           id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
           "number" varchar(100) NOT NULL,
           status varchar(50) NOT NULL,
           created_at timestamp DEFAULT now() NOT NULL,
           deleted_at timestamp NULL,
           CONSTRAINT orders_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX orders_number_idx ON public.orders USING btree (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
