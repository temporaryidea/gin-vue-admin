FROM golang:alpine as go-builder

WORKDIR /go/src/github.com/flipped-aurora/gin-vue-admin/server
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o server .

FROM python:3.12-slim

WORKDIR /app

# Install Go binary and dependencies
COPY --from=go-builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/server ./
COPY --from=go-builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/resource ./resource/
COPY --from=go-builder /go/src/github.com/flipped-aurora/gin-vue-admin/server/config.docker.yaml ./

# Set timezone and install dependencies
ENV TZ=Asia/Shanghai
RUN apt-get update && apt-get install -y tzdata \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install Python dependencies directly
RUN pip install fastapi==0.109.2 uvicorn==0.27.1 httpx==0.26.0

# Copy FastAPI app
COPY proxy.py .

EXPOSE 8000
EXPOSE 8888

# Start both services
CMD ./server -c config.docker.yaml & python proxy.py
