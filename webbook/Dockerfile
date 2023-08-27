# 指定基础镜像
FROM ubuntu:20.04
# 将编译好的可执行文件拷贝到镜像中
COPY webook /app/webook
# 指定工作目录
WORKDIR /app
# CMD 是执行命令
# 最佳
ENTRYPOINT ["/app/webook"]