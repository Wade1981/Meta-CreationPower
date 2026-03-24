# Test script to check container count issue

# Container statuses
$CONTAINER_STATUS_CREATED = "created"
$CONTAINER_STATUS_RUNNING = "running"
$CONTAINER_STATUS_STOPPED = "stopped"
$CONTAINER_STATUS_PAUSED = "paused"
$CONTAINER_STATUS_ERROR = "error"

# State file path
$STATE_FILE = "elr-state.json"

# Load state from file
function Load-State {
    if (Test-Path $STATE_FILE) {
        try {
            $state = Get-Content $STATE_FILE | ConvertFrom-Json
            $global:RUNTIME_STARTED = $state.RUNTIME_STARTED
            $global:RUNTIME_START_TIME = if ($state.RUNTIME_START_TIME) { [DateTime]::Parse($state.RUNTIME_START_TIME) } else { $null }
            $global:CONTAINERS = @()
            foreach ($container in $state.CONTAINERS) {
                $containerObj = @{
                    ID = $container.ID
                    Name = $container.Name
                    Image = $container.Image
                    Status = $container.Status
                    Created = [DateTime]::Parse($container.Created)
                }
                if ($container.Started) {
                    $containerObj.Started = [DateTime]::Parse($container.Started)
                }
                if ($container.Stopped) {
                    $containerObj.Stopped = [DateTime]::Parse($container.Stopped)
                }
                $global:CONTAINERS += $containerObj
            }
        } catch {
            Write-Host "Error loading state: $_"
        }
    } else {
        Write-Host "State file not found"
    }
}

# Check container status
function Check-Container-Status {
    Write-Host "===================================="
    Write-Host "Container Count Test"
    Write-Host "===================================="
    Write-Host "Total containers: $($global:CONTAINERS.Count)"
    Write-Host ""
    Write-Host "Container details:"
    foreach ($container in $global:CONTAINERS) {
        Write-Host "  ID: $($container.ID)"
        Write-Host "  Name: $($container.Name)"
        Write-Host "  Status: $($container.Status)"
        Write-Host "  Image: $($container.Image)"
        Write-Host "  Created: $($container.Created.ToString('yyyy-MM-dd HH:mm:ss'))"
        if ($container.Started) {
            Write-Host "  Started: $($container.Started.ToString('yyyy-MM-dd HH:mm:ss'))"
        }
        Write-Host ""
    }
    
    # Count running containers
    $runningCount = 0
    foreach ($container in $global:CONTAINERS) {
        Write-Host "Checking container $($container.Name) with status: $($container.Status)"
        if ($container.Status -eq $CONTAINER_STATUS_RUNNING) {
            Write-Host "  -> This container is running!"
            $runningCount++
        }
    }
    
    Write-Host ""
    Write-Host "Running containers count: $runningCount"
    Write-Host "===================================="
}

# Main script
Load-State
Check-Container-Status
