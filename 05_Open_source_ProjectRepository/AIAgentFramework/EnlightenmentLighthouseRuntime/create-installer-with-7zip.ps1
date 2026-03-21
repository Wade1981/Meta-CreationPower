# Create ELR Binary Installer using 7-Zip
# This script creates a self-extracting binary installer for ELR

Write-Host "===================================="
Write-Host "Creating ELR Binary Installer with 7-Zip"
Write-Host "===================================="

# Set variables
$installerName = "elr-installer.exe"
$outputDir = "output"
$tempDir = "temp-installer"
$currentDir = Get-Location

# Create output directory if it doesn't exist
if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

# Create temporary directory for installer files
if (Test-Path $tempDir) {
    Remove-Item -Path $tempDir -Recurse -Force
}
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

# Copy ELR files to temporary directory
Write-Host "Copying ELR files..."

# Create directory structure
$binDir = Join-Path $tempDir "bin"
$libDir = Join-Path $tempDir "lib"
$configDir = Join-Path $tempDir "config"
$modelsDir = Join-Path $tempDir "models"
$containersDir = Join-Path $tempDir "containers"

New-Item -ItemType Directory -Path $binDir -Force | Out-Null
New-Item -ItemType Directory -Path $libDir -Force | Out-Null
New-Item -ItemType Directory -Path $configDir -Force | Out-Null
New-Item -ItemType Directory -Path $modelsDir -Force | Out-Null
New-Item -ItemType Directory -Path $containersDir -Force | Out-Null

# Copy main ELR scripts
if (Test-Path "$currentDir\elr.ps1") {
    Copy-Item "$currentDir\elr.ps1" $binDir -Force
    Write-Host "Copied elr.ps1"
}

if (Test-Path "$currentDir\elr.bat") {
    Copy-Item "$currentDir\elr.bat" $binDir -Force
    Write-Host "Copied elr.bat"
}

# Copy micro_model directory
if (Test-Path "$currentDir\micro_model") {
    Copy-Item "$currentDir\micro_model" $libDir -Recurse -Force
    Write-Host "Copied micro_model"
}

# Copy models directory
if (Test-Path "$currentDir\models") {
    Copy-Item "$currentDir\models" $modelsDir -Recurse -Force
    Write-Host "Copied models"
}

# Copy api server
if (Test-Path "$currentDir\elr_api_server.py") {
    Copy-Item "$currentDir\elr_api_server.py" $binDir -Force
    Write-Host "Copied elr_api_server.py"
}

# Create a wrapper script to ensure ELR runs from the correct directory
$wrapperContent = '@echo off
set "ELR_HOME=%~dp0.."
cd /d "%ELR_HOME%"
powershell -ExecutionPolicy Bypass -File "%ELR_HOME%\bin\elr.ps1" %*'

$wrapperPath = Join-Path $binDir "elr.cmd"
$wrapperContent | Set-Content $wrapperPath -Force
Write-Host "Created wrapper script"

# Create installation script
$installScriptContent = '# ELR Container Installer
# Self-extracting binary installer

param(
    [string]$InstallDir = "$env:USERPROFILE\ELR"
)

Write-Host "===================================="
Write-Host "Enlightenment Lighthouse Runtime (ELR)"
Write-Host "Binary Installer"
Write-Host "===================================="

# Function to create directory if it doesn't exist
function Ensure-DirectoryExists {
    param([string]$directory)
    if (-not (Test-Path $directory)) {
        Write-Host "Creating directory: $directory"
        New-Item -ItemType Directory -Path $directory -Force | Out-Null
    }
}

# Function to add directory to PATH
function Add-ToPath {
    param([string]$directory)
    $path = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($path -notlike "*$directory*") {
        Write-Host "Adding $directory to PATH"
        $newPath = "$path;$directory"
        [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")
        Write-Host "PATH updated. You may need to restart your terminal for changes to take effect."
    } else {
        Write-Host "$directory is already in PATH"
    }
}

# Main installation process
try {
    # Get installation directory from user if not provided
    if (-not $InstallDir) {
        $InstallDir = Read-Host "Enter installation directory (default: $env:USERPROFILE\ELR)"
        if ([string]::IsNullOrEmpty($InstallDir)) {
            $InstallDir = "$env:USERPROFILE\ELR"
        }
    }
    
    # Create installation directory
    Ensure-DirectoryExists $InstallDir
    
    # Create subdirectories
    $binDir = Join-Path $InstallDir "bin"
    $libDir = Join-Path $InstallDir "lib"
    $configDir = Join-Path $InstallDir "config"
    $modelsDir = Join-Path $InstallDir "models"
    $containersDir = Join-Path $InstallDir "containers"
    
    Ensure-DirectoryExists $binDir
    Ensure-DirectoryExists $libDir
    Ensure-DirectoryExists $configDir
    Ensure-DirectoryExists $modelsDir
    Ensure-DirectoryExists $containersDir
    
    # Copy files from temporary directory
    $tempDir = Split-Path -Parent $MyInvocation.MyCommand.Path
    Write-Host "Copying ELR files from temporary directory to $InstallDir"
    
    # Copy main ELR scripts
    Copy-Item "$tempDir\bin\*" $binDir -Force
    Write-Host "Copied bin files"
    
    # Copy micro_model directory
    Copy-Item "$tempDir\lib\*" $libDir -Recurse -Force
    Write-Host "Copied lib files"
    
    # Copy models directory
    Copy-Item "$tempDir\models\*" $modelsDir -Recurse -Force
    Write-Host "Copied models"
    
    # Add bin directory to PATH
    Add-ToPath $binDir
    
    # Create a README file
    $readmeContent = "# Enlightenment Lighthouse Runtime (ELR)

## Installation
ELR has been successfully installed to: $InstallDir

## Usage

### Basic Commands
- `elr start` - Start the ELR runtime
- `elr stop` - Stop the ELR runtime
- `elr status` - Check runtime status
- `elr list` - List all containers
- `elr help` - Show help information

### Advanced Commands
- `elr create --name <name> --image <image>` - Create a new container
- `elr run --name <name> --image <image>` - Create and start a new container
- `elr start-container --id <container-id>` - Start a container
- `elr stop-container --id <container-id>` - Stop a container
- `elr delete --id <container-id>` - Delete a container
- `elr exec --id <container-id> --command <command>` - Execute a command in a container

### Model Commands
- `elr run-python --source <script.py>` - Run a Python script
- `elr run-python --code '<python code>'` - Run Python code directly
- `elr chat` - Start interactive chat with default local model
- `elr chat --model <model.py>` - Start chat with custom model

## Configuration
Configuration files are stored in: $configDir

## Models
Model files are stored in: $modelsDir

## Containers
Container data is stored in: $containersDir

## Troubleshooting
- If you encounter issues with Python, ensure Python 3.8+ is installed and in PATH
- If you encounter issues with C compilation, ensure GCC is installed
- For network issues, check your firewall settings

## Updates
To update ELR, simply run the installer again with the latest version.'
    
    $readmePath = Join-Path $InstallDir "README.md"
    $readmeContent | Set-Content $readmePath -Force
    Write-Host "Created README.md"
    
    # Test the installation
    Write-Host "===================================="
    Write-Host "Testing ELR installation..."
    Write-Host "===================================="
    
    # Change to installation directory and test
    Push-Location $InstallDir
    try {
        # Test ELR version
        Write-Host "Testing ELR version..."
        & powershell -ExecutionPolicy Bypass -File "$binDir\elr.ps1" version
        
        # Test ELR help
        Write-Host "\nTesting ELR help..."
        & powershell -ExecutionPolicy Bypass -File "$binDir\elr.ps1" help
    } finally {
        Pop-Location
    }
    
    Write-Host "===================================="
    Write-Host "ELR installation completed successfully!"
    Write-Host "===================================="
    Write-Host "Installation directory: $InstallDir"
    Write-Host "Binary directory: $binDir"
    Write-Host ""
    Write-Host "You can now use 'elr' command from anywhere in your terminal."
    Write-Host "Example: elr start"
    Write-Host ""
    Write-Host "For more information, see the README.md file in the installation directory."
    
    # Pause to allow user to see the output
    Read-Host "Press Enter to exit..."
    
} catch {
    Write-Host "Error during installation: $_"
    Read-Host "Press Enter to exit..."
    exit 1
}"

$installScriptPath = Join-Path $tempDir "install.ps1"
$installScriptContent | Set-Content $installScriptPath -Force
Write-Host "Created installation script"

# Create a batch file to run the installation
$batchContent = '@echo off
powershell -ExecutionPolicy Bypass -File "%~dp0install.ps1" %*'

$batchPath = Join-Path $tempDir "install.bat"
$batchContent | Set-Content $batchPath -Force
Write-Host "Created batch file"

# Check if 7-Zip is available
Write-Host "Checking for 7-Zip..."
$7zipPath = Get-Command 7z -ErrorAction SilentlyContinue
if ($7zipPath) {
    Write-Host "Found 7-Zip at: $($7zipPath.Source)"
    
    # Create self-extracting archive using 7-Zip
    Write-Host "Creating self-extracting installer with 7-Zip...