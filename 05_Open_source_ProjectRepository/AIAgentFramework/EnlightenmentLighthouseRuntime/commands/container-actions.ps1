# Container actions command module

function Execute-ContainerCommand {
    param(
        [string]$command = ""
    )
    
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    
    if ([string]::IsNullOrEmpty($command)) {
        Write-Host "Error: No command specified"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Executing command in ELR container:"
    Write-Host "$command"
    Write-Host "===================================="
    
    try {
        # 直接在本地执行命令，模拟容器执行
        $output = Invoke-Expression $command 2>&1
        Write-Host $output
        Write-Host "===================================="
        Write-Host "Command executed successfully!"
    } catch {
        Write-Host "Error executing command: $($_.Exception.Message)"
    }
    Write-Host "===================================="
}

function Upload-FileToContainer {
    param(
        [string]$filePath = ""
    )
    
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    
    if ([string]::IsNullOrEmpty($filePath)) {
        Write-Host "Error: No file path specified"
        return
    }
    
    if (-not (Test-Path $filePath)) {
        Write-Host "Error: File not found: $filePath"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Uploading file to ELR container:"
    Write-Host "$filePath"
    Write-Host "===================================="
    
    try {
        # 模拟文件上传，实际项目中可以实现真实的文件传输
        $fileName = Split-Path $filePath -Leaf
        $destination = "$PSScriptRoot\..\elr\uploads\$fileName"
        
        # 创建上传目录
        if (-not (Test-Path "$PSScriptRoot\..\elr\uploads")) {
            New-Item -ItemType Directory -Path "$PSScriptRoot\..\elr\uploads" -Force | Out-Null
        }
        
        # 复制文件到上传目录
        Copy-Item -Path $filePath -Destination $destination -Force
        
        Write-Host "File uploaded successfully!"
        Write-Host "Destination: $destination"
    } catch {
        Write-Host "Error uploading file: $($_.Exception.Message)"
    }
    Write-Host "===================================="
}