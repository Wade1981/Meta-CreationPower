#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
生成 ELR 图标
"""

from PIL import Image, ImageDraw
import os

def create_elr_icon():
    # 确保 icons 目录存在
    if not os.path.exists('icons'):
        os.makedirs('icons')
    
    # 创建一个 64x64 的图像
    width, height = 64, 64
    image = Image.new('RGB', (width, height), color='#1a1a2e')
    draw = ImageDraw.Draw(image)
    
    # 绘制灯塔主体
    # 灯塔底座
    draw.rectangle([(24, 50), (40, 64)], fill='#e6e6fa')
    # 灯塔主体
    draw.rectangle([(28, 20), (36, 50)], fill='#e6e6fa')
    # 灯塔顶部
    draw.polygon([(24, 20), (32, 10), (40, 20)], fill='#e6e6fa')
    # 灯塔灯光
    draw.ellipse([(30, 22), (34, 26)], fill='#00ffff')
    
    # 绘制放射状光线
    for i in range(8):
        angle = i * 45
        if angle % 90 != 0:  # 只绘制对角线方向的光线
            draw.line([(32, 32), (width, height)], fill='#00ffff', width=1)
        image = image.rotate(45)
    
    # 保存图标
    icon_path = os.path.join('icons', 'elr_icon.png')
    image.save(icon_path)
    print(f"图标已保存到: {icon_path}")

if __name__ == "__main__":
    create_elr_icon()
