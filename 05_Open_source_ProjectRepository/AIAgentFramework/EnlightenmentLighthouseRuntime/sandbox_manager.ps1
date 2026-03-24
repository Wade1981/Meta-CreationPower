# Sandbox Manager Script
# Simplified version for testing model loading and unloading

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
    
    # Create a simple object structure
    $saveState = @{
        sandboxes = @(
            @{
                id = "sandbox-789012"
                container = "running-container"
                status = "running"
                created = "2026-03-23 02:48:10"
                models = @()
            }
        )
    }
    
    # Add models
    foreach ($sandbox in $state.sandboxes) {
        if ($sandbox.id -eq "sandbox-789012") {
            foreach ($model in $sandbox.models) {
                # Create a simple hashtable for each model
                $modelHash = @{
                    id = $model.id
                    name = $model.name
                    description = $model.description
                    status = $model.status
                    resources = $model.resources
                }
                $saveState.sandboxes[0].models += $modelHash
            }
        }
    }
    
    # Convert to JSON with sufficient depth
    $saveState | ConvertTo-Json -Depth 5 | Out-File $sandboxStateFile -Encoding UTF8
}

# Load model to sandbox
function Load-Model {
    param(
        [string]$SandboxId,
        [string]$ModelId
    )

    Write-Host "===================================="
    Write-Host "Loading Model to Sandbox"
    Write-Host "===================================="

    $state = Load-SandboxState

    # Find sandbox
    $sandbox = $state.sandboxes | Where-Object { $_.id -eq $SandboxId } | Select-Object -First 1
    
    if (-not $sandbox) {
        Write-Host "Error: Sandbox $SandboxId not found"
        return
    }

    # Check if model is already loaded
    $existingModel = $sandbox.models | Where-Object { $_.id -eq $ModelId }
    if ($existingModel) {
        Write-Host "Error: Model $ModelId is already loaded in sandbox $SandboxId"
        return
    }

    # Model information
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

    # Create a hashtable for the model
    $modelInfo = @{
        id = $ModelId
        name = $modelName
        description = $modelDescription
        status = "running"
        resources = "CPU: 10%, Memory: 256MB"
    }

    # Add model to sandbox
    $newModels = @()
    if ($sandbox.models) {
        if ($sandbox.models -is [array]) {
            foreach ($item in $sandbox.models) {
                $newModels += $item
            }
        }
    }
    $newModels += $modelInfo

    # Update sandbox
    $sandbox.models = $newModels

    # Save state
    $state.sandboxes = $state.sandboxes | ForEach-Object {
        if ($_.id -eq $SandboxId) {
            $sandbox
        } else {
            $_
        }
    }

    Save-SandboxState $state

    Write-Host "Loading model $ModelId to sandbox $SandboxId..."
    Start-Sleep -Milliseconds 500
    Write-Host "Model loaded successfully!"
    Write-Host "===================================="
}

# Unload model from sandbox
function Unload-Model {
    param(
        [string]$SandboxId,
        [string]$ModelId
    )

    Write-Host "===================================="
    Write-Host "Unloading Model from Sandbox"
    Write-Host "===================================="

    $state = Load-SandboxState

    # Find sandbox
    $sandbox = $state.sandboxes | Where-Object { $_.id -eq $SandboxId } | Select-Object -First 1
    
    if (-not $sandbox) {
        Write-Host "Error: Sandbox $SandboxId not found"
        return
    }

    # Check if model is loaded
    $modelFound = $false
    foreach ($model in $sandbox.models) {
        if ($model.id -eq $ModelId) {
            $modelFound = $true
            break
        }
    }
    if (-not $modelFound) {
        Write-Host "Error: Model $ModelId is not loaded in sandbox $SandboxId"
        return
    }

    # Remove model
    $newModels = @()
    foreach ($model in $sandbox.models) {
        if ($model.id -ne $ModelId) {
            $newModels += $model
        }
    }

    # Update sandbox
    $sandbox.models = $newModels

    # Save state
    $state.sandboxes = $state.sandboxes | ForEach-Object {
        if ($_.id -eq $SandboxId) {
            $sandbox
        } else {
            $_
        }
    }

    Save-SandboxState $state

    Write-Host "Unloading model $ModelId from sandbox $SandboxId..."
    Start-Sleep -Milliseconds 500
    Write-Host "Model unloaded successfully!"
    Write-Host "===================================="
}

# List models in sandbox
function List-Models {
    param(
        [string]$SandboxId
    )

    Write-Host "===================================="
    Write-Host "Models in Sandbox: $SandboxId"
    Write-Host "===================================="

    $state = Load-SandboxState

    # Find sandbox
    $sandbox = $state.sandboxes | Where-Object { $_.id -eq $SandboxId } | Select-Object -First 1
    
    if (-not $sandbox) {
        Write-Host "Error: Sandbox $SandboxId not found"
        return
    }

    if (-not $sandbox.models -or $sandbox.models.Count -eq 0) {
        Write-Host "No models loaded in sandbox $SandboxId"
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

# Main function
if ($args.Length -lt 2) {
    Write-Host "Usage: sandbox_manager.ps1 <command> <sandbox-id> [model-id]"
    Write-Host "Commands: load-model, unload-model, models"
    exit 1
}

$command = $args[0]
$sandboxId = $args[1]

switch ($command) {
    "load-model" {
        if ($args.Length -lt 3) {
            Write-Host "Error: Model ID is required"
            exit 1
        }
        $modelId = $args[2]
        Load-Model -SandboxId $sandboxId -ModelId $modelId
    }
    "unload-model" {
        if ($args.Length -lt 3) {
            Write-Host "Error: Model ID is required"
            exit 1
        }
        $modelId = $args[2]
        Unload-Model -SandboxId $sandboxId -ModelId $modelId
    }
    "models" {
        List-Models -SandboxId $sandboxId
    }
    default {
        Write-Host "Unknown command: $command"
        Write-Host "Available commands: load-model, unload-model, models"
    }
}