package client

import (
	"github.com/northmule/gophermart/internal/app/services/logger"
	"testing"
)

func TestAccrualClient_SendOrderNumber(t *testing.T) {
	t.Run("Сервис_accrual_запуск_и_проверка", func(t *testing.T) {

		chStop := make(chan struct{})
		//go func() {
		//	currentDir, err := os.Getwd()
		//	if err != nil {
		//		t.Fatal(err)
		//		return
		//	}
		//	cmd := exec.Command(currentDir + "/../../../../cmd/accrual/accrual_linux_amd64")
		//	err = cmd.Start()
		//
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//	<-chStop
		//	err = cmd.Process.Signal(os.Interrupt)
		//	if err != nil {
		//		t.Fatal(err)
		//	}
		//}()

		logger.NewLogger("info")
		ac := &AccrualClient{
			serviceURL: "http://localhost:8081",
			logger:     logger.LogSugar,
		}
		_, err := ac.SendOrderNumber("18")

		if err != nil {
			chStop <- struct{}{}
			t.Error(err)
		}
		chStop <- struct{}{}
	})

}
