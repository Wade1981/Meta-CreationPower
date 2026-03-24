# Tray application command module
function Start-TrayApplication {
    Write-Host "===================================="
    Write-Host "Starting ELR Tray Application..."
    Write-Host "===================================="
    
    try {
        # 检查ELR-Tray-App.ps1是否存在
        $trayAppPath = "$PSScriptRoot\..\ELR-Tray-App.ps1"
        
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