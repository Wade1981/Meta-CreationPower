# Lumina Runtime Container (LRC) - 多语言支持容器

# 使用Ubuntu作为基础镜像
FROM ubuntu:20.04

# 设置非交互式安装
ENV DEBIAN_FRONTEND=noninteractive

# 维护者信息
LABEL maintainer="Code Weaver <270586352@qq.com>"
LABEL version="1.0.0"
LABEL description="Lumina Runtime Container - 多语言支持的容器运行环境"

# 设置环境变量
ENV LRC_VERSION=1.0.0
ENV LRC_HOME=/app
ENV LRC_LANG=python
ENV LRC_MEMORY_LIMIT=4096
ENV LRC_CPU_LIMIT=2
ENV LRC_GPU_ENABLED=false
ENV LRC_LOG_LEVEL=info

# 设置工作目录
WORKDIR $LRC_HOME

# 更新包管理器并安装基础依赖
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    curl \
    wget \
    unzip \
    tar \
    gzip \
    bzip2 \
    libssl-dev \
    libffi-dev \
    zlib1g-dev \
    libbz2-dev \
    libreadline-dev \
    libsqlite3-dev \
    libncurses5-dev \
    libncursesw5-dev \
    libffi-dev \
    liblzma-dev \
    tzdata \
    locales \
    && rm -rf /var/lib/apt/lists/*

# 设置时区和语言
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    locale-gen en_US.UTF-8 && \
    update-locale LANG=en_US.UTF-8

ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US:en
ENV LC_ALL=en_US.UTF-8

# 安装C/C++开发环境
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    gdb \
    make \
    && rm -rf /var/lib/apt/lists/*

# 安装Python 3.9
RUN apt-get update && apt-get install -y \
    python3 \
    python3-pip \
    python3-dev \
    python3-venv \
    && rm -rf /var/lib/apt/lists/*

# 升级pip
RUN pip3 install --no-cache-dir --upgrade pip setuptools wheel

# 安装Python AI框架和常用库
RUN pip3 install --no-cache-dir \
    tensorflow \
    torch torchvision torchaudio \
    onnxruntime \
    numpy \
    pandas \
    scikit-learn \
    matplotlib \
    requests \
    flask \
    fastapi \
    uvicorn \
    pymongo \
    redis \
    sqlalchemy

# 安装Node.js 16
RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash - && \
    apt-get install -y nodejs && \
    npm install -g npm@latest && \
    rm -rf /var/lib/apt/lists/*

# 安装常用Node.js包
RUN npm install -g \
    express \
    koa \
    react \
    vue \
    angular \
    webpack \
    typescript \
    ts-node \
    nodemon

# 安装OpenJDK 11
RUN apt-get update && apt-get install -y \
    openjdk-11-jdk \
    maven \
    gradle \
    && rm -rf /var/lib/apt/lists/*

# 安装Go 1.17
RUN wget -O go.tar.gz https://golang.org/dl/go1.17.13.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV GOROOT=/usr/local/go
ENV PATH=$PATH:$GOPATH/bin

# 创建Go工作目录
RUN mkdir -p $GOPATH/src $GOPATH/bin

# 安装常用Go包
RUN go get -u \
    github.com/gin-gonic/gin \
    github.com/go-redis/redis/v8 \
    github.com/mongodb/mongo-go-driver \
    github.com/jinzhu/gorm \
    github.com/golang-migrate/migrate/v4 \
    golang.org/x/crypto \
    google.golang.org/grpc \
    github.com/spf13/cobra

# 安装Docker CLI (用于在容器内管理其他容器)
RUN apt-get update && apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common && \
    curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
    apt-get update && apt-get install -y docker-ce-cli && \
    rm -rf /var/lib/apt/lists/*

# 安装Kubernetes工具
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl && \
    rm kubectl

# 安装Helm
RUN curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 && \
    chmod 700 get_helm.sh && \
    ./get_helm.sh && \
    rm get_helm.sh

# 安装监控和日志工具
RUN apt-get update && apt-get install -y \
    prometheus-client \
    fluentd \
    && rm -rf /var/lib/apt/lists/*

# 创建项目目录结构
RUN mkdir -p \
    $LRC_HOME/src/core \
    $LRC_HOME/src/languages \
    $LRC_HOME/src/services \
    $LRC_HOME/config \
    $LRC_HOME/docs \
    $LRC_HOME/examples \
    $LRC_HOME/scripts \
    $LRC_HOME/logs \
    $LRC_HOME/data

# 复制配置文件
COPY config/ $LRC_HOME/config/

# 复制示例代码
COPY examples/ $LRC_HOME/examples/

# 复制脚本
COPY scripts/ $LRC_HOME/scripts/

# 复制源码
COPY src/ $LRC_HOME/src/

# 给脚本添加执行权限
RUN chmod +x $LRC_HOME/scripts/*.sh

# 设置环境变量
ENV PATH=$PATH:$LRC_HOME/scripts

# 暴露端口
EXPOSE 8080 9090 22

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# 默认命令
CMD ["bash", "-c", "echo 'Lumina Runtime Container started' && tail -f /dev/null"]
