# Enlightenment Lighthouse Runtime (ELR)
# Sandbox management functions

# Sandbox state file
$sandboxStateFile = Join-Path -Path $PSScriptRoot -ChildPath "elr\sandbox-state.json"

# Ensure sandbox state file exists
function Ensure-SandboxStateFile {
    if (-not (Test-Path $sandboxStateFile)) {
        $initialState = @{
            sandboxes = @()
        }
        $initialState | ConvertTo-Json | Out-File $sandboxStateFile -Encoding UTF8
    }
}

# Load sandbox state
function Load-SandboxState {
    Ensure-SandboxStateFile
    return Get-Content $sandboxStateFile | ConvertFrom-Json
}

# Save sandbox state
function Save-SandboxState {
    param(
        [object]$state
    )
    $state | ConvertTo-Json | Out-File $sandboxStateFile -Encoding UTF8
}

# Create sandbox
function Create-Sandbox {
    param(
        [string]$Container = "running-container"
    )
    
    Write-Host "===================================="
    Write-Host "Creating Sandbox"
    Write-Host "===================================="
    
    $state = Load-SandboxState
    
    $sandboxId = "sandbox-$(Get-Random -Minimum 100000 -Maximum 999999)"
    $createdTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    
    $newSandbox = New-Object PSObject -Property @{
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
}

# Start sandbox
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
    
    $state = Load-SandboxState
    
    # 找到沙箱索引
    $sandboxIndex = -1
    for ($i = 0; $i -lt $state.sandboxes.Count; $i++) {
        if ($state.sandboxes[$i].id -eq $Id) {
            $sandboxIndex = $i
            break
        }
    }
    
    if ($sandboxIndex -eq -1) {
        Write-Host "Error: Sandbox $Id not found"
        return
    }
    
    $sandbox = $state.sandboxes[$sandboxIndex]
    $sandbox.status = "running"
    $state.sandboxes[$sandboxIndex] = $sandbox
    Save-SandboxState $state
    
    Write-Host "Starting sandbox..."
    Start-Sleep -Milliseconds 500
    Write-Host "Sandbox started successfully!"
    Write-Host "Status: running"
    
    Write-Host "===================================="
}

# Load model to sandbox
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
    
    $state = Load-SandboxState
    
    # 找到沙箱索引
    $sandboxIndex = -1
    for ($i = 0; $i -lt $state.sandboxes.Count; $i++) {
        if ($state.sandboxes[$i].id -eq $SandboxId) {
            $sandboxIndex = $i
            break
        }
    }
    
    if ($sandboxIndex -eq -1) {
        Write-Host "Error: Sandbox $SandboxId not found"
        return
    }
    
    $sandbox = $state.sandboxes[$sandboxIndex]
    
    if ($sandbox.status -ne "running") {
        Write-Host "Error: Sandbox $SandboxId is not running"
        return
    }
    
    # Check if model is already loaded
    $existingModel = $sandbox.models | Where-Object { $_.id -eq $ModelId }
    if ($existingModel) {
        Write-Host "Error: Model $ModelId is already loaded in sandbox $SandboxId"
        return
    }
    
    # Simulated model information
    $modelInfo = New-Object PSObject -Property @{
        id = $ModelId
        name = if ($ModelId -eq "elr-chat") { "ELR Chat Model" } elseif ($ModelId -eq "fish-speech") { "Fish Speech Model" } elseif ($ModelId -eq "elr-cscc") { "EL-CSCC Archive" } else { "Unknown Model" }
        description = if ($ModelId -eq "elr-chat") { "Chat model for ELR" } elseif ($ModelId -eq "fish-speech") { "Text-to-speech model" } elseif ($ModelId -eq "elr-cscc") { "EL-CSCC Archive model" } else { "Unknown model description" }
        status = "running"
        resources = "CPU: 10%, Memory: 256MB"
    }
    
    # 创建新的模型数组并添加模型
    $newModels = @()
    if ($sandbox.models) {
        # 确保 $sandbox.models 是一个数组
        if ($sandbox.models -is [array]) {
            $newModels = $sandbox.models
        } else {
            # 如果不是数组，创建一个新数组
            $newModels = @()
        }
    }
    $newModels += $modelInfo
    
    # 更新沙箱对象
    $sandbox.models = $newModels
    $state.sandboxes[$sandboxIndex] = $sandbox
    
    Save-SandboxState $state
    
    Write-Host "Loading model $ModelId to sandbox $SandboxId..."
    Start-Sleep -Milliseconds 500
    Write-Host "Model loaded successfully!"
    Write-Host "Model: $ModelId"
    Write-Host "Status: running"
    
    Write-Host "===================================="
}

# Unload model from sandbox
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
    
    $state = Load-SandboxState
    
    # 找到沙箱索引
    $sandboxIndex = -1
    for ($i = 0; $i -lt $state.sandboxes.Count; $i++) {
        if ($state.sandboxes[$i].id -eq $SandboxId) {
            $sandboxIndex = $i
            break
        }
    }
    
    if ($sandboxIndex -eq -1) {
        Write-Host "Error: Sandbox $SandboxId not found"
        return
    }
    
    $sandbox = $state.sandboxes[$sandboxIndex]
    
    # Check if model is loaded
    $existingModel = $sandbox.models | Where-Object { $_.id -eq $ModelId }
    if (-not $existingModel) {
        Write-Host "Error: Model $ModelId is not loaded in sandbox $SandboxId"
        return
    }
    
    # Remove model from sandbox
    $newModels = $sandbox.models | Where-Object { $_.id -ne $ModelId }
    
    # 更新沙箱对象
    $sandbox.models = $newModels
    $state.sandboxes[$sandboxIndex] = $sandbox
    
    Save-SandboxState $state
    
    Write-Host "Unloading model $ModelId from sandbox $SandboxId..."
    Start-Sleep -Milliseconds 500
    Write-Host "Model unloaded successfully!"
    
    Write-Host "===================================="
}

# Get sandbox models
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
}

# List all sandboxes
function List-Sandboxes {
    Write-Host "===================================="
    Write-Host "ELR Sandboxes"
    Write-Host "===================================="
    
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
}

# Main function
if ($args.Length -lt 1) {
    Write-Host "Usage: elr_sandbox.ps1 <command> [options]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  create            Create sandbox"
    Write-Host "  start             Start sandbox"
    Write-Host "  load-model        Load model to sandbox"
    Write-Host "  unload-model      Unload model from sandbox"
    Write-Host "  models            List models in sandbox"
    Write-Host "  list              List all sandboxes"
    exit 1
}

$command = $args[0]
switch ($command) {
    "create" {
        Create-Sandbox
    }
    "start" {
        if ($args.Length -lt 2) {
            Write-Host "Error: Sandbox ID is required"
            exit 1
        }
        Start-Sandbox -Id $args[1]
    }
    "load-model" {
        if ($args.Length -lt 3) {
            Write-Host "Error: Sandbox ID and Model ID are required"
            exit 1
        }
        Load-ModelToSandbox -SandboxId $args[1] -ModelId $args[2]
    }
    "unload-model" {
        if ($args.Length -lt 3) {
            Write-Host "Error: Sandbox ID and Model ID are required"
            exit 1
        }
        Unload-ModelFromSandbox -SandboxId $args[1] -ModelId $args[2]
    }
    "models" {
        if ($args.Length -lt 2) {
            Write-Host "Error: Sandbox ID is required"
            exit 1
        }
        Get-SandboxModels -Id $args[1]
    }
    "list" {
        List-Sandboxes
    }
    default {
        Write-Host "Unknown command: $command"
        Write-Host "Use 'elr_sandbox.ps1' for usage information"
        exit 1
    }
}
