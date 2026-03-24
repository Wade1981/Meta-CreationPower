# Container management commands module

# Function: Create a new container
function Create-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerName = ""
    $containerImage = "ubuntu:latest"

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--name" -and $i + 1 -lt $args.Length) {
            $containerName = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--image" -and $i + 1 -lt $args.Length) {
            $containerImage = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerName)) {
        $containerName = "container-$(Get-Date -Format 'HHmmss')"
    }

    $containerID = "elr-$(Get-Date -Format 'HHmmssfff')"
    $containerStatus = $CONTAINER_STATUS_CREATED
    $containerCreated = Get-Date

    # Add container to global list
    $newContainer = @{
        ID = $containerID
        Name = $containerName
        Image = $containerImage
        Status = $containerStatus
        Created = $containerCreated
    }

    $global:CONTAINERS += $newContainer

    Write-Host "===================================="
    Write-Host "Created container: $containerID ($containerName)"
    Write-Host "Image: $containerImage"
    Write-Host "Status: $containerStatus"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: Run a new container
function Run-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerName = ""
    $containerImage = "ubuntu:latest"

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--name" -and $i + 1 -lt $args.Length) {
            $containerName = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--image" -and $i + 1 -lt $args.Length) {
            $containerImage = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerName)) {
        $containerName = "container-$(Get-Date -Format 'HHmmss')"
    }

    $containerID = "elr-$(Get-Date -Format 'HHmmssfff')"
    $containerStatus = $CONTAINER_STATUS_RUNNING
    $containerCreated = Get-Date
    $containerStarted = Get-Date

    # Add container to global list
    $newContainer = @{
        ID = $containerID
        Name = $containerName
        Image = $containerImage
        Status = $containerStatus
        Created = $containerCreated
        Started = $containerStarted
    }

    $global:CONTAINERS += $newContainer

    Write-Host "===================================="
    Write-Host "Running container: $containerID ($containerName)"
    Write-Host "Image: $containerImage"
    Write-Host "Status: $containerStatus"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: Start a container
function Start-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerID = ""

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerID)) {
        Write-Host "Error: Container ID is required"
        return
    }

    # Find container
    $container = $null
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -eq $containerID) {
            $container = $c
            break
        }
    }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $containerID not found"
        return
    }

    # Update container status
    $container.Status = $CONTAINER_STATUS_RUNNING
    $container.Started = Get-Date

    Write-Host "===================================="
    Write-Host "Started container: $containerID"
    Write-Host "Status: $CONTAINER_STATUS_RUNNING"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: Stop a container
function Stop-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerID = ""

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerID)) {
        Write-Host "Error: Container ID is required"
        return
    }

    # Find container
    $container = $null
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -eq $containerID) {
            $container = $c
            break
        }
    }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $containerID not found"
        return
    }

    # Update container status
    $container.Status = $CONTAINER_STATUS_STOPPED
    $container.Stopped = Get-Date

    Write-Host "===================================="
    Write-Host "Stopped container: $containerID"
    Write-Host "Status: $CONTAINER_STATUS_STOPPED"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: Delete a container
function Delete-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerID = ""

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerID)) {
        Write-Host "Error: Container ID is required"
        return
    }

    # Remove container from global list
    $newContainers = @()
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -ne $containerID) {
            $newContainers += $c
        }
    }
    $global:CONTAINERS = $newContainers

    Write-Host "===================================="
    Write-Host "Deleted container: $containerID"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: Inspect a container
function Inspect-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerID = ""

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerID)) {
        Write-Host "Error: Container ID is required"
        return
    }

    # Find container
    $container = $null
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -eq $containerID) {
            $container = $c
            break
        }
    }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $containerID not found"
        return
    }

    Write-Host "===================================="
    Write-Host "Container Details:"
    Write-Host "===================================="
    Write-Host "ID: $($container.ID)"
    Write-Host "Name: $($container.Name)"
    Write-Host "Image: $($container.Image)"
    Write-Host "Status: $($container.Status)"
    Write-Host "Created: $($container.Created.ToString('yyyy-MM-dd HH:mm:ss'))"
    if ($container.Started) {
        Write-Host "Started: $($container.Started.ToString('yyyy-MM-dd HH:mm:ss'))"
    }
    if ($container.Stopped) {
        Write-Host "Stopped: $($container.Stopped.ToString('yyyy-MM-dd HH:mm:ss'))"
    }
    Write-Host "===================================="
}

# Function: Execute a command in a container
function Exec-Container {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerID = ""
    $command = ""

    for ($i = 1; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--command" -and $i + 1 -lt $args.Length) {
            $command = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($containerID)) {
        Write-Host "Error: Container ID is required"
        return
    }

    if ([string]::IsNullOrEmpty($command)) {
        Write-Host "Error: Command is required"
        return
    }

    # Find container
    $container = $null
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -eq $containerID) {
            $container = $c
            break
        }
    }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $containerID not found"
        return
    }

    if ($container.Status -ne $CONTAINER_STATUS_RUNNING) {
        Write-Host "Error: Container is not running"
        return
    }

    Write-Host "===================================="
    Write-Host "Executing command in container: $containerID"
    Write-Host "Command: $command"
    Write-Host "===================================="
    
    # Execute the command
    try {
        Invoke-Expression $command
    } catch {
        Write-Host "Error executing command: $_"
    }
    
    Write-Host "===================================="
}

