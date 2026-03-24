#!/usr/bin/env powershell

# ELR API测试脚本
# 测试Desktop API、Public API和模型API的可用性

Write-Host "=== ELR API测试脚本 ==="
Write-Host "测试四个API的可用性和功能"
Write-Host "=============================================="

# 定义API地址和端口
$desktopAPI = "http://localhost:8081"
$publicAPI = "http://localhost:8080"
$modelServiceAPI = "http://localhost:8082"
$microModelServerAPI = "http://localhost:8083"

# 测试函数
function Test-API {
    param(
        [string]$apiName,
        [string]$apiUrl,
        [string]$endpoint,
        [string]$method = "GET"
    )
    
    Write-Host "\n测试 $apiName - $endpoint"
    Write-Host "URL: $apiUrl$endpoint"
    
    try {
        if ($method -eq "GET") {
            $response = Invoke-RestMethod -Uri "$apiUrl$endpoint" -Method GET -ErrorAction Stop
            Write-Host "✓ 成功: $apiName $endpoint 响应正常"
            Write-Host "响应:"
            $response | ConvertTo-Json -Depth 3
            return $true
        } else {
            $response = Invoke-RestMethod -Uri "$apiUrl$endpoint" -Method POST -ErrorAction Stop
            Write-Host "✓ 成功: $apiName $endpoint 响应正常"
            Write-Host "响应:"
            $response | ConvertTo-Json -Depth 3
            return $true
        }
    } catch {
        Write-Host "✗ 失败: $apiName $endpoint 无法访问"
        Write-Host "错误: $($_.Exception.Message)"
        return $false
    }
}

# 测试Desktop API
Write-Host "\n=== 测试 Desktop API ==="
Test-API -apiName "Desktop API" -apiUrl $desktopAPI -endpoint "/api/desktop/health"
Test-API -apiName "Desktop API" -apiUrl $desktopAPI -endpoint "/api/desktop/status"
Test-API -apiName "Desktop API" -apiUrl $desktopAPI -endpoint "/api/desktop/containers"
Test-API -apiName "Desktop API" -apiUrl $desktopAPI -endpoint "/api/desktop/resources"

# 测试Public API
Write-Host "\n=== 测试 Public API ==="
Test-API -apiName "Public API" -apiUrl $publicAPI -endpoint "/health"
Test-API -apiName "Public API" -apiUrl $publicAPI -endpoint "/api/status"
Test-API -apiName "Public API" -apiUrl $publicAPI -endpoint "/api/network/status"
Test-API -apiName "Public API" -apiUrl $publicAPI -endpoint "/api/container/list"
Test-API -apiName "Public API" -apiUrl $publicAPI -endpoint "/api/model/list"

# 测试Model Service API
Write-Host "\n=== 测试 Model Service API ==="
Test-API -apiName "Model Service API" -apiUrl $modelServiceAPI -endpoint "/health"
Test-API -apiName "Model Service API" -apiUrl $modelServiceAPI -endpoint "/api/models"

# 测试Micro Model Server API
Write-Host "\n=== 测试 Micro Model Server API ==="
Test-API -apiName "Micro Model Server API" -apiUrl $microModelServerAPI -endpoint "/health"
Test-API -apiName "Micro Model Server API" -apiUrl $microModelServerAPI -endpoint "/api/models"

# 推荐的测试命令
Write-Host "\n=== 推荐的测试命令 ==="
Write-Host "\n1. Desktop API 测试命令:"
Write-Host "   # 健康检查"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8081/api/desktop/health' -Method GET"
Write-Host "   # 获取状态"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8081/api/desktop/status' -Method GET"
Write-Host "   # 获取容器列表"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8081/api/desktop/containers' -Method GET"
Write-Host "   # 获取资源使用情况"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8081/api/desktop/resources' -Method GET"

Write-Host "\n2. Public API 测试命令:"
Write-Host "   # 健康检查"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8080/health' -Method GET"
Write-Host "   # 获取网络状态"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8080/api/network/status' -Method GET"
Write-Host "   # 获取容器列表"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8080/api/container/list' -Method GET"
Write-Host "   # 获取模型列表"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8080/api/model/list' -Method GET"

Write-Host "\n3. Model Service API 测试命令:"
Write-Host "   # 健康检查"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8082/health' -Method GET"
Write-Host "   # 获取模型列表"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8082/api/models' -Method GET"

Write-Host "\n4. Micro Model Server API 测试命令:"
Write-Host "   # 健康检查"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8083/health' -Method GET"
Write-Host "   # 获取模型列表"
Write-Host "   Invoke-RestMethod -Uri 'http://localhost:8083/api/models' -Method GET"

Write-Host "\n=== API测试完成 ==="
