package order

import "github.com/northmule/gophermart/internal/app/services/logger"

const (
	NumberValidateAlg = "luhn"
)

type OrderService struct {
	alg string
}

func NewOrderService() *OrderService {
	instance := &OrderService{
		alg: NumberValidateAlg,
	}
	return instance
}

func (os *OrderService) ValidateOrderNumber(number int) bool {
	switch os.alg {
	case NumberValidateAlg:
		logger.LogSugar.Infof("Проверка номера заказа %d по алгоритму Луна", number)
		return os.luhnValid(number)
	}
	logger.LogSugar.Errorf("Указанный алгоритм проверки номера заказа не реализован: %s", os.alg)
	return false
}

func (os *OrderService) luhnValid(number int) bool {
	isValid := (number%10+os.luhnChecksum(number/10))%10 == 0
	logger.LogSugar.Infof("Результат проверки: номер заказа: %v", isValid)
	return isValid
}

func (os *OrderService) luhnChecksum(number int) int {
	var luhn int

	for i := 0; number > 0; i++ {
		cur := number % 10

		if i%2 == 0 {
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		number = number / 10
	}
	return luhn % 10
}
