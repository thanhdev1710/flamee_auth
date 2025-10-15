# ---------- Stage 1: Build ----------
FROM golang:1.24-alpine3.20 AS builder

# Bật Go module và tắt CGO để build binary tĩnh (tương thích Alpine)
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app

# Copy go.mod và go.sum trước để tận dụng cache build layer
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code vào container
COPY . .

# Build binary (tên flamee_auth)
RUN go build -ldflags="-s -w" -o flamee_auth ./cmd/server/main.go


# ---------- Stage 2: Runtime ----------
FROM alpine:3.20

# Thêm chứng chỉ CA (cần thiết nếu app gọi HTTPS — ví dụ AWS SDK, API bên ngoài)
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser

WORKDIR /app

# Copy binary từ builder stage
COPY --from=builder /app/flamee_auth .

USER root

# Expose port app lắng nghe
EXPOSE 8081

# Run app
ENTRYPOINT ["./flamee_auth"]
