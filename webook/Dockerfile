# 基础镜像
FROM ubuntu:20.04

# 将编译后执行文件放到工作目录
COPY webook /app/webook

# 设定我的工作目录
WORKDIR /app

# 设置容器启动时执行的命令
ENTRYPOINT ["/app/webook"]

