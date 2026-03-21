# 测试ELR容器适配性改造
Write-Host "Testing ELR container adaptation changes..."

# 检查新增的文件是否存在
Write-Host "Checking for new files..."

$modelPropertiesFile = "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\model\model_properties.go"
$modelAdapterFile = "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\model\model_adapter.go"
$modelPropertiesJson = "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\examples\model_properties.json"
$fishSpeechProperties = "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\model\models\fish-speech\model_properties.json"

if (Test-Path $modelPropertiesFile) {
    Write-Host "✓ model_properties.go exists"
} else {
    Write-Host "✗ model_properties.go missing" -ForegroundColor Red
}

if (Test-Path $modelAdapterFile) {
    Write-Host "✓ model_adapter.go exists"
} else {
    Write-Host "✗ model_adapter.go missing" -ForegroundColor Red
}

if (Test-Path $modelPropertiesJson) {
    Write-Host "✓ model_properties.json exists"
} else {
    Write-Host "✗ model_properties.json missing" -ForegroundColor Red
}

if (Test-Path $fishSpeechProperties) {
    Write-Host "✓ fish-speech model_properties.json exists"
} else {
    Write-Host "✗ fish-speech model_properties.json missing" -ForegroundColor Red
}

# 检查修改的文件
Write-Host "`nChecking for modified files..."

$modelFile = "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\micro_model\model\model.go"

if (Test-Path $modelFile) {
    Write-Host "✓ model.go exists"
    # 检查文件是否包含新的代码
    $content = Get-Content $modelFile -Raw
    if ($content -match "ModelProperties") {
        Write-Host "✓ model.go contains ModelProperties"
    } else {
        Write-Host "✗ model.go missing ModelProperties" -ForegroundColor Red
    }
    if ($content -match "GetModelAdapter") {
        Write-Host "✓ model.go contains GetModelAdapter"
    } else {
        Write-Host "✗ model.go missing GetModelAdapter" -ForegroundColor Red
    }
} else {
    Write-Host "✗ model.go missing" -ForegroundColor Red
}

Write-Host "`nELR container adaptation test completed!"
