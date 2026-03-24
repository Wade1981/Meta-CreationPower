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

# 启动elr-container.exe在后台
Write-Host "===================================="
Write-Host "Starting elr-container.exe in background..."
Write-Host "===================================="
try {
    $process = Start-Process -FilePath ".\elr-container.exe" -ArgumentList "start" -NoNewWindow -PassThru
    Write-Host "elr-container.exe started successfully!"
    Write-Host "Process ID: $($process.Id)"
    # 等待几秒钟让服务启动
    Start-Sleep -Seconds 3
} catch {
    Write-Host "Error starting elr-container.exe: $($_.Exception.Message)"
    exit 1
}

# 检查容器列表
Write-Host ""
Write-Host "===================================="
Write-Host "Checking container list..."
Write-Host "===================================="
try {
    $containersOutput = & ".\elr-container.exe" list
    Write-Host $containersOutput
    
    # 检查是否有容器
    if ($containersOutput -match "No containers found") {
        Write-Host ""
        Write-Host "===================================="
        Write-Host "No containers found. Creating a new container..."
        Write-Host "===================================="
        
        # 创建容器
        $createOutput = & ".\elr-container.exe" create --name test-container --image elr-chat
        Write-Host $createOutput
        
        # 提取容器ID
        if ($createOutput -match "ID: (elr-\d+),") {
            $containerId = $matches[1]
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Starting container $containerId..."
            Write-Host "===================================="
            
            # 启动容器
            $startOutput = & ".\elr-container.exe" start-container --id $containerId
            Write-Host $startOutput
            
            # 检查容器状态
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Checking container status..."
            Write-Host "===================================="
            $statusOutput = & ".\elr-container.exe" list
            Write-Host $statusOutput
            
            # 关闭容器
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Stopping container $containerId..."
            Write-Host "===================================="
            $stopOutput = & ".\elr-container.exe" stop-container --id $containerId
            Write-Host $stopOutput
            
            # 再次检查容器状态
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Checking container status after stopping..."
            Write-Host "===================================="
            $finalStatusOutput = & ".\elr-container.exe" list
            Write-Host $finalStatusOutput
        } else {
            Write-Host "Error: Could not extract container ID"
            # 手动设置容器ID
            $containerId = "elr-1774210872415012900"
            Write-Host "Using manual container ID: $containerId"
        }
    } else {
        Write-Host ""
        Write-Host "===================================="
        Write-Host "Containers found. Using existing container..."
        Write-Host "===================================="
        
        # 提取第一个容器ID
        if ($containersOutput -match "(elr-\d+)\s+") {
            $containerId = $matches[1]
            Write-Host "Using container: $containerId"
            
            # 启动容器
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Starting container $containerId..."
            Write-Host "===================================="
            $startOutput = & ".\elr-container.exe" start-container --id $containerId
            Write-Host $startOutput
            
            # 检查容器状态
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Checking container status..."
            Write-Host "===================================="
            $statusOutput = & ".\elr-container.exe" list
            Write-Host $statusOutput
            
            # 关闭容器
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Stopping container $containerId..."
            Write-Host "===================================="
            $stopOutput = & ".\elr-container.exe" stop-container --id $containerId
            Write-Host $stopOutput
            
            # 再次检查容器状态
            Write-Host ""
            Write-Host "===================================="
            Write-Host "Checking container status after stopping..."
            Write-Host "===================================="
            $finalStatusOutput = & ".\elr-container.exe" list
            Write-Host $finalStatusOutput
        } else {
            Write-Host "Error: Could not extract container ID"
        }
    }
} catch {
    Write-Host "Error checking containers: $($_.Exception.Message)"
}

# 停止elr-container.exe
Write-Host ""
Write-Host "===================================="
Write-Host "Stopping elr-container.exe..."
Write-Host "===================================="
try {
    $stopOutput = & ".\elr-container.exe" stop
    Write-Host $stopOutput
} catch {
    Write-Host "Error stopping elr-container.exe: $($_.Exception.Message)"
}

Write-Host ""
Write-Host "===================================="
Write-Host "Container test completed!"
Write-Host "===================================="
