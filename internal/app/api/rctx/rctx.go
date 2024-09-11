package rctx

type key int

// Наименования контекста
const (
	UserCtxKey key = iota
	OrderUpload
)
