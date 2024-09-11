# 使用官方的 Go 语言镜像作为基础镜像
FROM golang:1.20-alpine

# 设置工作目录
WORKDIR /app

# 将当前目录的文件复制到工作目录中
COPY . .

# 下载所需的 Go 模块
RUN go mod tidy

# 构建可执行文件
RUN go build -o ai-host-server

# 运行 Go 应用程序
CMD ["./ai-host-server"]

# 暴露端口
EXPOSE 10005
