echo $GOFERMART_PROJECT_PATH
cd $GOFERMART_PROJECT_PATH/cmd/gophermart
go build -buildvcs=false -o gophermart
chmod +x gophermart