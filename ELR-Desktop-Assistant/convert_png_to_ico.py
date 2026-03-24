from PIL import Image
import os

# 定义路径
png_path = os.path.join('icons', 'elr_icon.png')
ico_path = os.path.join('icons', 'elr_icon.ico')

print(f"Converting {png_path} to {ico_path}")

# 打开PNG文件
img = Image.open(png_path)

# 保存为ICO文件，包含多个尺寸
img.save(ico_path, format='ICO', sizes=[(16,16), (32,32), (48,48), (256,256)])

print(f"Conversion completed successfully: {ico_path}")
