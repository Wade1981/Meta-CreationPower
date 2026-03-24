#!/usr/bin/env powershell

# ELR Chat Command
# 启动与elr-chat模型的对话

Write-Host "===================================="
Write-Host "ELR Chat Command"
Write-Host "===================================="

# 检查Model Service是否运行
Write-Host "Checking Model Service status..."
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8082/health" -Method Get
    Write-Host "Model Service is running"
} catch {
    Write-Host "Error: Model Service is not running"
    Write-Host "Please start Model Service first: .\elr.ps1 start-model"
    exit 1
}

# 检查elr-chat模型是否存在
Write-Host "Checking elr-chat model..."
try {
    $models = Invoke-RestMethod -Uri "http://localhost:8082/api/models/" -Method Get
    $elrChatModel = $models | Where-Object { $_.id -eq "elr-chat" }
    if (-not $elrChatModel) {
        Write-Host "Error: elr-chat model not found"
        exit 1
    }
    Write-Host "elr-chat model found: $($elrChatModel.name) v$($elrChatModel.version)"
} catch {
    Write-Host "Error: Failed to check elr-chat model: $($_.Exception.Message)"
    exit 1
}

Write-Host ""
Write-Host "ELR Chat started!"
Write-Host "Type 'exit' to quit"
Write-Host "===================================="

# 交互式对话
while ($true) {
    Write-Host -NoNewline "You: "
    $inputMessage = Read-Host
    
    if ($inputMessage -eq "exit") {
        break
    }
    
    try {
        # 发送聊天请求
        $chatInput = $inputMessage
        $jsonBody = '{"container_name": "test-container", "model_id": "elr-chat", "input": "' + $chatInput + '"}'
        $response = Invoke-RestMethod -Uri "http://localhost:8082/api/models/run" -Method Post -Body $jsonBody -ContentType "application/json"
        
        # 提取响应
        $output = $response.output
        Write-Host "ELR: $output"
    } catch {
        Write-Host "Error: Failed to send message: $($_.Exception.Message)"
    }
    
    Write-Host "===================================="
}

Write-Host "ELR Chat ended."
Write-Host "===================================="
