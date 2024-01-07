FROM debian:stable-slim

# 设置工作目录
WORKDIR /app

# 安装 curl
RUN apt-get update && \
    apt-get install -y curl && \
    # 下载所需文件
    curl -O http://192.168.3.8:3000/package/TDengine-client/TDengine-client-3.2.2.0-Linux-x64.tar.gz && \
    tar -xf TDengine-client-3.2.2.0-Linux-x64.tar.gz && \
    cd TDengine-client-3.2.2.0 && ./install_client.sh && \
    cd .. && rm -rf /app/TDengine-client-3.2.2.0 && rm -rf /app/TDengine-client-3.2.2.0-Linux-x64.tar.gz && \
    # 清理缓存
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 复制应用程序和配置文件到工作目录
COPY go-emqx-to-tdengine config.toml create.sql insert.sql /app/

# 容器启动命令
CMD ["./go-emqx-to-tdengine"]
