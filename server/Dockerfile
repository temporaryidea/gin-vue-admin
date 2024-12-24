FROM golang:1.22.0-alpine

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
COPY go.work go.work.sum ./

# Set timezone
ENV TZ=Asia/Shanghai

# Install dependencies
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    go env -w GO111MODULE=on && \
    go env -w CGO_ENABLED=0 && \
    go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o server .

EXPOSE 8888

# Command to run
CMD ["./server", "-c", "config.docker.yaml"]
