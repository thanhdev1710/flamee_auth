# Dockerfile

# Build stage
FROM golang:1.24-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o flamee_auth ./cmd/server/main.go

# Runtime stage
FROM alpine:3.20

WORKDIR /app

# Copy binary từ stage build
COPY --from=builder /app/flamee_auth .

# Expose port nếu cần (ví dụ 8081)
EXPOSE 8081

# Chạy binary
CMD ["./flamee_auth"]
