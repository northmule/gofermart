package constants

const (
	//OrderStatusNew заказ загружен в систему, но не попал в обработку
	OrderStatusNew = "NEW"
	//OrderStatusProcessing вознаграждение за заказ рассчитывается
	OrderStatusProcessing = "PROCESSING"
	//OrderStatusInvalid система расчёта вознаграждений отказала в расчёте
	OrderStatusInvalid = "INVALID"
	//OrderStatusProcessed данные по заказу проверены и информация о расчёте успешно получена
	OrderStatusProcessed = "PROCESSED"
)
