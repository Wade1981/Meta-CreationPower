#!/usr/bin/env powershell

# ELR Container Management Script
# This script provides a PowerShell interface to elr-container.exe

# Path to elr-container.exe
$ELR_CONTAINER_EXE = ".\elr-container.exe"

# Check if elr-container.exe exists
function Test-ELRContainerExe {
    return Test-Path $ELR_CONTAINER_EXE
}

# Start ELR container in background
function Start-ELRContainer {
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    Write-Host "Starting ELR container..."
    try {
        $process = Start-Process -FilePath $ELR_CONTAINER_EXE -ArgumentList "start" -NoNewWindow -PassThru
        Write-Host "ELR container started successfully!"
        Write-Host "Process ID: $($process.Id)"
    } catch {
        Write-Host "Error starting ELR container: $($_.Exception.Message)"
    }
}

# Stop ELR container
function Stop-ELRContainer {
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    Write-Host "Stopping ELR container..."
    try {
        # First try to stop using the command
        Start-Process -FilePath $ELR_CONTAINER_EXE -ArgumentList "stop" -NoNewWindow -Wait
        
        # Then kill any remaining processes
        $processes = Get-Process | Where-Object { $_.Name -eq "elr-container" }
        foreach ($process in $processes) {
            Write-Host "Stopping ELR container process: $($process.Id)"
            $process.Kill()
        }
        
        Write-Host "ELR container stopped successfully!"
    } catch {
        Write-Host "Error stopping ELR container: $($_.Exception.Message)"
    }
}

# List containers
function List-ELRContainers {
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    Write-Host "Listing containers..."
    try {
        & $ELR_CONTAINER_EXE list
    } catch {
        Write-Host "Error listing containers: $($_.Exception.Message)"
    }
}

# Create container
function New-ELRContainer {
    param(
        [string]$Name,
        [string]$Image
    )
    
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    if (-not $Name) {
        Write-Host "Error: Container name is required"
        return
    }
    
    if (-not $Image) {
        Write-Host "Error: Container image is required"
        return
    }
    
    Write-Host "Creating container $Name with image $Image..."
    try {
        & $ELR_CONTAINER_EXE create --name $Name --image $Image
    } catch {
        Write-Host "Error creating container: $($_.Exception.Message)"
    }
}

# Run container (create and start)
function Run-ELRContainer {
    param(
        [string]$Name,
        [string]$Image
    )
    
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    if (-not $Name) {
        Write-Host "Error: Container name is required"
        return
    }
    
    if (-not $Image) {
        Write-Host "Error: Container image is required"
        return
    }
    
    Write-Host "Running container $Name with image $Image..."
    try {
        & $ELR_CONTAINER_EXE run --name $Name --image $Image
    } catch {
        Write-Host "Error running container: $($_.Exception.Message)"
    }
}

# Start container
function Start-ELRContainerInstance {
    param(
        [string]$Id
    )
    
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    if (-not $Id) {
        Write-Host "Error: Container ID is required"
        return
    }
    
    Write-Host "Starting container $Id..."
    try {
        & $ELR_CONTAINER_EXE start-container --id $Id
    } catch {
        Write-Host "Error starting container: $($_.Exception.Message)"
    }
}

# Stop container
function Stop-ELRContainerInstance {
    param(
        [string]$Id
    )
    
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    if (-not $Id) {
        Write-Host "Error: Container ID is required"
        return
    }
    
    Write-Host "Stopping container $Id..."
    try {
        & $ELR_CONTAINER_EXE stop-container --id $Id
    } catch {
        Write-Host "Error stopping container: $($_.Exception.Message)"
    }
}

# Delete container
function Remove-ELRContainer {
    param(
        [string]$Id
    )
    
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    if (-not $Id) {
        Write-Host "Error: Container ID is required"
        return
    }
    
    Write-Host "Deleting container $Id..."
    try {
        & $ELR_CONTAINER_EXE delete --id $Id
    } catch {
        Write-Host "Error deleting container: $($_.Exception.Message)"
    }
}

# Inspect container
function Get-ELRContainerInfo {
    param(
        [string]$Id
    )
    
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    if (-not $Id) {
        Write-Host "Error: Container ID is required"
        return
    }
    
    Write-Host "Inspecting container $Id..."
    try {
        & $ELR_CONTAINER_EXE inspect --id $Id
    } catch {
        Write-Host "Error inspecting container: $($_.Exception.Message)"
    }
}

# Show version
function Get-ELRVersion {
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    try {
        & $ELR_CONTAINER_EXE version
    } catch {
        Write-Host "Error getting version: $($_.Exception.Message)"
    }
}

# Show help
function Show-ELRHelp {
    if (-not (Test-ELRContainerExe)) {
        Write-Host "Error: elr-container.exe not found"
        return
    }
    
    try {
        & $ELR_CONTAINER_EXE help
    } catch {
        Write-Host "Error showing help: $($_.Exception.Message)"
    }
}

# Main function
if ($args.Length -lt 1) {
    Write-Host "Usage: elr-container.ps1 [command] [options]"
    Write-Host ""
    Write-Host "Commands:"
    Write-Host "  start              Start ELR container"
    Write-Host "  stop               Stop ELR container"
    Write-Host "  list               List containers"
    Write-Host "  create             Create a new container"
    Write-Host "  run                Create and start a new container"
    Write-Host "  start-container    Start a container"
    Write-Host "  stop-container     Stop a container"
    Write-Host "  delete             Delete a container"
    Write-Host "  inspect            Inspect a container"
    Write-Host "  version            Show version"
    Write-Host "  help               Show help"
    exit 1
}

$command = $args[0]
switch ($command) {
    "start" {
        Start-ELRContainer
    }
    "stop" {
        Stop-ELRContainer
    }
    "list" {
        List-ELRContainers
    }
    "create" {
        # Parse --name and --image parameters
        $name = ""
        $image = ""
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--name" -and $i + 1 -lt $args.Length) {
                $name = $args[$i + 1]
            } elseif ($args[$i] -eq "--image" -and $i + 1 -lt $args.Length) {
                $image = $args[$i + 1]
            }
        }
        New-ELRContainer -Name $name -Image $image
    }
    "run" {
        # Parse --name and --image parameters
        $name = ""
        $image = ""
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--name" -and $i + 1 -lt $args.Length) {
                $name = $args[$i + 1]
            } elseif ($args[$i] -eq "--image" -and $i + 1 -lt $args.Length) {
                $image = $args[$i + 1]
            }
        }
        Run-ELRContainer -Name $name -Image $image
    }
    "start-container" {
        # Parse --id parameter
        $id = ""
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                $id = $args[$i + 1]
                break
            }
        }
        Start-ELRContainerInstance -Id $id
    }
    "stop-container" {
        # Parse --id parameter
        $id = ""
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                $id = $args[$i + 1]
                break
            }
        }
        Stop-ELRContainerInstance -Id $id
    }
    "delete" {
        # Parse --id parameter
        $id = ""
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                $id = $args[$i + 1]
                break
            }
        }
        Remove-ELRContainer -Id $id
    }
    "inspect" {
        # Parse --id parameter
        $id = ""
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
                $id = $args[$i + 1]
                break
            }
        }
        Get-ELRContainerInfo -Id $id
    }
    "version" {
        Get-ELRVersion
    }
    "help" {
        Show-ELRHelp
    }
    default {
        Write-Host "Unknown command: $command"
        Show-ELRHelp
    }
}
