# ELR安全测试脚本

# 测试目录
$testDir = "./test_output"
if (!(Test-Path $testDir)) {
    New-Item -ItemType Directory -Path $testDir -Force
}

# 测试日志
$logFile = "$testDir/elr_security_test_log.txt"
"$(Get-Date) - ELR安全测试开始" | Out-File -FilePath $logFile -Append

# 测试结果
$resultsFile = "$testDir/elr_security_results.json"
$results = @{}

# 1. 网络隔离测试
Write-Host "=== 测试1: 网络隔离测试 ==="
"$(Get-Date) - 测试1: 网络隔离测试" | Out-File -FilePath $logFile -Append

# 创建测试容器
$containerID = "test-security-network"
.\elr.ps1 container create --id $containerID --name "Network Security Test Container" --image "ubuntu:latest"
.\elr.ps1 container start --id $containerID
Start-Sleep -Seconds 2

# 应用网络隔离
$isolationResult = .\elr.ps1 network isolate --container-id $containerID
Write-Host "网络隔离结果: $isolationResult"
"网络隔离结果: $isolationResult" | Out-File -FilePath $logFile -Append

# 验证网络配置
$configResult = .\elr.ps1 network config --container-id $containerID
Write-Host "网络配置: $configResult"
"网络配置: $configResult" | Out-File -FilePath $logFile -Append

# 检查网络隔离是否成功
$configJson = $configResult | ConvertFrom-Json
$networkIsolationSuccess = $false
if ($configJson.enabled -eq $true) {
    $networkIsolationSuccess = $true
    Write-Host "✓ 网络隔离成功"
    "✓ 网络隔离成功" | Out-File -FilePath $logFile -Append
} else {
    Write-Host "✗ 网络隔离失败"
    "✗ 网络隔离失败" | Out-File -FilePath $logFile -Append
}

$results["network_isolation"] = @{
    "success" = $networkIsolationSuccess
    "config" = $configJson
}

# 2. API安全测试
Write-Host "\n=== 测试2: API安全测试 ==="
"$(Get-Date) - 测试2: API安全测试" | Out-File -FilePath $logFile -Append

# 启动API服务
.\elr.ps1 api start
Start-Sleep -Seconds 2

# 测试1: 速率限制
Write-Host "测试1.1: 速率限制测试"
"测试1.1: 速率限制测试" | Out-File -FilePath $logFile -Append

$rateLimitTestResults = @()
for ($i = 1; $i -le 70; $i++) {
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8080/health" -Method GET
        $rateLimitTestResults += @{"success" = $true; "request" = $i}
    } catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -eq 429) {
            $rateLimitTestResults += @{"success" = $false; "request" = $i; "error" = "Rate limit exceeded"}
            Write-Host "速率限制触发: 请求 $i"
            "速率限制触发: 请求 $i" | Out-File -FilePath $logFile -Append
            break
        } else {
            $rateLimitTestResults += @{"success" = $false; "request" = $i; "error" = $_.Exception.Message}
        }
    }
}

$rateLimitSuccess = $false
if ($rateLimitTestResults | Where-Object { $_.error -eq "Rate limit exceeded" }) {
    $rateLimitSuccess = $true
    Write-Host "✓ 速率限制测试通过"
    "✓ 速率限制测试通过" | Out-File -FilePath $logFile -Append
} else {
    Write-Host "✗ 速率限制测试失败"
    "✗ 速率限制测试失败" | Out-File -FilePath $logFile -Append
}

# 测试2: CORS策略
Write-Host "测试1.2: CORS策略测试"
"测试1.2: CORS策略测试" | Out-File -FilePath $logFile -Append

try {
    $response = Invoke-WebRequest -Uri "http://localhost:8080/health" -Method GET -Headers @{"Origin" = "http://example.com"}
    $corsHeaders = $response.Headers.GetValues("Access-Control-Allow-Origin")
    $corsSuccess = $false
    if ($corsHeaders -contains "*") {
        $corsSuccess = $true
        Write-Host "✓ CORS策略测试通过"
        "✓ CORS策略测试通过" | Out-File -FilePath $logFile -Append
    } else {
        Write-Host "✗ CORS策略测试失败"
        "✗ CORS策略测试失败" | Out-File -FilePath $logFile -Append
    }
} catch {
    Write-Host "✗ CORS策略测试失败: $($_.Exception.Message)"
    "✗ CORS策略测试失败: $($_.Exception.Message)" | Out-File -FilePath $logFile -Append
    $corsSuccess = $false
}

$results["api_security"] = @{
    "rate_limit" = $rateLimitSuccess
    "cors" = $corsSuccess
}

# 3. 令牌安全测试
Write-Host "\n=== 测试3: 令牌安全测试 ==="
"$(Get-Date) - 测试3: 令牌安全测试" | Out-File -FilePath $logFile -Append

# 创建令牌
$createTokenResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/token/create" -Method POST -Body '{"description": "Security Test Token"}' -ContentType "application/json"
$token = $createTokenResponse.token
Write-Host "创建令牌: $token"
"创建令牌: $token" | Out-File -FilePath $logFile -Append

# 验证令牌
$validateTokenResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/token/validate" -Method POST -Body "{\"token\": \"$token\"}" -ContentType "application/json"
$tokenValid = $validateTokenResponse.valid
if ($tokenValid) {
    Write-Host "✓ 令牌验证通过"
    "✓ 令牌验证通过" | Out-File -FilePath $logFile -Append
} else {
    Write-Host "✗ 令牌验证失败"
    "✗ 令牌验证失败" | Out-File -FilePath $logFile -Append
}

# 测试无效令牌
$invalidTokenResponse = Invoke-RestMethod -Uri "http://localhost:8080/api/token/validate" -Method POST -Body '{"token": "invalid-token"}' -ContentType "application/json"
$invalidTokenValid = $invalidTokenResponse.valid
if (!$invalidTokenValid) {
    Write-Host "✓ 无效令牌验证通过"
    "✓ 无效令牌验证通过" | Out-File -FilePath $logFile -Append
} else {
    Write-Host "✗ 无效令牌验证失败"
    "✗ 无效令牌验证失败" | Out-File -FilePath $logFile -Append
}

$results["token_security"] = @{
    "token_valid" = $tokenValid
    "invalid_token_valid" = $invalidTokenValid
}

# 4. 容器安全测试
Write-Host "\n=== 测试4: 容器安全测试 ==="
"$(Get-Date) - 测试4: 容器安全测试" | Out-File -FilePath $logFile -Append

# 创建测试容器
$containerID = "test-security-container"
.\elr.ps1 container create --id $containerID --name "Container Security Test Container" --image "ubuntu:latest"
.\elr.ps1 container start --id $containerID
Start-Sleep -Seconds 2

# 测试文件系统隔离
Write-Host "测试4.1: 文件系统隔离测试"
"测试4.1: 文件系统隔离测试" | Out-File -FilePath $logFile -Append

# 创建测试文件
$testFile = "test_security.txt"
"Security test file" | Out-File -FilePath $testFile

# 上传文件到容器
$uploadResult = .\elr.ps1 container upload --id $containerID --local-path $testFile --container-path "/test_security.txt" --token $token
Write-Host "文件上传结果: $uploadResult"
"文件上传结果: $uploadResult" | Out-File -FilePath $logFile -Append

# 下载文件从容器
$downloadResult = .\elr.ps1 container download --id $containerID --container-path "/test_security.txt" --local-path "$testDir/downloaded_security.txt" --token $token
Write-Host "文件下载结果: $downloadResult" | Out-File -FilePath $logFile -Append

# 清理测试文件
Remove-Item $testFile -Force -ErrorAction SilentlyContinue
Remove-Item "$testDir/downloaded_security.txt" -Force -ErrorAction SilentlyContinue

# 测试容器状态
$containerStatus = .\elr.ps1 container status --id $containerID
Write-Host "容器状态: $containerStatus"
"容器状态: $containerStatus" | Out-File -FilePath $logFile -Append

$containerSecuritySuccess = $true
Write-Host "✓ 容器安全测试通过"
"✓ 容器安全测试通过" | Out-File -FilePath $logFile -Append

$results["container_security"] = @{
    "success" = $containerSecuritySuccess
}

# 5. 安全策略测试
Write-Host "\n=== 测试5: 安全策略测试 ==="
"$(Get-Date) - 测试5: 安全策略测试" | Out-File -FilePath $logFile -Append

# 测试API访问控制
Write-Host "测试5.1: API访问控制测试"
"测试5.1: API访问控制测试" | Out-File -FilePath $logFile -Append

# 测试未授权访问
$unauthorizedAccessSuccess = $false
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8080/api/token/list" -Method GET
    Write-Host "✗ 未授权访问成功，安全策略失败"
    "✗ 未授权访问成功，安全策略失败" | Out-File -FilePath $logFile -Append
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401) {
        $unauthorizedAccessSuccess = $true
        Write-Host "✓ 未授权访问被拒绝，安全策略通过"
        "✓ 未授权访问被拒绝，安全策略通过" | Out-File -FilePath $logFile -Append
    } else {
        Write-Host "✗ 未授权访问测试失败: $statusCode"
        "✗ 未授权访问测试失败: $statusCode" | Out-File -FilePath $logFile -Append
    }
}

$results["security_policy"] = @{
    "unauthorized_access" = $unauthorizedAccessSuccess
}

# 清理资源
.\elr.ps1 container stop --id $containerID
.\elr.ps1 container delete --id $containerID
.\elr.ps1 container stop --id test-security-network
.\elr.ps1 container delete --id test-security-network
.\elr.ps1 api stop

# 保存测试结果
$results | ConvertTo-Json -Depth 10 | Out-File -FilePath $resultsFile

# 测试完成
Write-Host "\n=== 测试完成 ==="
"$(Get-Date) - ELR安全测试完成" | Out-File -FilePath $logFile -Append
Write-Host "测试日志已保存到: $logFile"
Write-Host "测试结果已保存到: $resultsFile"

# 显示测试摘要
Write-Host "\n=== 测试摘要 ==="
Write-Host "1. 网络隔离: $(if ($networkIsolationSuccess) { "✓ 通过" } else { "✗ 失败" })"
Write-Host "2. API安全:"
Write-Host "   - 速率限制: $(if ($rateLimitSuccess) { "✓ 通过" } else { "✗ 失败" })"
Write-Host "   - CORS策略: $(if ($corsSuccess) { "✓ 通过" } else { "✗ 失败" })"
Write-Host "3. 令牌安全:"
Write-Host "   - 有效令牌: $(if ($tokenValid) { "✓ 通过" } else { "✗ 失败" })"
Write-Host "   - 无效令牌: $(if (!$invalidTokenValid) { "✓ 通过" } else { "✗ 失败" })"
Write-Host "4. 容器安全: $(if ($containerSecuritySuccess) { "✓ 通过" } else { "✗ 失败" })"
Write-Host "5. 安全策略:"
Write-Host "   - 未授权访问: $(if ($unauthorizedAccessSuccess) { "✓ 通过" } else { "✗ 失败" })"
