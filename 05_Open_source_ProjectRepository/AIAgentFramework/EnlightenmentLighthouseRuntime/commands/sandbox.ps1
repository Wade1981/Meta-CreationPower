# Sandbox management command module

# 沙箱状态文件路径
$sandboxStateFile = "..\elr\sandbox-state.json"

# 确保状态文件存在
function Ensure-SandboxStateFile {
    if (-not (Test-Path $sandboxStateFile)) {
        $initialState = @{
            sandboxes = @()
        }
        $initialState | ConvertTo-Json | Out-File $sandboxStateFile -Encoding UTF8
    }
}

# 加载沙箱状态
function Load-SandboxState {
    Ensure-SandboxStateFile
    return Get-Content $sandboxStateFile | ConvertFrom-Json
}

# 保存沙箱状态
function Save-SandboxState {
    param(
        [object]$state
    )
    $state | ConvertTo-Json | Out-File $sandboxStateFile -Encoding UTF8
}

function List-Sandboxes {
    Write-Host "===================================="
    Write-Host "ELR Sandboxes"
    Write-Host "===================================="
    
    try {
        $state = Load-SandboxState
        
        if ($state.sandboxes.Count -eq 0) {
            Write-Host "No sandboxes found"
        } else {
            Write-Host "ID                STATUS    CONTAINER           CREATED                MODELS"
            Write-Host "--                ------    ---------           -------                ------"
            
            foreach ($sandbox in $state.sandboxes) {
                $modelCount = if ($sandbox.models) { $sandbox.models.Count } else { 0 }
                Write-Host "$($sandbox.id)    $($sandbox.status)   $($sandbox.container)   $($sandbox.created)    $modelCount"
            }
        }
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error listing sandboxes: $($_.Exception.Message)"
    }
}

function Get-Sandbox {
    param(
        [string]$Id
    )
    
    if ([string]::IsNullOrEmpty($Id)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Sandbox Details: $Id"
    Write-Host "===================================="
    
    try {
        # 这里应该调用沙箱管理API获取沙箱详情
        # 暂时返回模拟数据
        Write-Host "ID:                $Id"
        Write-Host "Status:            running"
        Write-Host "Container:         running-container"
        Write-Host "Created:           2026-03-23 10:00:00"
        Write-Host "Started:           2026-03-23 10:01:00"
        Write-Host "Uptime:            10m"
        Write-Host "Resources:"
        Write-Host "  CPU:            20.5%"
        Write-Host "  Memory:         512MB"
        Write-Host "  Disk:           1GB"
        Write-Host "Models:           2"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error getting sandbox: $($_.Exception.Message)"
    }
}

function Create-Sandbox {
    param(
        [string]$Container = "running-container"
    )
    
    Write-Host "===================================="
    Write-Host "Creating Sandbox"
    Write-Host "===================================="
    
    try {
        $state = Load-SandboxState
        
        $sandboxId = "sandbox-$(Get-Random -Minimum 100000 -Maximum 999999)"
        $createdTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
        
        $newSandbox = @{
            id = $sandboxId
            container = $Container
            status = "created"
            created = $createdTime
            models = @()
        }
        
        $state.sandboxes += $newSandbox
        Save-SandboxState $state
        
        Write-Host "Creating sandbox in container: $Container"
        Write-Host "Sandbox created successfully!"
        Write-Host "Sandbox ID: $sandboxId"
        Write-Host "Status: created"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error creating sandbox: $($_.Exception.Message)"
    }
}

function Start-Sandbox {
    param(
        [string]$Id
    )
    
    if ([string]::IsNullOrEmpty($Id)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Starting Sandbox: $Id"
    Write-Host "===================================="
    
    try {
        $state = Load-SandboxState
        
        $sandbox = $state.sandboxes | Where-Object { $_.id -eq $Id }
        if (-not $sandbox) {
            Write-Host "Error: Sandbox $Id not found"
            return
        }
        
        $sandbox.status = "running"
        Save-SandboxState $state
        
        Write-Host "Starting sandbox..."
        Start-Sleep -Milliseconds 500
        Write-Host "Sandbox started successfully!"
        Write-Host "Status: running"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error starting sandbox: $($_.Exception.Message)"
    }
}

function Stop-Sandbox {
    param(
        [string]$Id
    )
    
    if ([string]::IsNullOrEmpty($Id)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Stopping Sandbox: $Id"
    Write-Host "===================================="
    
    try {
        # 这里应该调用沙箱管理API停止沙箱
        # 暂时返回模拟数据
        Write-Host "Stopping sandbox..."
        Start-Sleep -Milliseconds 500
        Write-Host "Sandbox stopped successfully!"
        Write-Host "Status: stopped"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error stopping sandbox: $($_.Exception.Message)"
    }
}

function Delete-Sandbox {
    param(
        [string]$Id
    )
    
    if ([string]::IsNullOrEmpty($Id)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Deleting Sandbox: $Id"
    Write-Host "===================================="
    
    try {
        # 这里应该调用沙箱管理API删除沙箱
        # 暂时返回模拟数据
        Write-Host "Deleting sandbox..."
        Start-Sleep -Milliseconds 500
        Write-Host "Sandbox deleted successfully!"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error deleting sandbox: $($_.Exception.Message)"
    }
}

function Load-ModelToSandbox {
    param(
        [string]$SandboxId,
        [string]$ModelId
    )
    
    if ([string]::IsNullOrEmpty($SandboxId)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    if ([string]::IsNullOrEmpty($ModelId)) {
        Write-Host "Error: Model ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Loading Model to Sandbox"
    Write-Host "===================================="
    
    try {
        $state = Load-SandboxState
        
        $sandbox = $state.sandboxes | Where-Object { $_.id -eq $SandboxId }
        if (-not $sandbox) {
            Write-Host "Error: Sandbox $SandboxId not found"
            return
        }
        
        if ($sandbox.status -ne "running") {
            Write-Host "Error: Sandbox $SandboxId is not running"
            return
        }
        
        # 检查模型是否已加载
        $existingModel = $sandbox.models | Where-Object { $_.id -eq $ModelId }
        if ($existingModel) {
            Write-Host "Error: Model $ModelId is already loaded in sandbox $SandboxId"
            return
        }
        
        # 模拟模型信息
        $modelInfo = @{
            id = $ModelId
            name = if ($ModelId -eq "elr-chat") { "ELR Chat Model" } else { "Fish Speech Model" }
            description = if ($ModelId -eq "elr-chat") { "Chat model for ELR" } else { "Text-to-speech model" }
            status = "running"
            resources = "CPU: 10%, Memory: 256MB"
        }
        
        $sandbox.models += $modelInfo
        Save-SandboxState $state
        
        Write-Host "Loading model $ModelId to sandbox $SandboxId..."
        Start-Sleep -Milliseconds 500
        Write-Host "Model loaded successfully!"
        Write-Host "Model: $ModelId"
        Write-Host "Status: running"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error loading model: $($_.Exception.Message)"
    }
}

function Unload-ModelFromSandbox {
    param(
        [string]$SandboxId,
        [string]$ModelId
    )
    
    if ([string]::IsNullOrEmpty($SandboxId)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    if ([string]::IsNullOrEmpty($ModelId)) {
        Write-Host "Error: Model ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Unloading Model from Sandbox"
    Write-Host "===================================="
    
    try {
        $state = Load-SandboxState
        
        $sandbox = $state.sandboxes | Where-Object { $_.id -eq $SandboxId }
        if (-not $sandbox) {
            Write-Host "Error: Sandbox $SandboxId not found"
            return
        }
        
        # 检查模型是否已加载
        $existingModel = $sandbox.models | Where-Object { $_.id -eq $ModelId }
        if (-not $existingModel) {
            Write-Host "Error: Model $ModelId is not loaded in sandbox $SandboxId"
            return
        }
        
        # 从沙箱中移除模型
        $sandbox.models = $sandbox.models | Where-Object { $_.id -ne $ModelId }
        Save-SandboxState $state
        
        Write-Host "Unloading model $ModelId from sandbox $SandboxId..."
        Start-Sleep -Milliseconds 500
        Write-Host "Model unloaded successfully!"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error unloading model: $($_.Exception.Message)"
    }
}

function Get-SandboxModels {
    param(
        [string]$Id
    )
    
    if ([string]::IsNullOrEmpty($Id)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Models in Sandbox: $Id"
    Write-Host "===================================="
    
    try {
        $state = Load-SandboxState
        
        $sandbox = $state.sandboxes | Where-Object { $_.id -eq $Id }
        if (-not $sandbox) {
            Write-Host "Error: Sandbox $Id not found"
            return
        }
        
        if (-not $sandbox.models -or $sandbox.models.Count -eq 0) {
            Write-Host "No models loaded in sandbox $Id"
        } else {
            Write-Host "ID           NAME             DESCRIPTION                         STATUS    RESOURCES"
            Write-Host "--           ----             -----------                         ------    ---------"
            
            foreach ($model in $sandbox.models) {
                Write-Host "$($model.id)     $($model.name)   $($model.description)                         $($model.status)    $($model.resources)"
            }
            
            Write-Host "===================================="
            Write-Host "Total models: $($sandbox.models.Count)"
        }
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error getting sandbox models: $($_.Exception.Message)"
    }
}

function Get-SandboxStats {
    param(
        [string]$Id
    )
    
    if ([string]::IsNullOrEmpty($Id)) {
        Write-Host "Error: Sandbox ID is required"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Sandbox Statistics: $Id"
    Write-Host "===================================="
    
    try {
        # 这里应该调用沙箱管理API获取统计信息
        # 暂时返回模拟数据
        Write-Host "Status:            running"
        Write-Host "Uptime:            10m"
        Write-Host "Models Loaded:     2"
        Write-Host "Resource Usage:"
        Write-Host "  CPU:            25.5%"
        Write-Host "  Memory:         768MB"
        Write-Host "  Disk:           1.5GB"
        Write-Host "Model Status:"
        Write-Host "  elr-chat:      running"
        Write-Host "  fish-speech:   running"
        
        Write-Host "===================================="
    } catch {
        Write-Host "Error getting sandbox stats: $($_.Exception.Message)"
    }
}
