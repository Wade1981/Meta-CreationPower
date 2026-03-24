# Python execution module
function Run-Python {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

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
            return
        }
    }
    Write-Host "Debug: Found Python interpreter: $($pythonPath.Source)"

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
            Write-Host "Debug: Executing: $($pythonPath.Source) $sourceFile"
            & $pythonPath $sourceFile
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
            Write-Host "Debug: Executing: $($pythonPath.Source) -c $pythonCode"
            & $pythonPath -c $pythonCode
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

