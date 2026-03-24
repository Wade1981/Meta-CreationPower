from PIL import Image
import os

# 检查ico文件是否有效
ico_path = os.path.join('icons', 'elr_icon.ico')

print(f"Checking icon file: {ico_path}")

if os.path.exists(ico_path):
    print(f"Icon file exists: {ico_path}")
    try:
        img = Image.open(ico_path)
        print(f"Icon file is valid: {img.format}")
        print(f"Icon size: {img.size}")
    except Exception as e:
        print(f"Error opening icon file: {e}")
else:
    print(f"Icon file does not exist: {ico_path}")

# 检查png文件
png_path = os.path.join('icons', 'elr_icon.png')
if os.path.exists(png_path):
    print(f"PNG file exists: {png_path}")
    try:
        img = Image.open(png_path)
        print(f"PNG file is valid: {img.format}")
        print(f"PNG size: {img.size}")
    except Exception as e:
        print(f"Error opening PNG file: {e}")
else:
    print(f"PNG file does not exist: {png_path}")
