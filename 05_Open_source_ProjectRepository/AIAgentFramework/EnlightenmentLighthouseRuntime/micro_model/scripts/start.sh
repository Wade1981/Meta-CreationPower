#!/bin/bash

# 启动微模型服务器

echo "Starting micro model server..."

# 检查是否在正确的目录
if [ ! -f "main.go" ]; then
    echo "Error: main.go not found. Please run this script in the micro_model directory."
    exit 1
fi

# 构建并运行服务器
echo "Building server..."
go build -o micro_model_server .

if [ $? -eq 0 ]; then
    echo "Build successful. Starting server..."
    ./micro_model_server
else
    echo "Build failed."
    exit 1
fi
