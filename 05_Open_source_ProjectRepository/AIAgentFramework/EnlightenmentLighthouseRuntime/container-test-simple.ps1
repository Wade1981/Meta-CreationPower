#!/usr/bin/env powershell

# 关闭可能运行的elr-container.exe进程
try {
    $processes = Get-Process | Where-Object { $_.Name -eq "elr-container" }
    foreach ($process in $processes) {
        Write-Host "Stopping existing elr-container.exe process: $($process.Id)"
        $process.Kill()
    }
} catch {
    Write-Host "Error stopping existing processes: $($_.Exception.Message)"
}

Write-Host "===================================="
Write-Host "Testing ELR Container Management"
Write-Host "===================================="

# 1. 启动elr-container.exe
Write-Host ""
Write-Host "1. Starting elr-container.exe..."
try {
    $process = Start-Process -FilePath ".\elr-container.exe" -ArgumentList "start" -NoNewWindow -PassThru
    Write-Host "elr-container.exe started successfully! Process ID: $($process.Id)"
    Start-Sleep -Seconds 3
} catch {
    Write-Host "Error starting elr-container.exe: $($_.Exception.Message)"
    exit 1
}

# 2. 创建容器
Write-Host ""
Write-Host "2. Creating container..."
try {
    $createOutput = & ".\elr-container.exe" create --name test-container --image elr-chat
    Write-Host $createOutput
} catch {
    Write-Host "Error creating container: $($_.Exception.Message)"
}

# 3. 列出容器
Write-Host ""
Write-Host "3. Listing containers..."
try {
    $listOutput = & ".\elr-container.exe" list
    Write-Host $listOutput
} catch {
    Write-Host "Error listing containers: $($_.Exception.Message)"
}

# 4. 启动容器（使用固定的容器ID格式）
Write-Host ""
Write-Host "4. Starting container..."
try {
    # 尝试使用常见的容器ID格式
    $containerId = "elr-1774210872415012900"
    $startOutput = & ".\elr-container.exe" start-container --id $containerId
    Write-Host $startOutput
} catch {
    Write-Host "Error starting container: $($_.Exception.Message)"
}

# 5. 再次列出容器
Write-Host ""
Write-Host "5. Listing containers again..."
try {
    $listOutput = & ".\elr-container.exe" list
    Write-Host $listOutput
} catch {
    Write-Host "Error listing containers: $($_.Exception.Message)"
}

# 6. 停止容器
Write-Host ""
Write-Host "6. Stopping container..."
try {
    $stopOutput = & ".\elr-container.exe" stop-container --id $containerId
    Write-Host $stopOutput
} catch {
    Write-Host "Error stopping container: $($_.Exception.Message)"
}

# 7. 再次列出容器
Write-Host ""
Write-Host "7. Listing containers after stopping..."
try {
    $listOutput = & ".\elr-container.exe" list
    Write-Host $listOutput
} catch {
    Write-Host "Error listing containers: $($_.Exception.Message)"
}

# 8. 停止elr-container.exe
Write-Host ""
Write-Host "8. Stopping elr-container.exe..."
try {
    $stopOutput = & ".\elr-container.exe" stop
    Write-Host $stopOutput
} catch {
    Write-Host "Error stopping elr-container.exe: $($_.Exception.Message)"
}

Write-Host ""
Write-Host "===================================="
Write-Host "Container management test completed!"
Write-Host "===================================="
