include .env
export


DSN=host=localhost port=3306 user=${MYSQL_USER} password=${MYSQL_PASSWORD} dbname=${MYSQL_DB} sslmode=disable timezone=UTC connect_timeout=5
BINARY_NAME=myapp.exe
REDIS=127.0.0.1:6379

## build: builds all binaries
build:
	@go build -o ${BINARY_NAME} ./cmd/web
	@echo back end built!

run: build
	@echo Starting...
	@DSN=${DSN} REDIS=${REDIS} ./${BINARY_NAME} &
	@echo back end started!

clean:
	@echo Cleaning...
	@rm -f ${BINARY_NAME}
	@go clean
	@echo Cleaned!

start: build
	@echo Starting...
	@export DSN=${DSN} REDIS=${REDIS} && ./${BINARY_NAME} &
	@echo back end started!


stop:
	@echo "Stopping..."
	@pkill -f ${BINARY_NAME}
	@echo Stopped back end

restart: stop start

test:
	@echo "Testing..."
	go test -v ./...
