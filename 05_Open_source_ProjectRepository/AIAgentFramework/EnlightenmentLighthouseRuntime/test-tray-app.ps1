#!/usr/bin/env powershell
# Test script for ELR tray app

Write-Host "Testing ELR tray app functionality..."
Write-Host "=========================================="

# Test 1: Check if ELR script exists
$ELRPS1 = "$PSScriptRoot\elr.ps1"
if (Test-Path $ELRPS1) {
    Write-Host "✓ ELR script found: $ELRPS1"
} else {
    Write-Host "✗ ELR script not found"
    exit 1
}

# Test 2: Test exec command
Write-Host "\nTesting exec command..."
try {
    $result = & $ELRPS1 exec --command "echo 'Hello from ELR container'" 2>&1
    Write-Host "✓ Exec command executed successfully"
    Write-Host "Output: $result"
} catch {
    Write-Host "✗ Exec command failed: $($_.Exception.Message)"
}

# Test 3: Test upload command
Write-Host "\nTesting upload command..."
try {
    # Create a test file
    $testFile = "$PSScriptRoot\test-upload.txt"
    "This is a test file for upload" | Set-Content $testFile
    
    $result = & $ELRPS1 upload --file "$testFile" 2>&1
    Write-Host "✓ Upload command executed successfully"
    Write-Host "Output: $result"
    
    # Clean up test file
    Remove-Item $testFile -Force -ErrorAction SilentlyContinue
} catch {
    Write-Host "✗ Upload command failed: $($_.Exception.Message)"
}

# Test 4: Test status command
Write-Host "\nTesting status command..."
try {
    $status = & $ELRPS1 status 2>&1
    Write-Host "✓ Status command executed successfully"
    Write-Host "Status: $status"
} catch {
    Write-Host "✗ Status command failed: $($_.Exception.Message)"
}

# Test 5: Test list command
Write-Host "\nTesting list command..."
try {
    $containers = & $ELRPS1 list 2>&1
    Write-Host "✓ List command executed successfully"
    Write-Host "Containers: $containers"
} catch {
    Write-Host "✗ List command failed: $($_.Exception.Message)"
}

Write-Host "\n=========================================="
Write-Host "Test completed."
