# Stop runtime command module
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

    # Set runtime status
    $global:RUNTIME_STARTED = $false
    $global:RUNTIME_START_TIME = $null

    Write-Host "Enlightenment Lighthouse Runtime stopped successfully!"
    Write-Host "===================================="
    
    # Save state
    Save-State
}

