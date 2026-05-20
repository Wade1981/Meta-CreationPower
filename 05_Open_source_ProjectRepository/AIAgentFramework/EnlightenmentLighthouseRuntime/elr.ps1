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

# Forward all arguments to elr.exe
& $ElrExe $args
