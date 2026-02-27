# Test script for ELR model assembler

Write-Host "===================================="
Write-Host "Testing ELR Model Assembler"
Write-Host "===================================="

# Test 1: List available models
Write-Host "\nTest 1: Listing available models"
Write-Host "------------------------------------"
try {
    & .\model-assembler-simple.ps1 list
} catch {
    Write-Host "Error: $_"
}

# Test 2: Assemble a model
Write-Host "\nTest 2: Assembling a business model"
Write-Host "------------------------------------"
try {
    & .\model-assembler-simple.ps1 assemble --type business --container test-container
} catch {
    Write-Host "Error: $_"
}

# Test 3: Run a model
Write-Host "\nTest 3: Running finance model"
Write-Host "------------------------------------"
try {
    & .\model-assembler-simple.ps1 run --model finance-1.0 --input "quarterly financial data" --container test-container
} catch {
    Write-Host "Error: $_"
}

# Test 4: Generate code
Write-Host "\nTest 4: Generating code from model output"
Write-Host "------------------------------------"
try {
    & .\model-assembler-simple.ps1 generate --model codegen-1.0 --input "create a function to calculate factorial" --output test_generated_code
} catch {
    Write-Host "Error: $_"
}

Write-Host "\n===================================="
Write-Host "Test completed!"
Write-Host "===================================="
