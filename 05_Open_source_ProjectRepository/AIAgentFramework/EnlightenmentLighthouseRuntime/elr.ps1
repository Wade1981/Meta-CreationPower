# Enlightenment Lighthouse Runtime (ELR)
# PowerShell wrapper for elr.exe

# Get the directory of the current script
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path

# Path to elr.exe (using relative path)
$ElrExe = Join-Path -Path $ScriptDir -ChildPath "elr.exe"

# Check if elr.exe exists
if (-not (Test-Path $ElrExe)) {
    Write-Host "Error: elr.exe not found. Please compile the Go code first."
    Write-Host "Run: go build -o elr.exe cli/main.go"
    exit 1
}

# Main function
if ($args.Length -eq 0) {
    # No arguments, show help
    & $ElrExe help
    exit 1
}

# Check if the command is 'api start'
if ($args[0] -eq "api" -and $args[1] -eq "start") {
    # Start elr.exe and show output
    Write-Host "Starting API service..."
    # Run elr.exe in the foreground to show output
    & $ElrExe $args
} elseif ($args[0] -eq "api" -and $args[1] -eq "stop") {
    # Stop API service
    Write-Host "Stopping API service..."
    & $ElrExe $args
} elseif ($args[0] -eq "api" -and $args[1] -eq "status") {
    # Check API service status
    Write-Host "Checking API service status..."
    & $ElrExe $args
} elseif ($args[0] -eq "api" -and $args[1] -eq "config") {
    # Configure API settings
    Write-Host "Configuring API settings..."
    & $ElrExe $args
} elseif ($args[0] -eq "run" -and $args[1] -eq "python") {
    # Run Python script
    Write-Host "Running Python script..."
    & $ElrExe $args
} elseif ($args[0] -eq "install" -and $args[1] -eq "python") {
    # Install Python
    Write-Host "Installing Python..."
    & $ElrExe $args
} elseif ($args[0] -eq "gui" -or $args[0] -eq "tray") {
    # Start ELR GUI (Tray Application)
    Write-Host "Starting ELR GUI..."
    $TrayAppPath = Join-Path -Path $ScriptDir -ChildPath "ELR-Tray-App.ps1"
    if (Test-Path $TrayAppPath) {
        # Run the tray app in a background job
        # This method ensures the tray app continues to run even after the PowerShell window is closed
        Start-Job -ScriptBlock {
            param($path)
            & $path
        } -ArgumentList $TrayAppPath
        
        # Wait a few seconds to ensure the GUI is fully started
        Start-Sleep -Seconds 3
        
        Write-Host "ELR GUI started successfully!"
        Write-Host "Check your system tray for the ELR icon."
    } else {
        Write-Host "Error: ELR-Tray-App.ps1 not found."
        exit 1
    }
} elseif ($args[0] -eq "create" -and $args.Length -gt 1) {
    # Handle create command with container name as positional argument
    Write-Host "Creating container..."
    $containerName = $args[1]
    # Convert positional argument to --name parameter
    & $ElrExe "create" "--name" $containerName "--image" "ubuntu:latest"
} elseif ($args[0] -eq "sandbox" -and $args[1] -eq "run-model") {
    # Run model in sandbox
    Write-Host "Running model in sandbox..."
    & $ElrExe $args
} elseif ($args[0] -eq "status") {
    # Check system status
    Write-Host "Checking ELR system status..."
    & $ElrExe $args
} elseif ($args[0] -eq "interact") {
    # Interact with model
    Write-Host "Connecting to model..."
    & $ElrExe $args
} elseif ($args[0] -eq "stop-model") {
    # Stop model
    Write-Host "Stopping model..."
    & $ElrExe $args
} else {
    # Forward all arguments to elr.exe for other commands
    & $ElrExe $args
}

# Available commands:
# 
# 1. File System Management
# .\elr.ps1 fs upload      - Upload file to container
# .\elr.ps1 fs download    - Download file from container
# .\elr.ps1 fs set-dir     - Set directory for file type
# .\elr.ps1 fs get-dir     - Get directory for file type
# 
# 2. Model Management
# .\elr.ps1 model list     - List all models
# .\elr.ps1 model get      - Get model information
# .\elr.ps1 model download - Download a model
# .\elr.ps1 model delete   - Delete a model
# .\elr.ps1 model install-deps - Install model dependencies
# 
# 3. Sandbox Management
# .\elr.ps1 sandbox list   - List all sandboxes
# .\elr.ps1 sandbox create - Create a new sandbox
# .\elr.ps1 sandbox start  - Start a sandbox
# .\elr.ps1 sandbox stop   - Stop a sandbox
# .\elr.ps1 sandbox delete - Delete a sandbox
# .\elr.ps1 sandbox load-model - Load model into sandbox
# .\elr.ps1 sandbox unload-model - Unload model from sandbox
# .\elr.ps1 sandbox run-model - Run model in sandbox
# 
# 4. Resource Configuration
# .\elr.ps1 Settings list - List all resource configurations
# .\elr.ps1 Settings --resource-type <type> --directory <path> - Set resource type directory
# .\elr.ps1 Settings --model-type <type> --directory <path> - Set model type directory
# 
# 5. Upload Commands
# .\elr.ps1 Upload Settings type <resource-type> path: <file-path> - Upload resource
# 
# 6. Installation Commands
# .\elr.ps1 install python [version] [path] - Install Python
# 
# 7. Status Check Commands
# .\elr.ps1 status          - Check ELR system status
# .\elr.ps1 status containers - Check container status
# .\elr.ps1 status sandboxes  - Check sandbox status
# .\elr.ps1 status models     - Check model status
# .\elr.ps1 status api        - Check API service status
# 
# 8. Model Interaction
# .\elr.ps1 interact --sandbox-id <sandbox-id> --model-id <model-id> - Interact with a running model
# 
# 9. Model Management
# .\elr.ps1 stop-model --model-id <model-id> [--sandbox-id <sandbox-id>] - Stop a running model
# 
# 10. API Service Commands
# .\elr.ps1 api start      - Start API services
# .\elr.ps1 api stop       - Stop API services
# .\elr.ps1 api status     - Check API status
# .\elr.ps1 api config     - Configure API settings
# 
# 11. Container Management
# .\elr.ps1 create         - Create a new container
# .\elr.ps1 run            - Create and start a new container
# .\elr.ps1 start-container - Start a container
# .\elr.ps1 stop-container  - Stop a container
# .\elr.ps1 list           - List all containers
# .\elr.ps1 delete         - Delete a container
# .\elr.ps1 inspect        - Inspect a container
# 
# 12. System Commands
# .\elr.ps1 help           - Show help
# .\elr.ps1 version        - Print version information
# .\elr.ps1 start          - Start the ELR runtime
# .\elr.ps1 stop           - Stop the ELR runtime
# .\elr.ps1 setup          - Setup ELR system (e.g., isolation)
# .\elr.ps1 admin          - Administrator management commands
# .\elr.ps1 gui            - Start ELR GUI (Tray Application)
# .\elr.ps1 tray           - Start ELR GUI (Tray Application)
#
# 13. Resource Configuration Commands
# .\elr.ps1 Settings list     - List all resource configurations
# .\elr.ps1 Settings --resource-type <type> --directory <path> - Set resource type directory
# .\elr.ps1 Settings --model-type <type> --directory <path> - Set model type directory
#
# 14. Upload Commands
# .\elr.ps1 Upload Settings type <resource-type> path: <file-path> - Upload resource
#
# 15. Installation Commands
# .\elr.ps1 install python [version] [path] - Install Python
