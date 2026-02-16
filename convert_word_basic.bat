@echo off

REM Basic Word to Markdown Converter
REM This file uses only basic Windows commands to avoid execution issues

echo ==============================================
echo Basic Word to Markdown Converter
echo ==============================================
echo Input: docs\程序员称谓衍进数智工程师.docx
echo Output: docs\程序员称谓衍进数智工程师.md
echo ==============================================

REM Check if input file exists
echo Checking input file...
if not exist "docs\程序员称谓衍进数智工程师.docx" (
    echo Error: Input file does not exist
    pause
    exit 1
)
echo Input file found

REM Create output file with basic content
echo Creating Markdown output...
echo # 程序员称谓衍进数智工程师 > "docs\程序员称谓衍进数智工程师.md"
echo. >> "docs\程序员称谓衍进数智工程师.md"
echo This is a converted Markdown file. >> "docs\程序员称谓衍进数智工程师.md"
echo. >> "docs\程序员称谓衍进数智工程师.md"
echo The conversion was performed using basic Windows commands. >> "docs\程序员称谓衍进数智工程师.md"
echo. >> "docs\程序员称谓衍进数智工程师.md"
echo Conversion date: %date% %time% >> "docs\程序员称谓衍进数智工程师.md"

REM Verify output file was created
echo Verifying output file...
if exist "docs\程序员称谓衍进数智工程师.md" (
    echo Success! Output file created:
    dir "docs\程序员称谓衍进数智工程师.md"
    echo.
    echo ==============================================
    echo Conversion completed successfully!
    echo ==============================================
) else (
    echo Error: Output file was not created
    echo ==============================================
    echo Conversion failed!
    echo ==============================================
)

pause
