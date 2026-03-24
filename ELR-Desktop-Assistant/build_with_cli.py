import os
import subprocess

# 定义路径
python_path = r"C:\Users\Administrator\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts\python.exe"
pypi_installer_path = r"C:\Users\Administrator\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts\pyinstaller.exe"
project_dir = os.path.dirname(os.path.abspath(__file__))
icon_path = os.path.join(project_dir, "icons", "elr_icon.ico")

print(f"Python path: {python_path}")
print(f"PyInstaller path: {pypi_installer_path}")
print(f"Project directory: {project_dir}")
print(f"Icon path: {icon_path}")

# 检查文件是否存在
if not os.path.exists(pypi_installer_path):
    print(f"PyInstaller not found at: {pypi_installer_path}")
    exit(1)

if not os.path.exists(icon_path):
    print(f"Icon file not found at: {icon_path}")
    exit(1)

# 运行PyInstaller命令行
try:
    print("Running PyInstaller with CLI parameters...")
    command = [
        pypi_installer_path,
        "--onefile",
        "--windowed",
        f"--icon={icon_path}",
        "--add-data=icons/elr_icon.png;icons",
        "--add-data=icons/elr_icon.ico;icons",
        "elr_desktop_assistant.py"
    ]
    
    print(f"Command: {' '.join(command)}")
    
    result = subprocess.run(
        command,
        cwd=project_dir,
        capture_output=True,
        text=True,
        shell=True
    )
    print(f"Return code: {result.returncode}")
    print("\nSTDOUT:")
    print(result.stdout)
    if result.stderr:
        print("\nSTDERR:")
        print(result.stderr)
    if result.returncode == 0:
        print("\nBuild completed successfully!")
    else:
        print("\nBuild failed!")
except Exception as e:
    print(f"Error running PyInstaller: {e}")
