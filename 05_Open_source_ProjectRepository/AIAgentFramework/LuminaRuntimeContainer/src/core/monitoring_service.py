#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
Lumina Runtime Container - 监控服务
"""

import os
import sys
import time
import psutil
from flask import Flask
from prometheus_client import start_http_server, Gauge, Counter, Summary

# 创建Flask应用
app = Flask(__name__)

# 服务版本
SERVICE_VERSION = "1.0.0"

# 定义监控指标

# 系统资源指标
CPU_USAGE = Gauge('cpu_usage_percent', 'CPU usage percentage')
MEMORY_USAGE = Gauge('memory_usage_percent', 'Memory usage percentage')
DISK_USAGE = Gauge('disk_usage_percent', 'Disk usage percentage')

# 网络指标
NETWORK_BYTES_SENT = Counter('network_bytes_sent', 'Network bytes sent')
NETWORK_BYTES_RECEIVED = Counter('network_bytes_received', 'Network bytes received')

# 服务指标
API_REQUESTS = Counter('api_requests_total', 'Total API requests')
API_REQUEST_LATENCY = Summary('api_request_latency_seconds', 'API request latency')

# 容器指标
CONTAINER_UPTIME = Gauge('container_uptime_seconds', 'Container uptime in seconds')
CONTAINER_STATUS = Gauge('container_status', 'Container status (1=running, 0=stopped)')

# 语言运行时指标
LANGUAGE_RUNTIME_STATUS = Gauge('language_runtime_status', 'Language runtime status', ['language'])

# 启动时间
START_TIME = time.time()

# 函数：更新系统指标
def update_system_metrics():
    """更新系统指标"""
    # CPU使用率
    cpu_percent = psutil.cpu_percent(interval=1)
    CPU_USAGE.set(cpu_percent)
    
    # 内存使用率
    memory = psutil.virtual_memory()
    memory_percent = memory.percent
    MEMORY_USAGE.set(memory_percent)
    
    # 磁盘使用率
    disk = psutil.disk_usage('/')
    disk_percent = disk.percent
    DISK_USAGE.set(disk_percent)
    
    # 网络指标
    net_io = psutil.net_io_counters()
    NETWORK_BYTES_SENT.inc(net_io.bytes_sent)
    NETWORK_BYTES_RECEIVED.inc(net_io.bytes_recv)
    
    # 容器指标
    uptime = time.time() - START_TIME
    CONTAINER_UPTIME.set(uptime)
    CONTAINER_STATUS.set(1)
    
    # 语言运行时状态
    languages = ["cpp", "python", "nodejs", "java", "go"]
    for lang in languages:
        LANGUAGE_RUNTIME_STATUS.labels(language=lang).set(1)

# 定时更新指标
def metrics_updater():
    """定时更新指标"""
    while True:
        try:
            update_system_metrics()
        except Exception as e:
            print(f"Error updating metrics: {e}")
        time.sleep(15)  # 每15秒更新一次

# 根路径
@app.route('/', methods=['GET'])
def root():
    """根路径"""
    return f"Lumina Runtime Container Monitoring Service v{SERVICE_VERSION}"

# 主函数
if __name__ == '__main__':
    # 获取端口
    port = int(os.environ.get('MONITORING_PORT', 9090))
    
    # 启动Prometheus HTTP服务器
    print(f"Starting Prometheus metrics server on port {port}...")
    start_http_server(port)
    
    # 启动指标更新线程
    import threading
    updater_thread = threading.Thread(target=metrics_updater, daemon=True)
    updater_thread.start()
    
    # 启动Flask应用
    print(f"Starting monitoring service...")
    app.run(host='0.0.0.0', port=port+1, debug=False)
