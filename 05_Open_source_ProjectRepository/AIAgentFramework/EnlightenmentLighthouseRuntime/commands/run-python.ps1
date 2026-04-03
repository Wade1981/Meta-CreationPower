# Python execution module
function Run-Python {
    # if (-not $global:RUNTIME_STARTED) {
    #     Write-Host "Error: ELR runtime is not running"
    #     return
    # }

    # Parse arguments
    $sourceFile = ""
    $pythonCode = ""

    # Debug: Show all arguments
    Write-Host "Debug: All arguments: $($args -join ' '); Length: $($args.Length)"

    # Check if we have at least 2 arguments (command + option)
    if ($args.Length -lt 2) {
        Write-Host "Error: Not enough arguments"
        Write-Host "Usage: elr run-python --source <script.py>"
        Write-Host "       elr run-python --code '<python code>'"
        return
    }

    # Parse arguments - start from index 1 since $args[0] is the command name
    for ($i = 1; $i -lt $args.Length; $i++) {
        Write-Host "Debug: Checking argument ${i}: $($args[$i])"
        if ($args[$i] -eq "--source" -and $i + 1 -lt $args.Length) {
            $sourceFile = $args[$i + 1]
            Write-Host "Debug: Found --source: $sourceFile"
            $i++
        } elseif ($args[$i] -eq "--code" -and $i + 1 -lt $args.Length) {
            $pythonCode = $args[$i + 1]
            Write-Host "Debug: Found --code: $pythonCode"
            $i++
        }
    }

    Write-Host "Debug: Final sourceFile: '$sourceFile'"
    Write-Host "Debug: Final pythonCode: '$pythonCode'"

    if ([string]::IsNullOrEmpty($sourceFile) -and [string]::IsNullOrEmpty($pythonCode)) {
        Write-Host "Error: Either --source or --code is required"
        Write-Host "Usage: elr run-python --source <script.py>"
        Write-Host "       elr run-python --code '<python code>'"
        return
    }

    # Check if Python is available
    Write-Host "Debug: Checking for Python interpreter..."
    
    # First check config file for Python path
    $configPath = Join-Path -Path $PSScriptRoot -ChildPath "..\micro_model\config\config.yaml"
    $pythonExe = $null
    
    if (Test-Path $configPath) {
        $configContent = Get-Content $configPath -Raw
        if ($configContent -match 'path: "([^"]+)"') {
            $pythonPathFromConfig = $matches[1]
            if (-not [string]::IsNullOrEmpty($pythonPathFromConfig)) {
                $pythonExe = Join-Path -Path $pythonPathFromConfig -ChildPath "python.exe"
                if (Test-Path $pythonExe) {
                    Write-Host "Debug: Found Python from config: $pythonExe"
                } else {
                    Write-Host "Debug: Python path from config not found: $pythonExe"
                    $pythonExe = $null
                }
            }
        }
    }
    
    # If not found in config, search in PATH
    if ($null -eq $pythonExe) {
        $pythonPath = Get-Command python -ErrorAction SilentlyContinue
        if ($null -eq $pythonPath) {
            Write-Host "Debug: python not found, trying python3..."
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
                Write-Host ""
                Write-Host "Or set Python path using:"
                Write-Host "  elr setup python --path <python-path>"
                return
            }
            $pythonExe = $pythonPath.Source
        } else {
            $pythonExe = $pythonPath.Source
        }
    }
    Write-Host "Debug: Found Python interpreter: $pythonExe"

    # Check if it's a Windows Store placeholder
    if ($pythonExe -like "*Microsoft\WindowsApps\python.exe") {
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
            Write-Host "Debug: Executing: $pythonExe $sourceFile"
            & $pythonExe $sourceFile
            $exitCode = $LASTEXITCODE
            Write-Host "Debug: Execution completed with exit code: $exitCode"
            
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
            Write-Host "Debug: Executing: $pythonExe -c $pythonCode"
            & $pythonExe -c $pythonCode
            $exitCode = $LASTEXITCODE
            Write-Host "Debug: Execution completed with exit code: $exitCode"
            
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

