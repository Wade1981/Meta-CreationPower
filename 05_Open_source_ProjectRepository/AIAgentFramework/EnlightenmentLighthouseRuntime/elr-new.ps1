# Enlightenment Lighthouse Runtime (ELR)
# PowerShell implementation for Windows
# No external dependencies required

# Version information
$ELR_VERSION = "1.3.0"
$PLATFORM = "Windows"

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
            # If state file is corrupted, initialize with default values
            Initialize-DefaultState
        }
    } else {
        Initialize-DefaultState
    }
}

# Save state to file
function Save-State {
    # Convert containers to a format that can be properly serialized to JSON
    $serializableContainers = @()
    foreach ($container in $global:CONTAINERS) {
        $serializableContainer = @{
            ID = $container.ID
            Name = $container.Name
            Image = $container.Image
            Status = $container.Status
            Created = $container.Created.ToString('o')
        }
        if ($container.Started) {
            $serializableContainer.Started = $container.Started.ToString('o')
        }
        if ($container.Stopped) {
            $serializableContainer.Stopped = $container.Stopped.ToString('o')
        }
        $serializableContainers += $serializableContainer
    }
    
    $state = @{
        RUNTIME_STARTED = $global:RUNTIME_STARTED
        RUNTIME_START_TIME = if ($global:RUNTIME_START_TIME) { $global:RUNTIME_START_TIME.ToString('o') } else { $null }
        CONTAINERS = $serializableContainers
    }
    $state | ConvertTo-Json -Depth 3 | Set-Content $STATE_FILE
}

# Initialize default state
function Initialize-DefaultState {
    $global:RUNTIME_STARTED = $false
    $global:RUNTIME_START_TIME = $null
    $global:CONTAINERS = @(
        @{
            ID = "elr-1234567890"
            Name = "test-container"
            Image = "ubuntu:latest"
            Status = $CONTAINER_STATUS_CREATED
            Created = Get-Date
        },
        @{
            ID = "elr-0987654321"
            Name = "python-app"
            Image = "python:3.9"
            Status = $CONTAINER_STATUS_RUNNING
            Created = Get-Date
            Started = Get-Date
        }
    )
    Save-State
}

# Load initial state
Load-State

# Function: Print version information
function Print-Version {
    Write-Host "Enlightenment Lighthouse Runtime v$ELR_VERSION"
    Write-Host "Platform: $PLATFORM"
    Write-Host "PowerShell Implementation"
    Write-Host "No external dependencies required"
}

# Function: Print help information
function Print-Help {
    Write-Host "Enlightenment Lighthouse Runtime (ELR)"
    Write-Host "Usage: elr [command] [options]"
    Write-Host
    Write-Host "Commands:"
    Write-Host "  version           Print version information"
    Write-Host "  help              Print this help message"
    Write-Host "  start             Start the ELR runtime"
    Write-Host "  stop              Stop the ELR runtime"
    Write-Host "  status            Check the runtime status"
    Write-Host "  create            Create a new container"
    Write-Host "  run               Create and start a new container"
    Write-Host "  start-container   Start a container"
    Write-Host "  stop-container    Stop a container"
    Write-Host "  list              List all containers"
    Write-Host "  delete            Delete a container"
    Write-Host "  inspect           Inspect a container"
    Write-Host "  exec              Execute a command in a container"
    Write-Host "  run-c             Compile and run a C program"
    Write-Host "  run-python        Run a Python script or code"
    Write-Host "  chat              Interactive chat with model"
    Write-Host "  desktop-api       Start ELR Desktop API server"
    Write-Host "  token             Manage ELR container tokens"
    Write-Host "  network-status    Check network status"
    Write-Host
    Write-Host "Options:"
    Write-Host "  --name            Container name"
    Write-Host "  --image           Container image"
    Write-Host "  --id              Container ID"
    Write-Host "  --command         Command to execute"
    Write-Host "  --source          Source file for C program or Python script"
    Write-Host "  --output          Output file for compiled C program"
    Write-Host "  --args            Additional compile arguments"
    Write-Host "  --code            Python code to execute directly"
    Write-Host "  --model           Model file path"
    Write-Host "  --target          Target environment (local, container, sandbox)"
    Write-Host "  --action          Token action (create, validate, refresh, list, revoke)"
    Write-Host "  --token           Token value"
    Write-Host "  --description     Token description"
    Write-Host
    Write-Host "Examples:"
    Write-Host "  elr run-c --source hello.c"
    Write-Host "  elr run-c --source hello.c --output hello.exe"
    Write-Host "  elr run-c --source hello.c --args '-Wall -O2'"
    Write-Host "  elr exec --id elr-1234567890 --command 'ls -la'"
    Write-Host "  elr run-python --source script.py"
    Write-Host "  elr run-python --code 'print("Hello from Python!")'"
    Write-Host "  elr chat                           Start interactive chat with default local model"
    Write-Host "  elr chat --model path/to/model.py  Start chat with custom local model"
    Write-Host "  elr chat --target container --id elr-1234567890  Start chat with container"
    Write-Host "  elr chat --target sandbox          Start chat with sandbox model"
    Write-Host "  elr desktop-api                    Start ELR Desktop API server"
    Write-Host "  elr token --action create --description 'Admin Token'"
    Write-Host "  elr token --action validate --token 'token-value'"
    Write-Host "  elr token --action refresh --token 'old-token' --description 'Refreshed Token'"
    Write-Host "  elr token --action list"
    Write-Host "  elr network-status"
}

# Function: Check runtime status
function Check-Status {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    Write-Host "Enlightenment Lighthouse Runtime is RUNNING"
    Write-Host "Started: $($global:RUNTIME_START_TIME.ToString('yyyy-MM-dd HH:mm:ss'))"
    Write-Host "Containers: $($global:CONTAINERS.Count)"
    $runningContainers = $global:CONTAINERS | Where-Object { $_.Status -eq $CONTAINER_STATUS_RUNNING }
    Write-Host "Running containers: $($runningContainers.Count)"
}

# Function: Start ELR runtime
function Start-Runtime {
    if ($global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is already running"
        return
    }

    Write-Host "===================================="
    Write-Host "Starting Enlightenment Lighthouse Runtime v$ELR_VERSION"
    Write-Host "Platform: $PLATFORM"
    Write-Host "===================================="
    Write-Host "Initializing platform..."
    Start-Sleep -Milliseconds 500
    Write-Host "Loading plugins..."
    Start-Sleep -Milliseconds 500
    Write-Host "Loading containers..."
    Start-Sleep -Milliseconds 500
    Write-Host "===================================="

    # 启动新的ELR容器
    $elrContainerPath = "$PSScriptRoot\elr\elr-container.exe"
    if (Test-Path $elrContainerPath) {
        Write-Host "Starting ELR container..."
        Start-Process -FilePath $elrContainerPath -WorkingDirectory "$PSScriptRoot\elr"
        Start-Sleep -Seconds 2
        Write-Host "ELR container started successfully!"
    } else {
        Write-Host "Warning: ELR container executable not found at $elrContainerPath"
        Write-Host "Using PowerShell simulation mode instead"
    }

    # Set runtime status
    $global:RUNTIME_STARTED = $true
    $global:RUNTIME_START_TIME = Get-Date

    # Display container information
    foreach ($container in $global:CONTAINERS) {
        Write-Host "Created container: $($container.ID) ($($container.Name))"
    }

    Write-Host
    Write-Host "Enlightenment Lighthouse Runtime started successfully!"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: Stop ELR runtime
function Stop-Runtime {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    Write-Host "===================================="
    Write-Host "Stopping Enlightenment Lighthouse Runtime..."
    Write-Host "Stopping containers..."
    Start-Sleep -Milliseconds 500
    Write-Host "Cleaning up plugins..."
    Start-Sleep -Milliseconds 500
    Write-Host "Cleaning up platform..."
    Start-Sleep -Milliseconds 500
    Write-Host "===================================="

    # 停止ELR容器
    Write-Host "Stopping ELR container..."
    try {
        $elrProcesses = Get-Process | Where-Object {$_.Name -eq "elr-container"}
        foreach ($process in $elrProcesses) {
            Stop-Process -Id $process.Id -Force
            Write-Host "Stopped ELR container process: $($process.Id)"
        }
    } catch {
        Write-Host "Warning: Failed to stop ELR container: $($_.Exception.Message)"
    }

    # Set runtime status
    $global:RUNTIME_STARTED = $false
    $global:RUNTIME_START_TIME = $null

    Write-Host "Enlightenment Lighthouse Runtime stopped successfully!"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

# Function: List all containers
function List-Containers {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    Write-Host "===================================="
    Write-Host "Containers:"
    Write-Host "===================================="
    Write-Host "ID                 NAME            IMAGE           STATUS    CREATED"
    Write-Host "--                 ----            -----           ------    -------"

    foreach ($container in $global:CONTAINERS) {
        $id = $container.ID
        $name = $container.Name
        $image = $container.Image
        $status = $container.Status
        $created = $container.Created.ToString('yyyy-MM-dd HH:mm:ss')

        # Format output
        Write-Host "$($id.PadRight(17)) $($name.PadRight(14)) $($image.PadRight(15)) $($status.PadRight(8)) $created"
    }

    Write-Host "===================================="
}

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
    $container = $global:CONTAINERS | Where-Object { $_.ID -eq $containerID }
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
    $container = $global:CONTAINERS | Where-Object { $_.ID -eq $containerID }
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
    $global:CONTAINERS = $global:CONTAINERS | Where-Object { $_.ID -ne $containerID }

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
    $container = $global:CONTAINERS | Where-Object { $_.ID -eq $containerID }
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
    $container = $global:CONTAINERS | Where-Object { $_.ID -eq $containerID }
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

# Function: Run a C program
function Run-C-Program {
    param(
        [string]$source = "",
        [string]$output = "program.exe",
        [string]$args = ""
    )
    
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments from script args
    $sourceFile = $source
    $outputFile = $output
    $compileArgs = $args

    if ([string]::IsNullOrEmpty($sourceFile)) {
        Write-Host "Error: Source file is required"
        Write-Host "Usage: elr run-c --source file.c [--output output.exe] [--args compile_args]"
        return
    }

    # Check if source file exists
    if (-not (Test-Path $sourceFile)) {
        Write-Host "Error: Source file '$sourceFile' not found"
        return
    }

    Write-Host "===================================="
    Write-Host "Compiling C program..."
    Write-Host "Source: $sourceFile"
    Write-Host "Output: $outputFile"
    if ($compileArgs) {
        Write-Host "Compile args: $compileArgs"
    }
    Write-Host "===================================="

    # Check if gcc is available
    $gccPath = Get-Command gcc -ErrorAction SilentlyContinue
    if ($null -eq $gccPath) {
        Write-Host "Error: gcc compiler not found"
        Write-Host "Please install gcc or specify a different compiler"
        Write-Host ""
        Write-Host "For Windows, you can install gcc through:"
        Write-Host "  1. MinGW-w64: https://www.mingw-w64.org/"
        Write-Host "  2. MSYS2: https://www.msys2.org/"
        Write-Host "  3. Cygwin: https://www.cygwin.com/"
        return
    }

    # Compile the C program
    try {
        if ($compileArgs) {
            $compileCmd = "gcc $compileArgs $sourceFile -o $outputFile"
        } else {
            $compileCmd = "gcc $sourceFile -o $outputFile"
        }
        
        Write-Host "Executing: $compileCmd"
        $result = Invoke-Expression $compileCmd 2>&1
        
        if ($LASTEXITCODE -ne 0) {
            Write-Host "Error: Compilation failed"
            Write-Host $result
            return
        }
        
        Write-Host "===================================="
        Write-Host "Compilation successful!"
        Write-Host "===================================="
        
        # Run the compiled program
        if (Test-Path $outputFile) {
            Write-Host "Running program..."
            Write-Host "===================================="
            & .\$outputFile
            Write-Host "===================================="
            Write-Host "Program execution completed"
            Write-Host "===================================="
        } else {
            Write-Host "Error: Output file '$outputFile' not found"
        }
    } catch {
        Write-Host "Error: $_"
    }
}

# Function: Run a Python script or code
function Run-Python {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $sourceFile = ""
    $pythonCode = ""

    # Check if we have at least 2 arguments (command + option)
    if ($args.Length -lt 2) {
        Write-Host "Error: Not enough arguments"
        Write-Host "Usage: elr run-python --source script.py"
        Write-Host "       elr run-python --code 'python code'"
        return
    }

    # Parse arguments - start from index 1 since $args[0] is the command name
    for ($i = 1; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--source" -and $i + 1 -lt $args.Length) {
            $sourceFile = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--code" -and $i + 1 -lt $args.Length) {
            $pythonCode = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($sourceFile) -and [string]::IsNullOrEmpty($pythonCode)) {
        Write-Host "Error: Either --source or --code is required"
        Write-Host "Usage: elr run-python --source script.py"
        Write-Host "       elr run-python --code 'python code'"
        return
    }

    # Check if Python is available
    $pythonPath = Get-Command python -ErrorAction SilentlyContinue
    if ($null -eq $pythonPath) {
        $pythonPath = Get-Command python3 -ErrorAction SilentlyContinue
        if ($null -eq $pythonPath) {
            Write-Host "Error: Python interpreter not found"
            Write-Host "Please install Python 3.8 or higher"
            Write-Host ""
            Write-Host "You can download Python from:"
            Write-Host "  https://www.python.org/downloads/"
            Write-Host ""
            Write-Host "Or use Python portable version:"
            Write-Host "  https://www.python.org/downloads/windows/"
            Write-Host "  (Choose Windows embeddable package)"
            return
        }
    }

    # Check if it's a Windows Store placeholder
    if ($pythonPath.Source -like "*Microsoft\WindowsApps\python.exe") {
        Write-Host "Error: Found Windows Store Python placeholder, not actual Python interpreter"
        Write-Host "Please install Python from official website:"
        Write-Host "  https://www.python.org/downloads/"
        Write-Host ""
        Write-Host "Or use Python portable version:"
        Write-Host "  https://www.python.org/downloads/windows/"
        Write-Host "  (Choose Windows embeddable package)"
        return
    }

    Write-Host "===================================="
    Write-Host "Running Python..."
    
    if (-not [string]::IsNullOrEmpty($sourceFile)) {
        # Run Python script
        if (-not (Test-Path $sourceFile)) {
            Write-Host "Error: Source file '$sourceFile' not found"
            return
        }
        
        Write-Host "Script: $sourceFile"
        Write-Host "===================================="
        
        try {
            & $pythonPath $sourceFile
            $exitCode = $LASTEXITCODE
            
            if ($exitCode -ne 0) {
                Write-Host "Warning: Python execution returned non-zero exit code: $exitCode"
                Write-Host "This may indicate a problem with Python installation or script execution"
            }
            
            Write-Host "===================================="
            Write-Host "Python script execution completed"
            Write-Host "===================================="
        } catch {
            Write-Host "Error: $_"
        }
    } else {
        # Run Python code directly
        Write-Host "Code: $pythonCode"
        Write-Host "===================================="
        
        try {
            & $pythonPath -c $pythonCode
            $exitCode = $LASTEXITCODE
            
            if ($exitCode -ne 0) {
                Write-Host "Warning: Python execution returned non-zero exit code: $exitCode"
                Write-Host "This may indicate a problem with Python installation or code execution"
            }
            
            Write-Host "===================================="
            Write-Host "Python code execution completed"
            Write-Host "===================================="
        } catch {
            Write-Host "Error: $_"
        }
    }
}

# Function: Chat with model in container or sandbox
function Chat-With-Model {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $modelPath = "micro_model/examples/elr_chat_model.py"
    $containerID = ""
    $target = "local"  # local, container, sandbox

    for ($i = 1; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--model" -and $i + 1 -lt $args.Length) {
            $modelPath = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--target" -and $i + 1 -lt $args.Length) {
            $target = $args[$i + 1]
            $i++
        }
    }

    Write-Host "===================================="
    Write-Host "ELR Interactive Model Chat"
    Write-Host "===================================="
    Write-Host "Model: $modelPath"
    if (-not [string]::IsNullOrEmpty($containerID)) {
        Write-Host "Container: $containerID"
    }
    Write-Host "Target: $target"
    Write-Host "===================================="
    Write-Host "Welcome to ELR Interactive Model Chat!"
    Write-Host "You can chat with the model in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "===================================="

    # Check if Python is available
    $pythonAvailable = $false
    $pythonPath = $null
    $possiblePythonPaths = @(
        "python.exe",
        "python3.exe",
        "C:\Python39\python.exe",
        "C:\Python38\python.exe",
        "C:\Program Files\Python39\python.exe",
        "C:\Program Files\Python38\python.exe"
    )
    
    foreach ($path in $possiblePythonPaths) {
        try {
            $testPath = Get-Command $path -ErrorAction SilentlyContinue
            if ($testPath) {
                # Check if it's a Windows Store placeholder
                if (-not ($testPath.Source -like "*Microsoft\WindowsApps\python.exe" -or $testPath.Source -like "*Microsoft\WindowsApps\python3.exe")) {
                    $pythonPath = $testPath
                    $pythonAvailable = $true
                    break
                }
            }
        } catch {
            # Ignore errors
        }
    }
    
    if ($pythonAvailable) {
        Write-Host "Debug: Found Python at: $($pythonPath.Source)"
    } else {
        Write-Host "Warning: Python interpreter not found or only Windows Store placeholder available"
        Write-Host "Using PowerShell-based chat mode instead"
    }

    # Handle different targets
    switch ($target) {
        "container" {
            if ([string]::IsNullOrEmpty($containerID)) {
                Write-Host "Error: Container ID is required for container target"
                return
            }
            if ($pythonAvailable) {
                Chat-With-Container-Model -ContainerID $containerID -PythonPath $pythonPath
            } else {
                Chat-With-Container-Model-PowerShell -ContainerID $containerID
            }
        }
        "sandbox" {
            if ($pythonAvailable) {
                Chat-With-Sandbox-Model -ModelPath $modelPath -PythonPath $pythonPath
            } else {
                Chat-With-Sandbox-Model-PowerShell -ModelPath $modelPath
            }
        }
        default {
            # Local model
            if ($pythonAvailable) {
                Chat-With-Local-Model -ModelPath $modelPath -PythonPath $pythonPath
            } else {
                Chat-With-Local-Model-PowerShell
            }
        }
    }

    Write-Host "===================================="
    Write-Host "Chat session ended"
    Write-Host "===================================="
}

# Function: Chat with local model
function Chat-With-Local-Model {
    param(
        [string]$ModelPath,
        [string]$PythonPath
    )

    # Check if model file exists
    if (-not (Test-Path $ModelPath)) {
        # Try with micro_model path
        $fullModelPath = "micro_model/examples/elr_chat_model.py"
        if (Test-Path $fullModelPath) {
            $ModelPath = $fullModelPath
        } else {
            Write-Host "Error: Model file not found"
            return
        }
    }

    Write-Host "Starting interactive chat session with local model..."
    Write-Host "===================================="
    Write-Host "Using Python: $PythonPath"
    Write-Host "Model: $ModelPath"
    Write-Host "===================================="
    Write-Host "This feature is temporarily disabled due to PowerShell syntax issues."
    Write-Host "Please use the PowerShell-only chat mode instead."
    Write-Host "===================================="
}

# Function: Chat with model in container
function Chat-With-Container-Model {
    param(
        [string]$ContainerID,
        [string]$PythonPath
    )

    # Find container
    $container = $global:CONTAINERS | Where-Object { $_.ID -eq $ContainerID }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $ContainerID not found"
        return
    }

    if ($container.Status -ne $CONTAINER_STATUS_RUNNING) {
        Write-Host "Error: Container is not running"
        return
    }

    Write-Host "Starting interactive chat session with container..."
    Write-Host "===================================="
    Write-Host "Using Python: $PythonPath"
    Write-Host "Container: $ContainerID"
    Write-Host "===================================="
    Write-Host "This feature is temporarily disabled due to PowerShell syntax issues."
    Write-Host "Please use the PowerShell-only chat mode instead."
    Write-Host "===================================="
}

# Function: Chat with model in sandbox
function Chat-With-Sandbox-Model {
    param(
        [string]$ModelPath,
        [string]$PythonPath
    )

    Write-Host "Starting interactive chat session with sandbox model..."
    Write-Host "===================================="
    Write-Host "Using Python: $PythonPath"
    Write-Host "Model: $ModelPath"
    Write-Host "===================================="
    Write-Host "This feature is temporarily disabled due to PowerShell syntax issues."
    Write-Host "Please use the PowerShell-only chat mode instead."
    Write-Host "===================================="
}

# Function: Chat with local model using PowerShell only
function Chat-With-Local-Model-PowerShell {
    Write-Host "Starting interactive chat session with local model (PowerShell mode)..."
    Write-Host "===================================="
    
    # Define ELR chat model logic
    $modelName = "elr_chat_model"
    $modelVersion = "1.0"
    $modelDescription = "ELR container collaborative dialogue micro model, supporting multi-turn dialogue and ELR-related functions (PowerShell implementation)"
    
    Write-Host "Initializing model: $modelName v$modelVersion"
    Write-Host "Model loaded successfully!"
    Write-Host ""
    Write-Host "=== ELR Interactive Model Chat (PowerShell Mode) ==="
    Write-Host "Type your message in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "==================================="
    
    # Main chat loop
    while ($true) {
        try {
            $userInput = Read-Host "You"
            if ($userInput.ToLower() -in @('exit', 'quit', 'q')) {
                Write-Host "Model: Goodbye! Looking forward to chatting with you again."
                break
            } elseif ($userInput.ToLower() -eq 'help' -or $userInput.ToLower() -eq ',help') {
                Write-Host "Model: Available commands:"
                Write-Host "  ,help: Show available commands"
                Write-Host "  ,status: Check ELR container status"
                Write-Host "  ,models: List loaded models"
                Write-Host "  ,info: Show ELR container information"
                Write-Host "  ,clear: Clear conversation history"
                Write-Host "  ,exit: Exit conversation"
                Write-Host ""
                Write-Host "  Tip: You can directly ask questions about ELR container or use commands to get specific information."
            } elseif ($userInput.ToLower() -eq 'info' -or $userInput.ToLower() -eq ',info') {
                Write-Host "Model: ELR container information:"
                Write-Host "  - Version: v1.0"
                Write-Host "  - Run mode: Sandbox isolation"
                Write-Host "  - Dialogue model: $modelName v$modelVersion"
                Write-Host "  - Carbon-silicon synergy: Enabled"
            } elseif ($userInput.ToLower() -eq ',status') {
                Write-Host "Model: ELR container current status: Running"
                Write-Host "  - Model loading: Ready"
                Write-Host "  - Resource usage: Normal"
                Write-Host "  - Network connection: Available"
                Write-Host "  - Service status: Active"
            } elseif ($userInput.ToLower() -eq ',models') {
                Write-Host "Model: Loaded models:"
                Write-Host "  1. $modelName (Current dialogue model)"
                Write-Host "  2. Other models can be viewed through ELR container management interface"
            } elseif ($userInput.ToLower() -eq ',clear') {
                Write-Host "Model: Conversation history cleared"
            } else {
                # Process user input
                $lowerInput = $userInput.ToLower()
                $response = ""
                
                # Greetings
                $greetings = @('hello', 'hi', 'hey')
                foreach ($greeting in $greetings) {
                    if ($lowerInput -like "*$greeting*") {
                        $response = "Carbon-silicon synergy greeting! I'm $modelName, ELR container's dialogue assistant. Nice to serve you, how can I help you?"
                        break
                    }
                }
                
                # Questions
                if ([string]::IsNullOrEmpty($response)) {
                    $questions = @('how are you', 'how do you do', 'how are things', 'how is it going')
                    foreach ($question in $questions) {
                        if ($lowerInput -like "*$question*") {
                            $response = "Carbon-silicon synergy response! I'm $modelName, running status is good. ELR container is running normally, ready to serve you anytime."
                            break
                        }
                    }
                }
                
                # ELR related questions
                if ([string]::IsNullOrEmpty($response)) {
                    if ($lowerInput -like "*elr*" -and ($lowerInput -like "*what*" -or $lowerInput -like "*function*" -or $lowerInput -like "*capability*")) {
                        $response = "ELR container is the Enlightenment Lighthouse Runtime environment, main functions include: `n1. Model management and loading`n2. Sandbox isolated operation`n3. Resource monitoring and management`n4. Network communication and API services`n5. Container lifecycle management"
                    }
                }
                
                # Default response
                if ([string]::IsNullOrEmpty($response)) {
                    $response = "Carbon-silicon synergy response! I'm $modelName, received your message: '$userInput'. `n`nTip: You can input ',help' to see available commands, or ask questions about ELR container."
                }
                
                Write-Host "Model: $response"
            }
        } catch {
            Write-Host "Model: Error: $($_.Exception.Message)"
        }
    }
}

# Function: Chat with container model using PowerShell only
function Chat-With-Container-Model-PowerShell {
    param(
        [string]$ContainerID
    )
    
    # Find container
    $container = $global:CONTAINERS | Where-Object { $_.ID -eq $ContainerID }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $ContainerID not found"
        return
    }

    if ($container.Status -ne $CONTAINER_STATUS_RUNNING) {
        Write-Host "Error: Container is not running"
        return
    }
    
    Write-Host "Starting interactive chat session with container (PowerShell mode)..."
    Write-Host "===================================="
    Write-Host "Container: $ContainerID"
    Write-Host ""
    Write-Host "=== ELR Container Model Chat (PowerShell Mode) ==="
    Write-Host "Type your message in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "==================================="
    
    # Main chat loop
    while ($true) {
        try {
            $userInput = Read-Host "You"
            if ($userInput.ToLower() -in @('exit', 'quit', 'q')) {
                Write-Host "Goodbye!"
                break
            } elseif ($userInput.ToLower() -eq 'help') {
                Write-Host "Available commands:"
                Write-Host "  exit/quit/q - End the conversation"
                Write-Host "  help - Show this help"
                Write-Host "  info - Show container information"
            } elseif ($userInput.ToLower() -eq 'info') {
                Write-Host "Container Information:"
                Write-Host "  Container ID: $ContainerID"
                Write-Host "  Status: Running"
                Write-Host "  Type: ELR Container"
            } else {
                # Simulate container response
                $response = "[Container $ContainerID] Carbon-silicon synergy response! I've received your message: '$userInput'. Processing in container environment..."
                Write-Host "Container: $response"
            }
        } catch {
            Write-Host "Error: $($_.Exception.Message)"
        }
    }
}

# Function: Chat with sandbox model using PowerShell only
function Chat-With-Sandbox-Model-PowerShell {
    param(
        [string]$ModelPath
    )
    
    Write-Host "Starting interactive chat session with sandbox model (PowerShell mode)..."
    Write-Host "===================================="
    Write-Host "Model: $ModelPath"
    Write-Host ""
    
    # Generate sandbox ID
    $sandboxID = "sandbox-$(Get-Date -Format 'HHmmssfff')"
    
    Write-Host "=== ELR Sandbox Model Chat (PowerShell Mode) ==="
    Write-Host "Type your message in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "==================================="
    Write-Host "Sandbox initialized: $sandboxID"
    Write-Host "Model path: $ModelPath"
    
    # Main chat loop
    while ($true) {
        try {
            $userInput = Read-Host "You"
            if ($userInput.ToLower() -in @('exit', 'quit', 'q')) {
                Write-Host "Goodbye!"
                break
            } elseif ($userInput.ToLower() -eq 'help') {
                Write-Host "Available commands:"
                Write-Host "  exit/quit/q - End the conversation"
                Write-Host "  help - Show this help"
                Write-Host "  info - Show sandbox information"
            } elseif ($userInput.ToLower() -eq 'info') {
                Write-Host "Sandbox Information:"
                Write-Host "  Sandbox ID: $sandboxID"
                Write-Host "  Model Path: $ModelPath"
                Write-Host "  Status: Active"
                Write-Host "  Type: ELR Micro-Model Sandbox"
            } else {
                # Simulate sandbox model response
                $response = "[Sandbox $sandboxID] Carbon-silicon synergy response! I'm running in isolated sandbox environment. Your message: '$userInput' has been processed."
                Write-Host "Sandbox: $response"
            }
        } catch {
            Write-Host "Error: $($_.Exception.Message)"
        }
    }
    
    Write-Host "Sandbox closed: $sandboxID"
}

# Function: Start ELR Desktop API server
function Start-DesktopAPI {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    Write-Host "===================================="
    Write-Host "Starting ELR Desktop API Server"
    Write-Host "===================================="
    
    # 检查Desktop API启动脚本是否存在
    $desktopApiScript = "$PSScriptRoot\elr\start_desktop_api.ps1"
    if (-not (Test-Path $desktopApiScript)) {
        Write-Host "Error: Desktop API start script not found"
        Write-Host "Please ensure the script exists at: $desktopApiScript"
        return
    }
    
    Write-Host "Found Desktop API start script: $desktopApiScript"
    
    # 启动Desktop API服务器
    try {
        Write-Host "Starting Desktop API server..."
        Write-Host "This will open a new window for the Desktop API server"
        Write-Host "Press Ctrl+C to stop the server"
        Write-Host "===================================="
        
        # 在新窗口中启动Desktop API服务器
        Start-Process powershell.exe -ArgumentList "-ExecutionPolicy Bypass -File '$desktopApiScript'"
        
        Write-Host "Desktop API server started in a new window"
        Write-Host "Server address: http://localhost:8081"
        Write-Host "===================================="
    } catch {
        Write-Host "Error starting Desktop API server: $($_.Exception.Message)"
    }
}

# Function: Manage ELR container tokens
function Manage-Token {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # 解析参数
    $action = "help"
    $token = ""
    $description = "ELR Container Token"

    for ($i = 1; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--action" -and $i + 1 -lt $args.Length) {
            $action = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--token" -and $i + 1 -lt $args.Length) {
            $token = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--description" -and $i + 1 -lt $args.Length) {
            $description = $args[$i + 1]
            $i++
        }
    }

    # 检查令牌管理脚本是否存在
    $tokenManagerScript = "$PSScriptRoot\elr\token_manager.ps1"
    if (-not (Test-Path $tokenManagerScript)) {
        Write-Host "Error: Token manager script not found"
        Write-Host "Please ensure the script exists at: $tokenManagerScript"
        return
    }

    # 调用令牌管理脚本
    try {
        & $tokenManagerScript -Action $action -Token $token -Description $description
    } catch {
        Write-Host "Error managing token: $($_.Exception.Message)"
    }
}

# Function: Check network status
function Check-NetworkStatus {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # 检查网络状态脚本是否存在
    $tokenManagerScript = "$PSScriptRoot\elr\token_manager.ps1"
    if (-not (Test-Path $tokenManagerScript)) {
        Write-Host "Error: Token manager script not found"
        Write-Host "Please ensure the script exists at: $tokenManagerScript"
        return
    }

    # 调用网络状态查询功能
    try {
        & $tokenManagerScript -Action "network-status"
    } catch {
        Write-Host "Error checking network status: $($_.Exception.Message)"
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
        Print-Version
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
        Exec-Container @args
    }
    "run-c" {
        Run-C-Program @args
    }
    "run-python" {
        Run-Python @args
    }
    "chat" {
        Chat-With-Model @args
    }
    "desktop-api" {
        Start-DesktopAPI
    }
    "token" {
        Manage-Token @args
    }
    "network-status" {
        Check-NetworkStatus
    }
    default {
        Write-Host "Unknown command: $command"
        Print-Help
        exit 1
    }
}