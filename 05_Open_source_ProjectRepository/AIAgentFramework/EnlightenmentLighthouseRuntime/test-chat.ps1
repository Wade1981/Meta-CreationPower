#!/usr/bin/env powershell

# 测试聊天功能
Write-Host "=== 测试聊天功能 ==="

# 测试消息
$testMessage = "你好，ELR容器"

# 发送聊天请求
try {
    $jsonBody = '{"container_name": "test-container", "model_id": "elr-chat", "input": "' + $testMessage + '"}'
    $response = Invoke-RestMethod -Uri "http://localhost:8082/api/models/run" -Method Post -Body $jsonBody -ContentType "application/json"
    
    # 提取响应
    $output = $response.output
    Write-Host "输入: $testMessage"
    Write-Host "输出: $output"
} catch {
    Write-Host "Error: Failed to send message: $($_.Exception.Message)"
}

Write-Host "=== 测试完成 ==="
