FROM golang:alpine

# Set timezone
ENV TZ=Asia/Shanghai
RUN apk add --no-cache tzdata python3 py3-pip \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime \
    && echo $TZ > /etc/timezone

WORKDIR /app

# Copy Go files and build
COPY . .
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go mod tidy \
    && go build -o server

# Install Python dependencies directly
RUN pip3 install fastapi==0.109.2 uvicorn==0.27.1 httpx==0.26.0

EXPOSE 8000
EXPOSE 8888

# Start both services
CMD ./server -c config.docker.yaml & python3 proxy.py
