# Minimal test script

function Test-Function {
    Write-Host "Hello World"
    $testArray = @(1, 2, 3)
    Write-Host "Array count: $($testArray.Count)"
    foreach ($item in $testArray) {
        Write-Host "Item: $item"
    }
}

Test-Function
