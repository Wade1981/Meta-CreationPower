#!/bin/bash

# Lumina Runtime Container - 停止脚本

set -e

# 脚本版本
SCRIPT_VERSION="1.0.0"

# 工作目录
WORKDIR="/app"

# 日志目录
LOG_DIR="$WORKDIR/logs"
mkdir -p "$LOG_DIR"

# 日志文件
LOG_FILE="$LOG_DIR/lumina-stop.log"

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

# 主函数
main() {
    log "INFO" "Stopping Lumina Runtime Container v$SCRIPT_VERSION"
    log "INFO" "Working directory: $WORKDIR"
    
    # 检查必要的命令
    log "INFO" "Checking required commands..."
    check_command "docker"
    check_command "docker-compose"
    
    # 停止容器
    log "INFO" "Stopping containers..."
    cd "$WORKDIR"
    docker-compose down
    
    # 等待容器停止
    log "INFO" "Waiting for containers to stop..."
    sleep 5
    
    # 检查容器状态
    log "INFO" "Checking container status..."
    docker ps
    
    # 显示停止信息
    log "INFO" "===================================="
    log "INFO" "Lumina Runtime Container stopped successfully!"
    log "INFO" "===================================="
    log "INFO" "All containers have been stopped and removed."
    log "INFO" "To start the containers again, run: ./scripts/start.sh"
    log "INFO" "===================================="
}

# 执行主函数
main
