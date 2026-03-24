@echo off

REM 使用命令提示符构建ELR-Desktop-Assistant
REM 切换到命令提示符模式运行

cmd /c "python -m pyinstaller --name ELRDesktopAssistant --windowed --icon icons/elr_icon.png --onefile elr_desktop_assistant.py"

pause
