-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.orders ADD CONSTRAINT orders_number_unique UNIQUE ("number");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.orders DROP CONSTRAINT orders_number_unique;
-- +goose StatementEnd
