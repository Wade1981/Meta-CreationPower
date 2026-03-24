# Test sandbox management script

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
    $state = Load-SandboxState
    $sandboxId = "sandbox-789012"
    $createdTime = Get-Date -Format "yyyy-MM-dd HH:mm:ss"
    
    # Remove existing sandbox with the same ID
    $state.sandboxes = $state.sandboxes | Where-Object { $_.id -ne $sandboxId }
    
    # Create new sandbox
    $newSandbox = @{
        id = $sandboxId
        container = "running-container"
        status = "running"
        created = $createdTime
        models = @()
    }
    
    $state.sandboxes += $newSandbox
    Save-SandboxState $state
    
    Write-Host "Sandbox $sandboxId created and started successfully!"
}

# Load model to sandbox
function Load-Model {
    $state = Load-SandboxState
    $sandboxId = "sandbox-789012"
    $modelId = "elr-cscc"
    
    # Find sandbox
    $sandboxIndex = -1
    for ($i = 0; $i -lt $state.sandboxes.Count; $i++) {
        if ($state.sandboxes[$i].id -eq $sandboxId) {
            $sandboxIndex = $i
            break
        }
    }
    
    if ($sandboxIndex -eq -1) {
        Write-Host "Sandbox $sandboxId not found"
        return
    }
    
    # Get sandbox
    $sandbox = $state.sandboxes[$sandboxIndex]
    
    # Ensure models is an array
    if (-not $sandbox.models) {
        $sandbox.models = @()
    }
    
    # Check if model is already loaded
    $modelExists = $false
    foreach ($model in $sandbox.models) {
        if ($model.id -eq $modelId) {
            $modelExists = $true
            break
        }
    }
    
    if ($modelExists) {
        Write-Host "Model $modelId is already loaded"
        return
    }
    
    # Add model
    $modelInfo = @{
        id = $modelId
        name = "EL-CSCC Archive"
        description = "EL-CSCC Archive model"
        status = "running"
        resources = "CPU: 10%, Memory: 256MB"
    }
    
    $sandbox.models += $modelInfo
    $state.sandboxes[$sandboxIndex] = $sandbox
    Save-SandboxState $state
    
    Write-Host "Model $modelId loaded successfully!"
}

# Unload model from sandbox
function Unload-Model {
    $state = Load-SandboxState
    $sandboxId = "sandbox-789012"
    $modelId = "elr-cscc"
    
    # Find sandbox
    $sandboxIndex = -1
    for ($i = 0; $i -lt $state.sandboxes.Count; $i++) {
        if ($state.sandboxes[$i].id -eq $sandboxId) {
            $sandboxIndex = $i
            break
        }
    }
    
    if ($sandboxIndex -eq -1) {
        Write-Host "Sandbox $sandboxId not found"
        return
    }
    
    # Get sandbox
    $sandbox = $state.sandboxes[$sandboxIndex]
    
    # Ensure models is an array
    if (-not $sandbox.models) {
        $sandbox.models = @()
    }
    
    # Check if model is loaded
    $modelExists = $false
    foreach ($model in $sandbox.models) {
        if ($model.id -eq $modelId) {
            $modelExists = $true
            break
        }
    }
    
    if (-not $modelExists) {
        Write-Host "Model $modelId is not loaded"
        return
    }
    
    # Remove model
    $newModels = @()
    foreach ($model in $sandbox.models) {
        if ($model.id -ne $modelId) {
            $newModels += $model
        }
    }
    
    $sandbox.models = $newModels
    $state.sandboxes[$sandboxIndex] = $sandbox
    Save-SandboxState $state
    
    Write-Host "Model $modelId unloaded successfully!"
}

# List models in sandbox
function List-Models {
    $state = Load-SandboxState
    $sandboxId = "sandbox-789012"
    
    # Find sandbox
    $sandboxIndex = -1
    for ($i = 0; $i -lt $state.sandboxes.Count; $i++) {
        if ($state.sandboxes[$i].id -eq $sandboxId) {
            $sandboxIndex = $i
            break
        }
    }
    
    if ($sandboxIndex -eq -1) {
        Write-Host "Sandbox $sandboxId not found"
        return
    }
    
    # Get sandbox
    $sandbox = $state.sandboxes[$sandboxIndex]
    
    # Ensure models is an array
    if (-not $sandbox.models) {
        $sandbox.models = @()
    }
    
    Write-Host "Models in sandbox $sandboxId"
    if ($sandbox.models.Count -eq 0) {
        Write-Host "No models loaded"
    } else {
        foreach ($model in $sandbox.models) {
            Write-Host "- $($model.id): $($model.name)"
        }
    }
}

# Test the functions
Write-Host "Testing sandbox management..."
Write-Host "1. Creating sandbox..."
Create-Sandbox
Write-Host "2. Loading model..."
Load-Model
Write-Host "3. Listing models..."
List-Models
Write-Host "4. Unloading model..."
Unload-Model
Write-Host "5. Listing models again..."
List-Models
Write-Host "Test completed!"
