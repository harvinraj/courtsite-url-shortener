# Stage 1: Build binary
FROM golang:alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o url-shortener main.go

# Stage 2: Minimal runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/url-shortener .
EXPOSE 8080
CMD ["./url-shortener"]