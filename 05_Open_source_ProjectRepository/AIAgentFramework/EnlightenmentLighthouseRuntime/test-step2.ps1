# Step 2: Test with simpler running count calculation

function Test-RunningCount {
    $containers = @(
        @{ Name = "container1"; Status = "running" },
        @{ Name = "container2"; Status = "created" }
    )
    
    Write-Host "Containers: $($containers.Count)"
    foreach ($container in $containers) {
        Write-Host "  Container: $($container.Name), Status: $($container.Status)"
    }
    
    # Use a simpler approach to count running containers
    $runningCount = 0
    foreach ($container in $containers) {
        if ($container.Status -eq "running") {
            $runningCount++
        }
    }
    Write-Host "Running containers: $runningCount"
}

Test-RunningCount
