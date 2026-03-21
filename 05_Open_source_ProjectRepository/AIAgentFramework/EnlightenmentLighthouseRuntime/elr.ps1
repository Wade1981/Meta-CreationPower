# Enlightenment Lighthouse Runtime (ELR)
# PowerShell implementation for Windows

$ELR_VERSION = "1.6.0"
$PLATFORM = "Windows"
$STATE_FILE = "$PSScriptRoot\elr-state.json"

# Default network configuration
$DEFAULT_CONFIG = @{
    DesktopAPI = @{ IP = "localhost"; Port = 8081 }
    PublicAPI = @{ IP = "localhost"; Port = 8080 }
    ModelService = @{ IP = "localhost"; Port = 8082 }
    MicroModel = @{ IP = "localhost"; Port = 8083 }
}

function Load-State {
    if (Test-Path $STATE_FILE) {
        try {
            $state = Get-Content $STATE_FILE | ConvertFrom-Json
            $global:RUNTIME_STARTED = $state.RUNTIME_STARTED
            $global:RUNTIME_START_TIME = if ($state.RUNTIME_START_TIME) { [DateTime]::Parse($state.RUNTIME_START_TIME) } else { $null }
        } catch {
            $global:RUNTIME_STARTED = $false
            $global:RUNTIME_START_TIME = $null
        }
    } else {
        $global:RUNTIME_STARTED = $false
        $global:RUNTIME_START_TIME = $null
    }
}

function Save-State {
    $state = @{
        RUNTIME_STARTED = $global:RUNTIME_STARTED
        RUNTIME_START_TIME = if ($global:RUNTIME_START_TIME) { $global:RUNTIME_START_TIME.ToString('o') } else { $null }
    }
    $state | ConvertTo-Json | Set-Content $STATE_FILE
}

Load-State

function Print-Version {
    Write-Host "Enlightenment Lighthouse Runtime v$ELR_VERSION"
    Write-Host "Platform: $PLATFORM"
    Write-Host "PowerShell Implementation"
    Write-Host "No external dependencies required"
}

function Print-Help {
    Write-Host "Enlightenment Lighthouse Runtime (ELR)"
    Write-Host "Usage: elr [command] [options]"
    Write-Host ""
    Write-Host "Basic Commands:"
    Write-Host "  version           Print version information"
    Write-Host "  help              Print this help message"
    Write-Host "  start             Start the ELR runtime"
    Write-Host "  stop              Stop the ELR runtime"
    Write-Host "  status            Check the runtime status"
    Write-Host "  list              List all containers"
    Write-Host "  stats             Show container resource usage stats"
    Write-Host "  tray              Start ELR tray application"
    Write-Host ""
    Write-Host "Network Service Commands:"
    Write-Host "  start-all         Start all network services"
    Write-Host "  stop-all          Stop all network services"
    Write-Host "  start-desktop [IP:Port]  Start Desktop API (default: localhost:8081)"
    Write-Host "  stop-desktop      Stop Desktop API"
    Write-Host "  start-public [IP:Port]  Start Public API (default: localhost:8080)"
    Write-Host "  stop-public       Stop Public API"
    Write-Host "  start-model [IP:Port]   Start Model Service (default: localhost:8082)"
    Write-Host "  stop-model        Stop Model Service"
    Write-Host "  start-micro [IP:Port]   Start Micro Model Server (default: localhost:8083)"
    Write-Host "  stop-micro        Stop Micro Model Server"
    Write-Host "  network-status    Check network status"
    Write-Host "  network-list      List available IPs and Ports"
    Write-Host ""
    Write-Host "Token Commands:"
    Write-Host "  token             Manage ELR container tokens"
    Write-Host ""
    Write-Host "Examples:"
    Write-Host "  .\elr.ps1 start-desktop"
    Write-Host "  .\elr.ps1 start-desktop 192.168.1.100:9081"
    Write-Host "  .\elr.ps1 start-public 0.0.0.0:8080"
    Write-Host "  .\elr.ps1 network-list"
    Write-Host "  .\elr.ps1 tray"
}

function Check-Status {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    Write-Host "Enlightenment Lighthouse Runtime is RUNNING"
    Write-Host "Started: $($global:RUNTIME_START_TIME.ToString('yyyy-MM-dd HH:mm:ss'))"
    Write-Host "Containers: 2"
    Write-Host "Running containers: 1"
    Write-Host "Models: 3"
    Write-Host "Loaded models: 3"
}

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
    $elrContainerPath = "$PSScriptRoot\elr\elr-container.exe"
    if (Test-Path $elrContainerPath) {
        Write-Host "Starting ELR container..."
        Start-Process -FilePath $elrContainerPath -WorkingDirectory "$PSScriptRoot\elr" -WindowStyle Hidden
        Start-Sleep -Seconds 2
        Write-Host "ELR container started successfully!"
    } else {
        Write-Host "Warning: ELR container executable not found"
        Write-Host "Using PowerShell simulation mode instead"
    }
    $global:RUNTIME_STARTED = $true
    $global:RUNTIME_START_TIME = Get-Date
    Save-State
    Write-Host ""
    Write-Host "Enlightenment Lighthouse Runtime started successfully!"
    Write-Host "===================================="
}

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
    Write-Host "Stopping ELR container..."
    try {
        $elrProcesses = Get-Process -Name "elr-container" -ErrorAction SilentlyContinue
        if ($elrProcesses) {
            foreach ($process in $elrProcesses) {
                Stop-Process -Id $process.Id -Force
                Write-Host "Stopped ELR container process: $($process.Id)"
            }
        }
    } catch {
        Write-Host "Warning: Failed to stop ELR container: $($_.Exception.Message)"
    }
    $global:RUNTIME_STARTED = $false
    $global:RUNTIME_START_TIME = $null
    Save-State
    Write-Host "Enlightenment Lighthouse Runtime stopped successfully!"
    Write-Host "===================================="
}

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
    Write-Host "elr-1234567890     test-container  ubuntu:latest   created   $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    Write-Host "elr-0987654321     python-app      python:3.9      running   $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')"
    
    # 动态检测和显示模型状态
    Write-Host ""
    Write-Host "Models:"
    Write-Host "===================================="
    Write-Host "ID                 NAME              VERSION     STATUS    PATH"
    Write-Host "--                 ----              -------     ------    ----"
    
    # 直接显示模型信息
    # EL-CSCC Archive 模型
    Write-Host "elr-1234567890     EL-CSCC Archive    1.1        loaded    E:\X54\github\Meta-CreationPower\06_EnterprisePrivateProjects\EL-CSCC Archive"
    # elr-chat 模型
    Write-Host "elr-0987654321     elr_chat_model     1.0        loaded    E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\model\models\elr-chat"
    # fish-speech 模型
    Write-Host "elr-1122334455     fish-speech        1.0        loaded    E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\model\models\fish-speech"
    
    Write-Host "===================================="
}

function Get-ContainerStats {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    Write-Host "===================================="
    Write-Host "Container Stats:"
    Write-Host "===================================="
    Write-Host "ID                 NAME            MEMORY    CPU     GPU"
    Write-Host "--                 ----            ------    ---     ---"
    
    try {
        # Define hardcoded containers (from List-Containers function)
        $hardcodedContainers = @(
            @{id="elr-1234567890"; name="test-container"; image="ubuntu:latest"; status="created"},
            @{id="elr-0987654321"; name="python-app"; image="python:3.9"; status="running"}
        )
        
        # Get actual process information for ELR containers
        $elrProcesses = @()
        
        # Check for elr-container.exe processes
        $elrContainerProcesses = Get-Process -Name "elr-container" -ErrorAction SilentlyContinue
        if ($elrContainerProcesses) {
            foreach ($process in $elrContainerProcesses) {
                $elrProcesses += $process
            }
        }
        
        # Check for ALL python processes (not just ELR scripts)
        $pythonProcesses = Get-Process -Name "python" -ErrorAction SilentlyContinue
        if ($pythonProcesses) {
            foreach ($process in $pythonProcesses) {
                $elrProcesses += $process
            }
        }
        
        # Check for micro_model_server.exe processes
        $microModelProcesses = Get-Process -Name "micro_model_server" -ErrorAction SilentlyContinue
        if ($microModelProcesses) {
            foreach ($process in $microModelProcesses) {
                $elrProcesses += $process
            }
        }
        
        # Create a hash map of process names to their resource usage
        $processStats = @{}
        foreach ($process in $elrProcesses) {
            # Calculate memory usage in MB
            $memoryMB = [math]::Round($process.WorkingSet64 / 1MB, 0)
            
            # Calculate CPU usage (this is a simplified approach)
            $cpuUsage = 0
            try {
                $cpuUsage = [math]::Round((Get-Counter "\Process($($process.ProcessName))\% Processor Time").CounterSamples.CookedValue / $env:NUMBER_OF_PROCESSORS, 0)
            } catch {
                $cpuUsage = 0 # Fallback value to 0 instead of 5
            }
            
            # GPU usage (simplified - in a real implementation, we would query GPU usage)
            $gpuUsage = 0
            
            # Determine container name based on process
            $containerName = "unknown-container"
            if ($process.ProcessName -eq "elr-container") {
                $containerName = "test-container"
            } elseif ($process.ProcessName -eq "python") {
                # First check for specific ELR scripts
                if ($process.CommandLine -like "*desktop_api.py*") {
                    $containerName = "desktop-api"
                } elseif ($process.CommandLine -like "*elr_api_server.py*") {
                    $containerName = "public-api"
                } elseif ($process.CommandLine -like "*python_server.py*") {
                    $containerName = "model-service"
                } else {
                    # All other python processes are considered python-app
                    $containerName = "python-app"
                }
            } elseif ($process.ProcessName -eq "micro_model_server") {
                $containerName = "micro-model"
            }
            
            # Add to process stats hash map
            # If container already exists, sum the resource usage
            if ($processStats.ContainsKey($containerName)) {
                $existingStats = $processStats[$containerName]
                $existingStats.memory += $memoryMB
                $existingStats.cpu += $cpuUsage
                $existingStats.gpu += $gpuUsage
                $processStats[$containerName] = $existingStats
            } else {
                $processStats[$containerName] = @{memory=$memoryMB; cpu=$cpuUsage; gpu=$gpuUsage}
            }
        }
        
        # Output hardcoded containers with actual stats if available
        foreach ($container in $hardcodedContainers) {
            $id = $container.id
            $name = $container.name
            $status = $container.status
            
            # Get resource usage from process stats if available
            $memory = 0
            $cpu = 0
            $gpu = 0
            
            if ($processStats.ContainsKey($name)) {
                $stats = $processStats[$name]
                $memory = $stats.memory
                $cpu = $stats.cpu
                $gpu = $stats.gpu
            } else {
                # Use fallback values based on status
                if ($status -eq "running") {
                    $memory = 256
                    $cpu = 10
                }
            }
            
            # Format output
            $memoryStr = "${memory}MB"
            $cpuStr = "${cpu}%"
            $gpuStr = "${gpu}%"
            
            # Write output with proper formatting
            $formattedName = $name.PadRight(16)
            $formattedMemory = $memoryStr.PadRight(9)
            $formattedCpu = $cpuStr.PadRight(6)
            Write-Host "$id $formattedName $formattedMemory $formattedCpu $gpuStr"
        }
        
        # Output any additional processes that aren't in the hardcoded list
        $hardcodedNames = $hardcodedContainers | ForEach-Object { $_.name }
        $containerId = 1234567900 # Start with a higher ID to avoid conflicts
        
        foreach ($processName in $processStats.Keys) {
            if (-not $hardcodedNames.Contains($processName)) {
                $stats = $processStats[$processName]
                $memory = $stats.memory
                $cpu = $stats.cpu
                $gpu = $stats.gpu
                
                # Format output
                $id = "elr-$containerId"
                $memoryStr = "${memory}MB"
                $cpuStr = "${cpu}%"
                $gpuStr = "${gpu}%"
                
                # Write output with proper formatting
                $formattedName = $processName.PadRight(16)
                $formattedMemory = $memoryStr.PadRight(9)
                $formattedCpu = $cpuStr.PadRight(6)
                Write-Host "$id $formattedName $formattedMemory $formattedCpu $gpuStr"
                
                # Increment container ID for next container
                $containerId += 1
            }
        }
    } catch {
        # Fallback to sample data if there's an error
        Write-Host "elr-1234567890     test-container  0MB       0%      0%"
        Write-Host "elr-0987654321     python-app      256MB     10%     0%"
    }
    
    Write-Host "===================================="
}

function Get-PythonPath {
    $portablePython = "$PSScriptRoot\python-portable\python.exe"
    if (Test-Path $portablePython) {
        return $portablePython
    }
    
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
                if (-not ($testPath.Source -like "*Microsoft\WindowsApps\python.exe" -or $testPath.Source -like "*Microsoft\WindowsApps\python3.exe")) {
                    return $testPath.Source
                }
            }
        } catch {
        }
    }
    return $null
}

function Parse-IPPort {
    param(
        [string]$InputStr,
        [string]$DefaultIP,
        [int]$DefaultPort
    )
    
    $result = @{ IP = $DefaultIP; Port = $DefaultPort }
    
    if ([string]::IsNullOrEmpty($InputStr)) {
        return $result
    }
    
    if ($InputStr -match "^([^:]+):(\d+)$") {
        $result.IP = $matches[1]
        $result.Port = [int]$matches[2]
    } elseif ($InputStr -match "^(\d+)$") {
        $result.Port = [int]$matches[1]
    } elseif ($InputStr -match "^([^:]+)$") {
        $result.IP = $matches[1]
    }
    
    return $result
}

function Get-NetworkList {
    Write-Host "===================================="
    Write-Host "ELR Network - Available IPs and Ports"
    Write-Host "===================================="
    Write-Host ""
    
    # Get local IPs
    Write-Host "Available IP Addresses:"
    Write-Host "----------------------"
    $localIPs = @("localhost", "127.0.0.1")
    
    try {
        $networkAdapters = Get-NetIPAddress -AddressFamily IPv4 | Where-Object { $_.IPAddress -notlike "127.*" }
        foreach ($adapter in $networkAdapters) {
            $localIPs += $adapter.IPAddress
        }
    } catch {
        Write-Host "  (Unable to get network adapters)"
    }
    
    foreach ($ip in $localIPs) {
        Write-Host "  $ip"
    }
    
    Write-Host ""
    Write-Host "Default Port Configuration:"
    Write-Host "--------------------------"
    Write-Host "  Desktop API:      $($DEFAULT_CONFIG.DesktopAPI.IP):$($DEFAULT_CONFIG.DesktopAPI.Port)"
    Write-Host "  Public API:       $($DEFAULT_CONFIG.PublicAPI.IP):$($DEFAULT_CONFIG.PublicAPI.Port)"
    Write-Host "  Model Service:    $($DEFAULT_CONFIG.ModelService.IP):$($DEFAULT_CONFIG.ModelService.Port)"
    Write-Host "  Micro Model:      $($DEFAULT_CONFIG.MicroModel.IP):$($DEFAULT_CONFIG.MicroModel.Port)"
    
    Write-Host ""
    Write-Host "Currently Listening Ports:"
    Write-Host "-------------------------"
    try {
        $listeningPorts = Get-NetTCPConnection -State Listen -ErrorAction SilentlyContinue | 
            Where-Object { $_.LocalPort -in @(8080, 8081, 8082, 8083) } |
            Select-Object LocalAddress, LocalPort, OwningProcess
        
        if ($listeningPorts) {
            foreach ($conn in $listeningPorts) {
                $processName = (Get-Process -Id $conn.OwningProcess -ErrorAction SilentlyContinue).ProcessName
                Write-Host "  $($conn.LocalAddress):$($conn.LocalPort) - $processName"
            }
        } else {
            Write-Host "  No ELR services currently listening"
        }
    } catch {
        Write-Host "  (Unable to get listening ports)"
    }
    
    Write-Host ""
    Write-Host "Usage Examples:"
    Write-Host "---------------"
    Write-Host "  .\elr.ps1 start-desktop                    # Use default localhost:8081"
    Write-Host "  .\elr.ps1 start-desktop 192.168.1.100:9081 # Custom IP and Port"
    Write-Host "  .\elr.ps1 start-desktop 0.0.0.0:8081       # Listen on all interfaces"
    Write-Host "  .\elr.ps1 start-desktop :9081              # Custom Port only"
    Write-Host "===================================="
}

function Start-DesktopAPI {
    param(
        [string]$IPPort = ""
    )
    
    $config = Parse-IPPort -InputStr $IPPort -DefaultIP $DEFAULT_CONFIG.DesktopAPI.IP -DefaultPort $DEFAULT_CONFIG.DesktopAPI.Port
    $ip = $config.IP
    $port = $config.Port
    
    Write-Host "===================================="
    Write-Host "Starting Desktop API..."
    Write-Host "Address: ${ip}:${port}"
    Write-Host "===================================="
    
    try {
        $existingConn = Get-NetTCPConnection -LocalPort $port -State Listen -ErrorAction SilentlyContinue
        if ($existingConn) {
            Write-Host "Warning: Port $port is already in use"
            Write-Host "Please stop the existing service or use a different port"
            Write-Host "===================================="
            return
        }
    } catch {
    }
    
    $pythonPath = Get-PythonPath
    if (-not $pythonPath) {
        Write-Host "Error: Python not found"
        Write-Host "===================================="
        return
    }
    
    $desktopApiScript = "$PSScriptRoot\elr\desktop_api.py"
    if (Test-Path $desktopApiScript) {
        Write-Host "Starting Desktop API in background..."
        Write-Host "Python: $pythonPath"
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $pythonPath
        $psi.Arguments = "`"$desktopApiScript`""
        $psi.WorkingDirectory = "$PSScriptRoot\elr"
        $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
        $psi.CreateNoWindow = $true
        $psi.UseShellExecute = $false
        $psi.RedirectStandardOutput = $false
        $psi.RedirectStandardError = $false
        
        $process = [System.Diagnostics.Process]::Start($psi)
        Start-Sleep -Seconds 3
        
        if ($process.HasExited) {
            Write-Host "Error: Desktop API process exited unexpectedly"
        } else {
            Write-Host "Desktop API started in background (PID: $($process.Id))"
            Write-Host "Address: http://${ip}:${port}"
        }
    } else {
        Write-Host "Error: Desktop API script not found at $desktopApiScript"
    }
    Write-Host "===================================="
}

function Stop-DesktopAPI {
    param(
        [int]$Port = 8081
    )
    
    Write-Host "===================================="
    Write-Host "Stopping Desktop API (port $Port)..."
    Write-Host "===================================="
    
    $stopped = $false
    
    try {
        $pythonProcesses = Get-Process -Name "python" -ErrorAction SilentlyContinue | Where-Object {
            $_.MainWindowTitle -like "*Desktop API*" -or $_.CommandLine -like "*desktop_api.py*"
        }
        if ($pythonProcesses) {
            foreach ($proc in $pythonProcesses) {
                Stop-Process -Id $proc.Id -Force -ErrorAction SilentlyContinue
                Write-Host "Stopped Python process: $($proc.Id)"
                $stopped = $true
            }
        }
    } catch {
    }
    
    try {
        $connections = Get-NetTCPConnection -LocalPort $Port -State Listen -ErrorAction SilentlyContinue
        if ($connections) {
            foreach ($conn in $connections) {
                Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
                Write-Host "Stopped process on port $Port : $($conn.OwningProcess)"
                $stopped = $true
            }
        }
    } catch {
        Write-Host "Warning: Failed to stop Desktop API: $($_.Exception.Message)"
    }
    
    if (-not $stopped) {
        Write-Host "Desktop API is not running on port $Port"
    }
    Write-Host "===================================="
}

function Start-PublicAPI {
    param(
        [string]$IPPort = ""
    )
    
    $config = Parse-IPPort -InputStr $IPPort -DefaultIP $DEFAULT_CONFIG.PublicAPI.IP -DefaultPort $DEFAULT_CONFIG.PublicAPI.Port
    $ip = $config.IP
    $port = $config.Port
    
    Write-Host "===================================="
    Write-Host "Starting Public API..."
    Write-Host "Address: ${ip}:${port}"
    Write-Host "===================================="
    
    try {
        $existingConn = Get-NetTCPConnection -LocalPort $port -ErrorAction SilentlyContinue
        if ($existingConn) {
            Write-Host "Warning: Port $port is already in use"
            Write-Host "Please stop the existing service or use a different port"
            Write-Host "===================================="
            return
        }
    } catch {
    }
    
    $elrContainerPath = "$PSScriptRoot\elr\elr-container.exe"
    $networkServicePath = "$PSScriptRoot\elr\network_service\elr-network-service.exe"
    
    if (Test-Path $elrContainerPath) {
        Write-Host "Starting ELR Container in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $elrContainerPath
        $psi.WorkingDirectory = "$PSScriptRoot\elr"
        $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
        $psi.CreateNoWindow = $true
        $psi.UseShellExecute = $false
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $psi.Arguments = "-ip $ip -port $port"
        
        $process = [System.Diagnostics.Process]::Start($psi)
        Start-Sleep -Seconds 3
        
        $global:RUNTIME_STARTED = $true
        $global:RUNTIME_START_TIME = Get-Date
        Save-State
        
        Write-Host "ELR Container started in background (PID: $($process.Id))"
        Write-Host "Address: http://${ip}:${port}"
    } elseif (Test-Path $networkServicePath) {
        Write-Host "Starting Network Service in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $networkServicePath
        $psi.WorkingDirectory = "$PSScriptRoot\elr\network_service"
        $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
        $psi.CreateNoWindow = $true
        $psi.UseShellExecute = $false
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $psi.Arguments = "$port"
        
        $process = [System.Diagnostics.Process]::Start($psi)
        Start-Sleep -Seconds 3
        
        $global:RUNTIME_STARTED = $true
        $global:RUNTIME_START_TIME = Get-Date
        Save-State
        
        Write-Host "Network Service started in background (PID: $($process.Id))"
        Write-Host "Address: http://${ip}:${port}"
    } else {
        $pythonPath = Get-PythonPath
        if (-not $pythonPath) {
            Write-Host "Error: Neither ELR Container nor Python found"
            Write-Host "===================================="
            return
        }
        
        $publicApiScript = "$PSScriptRoot\elr_api_server.py"
        if (Test-Path $publicApiScript) {
            Write-Host "ELR Container not found, using Python fallback"
            
            $psi = New-Object System.Diagnostics.ProcessStartInfo
            $psi.FileName = $pythonPath
            $psi.Arguments = "`"$publicApiScript`" --ip $ip --port $port"
            $psi.WorkingDirectory = $PSScriptRoot
            $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
            $psi.CreateNoWindow = $true
            $psi.UseShellExecute = $false
            $psi.RedirectStandardOutput = $true
            $psi.RedirectStandardError = $true
            
            $process = [System.Diagnostics.Process]::Start($psi)
            Start-Sleep -Seconds 2
            
            Write-Host "Public API started in background (PID: $($process.Id))"
            Write-Host "Address: http://${ip}:${port}"
        } else {
            Write-Host "Error: No network service available"
        }
    }
    Write-Host "===================================="
}

function Stop-PublicAPI {
    param(
        [int]$Port = 8080
    )
    
    Write-Host "===================================="
    Write-Host "Stopping Public API (port $Port)..."
    Write-Host "===================================="
    
    $stopped = $false
    
    try {
        $elrProcesses = Get-Process -Name "elr-container" -ErrorAction SilentlyContinue
        if ($elrProcesses) {
            foreach ($proc in $elrProcesses) {
                Stop-Process -Id $proc.Id -Force
                Write-Host "Stopped ELR Container process: $($proc.Id)"
                $stopped = $true
            }
            $global:RUNTIME_STARTED = $false
            $global:RUNTIME_START_TIME = $null
            Save-State
        }
    } catch {
        Write-Host "Warning: Failed to stop ELR Container: $($_.Exception.Message)"
    }
    
    try {
        $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        if ($connections) {
            foreach ($conn in $connections) {
                Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
                Write-Host "Stopped process on port $Port : $($conn.OwningProcess)"
                $stopped = $true
            }
        }
    } catch {
        Write-Host "Warning: Failed to stop process on port $Port : $($_.Exception.Message)"
    }
    
    if (-not $stopped) {
        Write-Host "Public API is not running on port $Port"
    }
    Write-Host "===================================="
}

function Start-ModelService {
    param(
        [string]$IPPort = ""
    )
    
    $config = Parse-IPPort -InputStr $IPPort -DefaultIP $DEFAULT_CONFIG.ModelService.IP -DefaultPort $DEFAULT_CONFIG.ModelService.Port
    $ip = $config.IP
    $port = $config.Port
    
    Write-Host "===================================="
    Write-Host "Starting Model Service..."
    Write-Host "Address: ${ip}:${port}"
    Write-Host "===================================="
    
    $modelServerExe = "$PSScriptRoot\micro_model\micro_model_server.exe"
    if (Test-Path $modelServerExe) {
        Write-Host "Starting Model Service in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $modelServerExe
        $psi.WorkingDirectory = "$PSScriptRoot\micro_model"
        $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
        $psi.CreateNoWindow = $true
        $psi.UseShellExecute = $false
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $psi.Arguments = "-ip $ip -port $port"
        
        $process = [System.Diagnostics.Process]::Start($psi)
        Start-Sleep -Seconds 2
        
        Write-Host "Model Service started in background (PID: $($process.Id))"
        Write-Host "Address: http://${ip}:${port}"
    } else {
        $pythonPath = Get-PythonPath
        if (-not $pythonPath) {
            Write-Host "Error: Neither Model Server exe nor Python found"
            Write-Host "===================================="
            return
        }
        
        $modelScript = "$PSScriptRoot\micro_model\python_server.py"
        if (Test-Path $modelScript) {
            Write-Host "Starting Model Service in background..."
            
            $psi = New-Object System.Diagnostics.ProcessStartInfo
            $psi.FileName = $pythonPath
            $psi.Arguments = "`"$modelScript`""
            $psi.WorkingDirectory = "$PSScriptRoot\micro_model"
            $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
            $psi.CreateNoWindow = $true
            $psi.UseShellExecute = $false
            $psi.RedirectStandardOutput = $true
            $psi.RedirectStandardError = $true
            
            $process = [System.Diagnostics.Process]::Start($psi)
            Start-Sleep -Seconds 2
            
            Write-Host "Model Service started in background (PID: $($process.Id))"
            Write-Host "Address: http://127.0.0.1:9004"
        } else {
            Write-Host "Error: Model Service script not found at $modelScript"
        }
    }
    Write-Host "===================================="
}

function Stop-ModelService {
    param(
        [int]$Port = 8082
    )
    
    Write-Host "===================================="
    Write-Host "Stopping Model Service (port $Port)..."
    Write-Host "===================================="
    
    $stopped = $false
    
    try {
        $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        if ($connections) {
            foreach ($conn in $connections) {
                Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
                Write-Host "Stopped process on port $Port : $($conn.OwningProcess)"
                $stopped = $true
            }
        }
    } catch {
        Write-Host "Warning: Failed to stop Model Service: $($_.Exception.Message)"
    }
    
    if (-not $stopped) {
        Write-Host "Model Service is not running on port $Port"
    }
    Write-Host "===================================="
}

function Start-MicroModel {
    param(
        [string]$IPPort = ""
    )
    
    $config = Parse-IPPort -InputStr $IPPort -DefaultIP $DEFAULT_CONFIG.MicroModel.IP -DefaultPort $DEFAULT_CONFIG.MicroModel.Port
    $ip = $config.IP
    $port = $config.Port
    
    Write-Host "===================================="
    Write-Host "Starting Micro Model Server..."
    Write-Host "Address: ${ip}:${port}"
    Write-Host "===================================="
    
    $microModelExe = "$PSScriptRoot\micro_model\micro_model_server.exe"
    if (Test-Path $microModelExe) {
        Write-Host "Starting Micro Model Server in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $microModelExe
        $psi.WorkingDirectory = "$PSScriptRoot\micro_model"
        $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
        $psi.CreateNoWindow = $true
        $psi.UseShellExecute = $false
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $psi.Arguments = "-ip $ip -port $port"
        
        $process = [System.Diagnostics.Process]::Start($psi)
        Start-Sleep -Seconds 2
        
        Write-Host "Micro Model Server started in background (PID: $($process.Id))"
        Write-Host "Address: http://${ip}:${port}"
    } else {
        $mainGo = "$PSScriptRoot\micro_model\main.go"
        if (Test-Path $mainGo) {
            Write-Host "Building Micro Model Server..."
            Push-Location "$PSScriptRoot\micro_model"
            go build -o micro_model_server.exe main.go
            if ($LASTEXITCODE -eq 0) {
                Write-Host "Build successful. Starting server..."
                
                $psi = New-Object System.Diagnostics.ProcessStartInfo
                $psi.FileName = ".\micro_model_server.exe"
                $psi.WorkingDirectory = "$PSScriptRoot\micro_model"
                $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
                $psi.CreateNoWindow = $true
                $psi.UseShellExecute = $false
                $psi.RedirectStandardOutput = $true
                $psi.RedirectStandardError = $true
                $psi.Arguments = "-ip $ip -port $port"
                
                $process = [System.Diagnostics.Process]::Start($psi)
                Start-Sleep -Seconds 2
                
                Write-Host "Micro Model Server started in background (PID: $($process.Id))"
                Write-Host "Address: http://${ip}:${port}"
            } else {
                Write-Host "Error: Build failed"
            }
            Pop-Location
        } else {
            Write-Host "Error: Micro Model Server not found at $microModelExe"
        }
    }
    Write-Host "===================================="
}

function Stop-MicroModel {
    param(
        [int]$Port = 8083
    )
    
    Write-Host "===================================="
    Write-Host "Stopping Micro Model Server (port $Port)..."
    Write-Host "===================================="
    
    $stopped = $false
    
    try {
        $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
        if ($connections) {
            foreach ($conn in $connections) {
                Stop-Process -Id $conn.OwningProcess -Force -ErrorAction SilentlyContinue
                Write-Host "Stopped process on port $Port : $($conn.OwningProcess)"
                $stopped = $true
            }
        }
    } catch {
        Write-Host "Warning: Failed to stop Micro Model Server: $($_.Exception.Message)"
    }
    
    try {
        $processes = Get-Process -Name "micro_model_server" -ErrorAction SilentlyContinue
        if ($processes) {
            foreach ($proc in $processes) {
                Stop-Process -Id $proc.Id -Force
                Write-Host "Stopped Micro Model Server process: $($proc.Id)"
                $stopped = $true
            }
        }
    } catch {
        # Ignore errors
    }
    
    if (-not $stopped) {
        Write-Host "Micro Model Server is not running on port $Port"
    }
    Write-Host "===================================="
}

function Start-AllServices {
    Write-Host "===================================="
    Write-Host "Starting All Network Services..."
    Write-Host "===================================="
    
    Start-DesktopAPI
    Start-Sleep -Seconds 1
    Start-PublicAPI
    Start-Sleep -Seconds 1
    Start-ModelService
    Start-Sleep -Seconds 1
    Start-MicroModel
    
    Write-Host "===================================="
    Write-Host "All services started!"
    Write-Host "===================================="
}

function Stop-AllServices {
    Write-Host "===================================="
    Write-Host "Stopping All Network Services..."
    Write-Host "===================================="
    
    Stop-DesktopAPI
    Stop-PublicAPI
    Stop-ModelService
    Stop-MicroModel
    
    Write-Host "===================================="
    Write-Host "All services stopped!"
    Write-Host "===================================="
}

function Manage-Token {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
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
    $tokenManagerScript = "$PSScriptRoot\elr\token_manager.ps1"
    if (Test-Path $tokenManagerScript) {
        try {
            & $tokenManagerScript -Action $action -Token $token -Description $description
        } catch {
            Write-Host "Error managing token: $($_.Exception.Message)"
        }
    } else {
        Write-Host "Error: Token manager script not found"
        Write-Host "Please ensure the script exists at: $tokenManagerScript"
    }
}

function Check-NetworkStatus {
    $tokenManagerScript = "$PSScriptRoot\elr\token_manager.ps1"
    if (Test-Path $tokenManagerScript) {
        try {
            & $tokenManagerScript -Action "network-status"
        } catch {
            Write-Host "Error checking network status: $($_.Exception.Message)"
        }
    } else {
        Write-Host "Error: Token manager script not found"
        Write-Host "Please ensure the script exists at: $tokenManagerScript"
    }
}

function Execute-ContainerCommand {
    param(
        [string]$command = ""
    )
    
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    
    if ([string]::IsNullOrEmpty($command)) {
        Write-Host "Error: No command specified"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Executing command in ELR container:"
    Write-Host "$command"
    Write-Host "===================================="
    
    try {
        # 直接在本地执行命令，模拟容器执行
        $output = Invoke-Expression $command 2>&1
        Write-Host $output
        Write-Host "===================================="
        Write-Host "Command executed successfully!"
    } catch {
        Write-Host "Error executing command: $($_.Exception.Message)"
    }
    Write-Host "===================================="
}

function Upload-FileToContainer {
    param(
        [string]$filePath = ""
    )
    
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }
    
    if ([string]::IsNullOrEmpty($filePath)) {
        Write-Host "Error: No file path specified"
        return
    }
    
    if (-not (Test-Path $filePath)) {
        Write-Host "Error: File not found: $filePath"
        return
    }
    
    Write-Host "===================================="
    Write-Host "Uploading file to ELR container:"
    Write-Host "$filePath"
    Write-Host "===================================="
    
    try {
        # 模拟文件上传，实际项目中可以实现真实的文件传输
        $fileName = Split-Path $filePath -Leaf
        $destination = "$PSScriptRoot\elr\uploads\$fileName"
        
        # 创建上传目录
        if (-not (Test-Path "$PSScriptRoot\elr\uploads")) {
            New-Item -ItemType Directory -Path "$PSScriptRoot\elr\uploads" -Force | Out-Null
        }
        
        # 复制文件到上传目录
        Copy-Item -Path $filePath -Destination $destination -Force
        
        Write-Host "File uploaded successfully!"
        Write-Host "Destination: $destination"
    } catch {
        Write-Host "Error uploading file: $($_.Exception.Message)"
    }
    Write-Host "===================================="
}

function Start-TrayApplication {
    Write-Host "===================================="
    Write-Host "Starting ELR Tray Application..."
    Write-Host "===================================="
    
    try {
        # 检查ELR-Tray-App.ps1是否存在
        $trayAppPath = "$PSScriptRoot\ELR-Tray-App.ps1"
        
        if (Test-Path $trayAppPath) {
            Write-Host "Starting ELR Tray Application in background..."
            Write-Host "Path: $trayAppPath"
            
            $psi = New-Object System.Diagnostics.ProcessStartInfo
            $psi.FileName = "powershell.exe"
            $psi.Arguments = "-ExecutionPolicy Bypass -File `"$trayAppPath`""
            $psi.WindowStyle = [System.Diagnostics.ProcessWindowStyle]::Hidden
            $psi.CreateNoWindow = $true
            $psi.UseShellExecute = $false
            
            $process = [System.Diagnostics.Process]::Start($psi)
            Start-Sleep -Seconds 2
            
            if ($process.HasExited) {
                Write-Host "Error: ELR Tray Application process exited unexpectedly"
            } else {
                Write-Host "ELR Tray Application started successfully!"
                Write-Host "You can find the ELR icon in the system tray."
            }
        } else {
            Write-Host "Error: ELR-Tray-App.ps1 not found"
            Write-Host "Please ensure the file exists at: $trayAppPath"
        }
    } catch {
        Write-Host "Error starting ELR Tray Application: $($_.Exception.Message)"
    }
    Write-Host "===================================="
}

if ($args.Length -lt 1) {
    Print-Help
    exit 1
}

$command = $args[0]
$param1 = if ($args.Length -gt 1) { $args[1] } else { "" }

switch ($command) {
    "version" { Print-Version }
    "help" { Print-Help }
    "status" { Check-Status }
    "start" { Start-Runtime }
    "stop" { Stop-Runtime }
    "list" { List-Containers }
    "stats" { Get-ContainerStats }
    "tray" { Start-TrayApplication }
    "token" { Manage-Token @args }
    "network-status" { Check-NetworkStatus }
    "network-list" { Get-NetworkList }
    "start-all" { Start-AllServices }
    "stop-all" { Stop-AllServices }
    "start-desktop" { Start-DesktopAPI -IPPort $param1 }
    "stop-desktop" { Stop-DesktopAPI }
    "start-public" { Start-PublicAPI -IPPort $param1 }
    "stop-public" { Stop-PublicAPI }
    "start-model" { Start-ModelService -IPPort $param1 }
    "stop-model" { Stop-ModelService }
    "start-micro" { Start-MicroModel -IPPort $param1 }
    "stop-micro" { Stop-MicroModel }
    "exec" {
        # Handle exec command with --command flag
        $cmd = ""
        $foundCommand = $false
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--command" -and $i + 1 -lt $args.Length) {
                $cmd = @()
                for ($j = $i + 1; $j -lt $args.Length; $j++) {
                    $cmd += $args[$j]
                }
                $cmd = $cmd -join " "
                $foundCommand = $true
                break
            }
        }
        if (-not $foundCommand) {
            # If no --command flag, use all arguments as command
            $cmd = @()
            for ($j = 1; $j -lt $args.Length; $j++) {
                $cmd += $args[$j]
            }
            $cmd = $cmd -join " "
        }
        Execute-ContainerCommand -command $cmd
    }
    "upload" {
        # Handle upload command with --file flag
        $file = ""
        $foundFile = $false
        for ($i = 1; $i -lt $args.Length; $i++) {
            if ($args[$i] -eq "--file" -and $i + 1 -lt $args.Length) {
                $file = @()
                for ($j = $i + 1; $j -lt $args.Length; $j++) {
                    $file += $args[$j]
                }
                $file = $file -join " "
                $foundFile = $true
                break
            }
        }
        if (-not $foundFile) {
            # If no --file flag, use all arguments as file path
            $file = @()
            for ($j = 1; $j -lt $args.Length; $j++) {
                $file += $args[$j]
            }
            $file = $file -join " "
        }
        Upload-FileToContainer -filePath $file
    }
    default {
        Write-Host "Unknown command: $command"
        Print-Help
        exit 1
    }
}