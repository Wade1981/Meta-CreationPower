# Test Check-Status function

# Container statuses
$CONTAINER_STATUS_RUNNING = "running"

# Global variables
$global:RUNTIME_STARTED = $true
$global:RUNTIME_START_TIME = Get-Date
$global:CONTAINERS = @(
    @{ Name = "container1"; Status = "running" },
    @{ Name = "container2"; Status = "created" }
)

# Function: Check runtime status
function Check-Status {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    Write-Host "Enlightenment Lighthouse Runtime is RUNNING"
    Write-Host "Started: $($global:RUNTIME_START_TIME.ToString('yyyy-MM-dd HH:mm:ss'))"
    Write-Host "Containers: $($global:CONTAINERS.Count)"
    # 调试输出：显示每个容器的状态
    foreach ($container in $global:CONTAINERS) {
        Write-Host "  Container: $($container.Name), Status: $($container.Status)"
    }
    # 使用简单的foreach循环计算运行中的容器数量
    $runningCount = 0
    foreach ($container in $global:CONTAINERS) {
        if ($container.Status -eq $CONTAINER_STATUS_RUNNING) {
            $runningCount++
        }
    }
    Write-Host "Running containers: $runningCount"
}

# Call the function
Check-Status
