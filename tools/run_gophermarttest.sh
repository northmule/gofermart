clear
echo $GOFERMART_PROJECT_PATH
cd $GOFERMART_PROJECT_PATH/cmd/gophermart
go build -buildvcs=false -o gophermart
chmod +x gophermart
cd $GOFERMART_PROJECT_PATH

# Запуск тестов
gophermarttest \
            -test.v -test.run=^TestGophermart$ \
            -gophermart-binary-path=$GOFERMART_PROJECT_PATH/cmd/gophermart/gophermart \
            -gophermart-host=localhost \
            -gophermart-port=8081 \
            -gophermart-database-uri='postgres://postgres:123@localhost:5456/gofermart?sslmode=disable' \
            -accrual-binary-path=$GOFERMART_PROJECT_PATH/cmd/accrual/accrual_linux_amd64 \
            -accrual-host=localhost \
            -accrual-port=8091 \
            -accrual-database-uri='postgres://postgres:123@localhost:5456/gofermart?sslmode=disable'