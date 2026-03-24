# Step 1: Test basic function with foreach loop

function Test-Foreach {
    $containers = @(
        @{ Name = "container1"; Status = "running" },
        @{ Name = "container2"; Status = "created" }
    )
    
    Write-Host "Containers: $($containers.Count)"
    foreach ($container in $containers) {
        Write-Host "  Container: $($container.Name), Status: $($container.Status)"
    }
    
    $runningContainers = $containers | Where-Object { $_.Status -eq "running" }
    Write-Host "Running containers: $($runningContainers.Count)"
}

Test-Foreach
