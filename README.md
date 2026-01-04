# BizarreVPN API

This project is a template Go-based backend API for Telegram bot and Mini App based VPN service.

## Architecture
I use features-based layers DDD


## Setup

1. Clone the repository:

    ```shell
    git clone https://github.com/R4yya/bizarre-vpn-api.git
    cd bizarre-vpn-api
    ```

2. Install Go dependencies:

    ```shell
    go mod download
    ```

3. Set up environment variables:

   Rename `.env.example` to `.env` and configure necessary environment variables like:

    ```dotenv
    API_PORT=:5000
    DATABASE_PATH=./data/database.db
    SWAGGER_PATH=/docs
    ```

## Build and Run

### Running Locally

1. Generate Swagger documentation (optional for development):

    ```shell
    swag init --parseDependency --dir ./cmd/api,./internal/api/handlers --output ./docs
    ```

2. Start the API server:

    ```shell
    go run cmd/api/main.go
    ```

3. Start the Telegram bot:

    ```shell
    go run cmd/bot/main.go
    ```

### Building

To build the project as an executable:

```shell
go build -ldflags "-s -w" -o bin/api cmd/api/main.go
go build -ldflags "-s -w" -o bin/bot cmd/bot/main.go
```

This will generate `api` and `bot` binaries in the `bin` directory.

### Running the Executables

After building, run the binaries with:

```shell
./bin/api
./bin/bot
```

## Swagger Documentation

After running the API server, Swagger documentation is available at:

```plaintext
https://example-domain:<API_PORT>/<SWAGGER_PATH>/index.html
```

## Makefile

The Makefile includes targets for managing dependencies, building, and running the API and bot services.

Run this to see list of targets:
```shell
make help
```

## Logging

Logs are saved in the `logs` directory, with separate log files for the API and bot services.

---

Feel free to reach out or open issues for bugs or feature requests.
