# Minimal ELR script to test container count

# State file path
$STATE_FILE = "elr-state.json"

# Load state
if (Test-Path $STATE_FILE) {
    try {
        $state = Get-Content $STATE_FILE | ConvertFrom-Json
        $containers = $state.CONTAINERS
        Write-Host "Total containers: $($containers.Count)"
        
        # Count running containers
        $runningCount = 0
        foreach ($container in $containers) {
            Write-Host "Container: $($container.Name), Status: $($container.Status)"
            if ($container.Status -eq "running") {
                $runningCount++
            }
        }
        Write-Host "Running containers: $runningCount"
    } catch {
        Write-Host "Error loading state: $_"
    }
} else {
    Write-Host "State file not found"
}
