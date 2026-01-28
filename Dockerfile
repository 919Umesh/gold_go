# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install git and ssl certificates for build
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Run stage
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for external API calls/DB connections
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

CMD ["./main"]