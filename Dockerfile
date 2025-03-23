# 使用官方 Go 镜像作为构建阶段
FROM golang:1.24.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# **复制 config.yaml**
COPY configs/config.yaml /app/configs/config.yaml

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd/main.go

# 使用 Alpine 运行应用
FROM alpine:latest
WORKDIR /root/

# 安装 tzdata 并设置时区
RUN apk add --no-cache tzdata \
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

ENV TZ=Asia/Shanghai

# 复制二进制文件和配置文件
COPY --from=builder /app/app .
COPY --from=builder /app/configs/config.yaml ./configs/config.yaml

CMD ["./app"]

