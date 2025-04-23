# Build stage
FROM golang:latest AS builder

WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o go_swift ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/go_swift .
COPY .env .env

RUN apk add --no-cache libc6-compat
RUN chmod +x /app/go_swift

EXPOSE 8080

CMD ["./go_swift"]