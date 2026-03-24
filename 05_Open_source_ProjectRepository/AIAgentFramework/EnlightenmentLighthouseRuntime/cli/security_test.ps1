#!/usr/bin/env powershell

# Security test script for ELR container

$containerName = "test-security-$(Get-Date -Format 'HHmmss')"

Write-Host "=== ELR Container Security Test ==="
Write-Host "Testing container isolation..."
Write-Host "====================================="

# Create container
Write-Host "\n1. Creating container..."
$createOutput = .\elr-container.exe create --name $containerName --image ubuntu:latest
$containerID = $createOutput | Select-String -Pattern "ID: (elr-\d+)" | ForEach-Object { $_.Matches.Groups[1].Value }

if ([string]::IsNullOrEmpty($containerID)) {
    Write-Host "Error: Failed to get container ID"
    exit 1
}

Write-Host "Container ID: $containerID"

# Start container
Write-Host "\n2. Starting container..."
.\elr-container.exe start-container --id $containerID

# Check container directory structure
Write-Host "\n3. Checking container directory structure..."
$containerDir = "$env:USERPROFILE\.elr\data\containers\$containerID"
if (Test-Path $containerDir) {
    Write-Host "✓ Container directory exists: $containerDir"
    
    # Check rootfs directory
    $rootfsDir = "$containerDir\rootfs"
    if (Test-Path $rootfsDir) {
        Write-Host "✓ Rootfs directory exists"
        
        # Check standard directories
        $standardDirs = @("bin", "etc", "home", "lib", "lib64", "proc", "sys", "tmp", "usr", "var")
        foreach ($dir in $standardDirs) {
            $dirPath = "$rootfsDir\$dir"
            if (Test-Path $dirPath) {
                Write-Host "✓ $dir directory exists"
            } else {
                Write-Host "✗ $dir directory missing"
            }
        }
    } else {
        Write-Host "✗ Rootfs directory missing"
    }
} else {
    Write-Host "✗ Container directory missing"
}

# Check network configuration
Write-Host "\n4. Checking network configuration..."
$networkDir = "$containerDir\network"
if (Test-Path $networkDir) {
    Write-Host "✓ Network directory exists"
    
    $networkConfig = "$networkDir\config.json"
    if (Test-Path $networkConfig) {
        Write-Host "✓ Network config file exists"
        $configContent = Get-Content $networkConfig -Raw
        Write-Host "Network config: $configContent"
    } else {
        Write-Host "✗ Network config file missing"
    }
} else {
    Write-Host "✗ Network directory missing"
}

# Stop container
Write-Host "\n5. Stopping container..."
.\elr-container.exe stop-container --id $containerID

# Verify container directory structure exists
Write-Host "\n6. Verifying container directory structure..."
if (Test-Path $containerDir) {
    Write-Host "✓ Container directory exists: $containerDir"
    
    # List contents of container directory
    Write-Host "\nContents of container directory:"
    Get-ChildItem -Path $containerDir -Recurse | ForEach-Object {
        Write-Host "  $($_.FullName)"
    }
} else {
    Write-Host "✗ Container directory missing"
}

# Delete container
Write-Host "\n7. Deleting container..."
.\elr-container.exe delete --id $containerID

# Verify container deletion
Write-Host "\n8. Verifying container deletion..."
if (-not (Test-Path $containerDir)) {
    Write-Host "✓ Container directory deleted"
} else {
    Write-Host "✗ Container directory still exists"
}

Write-Host "\n=== Security Test Complete ==="
