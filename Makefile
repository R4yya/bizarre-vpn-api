# Environment variables
APP_NAME = api
BUILD_DIR = bin
SWAGGER_DIR = ./docs
API_SRC = cmd/api/main.go
CONFIG_FILE ?= config/local.yaml

# Install dependencies
.PHONY: deps
deps:
	go mod download

# Launch the API locally
.PHONY: run
run:
	go run $(API_SRC) --config=$(CONFIG_FILE)

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init --parseDependency --dir ./cmd/api,./internal/api/handlers --output $(SWAGGER_DIR)

# Build API for production
.PHONY: build
build:
	go build -ldflags "-s -w" -o $(BUILD_DIR)/$(APP_NAME_API) $(API_SRC)

# Launch the collected API binary
.PHONY: start
start:
	$(BUILD_DIR)/$(APP_NAME_API) --config=$(CONFIG_FILE)

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
	@echo "  run      - Launch API locally"
	@echo "  swagger      - Generate Swagger documentation"
	@echo "  build        - Build API for production"
	@echo "  start        - Launch the collected API binary"
	@echo "  clean        - Clear the collected files"
	@echo "  setup        - Install everything from scratch"
