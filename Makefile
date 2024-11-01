# Environment variables
APP_NAME_API = api
APP_NAME_BOT = bot
BUILD_DIR = bin
SWAGGER_DIR = ./docs
API_SRC = cmd/api/main.go
BOT_SRC = cmd/bot/main.go

# Install dependencies
.PHONY: deps
deps:
	go mod download

# Launch the API locally
.PHONY: run-api
run-api:
	go run $(API_SRC)

# Launch the bot locally
.PHONY: run-bot
run-bot:
	go run $(BOT_SRC)

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init --parseDependency --dir ./cmd/api,./internal/api/handlers --output $(SWAGGER_DIR)

# Build API and bot for production
.PHONY: build
build: build-api build-bot

build-api:
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME_API) $(API_SRC)

build-bot:
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME_BOT) $(BOT_SRC)

# Launch the collected API and bot binaries
.PHONY: start
start: start-api start-bot

start-api:
	$(BUILD_DIR)/$(APP_NAME_API)

start-bot:
	$(BUILD_DIR)/$(APP_NAME_BOT)

# Clear the collected files
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

# Installing everything from scratch
.PHONY: setup
setup: deps swagger build

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  deps         - Install dependencies"
	@echo "  run-api      - Launch API locally"
	@echo "  run-bot      - Launch bot locally"
	@echo "  swagger      - Generate Swagger documentation"
	@echo "  build        - Build API and bot for production"
	@echo "  start        - Launch the collected API and bot binaries"
	@echo "  clean        - Clear the collected files"
	@echo "  setup        - Install everything from scratch"
