# Run RootPulseOS in ELR Container

Write-Host "=====================================" -ForegroundColor Green
Write-Host "RootPulseOS ELR Container Runner" -ForegroundColor Green
Write-Host "=====================================" -ForegroundColor Green

# Check ELR path
$elrPath = "E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr.ps1"
if (-not (Test-Path $elrPath)) {
    Write-Host "ELR script not found. Please check the path." -ForegroundColor Red
    Write-Host "Expected path: $elrPath" -ForegroundColor Yellow
    exit 1
}

Write-Host "ELR script path: $elrPath" -ForegroundColor Cyan

# Navigate to project root
Set-Location "$(Split-Path -Parent $MyInvocation.MyCommand.Path)\.."
$projectRoot = Get-Location
Write-Host "Project root: $projectRoot" -ForegroundColor Cyan

# Start ELR runtime
Write-Host "Starting ELR runtime..." -ForegroundColor Yellow
try {
    & powershell -ExecutionPolicy RemoteSigned -File $elrPath start
    Write-Host "ELR runtime started successfully!" -ForegroundColor Green
} catch {
    Write-Host "Failed to start ELR runtime: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Wait for ELR runtime initialization
Write-Host "Waiting for ELR runtime initialization..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# Check ELR runtime status
Write-Host "Checking ELR runtime status..." -ForegroundColor Yellow
try {
    $status = & powershell -ExecutionPolicy RemoteSigned -File $elrPath status
    Write-Host "ELR runtime status: $status" -ForegroundColor Cyan
} catch {
    Write-Host "Failed to check ELR runtime status: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Run RootPulseOS in ELR container
Write-Host "Running RootPulseOS in ELR container..." -ForegroundColor Yellow
try {
    # Run Python application using ELR
    & powershell -ExecutionPolicy RemoteSigned -File $elrPath run `
        --name rootpulseos-container `
        --language python `
        --command "python main.py" `
        --port 8000:8000
    
    Write-Host "RootPulseOS started successfully in ELR container!" -ForegroundColor Green
} catch {
    Write-Host "Failed to run RootPulseOS: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Wait for container to start
Write-Host "Waiting for container to start..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Display access information
Write-Host "=====================================" -ForegroundColor Green
Write-Host "RootPulseOS is running in ELR container!" -ForegroundColor Green
Write-Host "Access URL: http://localhost:8000" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Green

# Instructions to stop container
Write-Host "To stop RootPulseOS container, run:" -ForegroundColor Yellow
Write-Host "powershell -ExecutionPolicy RemoteSigned -File $elrPath stop-container --id rootpulseos-container" -ForegroundColor Cyan

Write-Host "" -ForegroundColor White
Write-Host "Press any key to exit..." -ForegroundColor White
$null = $Host.UI.RawUI.ReadKey('NoEcho,IncludeKeyDown')
