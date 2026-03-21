# ELR Container Installer for Windows
# Similar to uv installation script
# Usage: powershell -ExecutionPolicy ByPass -c "irm https://example.com/elr/install.ps1 | iex"

Write-Host "===================================="
Write-Host "Enlightenment Lighthouse Runtime (ELR)"
Write-Host "One-line Installation Script for Windows"
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
    
    # Download ELR files from GitHub
    Write-Host "Downloading ELR files from GitHub..."
    
    # GitHub repository URL
    $githubRepo = "https://github.com/username/ELR/archive/refs/heads/main.zip"
    
    # Download the zip file
    $zipPath = Join-Path $env:TEMP "elr-main.zip"
    Download-File -url $githubRepo -outputPath $zipPath
    
    # Extract the zip file
    Write-Host "Extracting ELR files..."
    $extractPath = Join-Path $env:TEMP "elr-main"
    Expand-Archive -Path $zipPath -DestinationPath $extractPath -Force
    
    # Copy files from extracted directory
    $extractedELRDir = Get-ChildItem -Path $extractPath -Directory | Select-Object -First 1
    if ($extractedELRDir) {
        Write-Host "Copying ELR files from $($extractedELRDir.FullName) to $installDir"
        
        # Copy main ELR scripts
        $elrPs1 = Join-Path $extractedELRDir.FullName "elr.ps1"
        if (Test-Path $elrPs1) {
            Copy-Item $elrPs1 $binDir -Force
            Write-Host "Copied elr.ps1 to $binDir"
        }
        
        $elrBat = Join-Path $extractedELRDir.FullName "elr.bat"
        if (Test-Path $elrBat) {
            Copy-Item $elrBat $binDir -Force
            Write-Host "Copied elr.bat to $binDir"
        }
        
        # Copy micro_model directory
        $microModelDir = Join-Path $extractedELRDir.FullName "micro_model"
        if (Test-Path $microModelDir) {
            Copy-Item $microModelDir $libDir -Recurse -Force
            Write-Host "Copied micro_model to $libDir"
        }
        
        # Copy models directory
        $modelsSrcDir = Join-Path $extractedELRDir.FullName "models"
        if (Test-Path $modelsSrcDir) {
            Copy-Item $modelsSrcDir $modelsDir -Recurse -Force
            Write-Host "Copied models to $modelsDir"
        }
        
        # Copy api server
        $apiServer = Join-Path $extractedELRDir.FullName "elr_api_server.py"
        if (Test-Path $apiServer) {
            Copy-Item $apiServer $binDir -Force
            Write-Host "Copied elr_api_server.py to $binDir"
        }
    }
    
    # Clean up temporary files
    if (Test-Path $zipPath) {
        Remove-Item $zipPath -Force
    }
    if (Test-Path $extractPath) {
        Remove-Item $extractPath -Recurse -Force
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
To update ELR, simply run the installation command again:
`powershell -ExecutionPolicy ByPass -c "irm https://example.com/elr/install.ps1 | iex"`
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
