# Test script for model unloading functionality

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

# Test model unloading
function Test-ModelUnload {
    Write-Host "===================================="
    Write-Host "Testing Model Unloading"
    Write-Host "===================================="
    
    # Load current state
    $state = Load-SandboxState
    
    # Find sandbox with id sandbox-789012
    $sandbox = $state.sandboxes | Where-Object { $_.id -eq "sandbox-789012" } | Select-Object -First 1
    
    if (-not $sandbox) {
        Write-Host "Error: Sandbox sandbox-789012 not found"
        return
    }
    
    Write-Host "Found sandbox: $($sandbox.id)"
    Write-Host "Current models: $($sandbox.models.Count)"
    
    # Remove elr-cscc model
    $newModels = @()
    foreach ($model in $sandbox.models) {
        if ($model.id -ne "elr-cscc") {
            $newModels += $model
        } else {
            Write-Host "Removing model: $($model.id)"
        }
    }
    
    # Update sandbox
    $sandbox.models = $newModels
    
    # Save state
    $state.sandboxes = $state.sandboxes | ForEach-Object {
        if ($_.id -eq "sandbox-789012") {
            $sandbox
        } else {
            $_
        }
    }
    
    Save-SandboxState $state
    
    Write-Host "Model unloaded successfully!"
    Write-Host "Remaining models: $($newModels.Count)"
    
    # List remaining models
    if ($newModels.Count -gt 0) {
        Write-Host "===================================="
        Write-Host "Remaining Models:"
        Write-Host "===================================="
        foreach ($model in $newModels) {
            Write-Host "ID: $($model.id), Name: $($model.name)"
        }
    }
    
    Write-Host "===================================="
}

# Run test
Test-ModelUnload