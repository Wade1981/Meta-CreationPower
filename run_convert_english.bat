@echo off

REM Word to Markdown Converter Batch File
REM This file uses PowerShell commands to perform conversion with detailed logging

echo ==============================================
echo Word to Markdown Converter
echo ==============================================
echo Input: docs\programmer_evolution.docx
echo Output: docs\programmer_evolution.md
echo ==============================================

REM First, copy the Chinese-named file to English name to avoid encoding issues
echo Copying file to English name...
copy "docs\程序员称谓衍进数智工程师.docx" "docs\programmer_evolution.docx" /Y
if errorlevel 1 (
    echo Error copying file
    pause
    exit 1
)
echo File copied successfully

REM Use PowerShell to execute conversion with detailed logging
powershell -ExecutionPolicy Bypass -Command "&
    Write-Host 'Starting conversion...';
    $inputFile = 'docs\programmer_evolution.docx';
    $outputFile = 'docs\programmer_evolution.md';
    
    Write-Host 'Checking input file...';
    if (-not (Test-Path $inputFile)) {
        Write-Host 'Error: Input file does not exist';
        exit 1;
    }
    
    Write-Host 'Creating temporary directory...';
    $tempDir = 'temp_docx';
    if (Test-Path $tempDir) {
        Remove-Item $tempDir -Recurse -Force;
    }
    New-Item -ItemType Directory -Path $tempDir -Force | Out-Null;
    Write-Host 'Temporary directory: ' $tempDir;
    
    Write-Host 'Copying and extracting file...';
    $zipFile = $tempDir + '\temp.zip';
    Copy-Item $inputFile $zipFile -Force;
    Write-Host 'Copy completed';
    
    try {
        Expand-Archive -Path $zipFile -DestinationPath $tempDir -Force;
        Write-Host 'Extraction completed';
    } catch {
        Write-Host 'Extraction error: ' $_.Exception.Message;
        exit 1;
    }
    
    Write-Host 'Checking document.xml...';
    $xmlPath = $tempDir + '\word\document.xml';
    if (Test-Path $xmlPath) {
        Write-Host 'Found document.xml, starting text extraction...';
        try {
            $xmlContent = Get-Content $xmlPath -Encoding UTF8 -Raw;
            Write-Host 'XML read completed, size: ' $xmlContent.Length 'characters';
            
            Write-Host 'Extracting text...';
            $plainText = $xmlContent -replace '<[^>]+>', '';
            $plainText = $plainText -replace '\s+', ' ' -replace '^\s+|\s+$', '';
            Write-Host 'Text extraction completed, length: ' $plainText.Length 'characters';
            
            Write-Host 'Splitting into paragraphs...';
            $paragraphs = $plainText -split '\s{2,}' | Where-Object { $_.Trim() -ne '' };
            Write-Host 'Paragraphs count: ' $paragraphs.Length;
            
            Write-Host 'Writing Markdown file...';
            $paragraphs | Set-Content -Path $outputFile -Encoding UTF8;
            Write-Host 'Writing completed';
            
            Write-Host 'Checking output file...';
            if (Test-Path $outputFile) {
                $fileSize = (Get-Item $outputFile).Length;
                Write-Host 'Success: Conversion completed!';
                Write-Host 'Output file: ' $outputFile;
                Write-Host 'File size: ' $fileSize 'bytes';
            } else {
                Write-Host 'Error: Output file does not exist';
            }
        } catch {
            Write-Host 'XML processing error: ' $_.Exception.Message;
        }
    } else {
        Write-Host 'Error: document.xml not found';
        Write-Host 'Extracted directory structure:';
        Get-ChildItem -Path $tempDir -Recurse | Select-Object FullName;
    }
    
    Write-Host 'Cleaning temporary files...';
    Remove-Item $tempDir -Recurse -Force -ErrorAction SilentlyContinue;
    Write-Host 'Cleanup completed';
    
    Write-Host 'Conversion task completed';
"

echo ==============================================
echo Conversion task executed
echo Please check the output above for status
echo ==============================================

REM Verify the output file was created
echo Verifying output file...
if exist "docs\programmer_evolution.md" (
    echo Success! Output file created:
    dir "docs\programmer_evolution.md"
) else (
    echo Error: Output file was not created
)

pause
