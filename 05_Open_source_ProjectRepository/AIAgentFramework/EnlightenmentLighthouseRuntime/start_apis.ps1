#!/usr/bin/env powershell

# ELR API服务启动脚本
# 启动Desktop API、Public API和模型API服务

Write-Host "=== ELR API服务启动脚本 ==="
Write-Host "启动四个API服务"
Write-Host "=============================================="

# 启动Desktop API
Write-Host "`n1. 启动 Desktop API (端口: 8081)"
try {
    Start-Process python -ArgumentList "elr\desktop_api.py" -NoNewWindow -PassThru
    Write-Host "✓ Desktop API 服务已启动"
} catch {
    Write-Host "✗ 启动 Desktop API 失败: $($_.Exception.Message)"
}

# 启动Public API
Write-Host "`n2. 启动 Public API (端口: 8080)"
try {
    Start-Process go -ArgumentList "run", "main.go" -WorkingDirectory "elr\network_service" -NoNewWindow -PassThru
    Write-Host "✓ Public API 服务已启动"
} catch {
    Write-Host "✗ 启动 Public API 失败: $($_.Exception.Message)"
}

# 启动Model Service API
Write-Host "`n3. 启动 Model Service API (端口: 8082)"
try {
    Start-Process go -ArgumentList "run", "main.go", "8082" -WorkingDirectory "micro_model" -NoNewWindow -PassThru
    Write-Host "✓ Model Service API 服务已启动"
} catch {
    Write-Host "✗ 启动 Model Service API 失败: $($_.Exception.Message)"
}

# 启动Micro Model Server API
Write-Host "`n4. 启动 Micro Model Server API (端口: 8083)"
try {
    Start-Process go -ArgumentList "run", "main.go", "8083" -WorkingDirectory "micro_model" -NoNewWindow -PassThru
    Write-Host "✓ Micro Model Server API 服务已启动"
} catch {
    Write-Host "✗ 启动 Micro Model Server API 失败: $($_.Exception.Message)"
}

Write-Host "`n=== API服务启动完成 ==="
Write-Host "请等待几秒钟让服务完全启动，然后运行 test_apis.ps1 进行测试"
