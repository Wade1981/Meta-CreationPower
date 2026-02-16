@echo off

REM Basic ELR Test
REM This file tests basic ELR functionality

echo ==============================================
echo Basic ELR Test
echo ==============================================
echo Date: %date%
echo Time: %time%
echo ==============================================

REM Test 1: Check if ELR script exists
echo Test 1: Checking ELR script...
if exist "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr-with-python.ps1" (
    echo ✓ ELR script found
) else (
    echo ✗ ELR script not found
)

REM Test 2: Check if Python converter exists
echo Test 2: Checking Python converter...
if exist "true_docx_converter.py" (
    echo ✓ Python converter found
) else (
    echo ✗ Python converter not found
)

REM Test 3: Check if input file exists
echo Test 3: Checking input file...
if exist "docs\程序员称谓衍进数智工程师.docx" (
    echo ✓ Input file found
) else (
    echo ✗ Input file not found
)

REM Test 4: Try to run converter directly
echo Test 4: Running converter directly...
python true_docx_converter.py "docs\程序员称谓衍进数智工程师.docx" "docs\程序员称谓衍进数智工程师.md"

REM Test 5: Check if output was created
echo Test 5: Checking output file...
if exist "docs\程序员称谓衍进数智工程师.md" (
    echo ✓ Output file created
    dir "docs\程序员称谓衍进数智工程师.md"
) else (
    echo ✗ Output file not created
)

echo ==============================================
echo Basic ELR Test completed
echo ==============================================

pause
