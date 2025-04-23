# GO_SWIFT

A RESTful API for managing SWIFT codes using Go and MySQL, designed for easy deployment with Docker.

## Requirements

- Go 1.x
- MySQL 8.x+
- Docker (optional, for containerization)

---

## Quick Start

### 1. Run Locally

1. Configure the `.env` file with the required environment variables (see **Environment Variables** section).
2. Uncomment the `godotenv` section in `cmd/main.go` to load environment variables from `.env`:
   ```go
   if err := godotenv.Load(); err != nil {
       fmt.Println("Error loading .env file:", err)
       os.Exit(1)
   }
   ```
3. Run the application locally:
   ```sh
   go run cmd/main.go
   ```

---

### 2. Run Containers with Docker Compose

1. Configure the `.env` file with the required environment variables (see **Environment Variables** section).
2. Start the application and database containers:
   ```sh
   docker-compose up --build
   ```
3. The application will be available at `http://localhost:8080`.

---

### 3. Run the Application Container Only

1. Configure the `.env` file with the required environment variables (see **Environment Variables** section).
2. Build the Docker image:
   ```sh
   docker build -t go_swift .
   ```
3. Run the container:
   ```sh
   docker run --env-file .env -p 8080:8080 go_swift
   ```
4. The application will be available at `http://localhost:8080`.

---

### 4. Set Default Environment Variables in Docker

If you want to set default environment variables directly in the `docker-compose.yml` file, add the following section under the `environment` key:

```yaml
environment:
  DB_USER: root
  DB_PASS: ""
  DB_HOST: db
  DB_PORT: 3306
  DB_NAME: go_swift
```

Alternatively, you can pass environment variables directly when running the container:

```sh
docker run -e DB_USER=root -e DB_PASS= -e DB_HOST=db -e DB_PORT=3306 -e DB_NAME=go_swift -p 8080:8080 go_swift
```

---

## Environment Variables

Configure the following environment variables in a `.env` file or pass them directly:

```plaintext
DB_USER=root
DB_PASS=
DB_HOST=localhost
DB_PORT=3306
DB_NAME=go_swift
```

For Docker Compose, set `DB_HOST=db` to connect to the database container.

---

## API Endpoints

- **GET** `/v1/swift-codes/{swift-code}` - Retrieve details of a specific SWIFT code.
- **GET** `/v1/swift-codes/country/{countryISO2code}` - Retrieve SWIFT codes for a specific country.
- **POST** `/v1/swift-codes` - Add a new SWIFT code.
- **DELETE** `/v1/swift-codes/{swift-code}` - Delete a SWIFT code.

---

## Development

1. Fetch dependencies:
   ```sh
   go mod tidy
   ```
2. Run the application locally:
   ```sh
   go run cmd/main.go
   ```

---

## Tests

Run unit and integration tests:

1. Run all tests:
   ```sh
   go test ./... -v
   ```
2. Run tests for a specific module (e.g., `internal/service`):
   ```sh
   go test ./internal/service -v
   ```
3. Run tests with coverage:
   ```sh
   go test ./... -cover
   ```
4. Generate a coverage report:
   ```sh
   go test ./... -coverprofile=coverage.out
   go tool cover -html=coverage.out
   ```

---

## Notes

- To use local environment variables, ensure the `godotenv` section in `cmd/main.go` is uncommented.
- If you encounter database connection issues, verify that the database is correctly configured and accessible at the specified host and port.