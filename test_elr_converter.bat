@echo off

REM ELR Word to Markdown Converter Test
REM This batch file tests if ELR can execute the Python converter

echo ==============================================
echo ELR Word to Markdown Converter Test
echo ==============================================
echo Test Date: %date% %time%
echo ==============================================

REM Start ELR runtime
echo Starting ELR runtime...
powershell -ExecutionPolicy Bypass -File "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr-with-python.ps1" start
if errorlevel 1 (
    echo Error: Failed to start ELR runtime
    pause
    exit 1
)
echo ELR runtime started

REM Test Python execution
echo ==============================================
echo Testing Python execution with ELR...
echo ==============================================

REM Run the converter using ELR
powershell -ExecutionPolicy Bypass -File "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr-with-python.ps1" run-python --source "true_docx_converter.py" "docs\程序员称谓衍进数智工程师.docx" "docs\程序员称谓衍进数智工程师.md"

REM Check if conversion was successful
echo ==============================================
echo Checking conversion result...
echo ==============================================
if exist "docs\程序员称谓衍进数智工程师.md" (
    echo Success! Markdown file created:
    dir "docs\程序员称谓衍进数智工程师.md"
    echo.
    echo ==============================================
    echo ELR Python execution test PASSED!
    echo ==============================================
) else (
    echo Error: Markdown file was not created
    echo ==============================================
    echo ELR Python execution test FAILED!
    echo ==============================================
)

REM Stop ELR runtime
echo Stopping ELR runtime...
powershell -ExecutionPolicy Bypass -File "05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr-with-python.ps1" stop
echo ELR runtime stopped

echo ==============================================
echo Test completed
echo ==============================================

pause
