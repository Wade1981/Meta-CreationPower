# 极简版Word转Markdown转换器
param(
    [string]$InputFile = "docs\程序员称谓衍进数智工程师.docx",
    [string]$OutputFile = "docs\程序员称谓衍进数智工程师.md"
)

Write-Host "开始转换: $InputFile -> $OutputFile"

# 检查输入文件
if (-not (Test-Path $InputFile)) {
    Write-Host "错误: 输入文件不存在"
    exit 1
}

# 创建临时目录
$tempDir = "temp_docx"
New-Item -ItemType Directory -Path $tempDir -Force | Out-Null

# 复制并解压
$zipFile = "$tempDir\temp.zip"
Copy-Item $InputFile $zipFile -Force
Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force

# 读取document.xml
$xmlPath = "$tempDir\word\document.xml"
if (Test-Path $xmlPath) {
    Write-Host "找到document.xml，开始提取文本..."
    $xmlContent = Get-Content $xmlPath -Encoding UTF8 -Raw
    
    # 简单提取文本（移除XML标签）
    $plainText = $xmlContent -replace '<[^>]+>', ''
    $plainText = $plainText -replace '\s+', ' ' -replace '^\s+|\s+$', ''
    
    # 分割成段落
    $paragraphs = $plainText -split '\s{2,}' | Where-Object { $_.Trim() -ne '' }
    
    # 写入Markdown
    $paragraphs | Set-Content -Path $OutputFile -Encoding UTF8
    
    Write-Host "转换完成! 输出文件: $OutputFile"
    Write-Host "行数: $($paragraphs.Length)"
} else {
    Write-Host "错误: 无法找到document.xml"
}

# 清理
Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue

Write-Host "任务完成"
