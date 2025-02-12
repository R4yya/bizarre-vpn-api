# Environment variables
PROJECT_DIR = $(shell pwd)
PROJECT_BUILD = $(PROJECT_DIR)/bin
APP_NAME = bizarre-api
SWAGGER_DIR = ./docs
API_SRC = cmd/api/main.go
CONFIG_FILE ?= config/local.yaml

GOLANGCI_LINT = $(PROJECT_BUILD)/golangci-lint

# Install dependencies
.PHONY: deps
deps:
	go mod download

# Launch the API locally
.PHONY: run
run:
	make swagger
	go run $(API_SRC) --config=$(CONFIG_FILE)

# Generate Swagger documentation
.PHONY: swagger
swagger:
	swag init --parseDependency --dir ./cmd/api,./internal/api/handlers	--exclude ./internal/api/handlers/subscription_plan	--output $(SWAGGER_DIR)
	rm ./docs/docs.go
	rm ./docs/swagger.yaml
# Build API for production
.PHONY: build
build:
	go build -ldflags "-s -w" -o $(PROJECT_BUILD)/$(APP_NAME) $(API_SRC)

# Build API for production Linux
.PHONY: build
build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(PROJECT_BUILD)/$(APP_NAME) $(API_SRC)

# Launch the collected API binary
.PHONY: start
start:
	$(PROJECT_BUILD)/$(APP_NAME) --config=$(CONFIG_FILE)

# Clear the collected files
.PHONY: clean
clean:
	rm -rf $(PROJECT_BUILD)

# Installing everything from scratch
.PHONY: setup
setup: deps swagger build

# Install linter
.PHONY: .install-linter
.install-linter:
	@echo "INSTALL GOLANGCI-LINT"
	@if [ ! -f $(GOLANGCI_LINT) ]; then \
		echo "golangci-lint не найден. Скачиваем и устанавливаем..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BUILD) v1.63.4; \
		echo "Даём права на исполнение"; \
		chmod +x $(GOLANGCI_LINT); \
		echo "golangci-lint успешно установлен в $(GOLANGCI_LINT)"; \
	else \
		echo "golangci-lint уже установлен в $(GOLANGCI_LINT)"; \
	fi
# Run Linter
.PHONY: lint
lint: .install-linter
	### RUN GOLANGCI-LINT ###
	$(GOLANGCI_LINT) run ./... --config=./.golangci.yml

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
