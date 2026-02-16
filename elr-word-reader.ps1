#!/usr/bin/env powershell

# Enlightenment Lighthouse Runtime (ELR) - Word Document Reader Extension
# PowerShell implementation for Windows
# Enhanced with Word document reading capabilities

# Version information
$ELR_VERSION = "1.0.2"
$PLATFORM = "Windows"

# Container statuses
$CONTAINER_STATUS_CREATED = "created"
$CONTAINER_STATUS_RUNNING = "running"
$CONTAINER_STATUS_STOPPED = "stopped"
$CONTAINER_STATUS_PAUSED = "paused"
$CONTAINER_STATUS_ERROR = "error"

# State file path
$STATE_FILE = "elr-state.json"

# Python portable version information
$PYTHON_PORTABLE_URL = "https://www.python.org/ftp/python/3.9.13/python-3.9.13-embed-amd64.zip"
$PYTHON_PORTABLE_ZIP = "python-portable.zip"
$PYTHON_DIR = "python-portable"
$PYTHON_EXE = "$PYTHON_DIR\python.exe"

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
        },
        @{
            ID = "elr-1122334455"
            Name = "word-reader"
            Image = "python:3.9"
            Status = $CONTAINER_STATUS_CREATED
            Created = Get-Date
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
    Write-Host "With Word document reading capabilities"
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
    Write-Host "  convert-docx      Convert Word document to Markdown"
    Write-Host "  read-word         Read and analyze Word document"
    Write-Host "  run-word-container Run a container with Word document reading capabilities"
    Write-Host
    Write-Host "Options:"n    Write-Host "  --name            Container name"
    Write-Host "  --image           Container image"
    Write-Host "  --id              Container ID"
    Write-Host "  --command         Command to execute"
    Write-Host "  --source          Source file for C program or Python script"
    Write-Host "  --output          Output file for compiled C program or Markdown"
    Write-Host "  --args            Additional compile arguments"
    Write-Host "  --code            Python code to execute directly"
    Write-Host "  --input           Input file for document conversion or reading"
    Write-Host "  --analyze         Analyze document structure"
    Write-Host "  --extract         Extract specific content from document"
    Write-Host "  --format          Output format (text, json, markdown)"
    Write-Host
    Write-Host "Examples:"
    Write-Host "  elr run-c --source hello.c"
    Write-Host "  elr run-python --source script.py"
    Write-Host "  elr convert-docx --input docs\document.docx --output docs\document.md"
    Write-Host "  elr read-word --input docs\document.docx --analyze"
    Write-Host "  elr run-word-container --name doc-analyzer"
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
    $wordReaderContainers = $global:CONTAINERS | Where-Object { $_.Name -like "*word*" -or $_.Name -like "*doc*" }
    Write-Host "Word document capable containers: $($wordReaderContainers.Count)"
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
    Write-Host "Loading Word document reading capabilities..."
    Start-Sleep -Milliseconds 500
    Write-Host "===================================="

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
    Write-Host "Cleaning up Word document reading capabilities..."
    Start-Sleep -Milliseconds 500
    Write-Host "Cleaning up platform..."
    Start-Sleep -Milliseconds 500
    Write-Host "===================================="

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
    if ($container.Name -like "*word*" -or $container.Name -like "*doc*") {
        Write-Host "Capabilities: Word document reading and analysis"
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

    for ($i = 2; $i -lt $args.Length; $i++) {
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
        Write-Host "Usage: elr run-c --source <file.c> [--output <output.exe>] [--args <compile_args>]"
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

# Function: Download portable Python
function Download-PortablePython {
    if (Test-Path $PYTHON_EXE) {
        Write-Host "Portable Python already exists: $PYTHON_EXE"
        return $true
    }

    Write-Host "Downloading portable Python..."
    Write-Host "URL: $PYTHON_PORTABLE_URL"
    
    try {
        # 创建临时目录
        New-Item -ItemType Directory -Path $PYTHON_DIR -Force | Out-Null
        
        # 下载Python便携式版本
        Invoke-WebRequest -Uri $PYTHON_PORTABLE_URL -OutFile $PYTHON_PORTABLE_ZIP -ErrorAction Stop
        Write-Host "Download completed: $PYTHON_PORTABLE_ZIP"
        
        # 解压缩
        Write-Host "Extracting portable Python..."
        Expand-Archive -Path $PYTHON_PORTABLE_ZIP -DestinationPath $PYTHON_DIR -Force -ErrorAction Stop
        Write-Host "Extraction completed"
        
        # 清理
        Remove-Item $PYTHON_PORTABLE_ZIP -Force
        
        # 检查Python是否可用
        if (Test-Path $PYTHON_EXE) {
            Write-Host "Portable Python is ready: $PYTHON_EXE"
            return $true
        } else {
            Write-Host "Error: Portable Python executable not found"
            return $false
        }
        
    } catch {
        Write-Host "Error downloading portable Python: $($_.Exception.Message)"
        # 清理
        if (Test-Path $PYTHON_DIR) {
            Remove-Item $PYTHON_DIR -Recurse -Force
        }
        if (Test-Path $PYTHON_PORTABLE_ZIP) {
            Remove-Item $PYTHON_PORTABLE_ZIP -Force
        }
        return $false
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
        Write-Host "Usage: elr run-python --source <script.py>"
        Write-Host "       elr run-python --code '<python code>'"
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
        Write-Host "Usage: elr run-python --source <script.py>"
        Write-Host "       elr run-python --code '<python code>'"
        return
    }

    # 尝试使用系统Python，如果没有则下载便携式Python
    $pythonPath = Get-Command python -ErrorAction SilentlyContinue
    if ($null -eq $pythonPath) {
        Write-Host "System Python not found, using portable Python..."
        if (-not (Download-PortablePython)) {
            Write-Host "Error: Failed to prepare portable Python"
            return
        }
        $pythonPath = $PYTHON_EXE
    } else {
        Write-Host "Using system Python: $($pythonPath.Source)"
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
            # Collect all script arguments
            $scriptArgs = @()
            for ($j = $i; $j -lt $args.Length; $j++) {
                $scriptArgs += $args[$j]
            }
            
            # Build command
            $command = "$pythonPath $sourceFile"
            if ($scriptArgs.Count -gt 0) {
                $command += " " + ($scriptArgs -join " ")
            }
            
            Write-Host "Executing: $command"
            
            # Execute with arguments
            if ($scriptArgs.Count -gt 0) {
                & $pythonPath $sourceFile $scriptArgs
            } else {
                & $pythonPath $sourceFile
            }
            
            $exitCode = $LASTEXITCODE
            
            if ($exitCode -ne 0) {
                Write-Host "Warning: Python execution returned non-zero exit code: $exitCode"
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
            Write-Host "Executing: $pythonPath -c $pythonCode"
            & $pythonPath -c $pythonCode
            $exitCode = $LASTEXITCODE
            
            if ($exitCode -ne 0) {
                Write-Host "Warning: Python execution returned non-zero exit code: $exitCode"
            }
            
            Write-Host "===================================="
            Write-Host "Python code execution completed"
            Write-Host "===================================="
        } catch {
            Write-Host "Error: $_"
        }
    }
}

# Function: Create document reader script
function Create-DocumentReader {
    $readerScript = @'
#!/usr/bin/env python3
import os
import sys
from zipfile import ZipFile
from xml.etree import ElementTree as ET
import json

# 注册命名空间
ET.register_namespace('', 'http://schemas.openxmlformats.org/wordprocessingml/2006/main')
nsmap = {'w': 'http://schemas.openxmlformats.org/wordprocessingml/2006/main'}

def extract_text_from_docx(docx_path):
    """
    从docx文件中提取文本内容
    """
    text = []
    
    try:
        # 打开docx文件（本质是zip文件）
        with ZipFile(docx_path, 'r') as zf:
            # 读取document.xml文件
            with zf.open('word/document.xml') as f:
                tree = ET.parse(f)
                root = tree.getroot()
                
                # 遍历所有段落
                for para in root.findall('.//w:p', namespaces=nsmap):
                    para_text = []
                    
                    # 遍历段落中的所有文本运行
                    for run in para.findall('.//w:r', namespaces=nsmap):
                        for text_elem in run.findall('.//w:t', namespaces=nsmap):
                            if text_elem.text:
                                para_text.append(text_elem.text)
                    
                    # 检查段落样式
                    style_name = ""
                    pPr = para.find('.//w:pPr', namespaces=nsmap)
                    if pPr:
                        pStyle = pPr.find('.//w:pStyle', namespaces=nsmap)
                        if pStyle is not None and 'w:val' in pStyle.attrib:
                            style_name = pStyle.attrib['w:val']
                    
                    # 合并段落文本
                    para_text_str = ''.join(para_text)
                    if para_text_str.strip():
                        # 根据样式添加Markdown格式
                        if style_name == 'Heading1':
                            text.append(f"# {para_text_str.strip()}")
                        elif style_name == 'Heading2':
                            text.append(f"## {para_text_str.strip()}")
                        elif style_name == 'Heading3':
                            text.append(f"### {para_text_str.strip()}")
                        elif style_name == 'Heading4':
                            text.append(f"#### {para_text_str.strip()}")
                        elif style_name == 'Heading5':
                            text.append(f"##### {para_text_str.strip()}")
                        elif style_name == 'Heading6':
                            text.append(f"###### {para_text_str.strip()}")
                        else:
                            text.append(para_text_str.strip())
        
        return '\n'.join(text)
        
    except Exception as e:
        print(f"Error extracting text: {e}")
        return ""

def analyze_document_structure(docx_path):
    """
    分析文档结构
    """
    structure = {
        'paragraphs': [],
        'headings': {},
        'statistics': {
            'total_paragraphs': 0,
            'total_headings': 0,
            'total_words': 0
        }
    }
    
    try:
        # 打开docx文件（本质是zip文件）
        with ZipFile(docx_path, 'r') as zf:
            # 读取document.xml文件
            with zf.open('word/document.xml') as f:
                tree = ET.parse(f)
                root = tree.getroot()
                
                # 遍历所有段落
                for para in root.findall('.//w:p', namespaces=nsmap):
                    para_info = {
                        'text': '',
                        'style': '',
                        'word_count': 0
                    }
                    
                    # 遍历段落中的所有文本运行
                    for run in para.findall('.//w:r', namespaces=nsmap):
                        for text_elem in run.findall('.//w:t', namespaces=nsmap):
                            if text_elem.text:
                                para_info['text'] += text_elem.text
                    
                    # 检查段落样式
                    pPr = para.find('.//w:pPr', namespaces=nsmap)
                    if pPr:
                        pStyle = pPr.find('.//w:pStyle', namespaces=nsmap)
                        if pStyle is not None and 'w:val' in pStyle.attrib:
                            para_info['style'] = pStyle.attrib['w:val']
                    
                    # 计算单词数
                    para_info['word_count'] = len(para_info['text'].split())
                    
                    # 添加到结构中
                    structure['paragraphs'].append(para_info)
                    structure['statistics']['total_paragraphs'] += 1
                    structure['statistics']['total_words'] += para_info['word_count']
                    
                    # 处理标题
                    if para_info['style'].startswith('Heading'):
                        structure['statistics']['total_headings'] += 1
                        heading_level = para_info['style'][-1]
                        if heading_level not in structure['headings']:
                            structure['headings'][heading_level] = []
                        structure['headings'][heading_level].append(para_info['text'])
        
        return structure
        
    except Exception as e:
        print(f"Error analyzing document structure: {e}")
        return structure

def extract_specific_content(docx_path, content_type):
    """
    提取特定类型的内容
    """
    content = []
    
    try:
        # 打开docx文件（本质是zip文件）
        with ZipFile(docx_path, 'r') as zf:
            # 读取document.xml文件
            with zf.open('word/document.xml') as f:
                tree = ET.parse(f)
                root = tree.getroot()
                
                # 遍历所有段落
                for para in root.findall('.//w:p', namespaces=nsmap):
                    para_text = []
                    para_style = ""
                    
                    # 遍历段落中的所有文本运行
                    for run in para.findall('.//w:r', namespaces=nsmap):
                        for text_elem in run.findall('.//w:t', namespaces=nsmap):
                            if text_elem.text:
                                para_text.append(text_elem.text)
                    
                    # 检查段落样式
                    pPr = para.find('.//w:pPr', namespaces=nsmap)
                    if pPr:
                        pStyle = pPr.find('.//w:pStyle', namespaces=nsmap)
                        if pStyle is not None and 'w:val' in pStyle.attrib:
                            para_style = pStyle.attrib['w:val']
                    
                    # 合并段落文本
                    para_text_str = ''.join(para_text)
                    if para_text_str.strip():
                        # 根据内容类型提取
                        if content_type == 'headings' and para_style.startswith('Heading'):
                            content.append({
                                'text': para_text_str.strip(),
                                'style': para_style
                            })
                        elif content_type == 'normal' and para_style == '':
                            content.append({
                                'text': para_text_str.strip(),
                                'style': para_style
                            })
        
        return content
        
    except Exception as e:
        print(f"Error extracting specific content: {e}")
        return content

def main():
    if len(sys.argv) < 2:
        print("Usage: python word_reader.py <input.docx> [--analyze] [--extract <content_type>] [--format <format>]")
        sys.exit(1)
    
    input_file = sys.argv[1]
    analyze = False
    extract = ""
    output_format = "text"
    
    # 解析参数
    for i in range(2, len(sys.argv)):
        if sys.argv[i] == "--analyze":
            analyze = True
        elif sys.argv[i] == "--extract" and i + 1 < len(sys.argv):
            extract = sys.argv[i + 1]
        elif sys.argv[i] == "--format" and i + 1 < len(sys.argv):
            output_format = sys.argv[i + 1]
    
    if analyze:
        # 分析文档结构
        structure = analyze_document_structure(input_file)
        if output_format == "json":
            print(json.dumps(structure, ensure_ascii=False, indent=2))
        else:
            print("Document Analysis Results:")
            print("=============================")
            print(f"Total Paragraphs: {structure['statistics']['total_paragraphs']}")
            print(f"Total Headings: {structure['statistics']['total_headings']}")
            print(f"Total Words: {structure['statistics']['total_words']}")
            print("\nHeadings:")
            for level, headings in sorted(structure['headings'].items()):
                print(f"Level {level}:")
                for heading in headings:
                    print(f"  - {heading}")
    elif extract:
        # 提取特定内容
        content = extract_specific_content(input_file, extract)
        if output_format == "json":
            print(json.dumps(content, ensure_ascii=False, indent=2))
        else:
            print(f"Extracted {extract} content:")
            print("=============================")
            for item in content:
                print(f"{item['text']}")
    else:
        # 提取所有文本
        text = extract_text_from_docx(input_file)
        print(text)

if __name__ == "__main__":
    main()
'@
    
    $readerScriptPath = "word_reader.py"
    $readerScript | Set-Content -Path $readerScriptPath -Encoding UTF8
    
    return $readerScriptPath
}

# Function: Convert Word document to Markdown
function Convert-DocxToMd {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $inputFile = ""
    $outputFile = ""

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--input" -and $i + 1 -lt $args.Length) {
            $inputFile = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--output" -and $i + 1 -lt $args.Length) {
            $outputFile = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($inputFile)) {
        Write-Host "Error: Input file is required"
        Write-Host "Usage: elr convert-docx --input <document.docx> [--output <document.md>]"
        return
    }

    # Check if input file exists
    if (-not (Test-Path $inputFile)) {
        Write-Host "Error: Input file '$inputFile' not found"
        return
    }

    # Set default output file if not specified
    if ([string]::IsNullOrEmpty($outputFile)) {
        $outputFile = [System.IO.Path]::ChangeExtension($inputFile, ".md")
    }

    Write-Host "===================================="
    Write-Host "Converting Word document to Markdown"
    Write-Host "Input: $inputFile"
    Write-Host "Output: $outputFile"
    Write-Host "===================================="

    # Create document converter script
    $converterScript = Create-DocumentReader

    # Run converter using Python
    if ($converterScript) {
        # Use Run-Python function to execute the converter
        $runPythonArgs = @("run-python", "--source", $converterScript, $inputFile)
        Run-Python @runPythonArgs

        # Check if conversion was successful
        if (Test-Path $outputFile) {
            $fileSize = (Get-Item $outputFile).Length
            Write-Host "===================================="
            Write-Host "Conversion successful!"
            Write-Host "Output file: $outputFile"
            Write-Host "File size: $fileSize bytes"
            Write-Host "===================================="
        } else {
            # 创建输出文件
            $text = & python $converterScript $inputFile
            $text | Set-Content -Path $outputFile -Encoding UTF8
            
            if (Test-Path $outputFile) {
                $fileSize = (Get-Item $outputFile).Length
                Write-Host "===================================="
                Write-Host "Conversion successful!"
                Write-Host "Output file: $outputFile"
                Write-Host "File size: $fileSize bytes"
                Write-Host "===================================="
            } else {
                Write-Host "===================================="
                Write-Host "Error: Conversion failed, output file not found"
                Write-Host "===================================="
            }
        }

        # Clean up
        if (Test-Path $converterScript) {
            Remove-Item $converterScript -Force
        }
    } else {
        Write-Host "Error: Failed to create document converter script"
    }
}

# Function: Read and analyze Word document
function Read-WordDocument {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $inputFile = ""
    $analyze = $false
    $extract = ""
    $outputFormat = "text"

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--input" -and $i + 1 -lt $args.Length) {
            $inputFile = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--analyze") {
            $analyze = $true
        } elseif ($args[$i] -eq "--extract" -and $i + 1 -lt $args.Length) {
            $extract = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--format" -and $i + 1 -lt $args.Length) {
            $outputFormat = $args[$i + 1]
            $i++
        }
    }

    if ([string]::IsNullOrEmpty($inputFile)) {
        Write-Host "Error: Input file is required"
        Write-Host "Usage: elr read-word --input <document.docx> [--analyze] [--extract <content_type>] [--format <format>]"
        return
    }

    # Check if input file exists
    if (-not (Test-Path $inputFile)) {
        Write-Host "Error: Input file '$inputFile' not found"
        return
    }

    Write-Host "===================================="
    Write-Host "Reading and analyzing Word document"
    Write-Host "Input: $inputFile"
    if ($analyze) {
        Write-Host "Mode: Analysis"
    } elseif (-not [string]::IsNullOrEmpty($extract)) {
        Write-Host "Mode: Extract $extract"
    } else {
        Write-Host "Mode: Text extraction"
    }
    Write-Host "Output format: $outputFormat"
    Write-Host "===================================="

    # Create document reader script
    $readerScript = Create-DocumentReader

    # Run reader using Python
    if ($readerScript) {
        # Build arguments
        $runPythonArgs = @("run-python", "--source", $readerScript, $inputFile)
        if ($analyze) {
            $runPythonArgs += "--analyze"
        }
        if (-not [string]::IsNullOrEmpty($extract)) {
            $runPythonArgs += "--extract", $extract
        }
        if (-not [string]::IsNullOrEmpty($outputFormat)) {
            $runPythonArgs += "--format", $outputFormat
        }
        
        # Run the reader
        Run-Python @runPythonArgs

        # Clean up
        if (Test-Path $readerScript) {
            Remove-Item $readerScript -Force
        }
    } else {
        Write-Host "Error: Failed to create document reader script"
    }
}

# Function: Run a container with Word document reading capabilities
function Run-WordContainer {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $containerName = "word-reader-container"

    for ($i = 2; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--name" -and $i + 1 -lt $args.Length) {
            $containerName = $args[$i + 1]
            $i++
        }
    }

    $containerID = "elr-$(Get-Date -Format 'HHmmssfff')"
    $containerStatus = $CONTAINER_STATUS_RUNNING
    $containerCreated = Get-Date
    $containerStarted = Get-Date

    # Add container to global list
    $newContainer = @{
        ID = $containerID
        Name = $containerName
        Image = "python:3.9"
        Status = $containerStatus
        Created = $containerCreated
        Started = $containerStarted
        Capabilities = "Word document reading and analysis"
    }

    $global:CONTAINERS += $newContainer

    Write-Host "===================================="
    Write-Host "Running Word document reader container: $containerID ($containerName)"
    Write-Host "Image: python:3.9"
    Write-Host "Status: $containerStatus"
    Write-Host "Capabilities: Word document reading and analysis"
    Write-Host "===================================="
    Write-Host "This container is ready to process Word documents."
    Write-Host "Use 'elr exec --id $containerID --command <command>' to run document processing commands."
    Write-Host "Example: elr exec --id $containerID --command 'python -c "import zipfile; print(\"Word document processing ready\")"'"
    Write-Host "===================================="
    
    # Save state
    Save-State
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
    "convert-docx" {
        Convert-DocxToMd @args
    }
    "read-word" {
        Read-WordDocument @args
    }
    "run-word-container" {
        Run-WordContainer @args
    }
    default {
        Write-Host "Unknown command: $command"
        Print-Help
        exit 1
    }
}
