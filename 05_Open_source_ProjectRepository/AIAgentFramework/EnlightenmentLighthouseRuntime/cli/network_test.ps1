#!/usr/bin/env powershell

# Network configuration test script for ELR container

$containerName1 = "test-network-1-$(Get-Date -Format 'HHmmss')"
$containerName2 = "test-network-2-$(Get-Date -Format 'HHmmss')"

Write-Host "=== ELR Container Network Configuration Test ==="
Write-Host "Testing network configuration for multiple containers..."
Write-Host "=============================================="

# Create first container
Write-Host "\n1. Creating first container..."
$createOutput1 = .\elr-container.exe create --name $containerName1 --image ubuntu:latest
$containerID1 = $createOutput1 | Select-String -Pattern "ID: (elr-\d+)" | ForEach-Object { $_.Matches.Groups[1].Value }

if ([string]::IsNullOrEmpty($containerID1)) {
    Write-Host "Error: Failed to get container ID for first container"
    exit 1
}

Write-Host "First container ID: $containerID1"

# Create second container
Write-Host "\n2. Creating second container..."
$createOutput2 = .\elr-container.exe create --name $containerName2 --image ubuntu:latest
$containerID2 = $createOutput2 | Select-String -Pattern "ID: (elr-\d+)" | ForEach-Object { $_.Matches.Groups[1].Value }

if ([string]::IsNullOrEmpty($containerID2)) {
    Write-Host "Error: Failed to get container ID for second container"
    exit 1
}

Write-Host "Second container ID: $containerID2"

# Check network configuration for first container
Write-Host "\n3. Checking network configuration for first container..."
$containerDir1 = "$env:USERPROFILE\.elr\data\containers\$containerID1"
$networkDir1 = "$containerDir1\network"

if (Test-Path $networkDir1) {
    Write-Host "✓ Network directory exists for first container"
    
    $networkConfig1 = "$networkDir1\config.json"
    if (Test-Path $networkConfig1) {
        Write-Host "✓ Network config file exists for first container"
        $configContent1 = Get-Content $networkConfig1 -Raw
        Write-Host "First container network config: $configContent1"
    } else {
        Write-Host "✗ Network config file missing for first container"
    }
} else {
    Write-Host "✗ Network directory missing for first container"
}

# Check network configuration for second container
Write-Host "\n4. Checking network configuration for second container..."
$containerDir2 = "$env:USERPROFILE\.elr\data\containers\$containerID2"
$networkDir2 = "$containerDir2\network"

if (Test-Path $networkDir2) {
    Write-Host "✓ Network directory exists for second container"
    
    $networkConfig2 = "$networkDir2\config.json"
    if (Test-Path $networkConfig2) {
        Write-Host "✓ Network config file exists for second container"
        $configContent2 = Get-Content $networkConfig2 -Raw
        Write-Host "Second container network config: $configContent2"
    } else {
        Write-Host "✗ Network config file missing for second container"
    }
} else {
    Write-Host "✗ Network directory missing for second container"
}

# Verify network configurations are different
Write-Host "\n5. Verifying network configurations are different..."
if ($configContent1 -ne $configContent2) {
    Write-Host "✓ Network configurations are different for each container"
} else {
    Write-Host "✗ Network configurations are the same for both containers"
}

# Start both containers
Write-Host "\n6. Starting both containers..."
.\elr-container.exe start-container --id $containerID1
.\elr-container.exe start-container --id $containerID2

# Check network service status
Write-Host "\n7. Checking network service status..."
Start-Sleep -Seconds 2
$response = Invoke-RestMethod -Uri "http://localhost:8080/api/network/status" -Method Get
Write-Host "Network service status:"
$response | ConvertTo-Json -Depth 3

# Stop both containers
Write-Host "\n8. Stopping both containers..."
.\elr-container.exe stop-container --id $containerID1
.\elr-container.exe stop-container --id $containerID2

# Delete both containers
Write-Host "\n9. Deleting both containers..."
.\elr-container.exe delete --id $containerID1
.\elr-container.exe delete --id $containerID2

Write-Host "\n=== Network Configuration Test Complete ==="
