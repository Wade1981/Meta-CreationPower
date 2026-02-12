#!/bin/bash

# Lumina Runtime Container - 启动脚本

set -e

# 脚本版本
SCRIPT_VERSION="1.0.0"

# 工作目录
WORKDIR="/app"

# 日志目录
LOG_DIR="$WORKDIR/logs"
mkdir -p "$LOG_DIR"

# 日志文件
LOG_FILE="$LOG_DIR/lumina-start.log"

# 函数：记录日志
log() {
    local timestamp=$(date "%Y-%m-%d %H:%M:%S")
    local level="$1"
    local message="$2"
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

# 函数：检查命令是否存在
check_command() {
    if ! command -v "$1" &> /dev/null; then
        log "ERROR" "Command '$1' not found"
        return 1
    fi
    return 0
}

# 函数：检查端口是否可用
check_port() {
    local port="$1"
    if lsof -i:$port &> /dev/null; then
        log "ERROR" "Port $port is already in use"
        return 1
    fi
    return 0
}

# 函数：启动服务
start_service() {
    local service_name="$1"
    local service_command="$2"
    local service_log="$LOG_DIR/${service_name}.log"
    
    log "INFO" "Starting $service_name..."
    
    # 检查命令是否存在
    check_command "$(echo "$service_command" | awk '{print $1}')"
    
    # 启动服务
    nohup $service_command > "$service_log" 2>&1 &
    
    # 等待服务启动
    sleep 2
    
    # 检查服务是否启动成功
    if ps aux | grep "$service_command" | grep -v grep &> /dev/null; then
        log "INFO" "$service_name started successfully"
        return 0
    else
        log "ERROR" "Failed to start $service_name"
        log "ERROR" "Check $service_log for details"
        return 1
    fi
}

# 主函数
main() {
    log "INFO" "Starting Lumina Runtime Container v$SCRIPT_VERSION"
    log "INFO" "Working directory: $WORKDIR"
    
    # 检查必要的命令
    log "INFO" "Checking required commands..."
    check_command "docker"
    check_command "docker-compose"
    
    # 检查端口
    log "INFO" "Checking ports..."
    check_port 8080
    check_port 9090
    check_port 2222
    check_port 27017
    check_port 6379
    
    # 启动容器
    log "INFO" "Starting containers..."
    cd "$WORKDIR"
    docker-compose up -d
    
    # 等待容器启动
    log "INFO" "Waiting for containers to start..."
    sleep 10
    
    # 检查容器状态
    log "INFO" "Checking container status..."
    docker ps
    
    # 显示访问信息
    log "INFO" "===================================="
    log "INFO" "Lumina Runtime Container started successfully!"
    log "INFO" "===================================="
    log "INFO" "API Service: http://localhost:8080"
    log "INFO" "Monitoring Service: http://localhost:9090"
    log "INFO" "SSH Access: ssh -p 2222 root@localhost"
    log "INFO" "MongoDB: mongodb://localhost:27017"
    log "INFO" "Redis: redis://localhost:6379"
    log "INFO" "===================================="
    log "INFO" "To stop the containers, run: docker-compose down"
    log "INFO" "To view logs, run: docker-compose logs"
    log "INFO" "===================================="
}

# 执行主函数
main
