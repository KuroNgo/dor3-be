# Sử dụng hình ảnh Alpine Go
FROM golang:1.20-alpine

# Thiết lập thư mục làm việc là /app
WORKDIR /app

# Sao chép mã nguồn của ứng dụng Go vào thư mục /app trong container
COPY .. .

# Tải và cài đặt các phụ thuộc (nếu cần)
RUN go mod download

# Biên dịch ứng dụng Go
RUN go build -o main .

# Expose port 8000 cho ứng dụng Go
EXPOSE 8080
