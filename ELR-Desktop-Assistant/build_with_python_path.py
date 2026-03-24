import os
import subprocess

# 定义路径
python_path = r"C:\Users\Administrator\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts\python.exe"
pypi_installer_path = r"C:\Users\Administrator\AppData\Local\Packages\PythonSoftwareFoundation.Python.3.13_qbz5n2kfra8p0\LocalCache\local-packages\Python313\Scripts\pyinstaller.exe"
project_dir = os.path.dirname(os.path.abspath(__file__))
spec_file = os.path.join(project_dir, "ELRDesktopAssistant.spec")

print(f"Python path: {python_path}")
print(f"PyInstaller path: {pypi_installer_path}")
print(f"Project directory: {project_dir}")
print(f"Spec file: {spec_file}")

# 检查文件是否存在
if not os.path.exists(pypi_installer_path):
    print(f"PyInstaller not found at: {pypi_installer_path}")
    exit(1)

if not os.path.exists(spec_file):
    print(f"Spec file not found at: {spec_file}")
    exit(1)

# 运行PyInstaller
try:
    print("Running PyInstaller...")
    result = subprocess.run(
        [pypi_installer_path, spec_file],
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
