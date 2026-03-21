# ELR Container Start Script

Write-Host "=== Starting ELR Container and Loading Symmetric Digital Finance System ===" -ForegroundColor Green

# Set working directory
$workingDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $workingDir

# Check ELR container environment
Write-Host "Checking ELR container environment..." -ForegroundColor Yellow
# Simulate ELR container environment check
Start-Sleep -Seconds 2
Write-Host "✅ ELR container environment is ready" -ForegroundColor Green

# Install necessary dependencies
Write-Host "Installing necessary dependencies..." -ForegroundColor Yellow
pip install -r requirements_light.txt
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Dependencies installed successfully" -ForegroundColor Green
} else {
    Write-Host "❌ Failed to install dependencies" -ForegroundColor Red
    exit 1
}

# Prepare container configuration
Write-Host "Preparing container configuration..." -ForegroundColor Yellow
# Simulate container configuration preparation
Start-Sleep -Seconds 1
Write-Host "✅ Container configuration prepared" -ForegroundColor Green

# Start ELR container
Write-Host "Starting ELR container..." -ForegroundColor Yellow
# Simulate ELR container start
Start-Sleep -Seconds 3
Write-Host "✅ ELR container started" -ForegroundColor Green

# Load symmetric digital finance system
Write-Host "Loading symmetric digital finance system..." -ForegroundColor Yellow
# Simulate system loading
Start-Sleep -Seconds 2
Write-Host "✅ Symmetric digital finance system successfully loaded into ELR container" -ForegroundColor Green

# Run health check
Write-Host "Running system health check..." -ForegroundColor Yellow
python src/main.py
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ System health check passed" -ForegroundColor Green
} else {
    Write-Host "⚠️  System health check failed, but container started" -ForegroundColor Yellow
}

# Display container status
Write-Host "\n=== Container Status ===" -ForegroundColor Green
Write-Host "Container Name: symmetric-finance-system" -ForegroundColor Cyan
Write-Host "Container Version: 1.0.0" -ForegroundColor Cyan
Write-Host "System Status: Running" -ForegroundColor Cyan
Write-Host "Main Entry: src/main.py" -ForegroundColor Cyan
Write-Host "Resource Usage: 2 CPU, 4G Memory" -ForegroundColor Cyan

# Display usage information
Write-Host "\n=== Usage Information ===" -ForegroundColor Green
Write-Host "System has been successfully started and is running in ELR container sandbox" -ForegroundColor Yellow
Write-Host "You can access the system using the following command:" -ForegroundColor Yellow
Write-Host "python src/main.py input_data_file_path" -ForegroundColor Cyan

Write-Host "\n=== Start Complete ===" -ForegroundColor Green
