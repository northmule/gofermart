cd $GOFERMART_PROJECT_PATH
./cmd/goose/goose  -dir db/migrations postgres "postgres://postgres:123@localhost:5456/gofermart?sslmode=disable" up