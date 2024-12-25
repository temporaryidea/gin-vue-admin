FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY go.work ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/server .
# Copy config files
COPY --from=builder /app/config.docker.yaml ./config.yaml
COPY --from=builder /app/config/alipay.yaml ./config/

EXPOSE 8888

CMD ["./server"]
