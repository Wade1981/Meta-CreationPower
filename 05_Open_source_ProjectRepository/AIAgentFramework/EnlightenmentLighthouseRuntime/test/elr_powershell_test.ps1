# ELR PowerShell测试脚本

# 测试目录
$testDir = "./test_output"
if (!(Test-Path $testDir)) {
    New-Item -ItemType Directory -Path $testDir -Force
}

# 测试日志
$logFile = "$testDir/elr_test_log.txt"
"$(Get-Date) - ELR PowerShell测试开始" | Out-File -FilePath $logFile -Append

# 测试1: 版本信息测试
Write-Host "=== 测试1: 版本信息测试 ==="
"$(Get-Date) - 测试1: 版本信息测试" | Out-File -FilePath $logFile -Append
try {
    $version = .\elr.ps1 version
    Write-Host "版本信息: $version"
    "版本信息: $version" | Out-File -FilePath $logFile -Append
    Write-Host "✓ 版本信息测试通过"
} catch {
    Write-Host "✗ 版本信息测试失败: $($_.Exception.Message)"
    "版本信息测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试2: 容器列表测试
Write-Host "\n=== 测试2: 容器列表测试 ==="
"$(Get-Date) - 测试2: 容器列表测试" | Out-File -FilePath $logFile -Append
try {
    $containers = .\elr.ps1 container list
    Write-Host "容器列表: $containers"
    "容器列表: $containers" | Out-File -FilePath $logFile -Append
    Write-Host "✓ 容器列表测试通过"
} catch {
    Write-Host "✗ 容器列表测试失败: $($_.Exception.Message)"
    "容器列表测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试3: 网络状态测试
Write-Host "\n=== 测试3: 网络状态测试 ==="
"$(Get-Date) - 测试3: 网络状态测试" | Out-File -FilePath $logFile -Append
try {
    $networkStatus = .\elr.ps1 network status
    Write-Host "网络状态: $networkStatus"
    "网络状态: $networkStatus" | Out-File -FilePath $logFile -Append
    Write-Host "✓ 网络状态测试通过"
} catch {
    Write-Host "✗ 网络状态测试失败: $($_.Exception.Message)"
    "网络状态测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试4: API服务启动测试
Write-Host "\n=== 测试4: API服务启动测试 ==="
"$(Get-Date) - 测试4: API服务启动测试" | Out-File -FilePath $logFile -Append
try {
    $apiStart = .\elr.ps1 api start
    Write-Host "API服务启动: $apiStart"
    "API服务启动: $apiStart" | Out-File -FilePath $logFile -Append
    Write-Host "✓ API服务启动测试通过"
    
    # 等待服务启动
    Start-Sleep -Seconds 2
    
    # 测试API服务状态
    $apiStatus = .\elr.ps1 api status
    Write-Host "API服务状态: $apiStatus"
    "API服务状态: $apiStatus" | Out-File -FilePath $logFile -Append
    
    # 停止API服务
    $apiStop = .\elr.ps1 api stop
    Write-Host "API服务停止: $apiStop"
    "API服务停止: $apiStop" | Out-File -FilePath $logFile -Append
} catch {
    Write-Host "✗ API服务测试失败: $($_.Exception.Message)"
    "API服务测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试5: 容器创建和管理测试
Write-Host "\n=== 测试5: 容器创建和管理测试 ==="
"$(Get-Date) - 测试5: 容器创建和管理测试" | Out-File -FilePath $logFile -Append
try {
    # 创建测试容器
    $containerID = "test-container-ps"
    $createResult = .\elr.ps1 container create --id $containerID --name "PowerShell Test Container" --image "ubuntu:latest"
    Write-Host "容器创建: $createResult"
    "容器创建: $createResult" | Out-File -FilePath $logFile -Append
    
    # 启动容器
    $startResult = .\elr.ps1 container start --id $containerID
    Write-Host "容器启动: $startResult"
    "容器启动: $startResult" | Out-File -FilePath $logFile -Append
    
    # 查看容器状态
    $statusResult = .\elr.ps1 container status --id $containerID
    Write-Host "容器状态: $statusResult"
    "容器状态: $statusResult" | Out-File -FilePath $logFile -Append
    
    # 停止容器
    $stopResult = .\elr.ps1 container stop --id $containerID
    Write-Host "容器停止: $stopResult"
    "容器停止: $stopResult" | Out-File -FilePath $logFile -Append
    
    # 删除容器
    $deleteResult = .\elr.ps1 container delete --id $containerID
    Write-Host "容器删除: $deleteResult"
    "容器删除: $deleteResult" | Out-File -FilePath $logFile -Append
    
    Write-Host "✓ 容器创建和管理测试通过"
} catch {
    Write-Host "✗ 容器创建和管理测试失败: $($_.Exception.Message)"
    "容器创建和管理测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试6: 网络隔离测试
Write-Host "\n=== 测试6: 网络隔离测试 ==="
"$(Get-Date) - 测试6: 网络隔离测试" | Out-File -FilePath $logFile -Append
try {
    # 创建测试容器
    $containerID = "test-container-network"
    $createResult = .\elr.ps1 container create --id $containerID --name "Network Test Container" --image "ubuntu:latest"
    
    # 启动容器
    .\elr.ps1 container start --id $containerID
    
    # 应用网络隔离
    $isolateResult = .\elr.ps1 network isolate --container-id $containerID
    Write-Host "网络隔离: $isolateResult"
    "网络隔离: $isolateResult" | Out-File -FilePath $logFile -Append
    
    # 获取网络配置
    $configResult = .\elr.ps1 network config --container-id $containerID
    Write-Host "网络配置: $configResult"
    "网络配置: $configResult" | Out-File -FilePath $logFile -Append
    
    # 移除网络隔离
    $unisolateResult = .\elr.ps1 network unisolate --container-id $containerID
    Write-Host "移除网络隔离: $unisolateResult"
    "移除网络隔离: $unisolateResult" | Out-File -FilePath $logFile -Append
    
    # 清理容器
    .\elr.ps1 container stop --id $containerID
    .\elr.ps1 container delete --id $containerID
    
    Write-Host "✓ 网络隔离测试通过"
} catch {
    Write-Host "✗ 网络隔离测试失败: $($_.Exception.Message)"
    "网络隔离测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试7: 模型管理测试
Write-Host "\n=== 测试7: 模型管理测试 ==="
"$(Get-Date) - 测试7: 模型管理测试" | Out-File -FilePath $logFile -Append
try {
    # 列出模型
    $models = .\elr.ps1 model list
    Write-Host "模型列表: $models"
    "模型列表: $models" | Out-File -FilePath $logFile -Append
    
    Write-Host "✓ 模型管理测试通过"
} catch {
    Write-Host "✗ 模型管理测试失败: $($_.Exception.Message)"
    "模型管理测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
}

# 测试完成
Write-Host "\n=== 测试完成 ==="
"$(Get-Date) - ELR PowerShell测试完成" | Out-File -FilePath $logFile -Append
Write-Host "测试日志已保存到: $logFile"
