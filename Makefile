BINARY_NAME=SkywalkerApp

build:
	@go mod vendor
	@echo "Building Skywalker..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "Skywalker built!"

run: build
	@echo "Starting Skywalker..."
	@./tmp/${BINARY_NAME} &
	@echo "Skywalker started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

start_compose:
	docker-compose up -d

stop_compose:
	docker-compose down

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping Skywalker..."
	@-pkill -SIGTERM -f "./tmp/${BINARY_NAME}"
	@echo "Stopped Skywalker!"

restart: stop start