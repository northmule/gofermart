package order

import (
	"github.com/northmule/gophermart/internal/app/services/logger"
	"regexp"
	"strconv"
)

const (
	NumberValidateAlg = "luhn"
)

type OrderService struct {
	alg              string
	regexOrderNumber *regexp.Regexp
}

func NewOrderService() *OrderService {
	instance := &OrderService{
		alg:              NumberValidateAlg,
		regexOrderNumber: regexp.MustCompile(`\d+`),
	}
	return instance
}

func (os *OrderService) ValidateOrderNumber(number string) bool {

	if !os.regexOrderNumber.MatchString(number) {
		return false
	}

	orderInt, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		logger.LogSugar.Errorf(err.Error())
		return false
	}

	switch os.alg {
	case NumberValidateAlg:
		logger.LogSugar.Infof("Проверка номера заказа %d по алгоритму Луна", orderInt)
		return os.luhnValid(int(orderInt))
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
