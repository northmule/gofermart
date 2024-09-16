-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.jobs_order (
               id int8 GENERATED ALWAYS AS IDENTITY NOT NULL,
               created_at timestamp DEFAULT now() NOT NULL,
               updated_at timestamp DEFAULT now() NOT NULL,
               next_run timestamp DEFAULT now() NOT NULL,
               run_cnt int8 DEFAULT 1 NOT NULL,
               order_number varchar(100) NOT NULL,
               is_work bool DEFAULT false NOT NULL,
               CONSTRAINT jobs_order_order_number_unique UNIQUE (order_number),
               CONSTRAINT jobs_order_pk PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS jobs_order;
-- +goose StatementEnd
