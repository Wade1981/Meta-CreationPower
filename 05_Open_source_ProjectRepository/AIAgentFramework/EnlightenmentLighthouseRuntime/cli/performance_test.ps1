#!/usr/bin/env powershell

# Performance test script for ELR container

$containerName = "test-performance-$(Get-Date -Format 'HHmmss')"
$testCount = 5
$startTimes = @()

Write-Host "=== ELR Container Performance Test ==="
Write-Host "Testing container startup time..."
Write-Host "Test count: $testCount"
Write-Host "====================================="

for ($i = 1; $i -le $testCount; $i++) {
    Write-Host "\nTest $i of $testCount"
    
    # Create container
    $createStart = Get-Date
    $createOutput = .\elr-container.exe create --name $containerName-$i --image ubuntu:latest
    $createEnd = Get-Date
    $createTime = ($createEnd - $createStart).TotalMilliseconds
    Write-Host "Create time: $createTime ms"
    
    # Extract container ID from output
    $containerID = $createOutput | Select-String -Pattern "ID: (elr-\d+)" | ForEach-Object { $_.Matches.Groups[1].Value }
    
    if ([string]::IsNullOrEmpty($containerID)) {
        Write-Host "Error: Failed to get container ID"
        continue
    }
    
    Write-Host "Container ID: $containerID"
    
    # Start container
    $startStart = Get-Date
    .\elr-container.exe start-container --id $containerID
    $startEnd = Get-Date
    $startTime = ($startEnd - $startStart).TotalMilliseconds
    $startTimes += $startTime
    Write-Host "Start time: $startTime ms"
    
    # Stop container
    $stopStart = Get-Date
    .\elr-container.exe stop-container --id $containerID
    $stopEnd = Get-Date
    $stopTime = ($stopEnd - $stopStart).TotalMilliseconds
    Write-Host "Stop time: $stopTime ms"
    
    # Delete container
    $deleteStart = Get-Date
    .\elr-container.exe delete --id $containerID
    $deleteEnd = Get-Date
    $deleteTime = ($deleteEnd - $deleteStart).TotalMilliseconds
    Write-Host "Delete time: $deleteTime ms"
}

# Calculate average startup time
$averageStartTime = ($startTimes | Measure-Object -Average).Average

Write-Host "\n=== Performance Test Results ==="
Write-Host "Average startup time: $averageStartTime ms"
Write-Host "================================"
