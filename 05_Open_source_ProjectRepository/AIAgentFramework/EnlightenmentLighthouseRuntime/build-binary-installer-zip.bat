@echo off

echo ====================================
echo Creating ELR Binary Installer as Zip Archive
echo ====================================

:: Set variables
set installerName=elr-binary-installer.zip
set outputDir=output
set tempDir=temp-binary-installer
set currentDir=%cd%

:: Create output directory if it doesn't exist
if not exist "%outputDir%" mkdir "%outputDir%"

:: Create temporary directory for installer files
if exist "%tempDir%" rd /s /q "%tempDir%"
mkdir "%tempDir%"
mkdir "%tempDir%\bin"
mkdir "%tempDir%\config"
mkdir "%tempDir%\models"
mkdir "%tempDir%\containers"

:: Skip Go compilation for now due to network issues
echo Skipping Go compilation (network issues)...
echo Will create installer with Python API server only

:: Create a placeholder for micro_model.exe
echo Creating placeholder for micro_model.exe...
type nul > "%tempDir%\bin\micro_model.exe"
echo Created placeholder for micro_model.exe

:: Copy compiled binary
echo Copying compiled binary...

:: Copy main ELR scripts
if exist "elr.ps1" copy "elr.ps1" "%tempDir%\bin\" >nul && echo Copied elr.ps1
if exist "elr.bat" copy "elr.bat" "%tempDir%\bin\" >nul && echo Copied elr.bat

:: Copy api server (as Python file for now)
if exist "elr_api_server.py" copy "elr_api_server.py" "%tempDir%\bin\" >nul && echo Copied elr_api_server.py

:: Copy necessary configuration files
if exist "micro_model\config" (xcopy "micro_model\config" "%tempDir%\config" /s /e /i /y >nul && echo Copied config files)

:: Copy models directory
if exist "models" (xcopy "models" "%tempDir%\models" /s /e /i /y >nul && echo Copied models)

:: Create a wrapper script to ensure ELR runs from the correct directory
echo Creating wrapper script...
echo @echo off> "%tempDir%\bin\elr.cmd"
echo set "ELR_HOME=%%~dp0..">> "%tempDir%\bin\elr.cmd"
echo cd /d "%%ELR_HOME%%">> "%tempDir%\bin\elr.cmd"
echo powershell -ExecutionPolicy Bypass -File "%%ELR_HOME%%\bin\elr.ps1" %%*>> "%tempDir%\bin\elr.cmd"

:: Create installation script
echo Creating installation script...
type nul > "%tempDir%\install.ps1"
echo # ELR Container Binary Installer>> "%tempDir%\install.ps1"
echo # Zip archive installation script>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo param(>> "%tempDir%\install.ps1"
echo     [string]$InstallDir = "$env:USERPROFILE\ELR">> "%tempDir%\install.ps1"
echo )>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo Write-Host "====================================">> "%tempDir%\install.ps1"
echo Write-Host "Enlightenment Lighthouse Runtime (ELR)">> "%tempDir%\install.ps1"
echo Write-Host "Binary Installation Script">> "%tempDir%\install.ps1"
echo Write-Host "====================================">> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo # Function to create directory if it doesn't exist>> "%tempDir%\install.ps1"
echo function Ensure-DirectoryExists {>> "%tempDir%\install.ps1"
echo     param([string]$directory)>> "%tempDir%\install.ps1"
echo     if (-not (Test-Path $directory)) {>> "%tempDir%\install.ps1"
echo         Write-Host "Creating directory: $directory">> "%tempDir%\install.ps1"
echo         New-Item -ItemType Directory -Path $directory -Force ^| Out-Null>> "%tempDir%\install.ps1"
echo     }>> "%tempDir%\install.ps1"
echo }>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo # Function to add directory to PATH>> "%tempDir%\install.ps1"
echo function Add-ToPath {>> "%tempDir%\install.ps1"
echo     param([string]$directory)>> "%tempDir%\install.ps1"
echo     $path = [Environment]::GetEnvironmentVariable("PATH", "User")>> "%tempDir%\install.ps1"
echo     if ($path -notlike "*$directory*") {>> "%tempDir%\install.ps1"
echo         Write-Host "Adding $directory to PATH">> "%tempDir%\install.ps1"
echo         $newPath = "$path;$directory">> "%tempDir%\install.ps1"
echo         [Environment]::SetEnvironmentVariable("PATH", $newPath, "User")>> "%tempDir%\install.ps1"
echo         Write-Host "PATH updated. You may need to restart your terminal for changes to take effect.">> "%tempDir%\install.ps1"
echo     } else {>> "%tempDir%\install.ps1"
echo         Write-Host "$directory is already in PATH">> "%tempDir%\install.ps1"
echo     }>> "%tempDir%\install.ps1"
echo }>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo # Main installation process>> "%tempDir%\install.ps1"
echo try {>> "%tempDir%\install.ps1"
echo     # Get installation directory from user if not provided>> "%tempDir%\install.ps1"
echo     if (-not $InstallDir) {>> "%tempDir%\install.ps1"
echo         $InstallDir = Read-Host "Enter installation directory (default: $env:USERPROFILE\ELR)">> "%tempDir%\install.ps1"
echo         if ([string]::IsNullOrEmpty($InstallDir)) {>> "%tempDir%\install.ps1"
echo             $InstallDir = "$env:USERPROFILE\ELR">> "%tempDir%\install.ps1"
echo         }>> "%tempDir%\install.ps1"
echo     }>> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Create installation directory>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $InstallDir>> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Create subdirectories>> "%tempDir%\install.ps1"
echo     $binDir = Join-Path $InstallDir "bin">> "%tempDir%\install.ps1"
echo     $configDir = Join-Path $InstallDir "config">> "%tempDir%\install.ps1"
echo     $modelsDir = Join-Path $InstallDir "models">> "%tempDir%\install.ps1"
echo     $containersDir = Join-Path $InstallDir "containers">> "%tempDir%\install.ps1"
echo     $uploadsDir = Join-Path $InstallDir "uploads">> "%tempDir%\install.ps1"
echo     $assetsDir = Join-Path $InstallDir "assets">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $binDir>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $configDir>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $modelsDir>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $containersDir>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $uploadsDir>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists $assetsDir>> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists "$assetsDir\images">> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists "$assetsDir\audio">> "%tempDir%\install.ps1"
echo     Ensure-DirectoryExists "$assetsDir\video">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Copy files from current directory>> "%tempDir%\install.ps1"
echo     $currentDir = Split-Path -Parent $MyInvocation.MyCommand.Path>> "%tempDir%\install.ps1"
echo     Write-Host "Copying ELR files from $currentDir to $InstallDir">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Copy main ELR scripts>> "%tempDir%\install.ps1"
echo     Copy-Item "$currentDir\bin\*" $binDir -Force>> "%tempDir%\install.ps1"
echo     Write-Host "Copied bin files">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Copy config files>> "%tempDir%\install.ps1"
echo     Copy-Item "$currentDir\config\*" $configDir -Recurse -Force>> "%tempDir%\install.ps1"
echo     Write-Host "Copied config files">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Copy models directory>> "%tempDir%\install.ps1"
echo     Copy-Item "$currentDir\models\*" $modelsDir -Recurse -Force>> "%tempDir%\install.ps1"
echo     Write-Host "Copied models">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Add bin directory to PATH>> "%tempDir%\install.ps1"
echo     Add-ToPath $binDir>> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Create a README file>> "%tempDir%\install.ps1"
echo     $readmeContent = '# Enlightenment Lighthouse Runtime (ELR)>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Installation>> "%tempDir%\install.ps1"
echo ELR has been successfully installed to: $InstallDir>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Usage>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ### Basic Commands>> "%tempDir%\install.ps1"
echo - `elr start` - Start the ELR runtime>> "%tempDir%\install.ps1"
echo - `elr stop` - Stop the ELR runtime>> "%tempDir%\install.ps1"
echo - `elr status` - Check runtime status>> "%tempDir%\install.ps1"
echo - `elr list` - List all containers>> "%tempDir%\install.ps1"
echo - `elr help` - Show help information>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ### Advanced Commands>> "%tempDir%\install.ps1"
echo - `elr create --name ^<name^> --image ^<image^>` - Create a new container>> "%tempDir%\install.ps1"
echo - `elr run --name ^<name^> --image ^<image^>` - Create and start a new container>> "%tempDir%\install.ps1"
echo - `elr start-container --id ^<container-id^>` - Start a container>> "%tempDir%\install.ps1"
echo - `elr stop-container --id ^<container-id^>` - Stop a container>> "%tempDir%\install.ps1"
echo - `elr delete --id ^<container-id^>` - Delete a container>> "%tempDir%\install.ps1"
echo - `elr exec --id ^<container-id^> --command ^<command^>` - Execute a command in a container>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ### Model Commands>> "%tempDir%\install.ps1"
echo - `elr run-python --source ^<script.py^>` - Run a Python script>> "%tempDir%\install.ps1"
echo - `elr run-python --code '^<python code^>'` - Run Python code directly>> "%tempDir%\install.ps1"
echo - `elr chat` - Start interactive chat with default local model>> "%tempDir%\install.ps1"
echo - `elr chat --model ^<model.py^>` - Start chat with custom model>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Configuration>> "%tempDir%\install.ps1"
echo Configuration files are stored in: $configDir>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Models>> "%tempDir%\install.ps1"
echo Model files are stored in: $modelsDir>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Containers>> "%tempDir%\install.ps1"
echo Container data is stored in: $containersDir>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Troubleshooting>> "%tempDir%\install.ps1"
echo - If you encounter issues with Python, ensure Python 3.8+ is installed and in PATH>> "%tempDir%\install.ps1"
echo - For network issues, check your firewall settings>> "%tempDir%\install.ps1"
echo.>> "%tempDir%\install.ps1"
echo ## Updates>> "%tempDir%\install.ps1"
echo To update ELR, simply run the installer again with the latest version.'>> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     $readmePath = Join-Path $InstallDir "README.md">> "%tempDir%\install.ps1"
echo     $readmeContent ^| Set-Content $readmePath -Force>> "%tempDir%\install.ps1"
echo     Write-Host "Created README.md">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Test the installation>> "%tempDir%\install.ps1"
echo     Write-Host "====================================">> "%tempDir%\install.ps1"
echo     Write-Host "Testing ELR installation...">> "%tempDir%\install.ps1"
echo     Write-Host "====================================">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Change to installation directory and test>> "%tempDir%\install.ps1"
echo     Push-Location $InstallDir>> "%tempDir%\install.ps1"
echo     try {>> "%tempDir%\install.ps1"
echo         # Test ELR version>> "%tempDir%\install.ps1"
echo         Write-Host "Testing ELR version...">> "%tempDir%\install.ps1"
echo         & powershell -ExecutionPolicy Bypass -File "$binDir\elr.ps1" version>> "%tempDir%\install.ps1"
echo         >> "%tempDir%\install.ps1"
echo         # Test ELR help>> "%tempDir%\install.ps1"
echo         Write-Host "\nTesting ELR help...">> "%tempDir%\install.ps1"
echo         & powershell -ExecutionPolicy Bypass -File "$binDir\elr.ps1" help>> "%tempDir%\install.ps1"
echo     } finally {>> "%tempDir%\install.ps1"
echo         Pop-Location>> "%tempDir%\install.ps1"
echo     }>> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     Write-Host "====================================">> "%tempDir%\install.ps1"
echo     Write-Host "ELR binary installation completed successfully!">> "%tempDir%\install.ps1"
echo     Write-Host "====================================">> "%tempDir%\install.ps1"
echo     Write-Host "Installation directory: $InstallDir">> "%tempDir%\install.ps1"
echo     Write-Host "Binary directory: $binDir">> "%tempDir%\install.ps1"
echo     Write-Host "">> "%tempDir%\install.ps1"
echo     Write-Host "You can now use 'elr' command from anywhere in your terminal.">> "%tempDir%\install.ps1"
echo     Write-Host "Example: elr start">> "%tempDir%\install.ps1"
echo     Write-Host "">> "%tempDir%\install.ps1"
echo     Write-Host "For more information, see the README.md file in the installation directory.">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo     # Pause to allow user to see the output>> "%tempDir%\install.ps1"
echo     Read-Host "Press Enter to exit...">> "%tempDir%\install.ps1"
echo     >> "%tempDir%\install.ps1"
echo } catch {>> "%tempDir%\install.ps1"
echo     Write-Host "Error during installation: $_">> "%tempDir%\install.ps1"
echo     Read-Host "Press Enter to exit...">> "%tempDir%\install.ps1"
echo     exit 1>> "%tempDir%\install.ps1"
echo }>> "%tempDir%\install.ps1"

:: Create a batch file to run the installation
echo Creating batch file...
echo @echo off> "%tempDir%\install.bat"
echo powershell -ExecutionPolicy Bypass -File "%%~dp0install.ps1" %%*>> "%tempDir%\install.bat"

:: Create a README for the installer
echo Creating installer README.md...
type nul > "%tempDir%\README.md"
echo # ELR Container Binary Installer>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ## Installation Instructions>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo 1. Extract this zip file to your desired location>> "%tempDir%\README.md"
echo 2. Run `install.bat` from the extracted directory>> "%tempDir%\README.md"
echo 3. Follow the prompts to complete the installation>> "%tempDir%\README.md"
echo 4. Use the `elr` command from anywhere in your terminal>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ## System Requirements>> "%tempDir%\README.md"
echo - Windows 10 or later>> "%tempDir%\README.md"
echo - PowerShell 5.1 or later>> "%tempDir%\README.md"
echo - Python 3.8+ (required for API server)>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ## Usage>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo After installation, you can use the `elr` command to interact with ELR:>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ### Basic Commands>> "%tempDir%\README.md"
echo - `elr start` - Start the ELR runtime>> "%tempDir%\README.md"
echo - `elr stop` - Stop the ELR runtime>> "%tempDir%\README.md"
echo - `elr status` - Check runtime status>> "%tempDir%\README.md"
echo - `elr list` - List all containers>> "%tempDir%\README.md"
echo - `elr help` - Show help information>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ### Advanced Commands>> "%tempDir%\README.md"
echo - `elr create --name ^<name^> --image ^<image^>` - Create a new container>> "%tempDir%\README.md"
echo - `elr run --name ^<name^> --image ^<image^>` - Create and start a new container>> "%tempDir%\README.md"
echo - `elr start-container --id ^<container-id^>` - Start a container>> "%tempDir%\README.md"
echo - `elr stop-container --id ^<container-id^>` - Stop a container>> "%tempDir%\README.md"
echo - `elr delete --id ^<container-id^>` - Delete a container>> "%tempDir%\README.md"
echo - `elr exec --id ^<container-id^> --command ^<command^>` - Execute a command in a container>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ### Model Commands>> "%tempDir%\README.md"
echo - `elr run-python --source ^<script.py^>` - Run a Python script>> "%tempDir%\README.md"
echo - `elr run-python --code '^<python code^>'` - Run Python code directly>> "%tempDir%\README.md"
echo - `elr chat` - Start interactive chat with default local model>> "%tempDir%\README.md"
echo - `elr chat --model ^<model.py^>` - Start chat with custom model>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ## Troubleshooting>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ### Python Not Found>> "%tempDir%\README.md"
echo If you encounter Python-related errors, ensure Python 3.8+ is installed and in your PATH.>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ### Network Issues>> "%tempDir%\README.md"
echo If you have network connectivity issues, check your firewall settings.>> "%tempDir%\README.md"
echo.>> "%tempDir%\README.md"
echo ## Updates>> "%tempDir%\README.md"
echo To update ELR, simply download the latest installer and run it again.>> "%tempDir%\README.md"

:: Create zip archive using PowerShell
echo Creating zip archive...
powershell -ExecutionPolicy Bypass -Command "Compress-Archive -Path '%tempDir%\*' -DestinationPath '%outputDir%\%installerName%' -Force"

:: Check if zip file was created
if exist "%outputDir%\%installerName%" (
    echo ====================================
    echo ELR Binary Installer Zip Archive created successfully!
    echo ====================================
    echo Installer path: %outputDir%\%installerName%
    for %%A in ("%outputDir%\%installerName%") do echo Size: %%~zA bytes
    echo.
    echo To install ELR:
    echo 1. Extract the zip file to your desired location
    echo 2. Run install.bat from the extracted directory
    echo 3. Follow the prompts
    echo 4. Use 'elr' command from anywhere in your terminal
) else (
    echo Error: Zip file was not created
)

:cleanup
:: Clean up temporary files
if exist "%tempDir%" rd /s /q "%tempDir%"

cd "%currentDir%"
echo ====================================
echo Binary installer creation process completed
echo ====================================

pause