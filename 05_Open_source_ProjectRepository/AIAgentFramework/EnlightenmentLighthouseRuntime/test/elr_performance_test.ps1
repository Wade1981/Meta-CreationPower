# ELR性能测试脚本

# 测试目录
$testDir = "./test_output"
if (!(Test-Path $testDir)) {
    New-Item -ItemType Directory -Path $testDir -Force
}

# 测试日志
$logFile = "$testDir/elr_performance_test_log.txt"
"$(Get-Date) - ELR性能测试开始" | Out-File -FilePath $logFile -Append

# 测试结果
$resultsFile = "$testDir/elr_performance_results.json"
$results = @{}

# 1. 容器启动性能测试
Write-Host "=== 测试1: 容器启动性能测试 ==="
"$(Get-Date) - 测试1: 容器启动性能测试" | Out-File -FilePath $logFile -Append

$startupTimes = @()
for ($i = 1; $i -le 5; $i++) {
    $containerID = "test-performance-$i"
    
    # 创建容器
    $createStart = Get-Date
    $createResult = .\elr.ps1 container create --id $containerID --name "Performance Test Container $i" --image "ubuntu:latest"
    $createEnd = Get-Date
    $createTime = ($createEnd - $createStart).TotalSeconds
    
    # 启动容器
    $startStart = Get-Date
    $startResult = .\elr.ps1 container start --id $containerID
    $startEnd = Get-Date
    $startTime = ($startEnd - $startStart).TotalSeconds
    
    # 总时间
    $totalTime = $createTime + $startTime
    $startupTimes += $totalTime
    
    Write-Host "容器 $i: 创建时间 = $createTime 秒, 启动时间 = $startTime 秒, 总时间 = $totalTime 秒"
    "容器 $i: 创建时间 = $createTime 秒, 启动时间 = $startTime 秒, 总时间 = $totalTime 秒" | Out-File -FilePath $logFile -Append
    
    # 停止并删除容器
    .\elr.ps1 container stop --id $containerID
    .\elr.ps1 container delete --id $containerID
}

# 计算平均值
$avgStartupTime = ($startupTimes | Measure-Object -Average).Average
$minStartupTime = ($startupTimes | Measure-Object -Minimum).Minimum
$maxStartupTime = ($startupTimes | Measure-Object -Maximum).Maximum

Write-Host "平均启动时间: $avgStartupTime 秒"
Write-Host "最小启动时间: $minStartupTime 秒"
Write-Host "最大启动时间: $maxStartupTime 秒"
"平均启动时间: $avgStartupTime 秒" | Out-File -FilePath $logFile -Append
"最小启动时间: $minStartupTime 秒" | Out-File -FilePath $logFile -Append
"最大启动时间: $maxStartupTime 秒" | Out-File -FilePath $logFile -Append

$results["container_startup"] = @{
    "average_time" = $avgStartupTime
    "min_time" = $minStartupTime
    "max_time" = $maxStartupTime
    "times" = $startupTimes
}

# 2. API响应性能测试
Write-Host "\n=== 测试2: API响应性能测试 ==="
"$(Get-Date) - 测试2: API响应性能测试" | Out-File -FilePath $logFile -Append

# 启动API服务
.\elr.ps1 api start
Start-Sleep -Seconds 2

$apiResponseTimes = @()
for ($i = 1; $i -le 10; $i++) {
    $startTime = Get-Date
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
        $endTime = Get-Date
        $responseTime = ($endTime - $startTime).TotalMilliseconds
        $apiResponseTimes += $responseTime
        Write-Host "API请求 $i: $responseTime 毫秒"
        "API请求 $i: $responseTime 毫秒" | Out-File -FilePath $logFile -Append
    } catch {
        Write-Host "API请求 $i: 失败 - $($_.Exception.Message)"
        "API请求 $i: 失败 - $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
    }
}

# 计算平均值
$avgApiResponseTime = ($apiResponseTimes | Measure-Object -Average).Average
$minApiResponseTime = ($apiResponseTimes | Measure-Object -Minimum).Minimum
$maxApiResponseTime = ($apiResponseTimes | Measure-Object -Maximum).Maximum

Write-Host "平均API响应时间: $avgApiResponseTime 毫秒"
Write-Host "最小API响应时间: $minApiResponseTime 毫秒"
Write-Host "最大API响应时间: $maxApiResponseTime 毫秒"
"平均API响应时间: $avgApiResponseTime 毫秒" | Out-File -FilePath $logFile -Append
"最小API响应时间: $minApiResponseTime 毫秒" | Out-File -FilePath $logFile -Append
"最大API响应时间: $maxApiResponseTime 毫秒" | Out-File -FilePath $logFile -Append

$results["api_response"] = @{
    "average_time" = $avgApiResponseTime
    "min_time" = $minApiResponseTime
    "max_time" = $maxApiResponseTime
    "times" = $apiResponseTimes
}

# 停止API服务
.\elr.ps1 api stop

# 3. 资源使用测试
Write-Host "\n=== 测试3: 资源使用测试 ==="
"$(Get-Date) - 测试3: 资源使用测试" | Out-File -FilePath $logFile -Append

# 创建并启动测试容器
$containerID = "test-resource"
.\elr.ps1 container create --id $containerID --name "Resource Test Container" --image "ubuntu:latest"
.\elr.ps1 container start --id $containerID
Start-Sleep -Seconds 2

# 获取系统资源使用情况
$systemInfo = Get-CimInstance -ClassName Win32_OperatingSystem
$totalMemory = $systemInfo.TotalVisibleMemorySize / 1MB
$freeMemory = $systemInfo.FreePhysicalMemory / 1MB
$usedMemory = $totalMemory - $freeMemory
$memoryUsage = ($usedMemory / $totalMemory) * 100

$cpuUsage = (Get-Counter '\Processor(_Total)\% Processor Time').CounterSamples.CookedValue

Write-Host "内存使用: $usedMemory MB / $totalMemory MB ($memoryUsage%)"
Write-Host "CPU使用: $cpuUsage%"
"内存使用: $usedMemory MB / $totalMemory MB ($memoryUsage%)" | Out-File -FilePath $logFile -Append
"CPU使用: $cpuUsage%" | Out-File -FilePath $logFile -Append

$results["resource_usage"] = @{
    "memory_used" = $usedMemory
    "memory_total" = $totalMemory
    "memory_usage_percent" = $memoryUsage
    "cpu_usage_percent" = $cpuUsage
}

# 停止并删除容器
.\elr.ps1 container stop --id $containerID
.\elr.ps1 container delete --id $containerID

# 4. 并发测试
Write-Host "\n=== 测试4: 并发测试 ==="
"$(Get-Date) - 测试4: 并发测试" | Out-File -FilePath $logFile -Append

# 启动API服务
.\elr.ps1 api start
Start-Sleep -Seconds 2

$concurrency = 5
$concurrentResults = @()

# 并发请求
$jobs = @()
for ($i = 1; $i -le $concurrency; $i++) {
    $job = Start-Job -ScriptBlock {
        $startTime = Get-Date
        try {
            $response = Invoke-RestMethod -Uri "http://localhost:8080/api/container/list" -Method GET
            $endTime = Get-Date
            $responseTime = ($endTime - $startTime).TotalMilliseconds
            return @{"success" = $true; "time" = $responseTime}
        } catch {
            $endTime = Get-Date
            $responseTime = ($endTime - $startTime).TotalMilliseconds
            return @{"success" = $false; "time" = $responseTime; "error" = $_.Exception.Message}
        }
    }
    $jobs += $job
}

# 等待所有任务完成
foreach ($job in $jobs) {
    $result = Receive-Job -Job $job -Wait
    $concurrentResults += $result
}

# 计算并发测试结果
$successfulRequests = ($concurrentResults | Where-Object { $_.success -eq $true }).Count
$failedRequests = ($concurrentResults | Where-Object { $_.success -eq $false }).Count
$totalRequests = $concurrentResults.Count
$successRate = ($successfulRequests / $totalRequests) * 100

$concurrentResponseTimes = $concurrentResults | ForEach-Object { $_.time }
$avgConcurrentResponseTime = ($concurrentResponseTimes | Measure-Object -Average).Average

Write-Host "并发请求数: $totalRequests"
Write-Host "成功请求数: $successfulRequests"
Write-Host "失败请求数: $failedRequests"
Write-Host "成功率: $successRate%"
Write-Host "平均响应时间: $avgConcurrentResponseTime 毫秒"
"并发请求数: $totalRequests" | Out-File -FilePath $logFile -Append
"成功请求数: $successfulRequests" | Out-File -FilePath $logFile -Append
"失败请求数: $failedRequests" | Out-File -FilePath $logFile -Append
"成功率: $successRate%" | Out-File -FilePath $logFile -Append
"平均响应时间: $avgConcurrentResponseTime 毫秒" | Out-File -FilePath $logFile -Append

$results["concurrency"] = @{
    "total_requests" = $totalRequests
    "successful_requests" = $successfulRequests
    "failed_requests" = $failedRequests
    "success_rate" = $successRate
    "average_response_time" = $avgConcurrentResponseTime
}

# 停止API服务
.\elr.ps1 api stop

# 5. 网络隔离测试
Write-Host "\n=== 测试5: 网络隔离测试 ==="
"$(Get-Date) - 测试5: 网络隔离测试" | Out-File -FilePath $logFile -Append

# 创建测试容器
$containerID = "test-network-isolation"
.\elr.ps1 container create --id $containerID --name "Network Isolation Test Container" --image "ubuntu:latest"
.\elr.ps1 container start --id $containerID
Start-Sleep -Seconds 2

# 应用网络隔离
$isolationStart = Get-Date
$isolationResult = .\elr.ps1 network isolate --container-id $containerID
$isolationEnd = Get-Date
$isolationTime = ($isolationEnd - $isolationStart).TotalSeconds

Write-Host "网络隔离时间: $isolationTime 秒"
"网络隔离时间: $isolationTime 秒" | Out-File -FilePath $logFile -Append

# 获取网络配置
$configResult = .\elr.ps1 network config --container-id $containerID
Write-Host "网络配置: $configResult"
"网络配置: $configResult" | Out-File -FilePath $logFile -Append

# 移除网络隔离
$unisolationStart = Get-Date
$unisolationResult = .\elr.ps1 network unisolate --container-id $containerID
$unisolationEnd = Get-Date
$unisolationTime = ($unisolationEnd - $unisolationStart).TotalSeconds

Write-Host "移除网络隔离时间: $unisolationTime 秒"
"移除网络隔离时间: $unisolationTime 秒" | Out-File -FilePath $logFile -Append

$results["network_isolation"] = @{
    "isolation_time" = $isolationTime
    "unisolation_time" = $unisolationTime
}

# 清理容器
.\elr.ps1 container stop --id $containerID
.\elr.ps1 container delete --id $containerID

# 保存测试结果
$results | ConvertTo-Json -Depth 10 | Out-File -FilePath $resultsFile

# 测试完成
Write-Host "\n=== 测试完成 ==="
"$(Get-Date) - ELR性能测试完成" | Out-File -FilePath $logFile -Append
Write-Host "测试日志已保存到: $logFile"
Write-Host "测试结果已保存到: $resultsFile"

# 显示测试摘要
Write-Host "\n=== 测试摘要 ==="
Write-Host "1. 容器启动性能: 平均 $avgStartupTime 秒"
Write-Host "2. API响应性能: 平均 $avgApiResponseTime 毫秒"
Write-Host "3. 资源使用: 内存 $memoryUsage%, CPU $cpuUsage%"
Write-Host "4. 并发测试: 成功率 $successRate%, 平均响应时间 $avgConcurrentResponseTime 毫秒"
Write-Host "5. 网络隔离: 隔离时间 $isolationTime 秒, 移除时间 $unisolationTime 秒"
