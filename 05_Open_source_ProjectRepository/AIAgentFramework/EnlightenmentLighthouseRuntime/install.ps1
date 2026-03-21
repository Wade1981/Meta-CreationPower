# ELR Container Installer for Windows
# Similar to uv installation script

Write-Host "===================================="
Write-Host "Enlightenment Lighthouse Runtime (ELR)"
Write-Host "Installation Script for Windows"
Write-Host "===================================="

# Default installation directory
$DEFAULT_INSTALL_DIR = "$env:USERPROFILE\ELR"

# Function to prompt for installation directory
function Get-InstallationDirectory {
    $installDir = Read-Host "Enter installation directory (default: $DEFAULT_INSTALL_DIR)"
    if ([string]::IsNullOrEmpty($installDir)) {
        return $DEFAULT_INSTALL_DIR
    }
    return $installDir
}

# Function to create directory if it doesn't exist
function Ensure-DirectoryExists {
    param([string]$directory)
    if (-not (Test-Path $directory)) {
        Write-Host "Creating directory: $directory"
        New-Item -ItemType Directory -Path $directory -Force | Out-Null
    }
}

# Function to download files
function Download-File {
    param(
        [string]$url,
        [string]$outputPath
    )
    Write-Host "Downloading: $url"
    try {
        Invoke-WebRequest -Uri $url -OutFile $outputPath -ErrorAction Stop
        Write-Host "Download completed successfully"
    } catch {
        Write-Host "Error downloading file: $_"
        exit 1
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
    # Get installation directory
    $installDir = Get-InstallationDirectory
    
    # Create installation directory
    Ensure-DirectoryExists $installDir
    
    # Create subdirectories
    $binDir = Join-Path $installDir "bin"
    $libDir = Join-Path $installDir "lib"
    $configDir = Join-Path $installDir "config"
    $modelsDir = Join-Path $installDir "models"
    $containersDir = Join-Path $installDir "containers"
    
    Ensure-DirectoryExists $binDir
    Ensure-DirectoryExists $libDir
    Ensure-DirectoryExists $configDir
    Ensure-DirectoryExists $modelsDir
    Ensure-DirectoryExists $containersDir
    
    # Copy ELR files from current location
    $currentDir = Get-Location
    Write-Host "Copying ELR files from $currentDir to $installDir"
    
    # Copy main ELR scripts
    if (Test-Path "$currentDir\elr.ps1") {
        Copy-Item "$currentDir\elr.ps1" $binDir -Force
        Write-Host "Copied elr.ps1 to $binDir"
    }
    
    if (Test-Path "$currentDir\elr.bat") {
        Copy-Item "$currentDir\elr.bat" $binDir -Force
        Write-Host "Copied elr.bat to $binDir"
    }
    
    # Copy micro_model directory
    if (Test-Path "$currentDir\micro_model") {
        Copy-Item "$currentDir\micro_model" $libDir -Recurse -Force
        Write-Host "Copied micro_model to $libDir"
    }
    
    # Copy models directory
    if (Test-Path "$currentDir\models") {
        Copy-Item "$currentDir\models" $modelsDir -Recurse -Force
        Write-Host "Copied models to $modelsDir"
    }
    
    # Copy api server
    if (Test-Path "$currentDir\elr_api_server.py") {
        Copy-Item "$currentDir\elr_api_server.py" $binDir -Force
        Write-Host "Copied elr_api_server.py to $binDir"
    }
    
    # Create a wrapper script to ensure ELR runs from the correct directory
    $wrapperContent = @"
@echo off
set "ELR_HOME=%~dp0.."
cd /d "%ELR_HOME%"
powershell -ExecutionPolicy Bypass -File "%ELR_HOME%\bin\elr.ps1" %*
"@
    
    $wrapperPath = Join-Path $binDir "elr.cmd"
    $wrapperContent | Set-Content $wrapperPath -Force
    Write-Host "Created wrapper script: $wrapperPath"
    
    # Add bin directory to PATH
    Add-ToPath $binDir
    
    # Create a README file
    $readmeContent = @"
# Enlightenment Lighthouse Runtime (ELR)

## Installation
ELR has been successfully installed to: $installDir

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
To update ELR, simply run this installation script again with the latest version.
"@
    
    $readmePath = Join-Path $installDir "README.md"
    $readmeContent | Set-Content $readmePath -Force
    Write-Host "Created README.md at $readmePath"
    
    # Test the installation
    Write-Host "===================================="
    Write-Host "Testing ELR installation..."
    Write-Host "===================================="
    
    # Change to installation directory and test
    Push-Location $installDir
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
    Write-Host "Installation directory: $installDir"
    Write-Host "Binary directory: $binDir"
    Write-Host ""
    Write-Host "You can now use 'elr' command from anywhere in your terminal."
    Write-Host "Example: elr start"
    Write-Host ""
    Write-Host "For more information, see the README.md file in the installation directory."
    
} catch {
    Write-Host "Error during installation: $_"
    exit 1
}
