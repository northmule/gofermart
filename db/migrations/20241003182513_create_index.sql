-- +goose Up
-- +goose StatementBegin
CREATE UNIQUE INDEX user_balance_user_id_idx ON public.user_balance (user_id);
ALTER TABLE public.accruals ALTER COLUMN user_id SET NOT NULL;
CREATE INDEX accruals_user_id_idx ON public.accruals (user_id);
CREATE UNIQUE INDEX withdrawals_order_id_idx ON public.withdrawals (order_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX public.user_balance_user_id_idx;
DROP INDEX public.accruals_user_id_idx;
ALTER TABLE public.accruals ALTER COLUMN user_id DROP NOT NULL;
DROP INDEX public.withdrawals_order_id_idx;
-- +goose StatementEnd
