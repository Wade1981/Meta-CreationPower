@echo off

set "PYTHON_EXE=e:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\python-portable\python.exe"
set "SCRIPT_PATH=e:\X54\github\Meta-CreationPower\06_EnterprisePrivateProjects\EnlightenmentLighthouseOriginTeam\CommunicationChannel\tools\elr_taskbar_service.py"
set "DESKTOP_PATH=%USERPROFILE%\Desktop"
set "SHORTCUT_NAME=ELR任务栏服务.lnk"

echo 创建ELR任务栏服务桌面快捷方式...

:: 创建VBS脚本以创建快捷方式
set "VBS_SCRIPT=%TEMP%\CreateShortcut.vbs"
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%VBS_SCRIPT%"
echo sLinkFile = "%DESKTOP_PATH%\%SHORTCUT_NAME%" >> "%VBS_SCRIPT%"
echo Set oLink = oWS.CreateShortcut(sLinkFile) >> "%VBS_SCRIPT%"
echo oLink.TargetPath = "%PYTHON_EXE%" >> "%VBS_SCRIPT%"
echo oLink.Arguments = """%SCRIPT_PATH%""" >> "%VBS_SCRIPT%"
echo oLink.Description = "ELR任务栏服务" >> "%VBS_SCRIPT%"
echo oLink.IconLocation = "%PYTHON_EXE%, 0" >> "%VBS_SCRIPT%"
echo oLink.Save >> "%VBS_SCRIPT%"

:: 运行VBS脚本
cscript //nologo "%VBS_SCRIPT%"

:: 清理临时文件
del "%VBS_SCRIPT%"

echo 快捷方式创建成功！
echo.
echo 如何使用：
echo 1. 双击桌面图标 "%SHORTCUT_NAME" 启动ELR任务栏服务

echo 2. 服务启动后，会在系统托盘中显示一个灯塔图标

echo 3. 右键点击系统托盘中的图标，打开上下文菜单

echo 4. 从菜单中选择相应的操作，如启动/停止ELR运行时、管理容器等

echo.
echo 如何将图标固定到任务栏：
echo 1. 双击桌面图标 "%SHORTCUT_NAME" 启动ELR任务栏服务

echo 2. 在任务栏中找到ELR任务栏服务的图标

echo 3. 右键点击该图标，选择 "固定到任务栏"

echo.
echo 安装完成！
pause
