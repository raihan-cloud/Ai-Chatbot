# Gunakan Go 1.24 yang sudah kamu download tadi
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Pastikan git terinstall untuk fetch library
RUN apk add --no-cache git

COPY go.mod go.sum ./

# Daripada pakai @latest yang minta Go 1.25, 
# kita gunakan versi yang stabil untuk Go 1.24
RUN go mod download
RUN go get github.com/google/generative-ai-go/genai@v0.14.0
RUN go get google.golang.org/api/option@v0.200.0

COPY . .

# Build aplikasi
RUN go build -o main .

# Final stage
FROM alpine:latest
# Tambahkan ca-certificates agar bisa akses API Google (HTTPS)
RUN apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]