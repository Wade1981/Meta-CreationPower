#!/usr/bin/env pwsh
# ELR Desktop API启动脚本

param(
    [string]$IP = "localhost",
    [int]$Port = 8081
)

Write-Host "====================================="
Write-Host "ELR Desktop API 启动脚本"
Write-Host "====================================="

$pythonPath = $null

$portablePython = "$PSScriptRoot\..\python-portable\python.exe"
if (Test-Path $portablePython) {
    $pythonPath = $portablePython
    Write-Host "Found portable Python: $pythonPath"
} else {
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
                    $pythonPath = $testPath.Source
                    break
                }
            }
        } catch {
        }
    }
    
    if ($pythonPath) {
        Write-Host "Found Python: $pythonPath"
    }
}

if (-not $pythonPath) {
    Write-Host "Error: Python interpreter not found"
    Write-Host "Please install Python 3.8 or higher"
    exit 1
}

$desktopApiPath = "$PSScriptRoot\desktop_api.py"
if (-not (Test-Path $desktopApiPath)) {
    Write-Host "Error: desktop_api.py not found"
    Write-Host "Expected path: $desktopApiPath"
    exit 1
}

Write-Host "Found Desktop API: $desktopApiPath"

$uploadDir = "$PSScriptRoot\uploads"
$modelDir = "$PSScriptRoot\models"
$componentDir = "$PSScriptRoot\components"
$assetDir = "$PSScriptRoot\assets"

if (-not (Test-Path $uploadDir)) {
    New-Item -ItemType Directory -Path $uploadDir -Force | Out-Null
    Write-Host "Created upload directory: $uploadDir"
}

if (-not (Test-Path $modelDir)) {
    New-Item -ItemType Directory -Path $modelDir -Force | Out-Null
    Write-Host "Created model directory: $modelDir"
}

if (-not (Test-Path $componentDir)) {
    New-Item -ItemType Directory -Path $componentDir -Force | Out-Null
    Write-Host "Created component directory: $componentDir"
}

if (-not (Test-Path $assetDir)) {
    New-Item -ItemType Directory -Path $assetDir -Force | Out-Null
    New-Item -ItemType Directory -Path "$assetDir\images" -Force | Out-Null
    New-Item -ItemType Directory -Path "$assetDir\audio" -Force | Out-Null
    New-Item -ItemType Directory -Path "$assetDir\video" -Force | Out-Null
    Write-Host "Created asset directory: $assetDir"
}

Write-Host ""
Write-Host "====================================="
Write-Host "Starting ELR Desktop API Server..."
Write-Host "====================================="
Write-Host "Server address: http://${IP}:${Port}"
Write-Host ""
Write-Host "Available endpoints:"
Write-Host "  GET  /api/desktop/status      - Get ELR status"
Write-Host "  GET  /api/desktop/containers  - Get container list"
Write-Host "  POST /api/desktop/upload      - Upload file"
Write-Host "  GET  /api/desktop/files       - List uploaded files"
Write-Host "  DELETE /api/desktop/files/{name} - Delete file"
Write-Host "  GET  /api/desktop/resources   - Get system resources"
Write-Host "  GET  /api/desktop/health      - Health check"
Write-Host ""
Write-Host "Press Ctrl+C to stop"
Write-Host "====================================="

try {
    & $pythonPath $desktopApiPath
} catch {
    Write-Host "Error: $($_.Exception.Message)"
    exit 1
}
