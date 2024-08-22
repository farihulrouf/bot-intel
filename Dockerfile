# Gunakan image Go resmi sebagai base image
FROM golang:1.22.5-alpine AS builder

# Set environment variable
ENV GO111MODULE=on

# Buat direktori kerja untuk aplikasi
WORKDIR /app

# Copy go.mod dan go.sum untuk meng-cache dependency
COPY go.mod go.sum ./

# Download semua dependencies (akan di-cache jika tidak ada perubahan pada go.mod dan go.sum)
RUN go mod download

# Copy semua file sumber kode ke dalam container
COPY . .

# Build aplikasi
RUN go build -o /app/main .

# Multistage build untuk mengurangi ukuran image final
FROM alpine:3.18

# Copy binary dari stage sebelumnya
COPY --from=builder /app/main /app/main

# Set workdir dan jalankan aplikasi
WORKDIR /app

# Jalankan aplikasi Go
CMD ["/app/main"]
