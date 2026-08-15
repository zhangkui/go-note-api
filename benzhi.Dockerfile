# 官方 Go 镜像，自带完整工具链
FROM golang:1.22

WORKDIR /app

# 复制所有项目文件
COPY . .

# 预编译一次，把编译缓存留在镜像里
RUN go build ./...

# 容器启动后进入 shell，方便操作
CMD ["bash"]
