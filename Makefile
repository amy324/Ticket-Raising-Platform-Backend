
export

DSN=host=${DB_HOST} port=${DB_PORT} user=root password=${DB_PASSWORD} dbname=${DB_DATABASE} sslmode=disable timezone=UTC connect_timeout=5
BINARY_NAME=ticketplatform.exe

## build: builds all binaries
build:
	@ go build -o ${BINARY_NAME}

run: build
	@echo Starting...
	@export DSN=${DSN}; \
	 ./${BINARY_NAME} &
	@echo Backend started!


clean:
	@echo Cleaning...
	@rm -f cmd/web/${BINARY_NAME}
	@go clean
	@echo Cleaned!
	
start: build
	@echo Starting...
	@export DSN=${DSN}; \
	cd cmd/web && ./${BINARY_NAME} &
	@echo Backend started!

stop:
	@echo "Stopping..."
	@pkill -f cmd/web/${BINARY_NAME}
	@echo Stopped backend

restart: stop start

test:
	@echo "Testing..."
	cd cmd/web && go test -v ./...
