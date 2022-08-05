BINARY_NAME=skywalkerApp.exe

build:
	@go mod vendor
	@echo "Building Skywalker..."
	@go build -o tmp/${BINARY_NAME} .
	@echo "Skywalker built!"

run: build
	@echo "Starting Skywalker..."
	@start /min cmd /c tmp\${BINARY_NAME} &
	@echo "Skywalker started!"

clean:
	@echo "Skywalker..."
	@go clean
	@rm tmp/${BINARY_NAME}
	@echo "Cleaned!"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run

stop:
	@echo "Stopping Skywalker..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo "Stopped Skywalker!"

restart: stop start