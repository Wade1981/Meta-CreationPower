# Network service management module

# Default network configuration
$DEFAULT_CONFIG = @{
    DesktopAPI = @{ IP = "localhost"; Port = 8081 }
    PublicAPI = @{ IP = "localhost"; Port = 8080 }
    ModelService = @{ IP = "localhost"; Port = 8082 }
    MicroModel = @{ IP = "localhost"; Port = 8083 }
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
    
    $desktopApiScript = "$PSScriptRoot\..\elr\desktop_api.py"
    if (Test-Path $desktopApiScript) {
        Write-Host "Starting Desktop API in background..."
        Write-Host "Python: $pythonPath"
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $pythonPath
        $psi.Arguments = "`"$desktopApiScript`""
        $psi.WorkingDirectory = "$PSScriptRoot\..\elr"
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
    
    $elrContainerPath = "$PSScriptRoot\..\elr\elr-container.exe"
    $networkServicePath = "$PSScriptRoot\..\elr\network_service\elr-network-service.exe"
    
    if (Test-Path $elrContainerPath) {
        Write-Host "Starting ELR Container in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $elrContainerPath
        $psi.WorkingDirectory = "$PSScriptRoot\..\elr"
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
        $psi.WorkingDirectory = "$PSScriptRoot\..\elr\network_service"
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
        
        $publicApiScript = "$PSScriptRoot\..\elr_api_server.py"
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
    
    $modelServerExe = "$PSScriptRoot\..\micro_model\micro_model_server.exe"
    if (Test-Path $modelServerExe) {
        Write-Host "Starting Model Service in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $modelServerExe
        $psi.WorkingDirectory = "$PSScriptRoot\..\micro_model"
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
        
        $modelScript = "$PSScriptRoot\..\micro_model\python_server.py"
        if (Test-Path $modelScript) {
            Write-Host "Starting Model Service in background..."
            
            $psi = New-Object System.Diagnostics.ProcessStartInfo
            $psi.FileName = $pythonPath
            $psi.Arguments = "`"$modelScript`""
            $psi.WorkingDirectory = "$PSScriptRoot\..\micro_model"
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
    
    $microModelExe = "$PSScriptRoot\..\micro_model\micro_model_server.exe"
    if (Test-Path $microModelExe) {
        Write-Host "Starting Micro Model Server in background..."
        
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = $microModelExe
        $psi.WorkingDirectory = "$PSScriptRoot\..\micro_model"
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
        $mainGo = "$PSScriptRoot\..\micro_model\main.go"
        if (Test-Path $mainGo) {
            Write-Host "Building Micro Model Server..."
            Push-Location "$PSScriptRoot\..\micro_model"
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

function Check-NetworkStatus {
    $tokenManagerScript = "$PSScriptRoot\..\elr\token_manager.ps1"
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