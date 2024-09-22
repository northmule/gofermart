package rctx

type key int

// Наименования контекста
const (
	// UserCtxKey объект с пользователем
	UserCtxKey key = iota
	// OrderUpload объект заказа
	OrderUpload
	// TransactionCtxKey транзакция в рамках запроса
	TransactionCtxKey
)
