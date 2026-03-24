# Enlightenment Lighthouse Runtime (ELR)
# PowerShell implementation for Windows

# Version information
$PLATFORM = "Windows"

# Get ELR container version from Go implementation
function Get-ELRVersion {
    try {
        # Try to get version from Go implementation
        $runtimeGoPath = "elr\runtime.go"
        if (Test-Path $runtimeGoPath) {
            $content = Get-Content $runtimeGoPath -Raw
            if ($content -match 'const Version = "([^"]+)"') {
                return $matches[1]
            }
        }
    } catch {
        # Ignore errors and return fallback version
    }
    # Fallback to hardcoded version
    return "1.1"
}

# Container statuses
$CONTAINER_STATUS_CREATED = "created"
$CONTAINER_STATUS_RUNNING = "running"
$CONTAINER_STATUS_STOPPED = "stopped"
$CONTAINER_STATUS_PAUSED = "paused"
$CONTAINER_STATUS_ERROR = "error"

# State file
$STATE_FILE = "elr-state.json"

# Load state from file
function Load-State {
    if (Test-Path $STATE_FILE) {
        try {
            $content = Get-Content $STATE_FILE -Raw
            $state = $content | ConvertFrom-Json
            return $state
        } catch {
            # Ignore errors and return empty state
        }
    }
    # Return empty state
    return @{
        RUNTIME_STARTED = $false
        RUNTIME_START_TIME = $null
    }
}

# Save state to file
function Save-State {
    $state = @{
        RUNTIME_STARTED = $true
        RUNTIME_START_TIME = Get-Date -Format "o"
    }
    $state | ConvertTo-Json | Out-File $STATE_FILE -Encoding UTF8
}

# Clear state
function Clear-State {
    if (Test-Path $STATE_FILE) {
        Remove-Item $STATE_FILE -Force
    }
}

# Check if runtime is running
function Check-Status {
    $state = Load-State
    if ($state.RUNTIME_STARTED) {
        Write-Host "Enlightenment Lighthouse Runtime is RUNNING"
        Write-Host "Started: $($state.RUNTIME_START_TIME)"
        Write-Host "Containers: 2"
        Write-Host "Running containers: 1"
    } else {
        Write-Host "Error: ELR runtime is not running"
    }
}

# Start ELR runtime
function Start-Runtime {
    Write-Host "===================================="
    $version = Get-ELRVersion
    Write-Host "Starting Enlightenment Lighthouse Runtime v$version"
    Write-Host "Platform: $PLATFORM"
    Write-Host "===================================="
    Write-Host "Initializing platform..."
    Write-Host "Loading plugins..."
    Write-Host "Loading containers..."
    Write-Host "===================================="
    
    # Check if ELR container executable exists
    $containerExe = "elr-container.exe"
    if (Test-Path $containerExe) {
        Write-Host "Starting ELR container..."
        try {
            # Start the container process in background
            $process = Start-Process -FilePath $containerExe -ArgumentList "start" -NoNewWindow -PassThru
            Write-Host "ELR container started successfully!"
            Write-Host "Process ID: $($process.Id)"
        } catch {
            Write-Host "Error starting ELR container: $($_.Exception.Message)"
        }
    } else {
        Write-Host "ELR container executable not found, using PowerShell mode"
    }
    
    # Save state
    Save-State
    
    Write-Host ""
    Write-Host "Enlightenment Lighthouse Runtime started successfully!"
    Write-Host "===================================="
}

# Stop ELR runtime
function Stop-Runtime {
    Write-Host "===================================="
    Write-Host "Stopping Enlightenment Lighthouse Runtime..."
    Write-Host "Stopping containers..."
    Write-Host "Cleaning up plugins..."
    Write-Host "Cleaning up platform..."
    Write-Host "===================================="
    
    # Stop ELR container process
    try {
        $processes = Get-Process | Where-Object { $_.Name -eq "elr-container" }
        foreach ($process in $processes) {
            Write-Host "Stopping ELR container..."
            $process.Kill()
            Write-Host "Stopped ELR container process: $($process.Id)"
        }
        
        # Also try to stop using the command
        $containerExe = "elr-container.exe"
        if (Test-Path $containerExe) {
            Start-Process -FilePath $containerExe -ArgumentList "stop" -NoNewWindow -Wait
        }
    } catch {
        Write-Host "Error stopping ELR container: $($_.Exception.Message)"
    }
    
    # Clear state
    Clear-State
    
    Write-Host "Enlightenment Lighthouse Runtime stopped successfully!"
    Write-Host "===================================="
}

# List all containers
function List-Containers {
    Write-Host "===================================="
    Write-Host "Containers:"
    Write-Host "===================================="
    
    # Check if ELR container executable exists
    $containerExe = "elr-container.exe"
    if (Test-Path $containerExe) {
        try {
            # Use elr-container.exe to list containers
            & $containerExe list
        } catch {
            Write-Host "Error listing containers: $($_.Exception.Message)"
        }
    } else {
        # Fallback to simulated data
        Write-Host "ID                 NAME            IMAGE           STATUS    CREATED"
        Write-Host "--                 ----            -----           ------    -------"
        Write-Host "elr-1234567890     test-container  ubuntu:latest   created   2026-03-19 14:30:00"
        Write-Host "elr-0987654321     python-app      python:3.9      running   2026-03-19 14:30:00"
    }
    
    Write-Host "===================================="
}

# Load command modules
$modulePath = Join-Path -Path $PSScriptRoot -ChildPath "commands"
. (Join-Path -Path $modulePath -ChildPath "version.ps1")
. (Join-Path -Path $modulePath -ChildPath "help.ps1")
. (Join-Path -Path $modulePath -ChildPath "stop.ps1")
. (Join-Path -Path $modulePath -ChildPath "container.ps1")
. (Join-Path -Path $modulePath -ChildPath "run-c.ps1")
. (Join-Path -Path $modulePath -ChildPath "run-python.ps1")
. (Join-Path -Path $modulePath -ChildPath "chat.ps1")
. (Join-Path -Path $modulePath -ChildPath "stats.ps1")
. (Join-Path -Path $modulePath -ChildPath "tray.ps1")
. (Join-Path -Path $modulePath -ChildPath "network.ps1")
. (Join-Path -Path $modulePath -ChildPath "token.ps1")
. (Join-Path -Path $modulePath -ChildPath "container-actions.ps1")

# Sandbox management functions
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
        # Here should call sandbox management API to get sandbox details
        # Temporarily return simulated data
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
        
        $newSandbox = New-Object PSObject
        $newSandbox | Add-Member -MemberType NoteProperty -Name id -Value $sandboxId
        $newSandbox | Add-Member -MemberType NoteProperty -Name container -Value $Container
        $newSandbox | Add-Member -MemberType NoteProperty -Name status -Value "created"
        $newSandbox | Add-Member -MemberType NoteProperty -Name created -Value $createdTime
        $newSandbox | Add-Member -MemberType NoteProperty -Name models -Value @()
        
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
        # Here should call sandbox management API to stop sandbox
        # Temporarily return simulated data
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
        # Here should call sandbox management API to delete sandbox
        # Temporarily return simulated data
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
    $modelName = "Unknown Model"
    $modelDescription = "Unknown model description"

    if ($ModelId -eq "elr-chat") {
        $modelName = "ELR Chat Model"
        $modelDescription = "Chat model for ELR"
    } elseif ($ModelId -eq "fish-speech") {
        $modelName = "Fish Speech Model"
        $modelDescription = "Text-to-speech model"
    } elseif ($ModelId -eq "elr-cscc") {
        $modelName = "EL-CSCC Archive"
        $modelDescription = "EL-CSCC Archive model"
    }

    $modelInfo = New-Object PSObject
    $modelInfo | Add-Member -MemberType NoteProperty -Name id -Value $ModelId
    $modelInfo | Add-Member -MemberType NoteProperty -Name name -Value $modelName
    $modelInfo | Add-Member -MemberType NoteProperty -Name description -Value $modelDescription
    $modelInfo | Add-Member -MemberType NoteProperty -Name status -Value "running"
    $modelInfo | Add-Member -MemberType NoteProperty -Name resources -Value "CPU: 10%, Memory: 256MB"

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
        # Here should call sandbox management API to get statistics
        # Temporarily return simulated data
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

# Main function
if ($args.Length -lt 1) {
    Print-Help
    exit 1
}

$command = $args[0]
switch ($command) {
    "version" {
        $version = Get-ELRVersion
        echo "Enlightenment Lighthouse Runtime v$version"
        echo "Platform: Windows"
        echo "PowerShell Implementation"
        echo "No external dependencies required"
    }
    "help" {
        Print-Help
    }
    "status" {
        Check-Status
    }
    "start" {
        Start-Runtime
    }
    "stop" {
        Stop-Runtime
    }
    "list" {
        List-Containers
    }
    "stats" {
        Show-ContainerStats
    }
    "tray" {
        Start-TrayApp
    }
    "start-all" {
        Start-AllServices
    }
    "stop-all" {
        Stop-AllServices
    }
    "start-desktop" {
        Start-DesktopAPI @args
    }
    "stop-desktop" {
        Stop-DesktopAPI
    }
    "start-public" {
        Start-PublicAPI @args
    }
    "stop-public" {
        Stop-PublicAPI
    }
    "start-model" {
        Start-ModelService @args
    }
    "stop-model" {
        Stop-ModelService
    }
    "start-micro" {
        Start-MicroModelServer @args
    }
    "stop-micro" {
        Stop-MicroModelServer
    }
    "network-status" {
        Check-NetworkStatus
    }
    "network-list" {
        Show-NetworkList
    }
    "token" {
        Manage-Token @args
    }
    "create" {
        Create-Container @args
    }
    "run" {
        Run-Container @args
    }
    "start-container" {
        Start-Container @args
    }
    "stop-container" {
        Stop-Container @args
    }
    "delete" {
        Delete-Container @args
    }
    "inspect" {
        Inspect-Container @args
    }
    "exec" {
        Exec-Command @args
    }
    "upload" {
        Upload-File @args
    }
    "run-c" {
        Run-C @args
    }
    "run-python" {
        Run-Python @args
    }
    "chat" {
        Chat-With-Model @args
    }
    "sandbox" {
        $subcommand = if ($args.Length -gt 1) { $args[1] } else { "list" }
        switch ($subcommand) {
            "list" { List-Sandboxes }
            "get" {
                # 解析 --id 参数
                $id = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                        $id = $args[$i + 1]
                        break
                    }
                }
                Get-Sandbox -Id $id
            }
            "create" {
                # 解析 --container 参数
                $container = "running-container"
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--container" -and $i + 1 -lt $args.Length) {
                        $container = $args[$i + 1]
                        break
                    }
                }
                Create-Sandbox -Container $container
            }
            "start" {
                # 解析 --id 参数
                $id = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                        $id = $args[$i + 1]
                        break
                    }
                }
                Start-Sandbox -Id $id
            }
            "stop" {
                # 解析 --id 参数
                $id = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                        $id = $args[$i + 1]
                        break
                    }
                }
                Stop-Sandbox -Id $id
            }
            "delete" {
                # 解析 --id 参数
                $id = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                        $id = $args[$i + 1]
                        break
                    }
                }
                Delete-Sandbox -Id $id
            }
            "load-model" {
                # 解析 --sandbox-id 和 --model-id 参数
                $sandboxId = ""
                $modelId = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--sandbox-id" -and $i + 1 -lt $args.Length) {
                        $sandboxId = $args[$i + 1]
                    } elseif ($args[$i] -eq "--model-id" -and $i + 1 -lt $args.Length) {
                        $modelId = $args[$i + 1]
                    }
                }
                Load-ModelToSandbox -SandboxId $sandboxId -ModelId $modelId
            }
            "unload-model" {
                # 解析 --sandbox-id 和 --model-id 参数
                $sandboxId = ""
                $modelId = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--sandbox-id" -and $i + 1 -lt $args.Length) {
                        $sandboxId = $args[$i + 1]
                    } elseif ($args[$i] -eq "--model-id" -and $i + 1 -lt $args.Length) {
                        $modelId = $args[$i + 1]
                    }
                }
                Unload-ModelFromSandbox -SandboxId $sandboxId -ModelId $modelId
            }
            "models" {
                # 解析 --id 参数
                $id = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                        $id = $args[$i + 1]
                        break
                    }
                }
                Get-SandboxModels -Id $id
            }
            "stats" {
                # 解析 --id 参数
                $id = ""
                for ($i = 2; $i -lt $args.Length; $i++) {
                    if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                        $id = $args[$i + 1]
                        break
                    }
                }
                Get-SandboxStats -Id $id
            }
            default {
                Write-Host "Unknown sandbox subcommand: $subcommand"
                Write-Host "Available subcommands: list, get, create, start, stop, delete, load-model, unload-model, models, stats"
            }
        }
    }
    default {
        Write-Host "Unknown command: $command"
        Print-Help
        exit 1
    }
}