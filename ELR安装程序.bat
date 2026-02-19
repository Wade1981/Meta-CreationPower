@echo off

:: 设置标题
title ELR任务栏服务安装程序

:: 清屏
cls

:: 显示欢迎信息
echo ================================================================================
echo                          ELR任务栏服务安装程序
echo ================================================================================
echo
echo 欢迎使用ELR任务栏服务安装程序！
echo
echo 本安装程序将：
echo 1. 在桌面上创建ELR任务栏服务快捷方式
echo 2. 在开始菜单中创建ELR任务栏服务快捷方式
echo 3. 提供将ELR图标固定到任务栏的选项
echo 4. 安装必要的依赖项（如果尚未安装）
echo
echo 按任意键继续...
pause >nul

:: 清屏
cls

:: 检查Python环境
echo ================================================================================
echo                          检查系统环境
echo ================================================================================
echo
echo 检查Python环境...

set "PYTHON_EXE=e:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\python-portable\python.exe"
set "SCRIPT_PATH=e:\X54\github\Meta-CreationPower\06_EnterprisePrivateProjects\EnlightenmentLighthouseOriginTeam\CommunicationChannel\tools\elr_taskbar_service.py"

if exist "%PYTHON_EXE" (
    echo [✓] Python环境已找到
) else (
    echo [✗] Python环境未找到
    echo 请确保便携式Python已正确安装
    echo 按任意键退出...
    pause >nul
    exit 1
)

if exist "%SCRIPT_PATH" (
    echo [✓] ELR任务栏服务脚本已找到
) else (
    echo [✗] ELR任务栏服务脚本未找到
    echo 请确保脚本文件已正确放置
    echo 按任意键退出...
    pause >nul
    exit 1
)

echo
echo 系统环境检查完成！
echo
echo 按任意键继续...
pause >nul

:: 清屏
cls

:: 安装依赖项
echo ================================================================================
echo                          安装依赖项
echo ================================================================================
echo
echo 检查并安装必要的依赖项...

:: 检查pystray库
echo 检查pystray库...
%PYTHON_EXE% -c "import pystray" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] pystray库已安装
) else (
    echo [✗] pystray库未安装，正在安装...
    %PYTHON_EXE% -m pip install pystray >nul 2>&1
    if %errorlevel% equ 0 (
        echo [✓] pystray库安装成功
    ) else (
        echo [✗] pystray库安装失败
        echo 按任意键退出...
        pause >nul
        exit 1
    )
)

:: 检查Pillow库
echo 检查Pillow库...
%PYTHON_EXE% -c "import PIL" >nul 2>&1
if %errorlevel% equ 0 (
    echo [✓] Pillow库已安装
) else (
    echo [✗] Pillow库未安装，正在安装...
    %PYTHON_EXE% -m pip install Pillow >nul 2>&1
    if %errorlevel% equ 0 (
        echo [✓] Pillow库安装成功
    ) else (
        echo [✗] Pillow库安装失败
        echo 按任意键退出...
        pause >nul
        exit 1
    )
)

echo
echo 依赖项安装完成！
echo
echo 按任意键继续...
pause >nul

:: 清屏
cls

:: 创建快捷方式
echo ================================================================================
echo                          创建快捷方式
echo ================================================================================
echo
echo 创建桌面快捷方式...

set "DESKTOP_PATH=%USERPROFILE%\Desktop"
set "SHORTCUT_NAME=ELR任务栏服务.lnk"
set "START_MENU_PATH=%APPDATA%\Microsoft\Windows\Start Menu\Programs"

:: 创建VBS脚本以创建桌面快捷方式
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

echo [✓] 桌面快捷方式创建成功！
echo
echo 创建开始菜单快捷方式...

:: 确保开始菜单目录存在
if not exist "%START_MENU_PATH%" (
    mkdir "%START_MENU_PATH%"
)

:: 创建VBS脚本以创建开始菜单快捷方式
set "VBS_SCRIPT=%TEMP%\CreateStartMenuShortcut.vbs"
echo Set oWS = WScript.CreateObject("WScript.Shell") > "%VBS_SCRIPT%"
echo sLinkFile = "%START_MENU_PATH%\%SHORTCUT_NAME%" >> "%VBS_SCRIPT%"
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
echo [✓] 开始菜单快捷方式创建成功！
echo
echo 按任意键继续...
pause >nul

:: 清屏
cls

:: 固定到任务栏选项
echo ================================================================================
echo                          固定到任务栏
echo ================================================================================
echo
echo 是否将ELR图标固定到任务栏？
echo
echo 1. 是
2. 否
echo
echo 请输入选项编号（默认：2）：
set /p "CHOICE="
if "%CHOICE%"=="" set "CHOICE=2"

if "%CHOICE%"=="1" (
    echo
echo 正在将ELR图标固定到任务栏...
    
    :: 创建VBS脚本以固定到任务栏
    set "VBS_SCRIPT=%TEMP%\PinToTaskbar.vbs"
echo ' 固定应用程序到任务栏的VBS脚本 > "%VBS_SCRIPT%"
echo Set objShell = CreateObject("Shell.Application") >> "%VBS_SCRIPT%"
echo Set objFolder = objShell.Namespace("%DESKTOP_PATH%") >> "%VBS_SCRIPT%"
echo Set objFolderItem = objFolder.ParseName("%SHORTCUT_NAME%") >> "%VBS_SCRIPT%"
echo Set colVerbs = objFolderItem.Verbs >> "%VBS_SCRIPT%"
echo For Each objVerb in colVerbs >> "%VBS_SCRIPT%"
echo     If Replace(objVerb.name, "&", "") = "固定到任务栏" Then >> "%VBS_SCRIPT%"
echo         objVerb.DoIt >> "%VBS_SCRIPT%"
echo         Exit For >> "%VBS_SCRIPT%"
echo     End If >> "%VBS_SCRIPT%"
echo Next >> "%VBS_SCRIPT%"

    :: 运行VBS脚本
    cscript //nologo "%VBS_SCRIPT%"

    :: 清理临时文件
    del "%VBS_SCRIPT%"
    
    echo [✓] ELR图标已固定到任务栏！
) else (
    echo 跳过固定到任务栏步骤
)
echo
echo 按任意键继续...
pause >nul

:: 清屏
cls

:: 显示安装完成信息
echo ================================================================================
echo                          安装完成
echo ================================================================================
echo
echo [✓] ELR任务栏服务安装完成！
echo
echo 安装结果：
echo --------------------------------------------------------------------------------
echo [✓] 桌面快捷方式：已创建
echo [✓] 开始菜单快捷方式：已创建
echo [✓] 依赖项：已检查并安装
echo %IF_PIN% 任务栏固定：%PIN_STATUS%
echo --------------------------------------------------------------------------------
echo
echo 使用说明：
echo 1. 双击桌面图标 "%SHORTCUT_NAME" 启动ELR任务栏服务
echo 2. 服务启动后，会在系统托盘中显示一个灯塔图标
echo 3. 右键点击系统托盘中的图标，打开上下文菜单
echo 4. 从菜单中选择相应的操作，如启动/停止ELR运行时、管理容器等
echo
echo 按任意键退出安装程序...
pause >nul

:: 退出
exit 0
