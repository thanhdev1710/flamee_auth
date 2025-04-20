# Build stage
FROM golang:1.24-alpine3.20 AS builder

# Cài đặt thư viện cần thiết và môi trường làm việc
WORKDIR /app

# Copy go.mod và go.sum vào container và tải dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy phần còn lại của ứng dụng vào container
COPY . .

# Build ứng dụng Go
RUN go build -o flamee_auth ./cmd/server/main.go

# Runtime stage (Sử dụng image nhẹ cho sản phẩm cuối cùng)
FROM alpine:3.20

# Tạo non-root user cho ứng dụng (bảo mật tốt hơn)
RUN adduser -D -g '' appuser
USER appuser

# Cài đặt thư viện cần thiết cho môi trường runtime (nếu có)
# Cài thêm các package khác nếu cần thiết, ví dụ:
# RUN apk add --no-cache <package_name>

WORKDIR /app

# Copy binary từ build stage vào runtime container
COPY --from=builder /app/flamee_auth .

# Expose port mà ứng dụng sẽ chạy trên Render
EXPOSE 8081

# Chạy ứng dụng
CMD ["./flamee_auth"]
