-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.user_orders (
	id int8 GENERATED ALWAYS AS IDENTITY( INCREMENT BY 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1 NO CYCLE) NOT NULL,
	user_id int8 NOT NULL,
	order_id int8 NOT NULL,
	CONSTRAINT user_orders_pk PRIMARY KEY (id)
);
CREATE UNIQUE INDEX user_orders_user_id_idx ON public.user_orders USING btree (user_id, order_id);


-- public.user_orders внешние включи

ALTER TABLE public.user_orders ADD CONSTRAINT user_orders_orders_fk FOREIGN KEY (order_id) REFERENCES public.orders(id);
ALTER TABLE public.user_orders ADD CONSTRAINT user_orders_users_fk FOREIGN KEY (user_id) REFERENCES public.users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_orders;
-- +goose StatementEnd
