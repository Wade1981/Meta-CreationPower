@echo off

REM 构建ELR-Desktop-Assistant
REM 使用Python的绝对路径来确保正确的环境

set PYTHON_PATH=C:\Users\Administrator\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts

REM 尝试使用脚本路径中的pyinstaller
if exist "%PYTHON_PATH%\pyinstaller.exe" (
    echo 使用脚本路径中的pyinstaller
    "%PYTHON_PATH%\pyinstaller.exe" ELRDesktopAssistant.spec
) else (
    echo 使用Python模块方式
    python -m pyinstaller ELRDesktopAssistant.spec
)

pause
