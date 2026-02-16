@echo off

REM Word转Markdown转换批处理文件
REM 此文件使用PowerShell命令来执行转换，并输出详细日志

echo ==============================================
echo Word转Markdown转换器
echo ==============================================
echo 输入文件: docs\程序员称谓衍进数智工程师.docx
echo 输出文件: docs\程序员称谓衍进数智工程师.md
echo ==============================================

REM 使用PowerShell执行转换，并输出详细日志
powershell -ExecutionPolicy Bypass -Command "&
    Write-Host '开始执行转换...';
    $inputFile = 'docs\程序员称谓衍进数智工程师.docx';
    $outputFile = 'docs\程序员称谓衍进数智工程师.md';
    
    Write-Host '检查输入文件...';
    if (-not (Test-Path $inputFile)) {
        Write-Host '错误: 输入文件不存在';
        exit 1;
    }
    
    Write-Host '创建临时目录...';
    $tempDir = 'temp_docx';
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null;
    Write-Host '临时目录: ' $tempDir;
    
    Write-Host '复制并解压文件...';
    $zipFile = '$tempDir\temp.zip';
    Copy-Item $inputFile $zipFile -Force;
    Write-Host '复制完成';
    
    try {
        Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force;
        Write-Host '解压完成';
    } catch {
        Write-Host '错误解压: ' $_.Exception.Message;
        exit 1;
    }
    
    Write-Host '检查document.xml...';
    $xmlPath = '$tempDir\word\document.xml';
    if (Test-Path $xmlPath) {
        Write-Host '找到document.xml，开始提取文本...';
        try {
            $xmlContent = Get-Content $xmlPath -Encoding UTF8 -Raw;
            Write-Host '读取XML完成，大小: ' $xmlContent.Length '字符';
            
            Write-Host '提取文本...';
            $plainText = $xmlContent -replace '<[^>]+>', '';
            $plainText = $plainText -replace '\s+', ' ' -replace '^\s+|\s+$', '';
            Write-Host '文本提取完成，长度: ' $plainText.Length '字符';
            
            Write-Host '分割段落...';
            $paragraphs = $plainText -split '\s{2,}' | Where-Object { $_.Trim() -ne '' };
            Write-Host '段落数: ' $paragraphs.Length;
            
            Write-Host '写入Markdown文件...';
            $paragraphs | Set-Content -Path $outputFile -Encoding UTF8;
            Write-Host '写入完成';
            
            Write-Host '检查输出文件...';
            if (Test-Path $outputFile) {
                $fileSize = (Get-Item $outputFile).Length;
                Write-Host '成功: 转换完成!';
                Write-Host '输出文件: ' $outputFile;
                Write-Host '文件大小: ' $fileSize '字节';
            } else {
                Write-Host '错误: 输出文件不存在';
            }
        } catch {
            Write-Host '错误处理XML: ' $_.Exception.Message;
        }
    } else {
        Write-Host '错误: 无法找到document.xml';
        Write-Host '解压后的目录结构:';
        Get-ChildItem -Path $tempDir -Recurse | Select-Object FullName;
    }
    
    Write-Host '清理临时文件...';
    Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue;
    Write-Host '清理完成';
    
    Write-Host '转换任务完成';
"

echo ==============================================
echo 转换任务已执行完成
echo 请查看上面的输出信息以确认转换状态
echo ==============================================

pause
