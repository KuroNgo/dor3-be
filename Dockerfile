# Sử dụng hình ảnh Alpine Go
FROM golang:1.21-alpine

# Thiết lập thư mục làm việc là /app
WORKDIR /app

# Sao chép mã nguồn của ứng dụng Go vào thư mục /app trong container
COPY ../.. .

# Tải và cài đặt các phụ thuộc (nếu cần)
RUN go mod download

# Biên dịch ứng dụng Go
RUN go build -o main .

# Expose port 8000 cho ứng dụng Go
EXPOSE 8080

# Chạy ứng dụng Go khi container được khởi chạy và cung cấp các biến môi trường cho kết nối đến MongoDB
CMD ["./main"]