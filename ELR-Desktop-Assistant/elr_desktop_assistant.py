#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
ELR Desktop Assistant
基于 Python 开发的 ELR（Enlightenment Lighthouse Runtime）桌面助手
不依赖 Qt，使用 Python 标准库和轻量级依赖
"""

import tkinter as tk
from tkinter import ttk
from tkinter import filedialog
import threading
import time
import requests
import json
import sys
from PIL import Image, ImageTk
import pystray
from pystray import MenuItem as item

class ELRClient:
    """ELR 客户端类，用于与 ELR API 通信"""
    
    def __init__(self, api_base_url="http://localhost:8080/api"):
        self.api_base_url = api_base_url
        self.status = "未连接"
        self.containers = []
        self.callbacks = []
    
    def set_api_url(self, api_url):
        """设置 API 地址"""
        self.api_base_url = api_url
        self.get_elr_status()
        self.get_elr_containers()
        self.notify_callbacks()
    
    def get_elr_status(self):
        """获取 ELR 状态"""
        try:
            response = requests.get(f"{self.api_base_url}/status", timeout=2)
            if response.status_code == 200:
                data = response.json()
                if "status" in data:
                    self.status = data["status"]
                else:
                    self.status = "未知"
            else:
                self.status = "错误"
        except Exception:
            self.status = "未连接"
        return self.status
    
    def get_elr_containers(self):
        """获取 ELR 容器列表"""
        try:
            response = requests.get(f"{self.api_base_url}/containers", timeout=2)
            if response.status_code == 200:
                data = response.json()
                if isinstance(data, list):
                    self.containers = data
                else:
                    self.containers = []
            else:
                self.containers = []
        except Exception:
            self.containers = []
        return self.containers
    
    def register_callback(self, callback):
        """注册状态更新回调"""
        self.callbacks.append(callback)
    
    def notify_callbacks(self):
        """通知所有回调"""
        for callback in self.callbacks:
            callback()

class SettingsWindow:
    """设置窗口类"""
    
    def __init__(self, parent, elr_client):
        self.parent = parent
        self.elr_client = elr_client
        self.window = tk.Toplevel(parent)
        self.window.title("ELR 桌面助手设置")
        self.window.geometry("400x200")
        self.window.resizable(False, False)
        
        # 创建 UI 元素
        self.create_ui()
    
    def create_ui(self):
        """创建 UI 元素"""
        # 主框架
        main_frame = ttk.Frame(self.window, padding="20")
        main_frame.pack(fill=tk.BOTH, expand=True)
        
        # ELR API 地址
        api_frame = ttk.LabelFrame(main_frame, text="ELR API 配置")
        api_frame.pack(fill=tk.X, pady=(0, 20))
        
        # 地址标签和输入框
        address_frame = ttk.Frame(api_frame)
        address_frame.pack(fill=tk.X, pady=(10, 5))
        
        address_label = ttk.Label(address_frame, text="地址:", width=10)
        address_label.pack(side=tk.LEFT, padx=(10, 5))
        
        self.address_var = tk.StringVar()
        # 从当前 API 地址中提取主机地址
        current_url = self.elr_client.api_base_url
        if "http://" in current_url:
            host_part = current_url.split("/api")[0]
            self.address_var.set(host_part)
        else:
            self.address_var.set("http://localhost:8080")
        
        address_entry = ttk.Entry(address_frame, textvariable=self.address_var, width=30)
        address_entry.pack(side=tk.LEFT, fill=tk.X, expand=True)
        
        # 按钮框架
        button_frame = ttk.Frame(main_frame)
        button_frame.pack(fill=tk.X, pady=(0, 10))
        
        confirm_button = ttk.Button(button_frame, text="确认", command=self.confirm)
        confirm_button.pack(side=tk.LEFT, padx=(0, 10))
        
        cancel_button = ttk.Button(button_frame, text="取消", command=self.cancel)
        cancel_button.pack(side=tk.LEFT)
    
    def confirm(self):
        """确认设置"""
        # 显示连接提示
        message_window = tk.Toplevel(self.window)
        message_window.title("连接 ELR")
        message_window.geometry("300x100")
        message_window.resizable(False, False)
        message_window.transient(self.window)
        message_window.grab_set()
        
        # 居中显示
        screen_width = message_window.winfo_screenwidth()
        screen_height = message_window.winfo_screenheight()
        x = (screen_width - 300) // 2
        y = (screen_height - 100) // 2
        message_window.geometry(f"300x100+{x}+{y}")
        
        # 提示信息
        label = ttk.Label(message_window, text="正在连接 ELR...", font=("Microsoft YaHei", 10))
        label.pack(pady=20)
        
        # 处理确认
        def connect_elr():
            api_url = self.address_var.get().strip()
            if not api_url.endswith("/api"):
                api_url = api_url.rstrip("/") + "/api"
            
            # 连接 ELR
            old_status = self.elr_client.status
            self.elr_client.set_api_url(api_url)
            
            # 关闭提示窗口
            message_window.destroy()
            
            # 显示连接结果
            result_window = tk.Toplevel(self.window)
            result_window.title("连接结果")
            result_window.geometry("300x100")
            result_window.resizable(False, False)
            result_window.transient(self.window)
            result_window.grab_set()
            
            # 居中显示
            screen_width = result_window.winfo_screenwidth()
            screen_height = result_window.winfo_screenheight()
            x = (screen_width - 300) // 2
            y = (screen_height - 100) // 2
            result_window.geometry(f"300x100+{x}+{y}")
            
            # 显示结果信息
            if self.elr_client.status != "未连接" and self.elr_client.status != "错误":
                result_text = f"连接 ELR 成功！\n当前状态: {self.elr_client.status}"
            else:
                result_text = f"连接 ELR 失败！\n请检查地址是否正确。"
            
            label = ttk.Label(result_window, text=result_text, font=("Microsoft YaHei", 10))
            label.pack(pady=20)
            
            # 确定按钮
            button = ttk.Button(result_window, text="确定", command=lambda: [result_window.destroy(), self.window.destroy()])
            button.pack(pady=5)
        
        # 启动连接线程
        threading.Thread(target=connect_elr).start()
    
    def cancel(self):
        """取消设置"""
        self.window.destroy()

class DialogWindow:
    """ELR 容器对话窗口类"""
    
    def __init__(self, parent, elr_client):
        self.parent = parent
        self.elr_client = elr_client
        # 检查parent是否有desktop_widget属性，如果有则使用其root作为父窗口
        if hasattr(parent, 'desktop_widget') and hasattr(parent.desktop_widget, 'root'):
            self.window = tk.Toplevel(parent.desktop_widget.root)
        else:
            # 否则直接使用parent作为父窗口
            self.window = tk.Toplevel(parent)
        self.window.title("ELR 容器-构件对话")
        self.window.geometry("600x400")
        # 绑定主窗口关闭事件
        self.window.protocol("WM_DELETE_WINDOW", self.close)
        
        # 始终跟随主窗口，去掉跟随状态切换
        self.follow_main = True
        
        # 创建 UI 元素
        self.create_ui()
        
        # 使主窗口可拖动（仅标题栏）
        # 绑定标题栏点击事件
        self.window.bind("<Button-1>", self.on_window_title_click)
        # 绑定鼠标移动事件，但在drag_main方法中检查是否在标题栏
        self.window.bind("<B1-Motion>", self.drag_main)
        
        # 移除输出区的鼠标事件禁用，允许选中、复制、点击和右击操作
        # 只为输出区的直接父框架绑定事件，确保输出区本身不会触发拖动
        # 但保留主窗口边框区域的事件绑定，以便能够拖动主窗口
        if self.dialog_text.master:
            self.dialog_text.master.bind("<Button-1>", lambda e: "break")
            self.dialog_text.master.bind("<B1-Motion>", lambda e: "break")
    
    def create_ui(self):
        """创建 UI 元素"""
        # 主框架
        main_frame = ttk.Frame(self.window, padding="10")
        main_frame.pack(fill=tk.BOTH, expand=True)
        
        # 对话区域
        dialog_frame = ttk.Frame(main_frame)
        dialog_frame.pack(fill=tk.BOTH, expand=True, pady=(0, 10))
        
        # 对话内容
        self.dialog_text = tk.Text(dialog_frame, bg="black", fg="white", font=("Microsoft YaHei", 10))
        self.dialog_text.pack(fill=tk.BOTH, expand=True)
        
        # 初始文本
        self.dialog_text.insert(tk.END, "ELR 容器对话窗口\n\n")
        self.dialog_text.insert(tk.END, "欢迎使用 ELR 容器对话功能！\n")
        self.dialog_text.insert(tk.END, "您可以在这里与 ELR 容器进行交互。\n")
        self.dialog_text.insert(tk.END, "ELR 容器提示：点击此窗口的上边框或左右两边框都可以触发对话输入窗口跟随\n\n")

        # 浮动输入区域
        self.input_window = tk.Toplevel()  # 不设置父窗口，避免自动跟随
        self.input_window.title("ELR 对话输入")
        self.input_window.geometry("600x150")  # 宽度与主窗口一致
        
        # 设置输入窗口属性，确保始终在最上层
        self.input_window.attributes('-topmost', True)  # 始终在最上层
        
        # 限制输入窗口只能左右拉伸
        self.input_window.resizable(True, False)
        
        # 输入窗口关闭标志
        self.input_window_closed = False
        
        # 绑定输入窗口关闭事件
        self.input_window.protocol("WM_DELETE_WINDOW", self.on_input_window_close)
        
        # 绑定主窗口的点击事件，实现特定的焦点行为
        def on_main_window_click(e):
            # 获取窗口尺寸
            width = self.window.winfo_width()
            height = self.window.winfo_height()
            
            # 检查是否点击了标题栏或边框区域
            # 确保边框区域判定在窗口内部，不超出窗口范围
            if (e.y >= 0 and e.y < 20) or \
               (e.x >= 0 and e.x < 10) or \
               (e.x >= width - 10 and e.x < width) or \
               (e.y >= height - 10 and e.y < height):
                # 点击标题栏或边框区域，开始拖动
                self.start_drag_main(e)
                # 输入窗口失去焦点，主窗口获得焦点
                if not self.input_window_closed:
                    self.input_window.grab_release()
                self.window.focus_set()
            else:
                # 点击其他区域，输入窗口保持焦点
                # 注意：不要使用grab_set，因为这会干扰主窗口的拖动操作
                # 只设置焦点，不强制抓取
                if not self.input_window_closed:
                    self.input_window.focus_set()
        
        # 绑定主窗口的点击事件
        self.window.bind("<Button-1>", on_main_window_click)
        
        # 绑定主窗口的其他事件，实现特定的焦点行为
        def on_main_window_focus_in(e):
            # 检查主窗口是否获得焦点
            if e.widget == self.window:
                # 获取鼠标位置
                x = self.window.winfo_pointerx() - self.window.winfo_rootx()
                y = self.window.winfo_pointery() - self.window.winfo_rooty()
                # 获取窗口尺寸
                width = self.window.winfo_width()
                height = self.window.winfo_height()
                # 检查是否点击了标题栏或边框区域，确保在窗口内部
                if (y >= 0 and y < 20) or \
                   (x >= 0 and x < 5) or \
                   (x >= width - 5 and x < width) or \
                   (y >= height - 5 and y < height):
                    # 点击标题栏或边框区域，输入窗口失去焦点
                    if not self.input_window_closed:
                        self.input_window.grab_release()
                else:
                    # 点击其他区域，输入窗口保持焦点
                    # 注意：不要使用grab_set，因为这会干扰主窗口的拖动操作
                    # 只设置焦点，不强制抓取
                    if not self.input_window_closed:
                        self.input_window.focus_set()
        
        # 绑定主窗口的FocusIn事件
        self.window.bind("<FocusIn>", on_main_window_focus_in)
        
        # 计算输入窗口位置，使其位于主窗口下方
        main_x = self.window.winfo_x()
        main_y = self.window.winfo_y()
        main_height = self.window.winfo_height()
        if main_height == 1:  # 处理窗口还未完全初始化的情况
            main_height = 400
        # 直接位于主窗口下方，无间距
        self.input_window.geometry(f"600x150+{main_x}+{main_y + main_height}")
        
        # 使输入窗口可拖动（仅标题栏）
        self.input_window.bind("<Button-1>", self.on_input_window_title_click)
        self.input_window.bind("<B1-Motion>", self.drag_input)
        
        # 输入窗口框架
        input_frame = ttk.Frame(self.input_window, padding="10")
        input_frame.pack(fill=tk.BOTH, expand=True)
        
        # 多行输入框
        self.input_text = tk.Text(input_frame, height=3, wrap=tk.WORD, font=("Microsoft YaHei", 10))
        self.input_text.pack(fill=tk.X, pady=(0, 15))
        
        # 按钮框架
        buttons_frame = ttk.Frame(input_frame)
        buttons_frame.pack(fill=tk.X)
        
        # 左侧按钮（加号、构件模式、调试构件）
        left_buttons_frame = ttk.Frame(buttons_frame)
        left_buttons_frame.pack(side=tk.LEFT)
        
        # 加号按钮
        plus_button = ttk.Button(left_buttons_frame, text="+", width=3, command=self.add_input_field)
        plus_button.pack(side=tk.LEFT, padx=(0, 10))
        
        # 构件模式按钮
        component_button = ttk.Button(left_buttons_frame, text="构件模式 >")
        component_button.pack(side=tk.LEFT, padx=(0, 10))
        
        # 调试构件按钮
        debug_button = ttk.Button(left_buttons_frame, text="调试构件")
        debug_button.pack(side=tk.LEFT, padx=(0, 10))
        
        # 右侧按钮（发送）
        right_buttons_frame = ttk.Frame(buttons_frame)
        right_buttons_frame.pack(side=tk.RIGHT)
        
        # 发送按钮
        send_button = ttk.Button(right_buttons_frame, text="发送", command=self.send_message)
        send_button.pack(side=tk.RIGHT)
        
        # 在输入窗口标题栏添加取消跟随/跟随文字
        self.update_window_title()
        
        # 设置输入框提示文本
        self.set_input_placeholder("请输入你要跟ELR对话内容")
        
        # 绑定输入框事件
        self.input_text.bind("<FocusIn>", self.on_input_focus_in)
        self.input_text.bind("<FocusOut>", self.on_input_focus_out)
        self.input_text.bind("<Return>", self.on_enter_pressed)
    
    def start_drag_main(self, event):
        """开始拖动主窗口"""
        self.window.x = event.x
        self.window.y = event.y
    
    def on_window_title_click(self, event):
        """主窗口标题栏或边框点击事件"""
        # 获取窗口尺寸
        width = self.window.winfo_width()
        height = self.window.winfo_height()
        
        # 检查是否点击了标题栏或边框区域，确保在窗口内部
        if (event.y >= 0 and event.y < 20) or \
           (event.x >= 0 and event.x < 5) or \
           (event.x >= width - 5 and event.x < width) or \
           (event.y >= height - 5 and event.y < height):
            # 点击标题栏或边框区域，开始拖动
            self.start_drag_main(event)
        else:
            # 输出区点击，不开始拖动
            # 清除可能的拖动状态
            if hasattr(self.window, 'x'):
                delattr(self.window, 'x')
            if hasattr(self.window, 'y'):
                delattr(self.window, 'y')
            # 返回"break"来阻止事件继续传播，确保输出区的鼠标事件能够正常工作
            return "break"
            
    def drag_main(self, event):
        """拖动主窗口"""
        # 检查是否已经开始拖动（通过检查window.x是否存在）
        # 并且确保拖动是从真正的标题栏或边框区域开始的
        if hasattr(self.window, 'x') and hasattr(self.window, 'y'):
            # 获取窗口尺寸
            width = self.window.winfo_width()
            height = self.window.winfo_height()
            # 检查拖动是否从标题栏或边框区域开始，确保在窗口内部
            start_x = self.window.x
            start_y = self.window.y
            if (start_y >= 0 and start_y < 20) or \
               (start_x >= 0 and start_x < 10) or \
               (start_x >= width - 10 and start_x < width) or \
               (start_y >= height - 10 and start_y < height):
                deltax = event.x - self.window.x
                deltay = event.y - self.window.y
                x = self.window.winfo_x() + deltax
                y = self.window.winfo_y() + deltay
                
                # 移动主窗口
                self.window.geometry(f"600x400+{x}+{y}")
                
                # 根据跟随状态决定是否移动输入窗口
                if self.follow_main and not self.input_window_closed:
                    # 直接更新输入窗口位置
                    main_height = self.window.winfo_height()
                    if main_height == 1:  # 处理窗口还未完全初始化的情况
                        main_height = 400
                    new_input_x = x
                    # 直接位于主窗口下方，无间距
                    new_input_y = y + main_height
                    # 使用主窗口的宽度作为输入窗口的宽度
                    input_width = self.window.winfo_width()
                    input_height = 150
                    # 使用完整的geometry字符串
                    self.input_window.geometry(f"{input_width}x{input_height}+{new_input_x}+{new_input_y}")
                    # 强制更新输入窗口位置
                    self.input_window.update()
        else:
            # 如果拖动不是从标题栏开始的，清除拖动状态
            if hasattr(self.window, 'x'):
                delattr(self.window, 'x')
            if hasattr(self.window, 'y'):
                delattr(self.window, 'y')
    

    
    def toggle_follow(self):
        """切换跟随状态"""
        self.follow_main = not self.follow_main
        # 立即跟随主窗口
        if self.follow_main and not self.input_window_closed:
            x = self.window.winfo_x()
            y = self.window.winfo_y()
            main_height = self.window.winfo_height()
            if main_height == 1:  # 处理窗口还未完全初始化的情况
                main_height = 400
            new_input_x = x
            new_input_y = y + main_height
            self.input_window.geometry(f"600x150+{new_input_x}+{new_input_y}")
            # 强制更新输入窗口位置
            self.input_window.update_idletasks()
        # 更新窗口标题
        self.update_window_title()
    
    def update_window_title(self):
        """更新输入窗口标题"""
        # 只显示固定标题，去掉跟随状态
        self.input_window.title("ELR 对话输入")
    
    def start_drag_input(self, event):
        """开始拖动输入窗口"""
        self.input_window.x = event.x
        self.input_window.y = event.y
    
    def drag_input(self, event):
        """拖动输入窗口"""
        # 只有在拖动状态已经开始时才执行拖动操作
        if hasattr(self.input_window, 'x') and hasattr(self.input_window, 'y'):
            deltax = event.x - self.input_window.x
            deltay = event.y - self.input_window.y
            x = self.input_window.winfo_x() + deltax
            y = self.input_window.winfo_y() + deltay
            self.input_window.geometry(f"600x150+{x}+{y}")
    
    def on_input_window_title_click(self, event):
        """输入窗口标题栏点击事件"""
        # 只有在真正的标题栏区域（窗口顶部）点击时才开始拖动
        # 标题栏高度大约为30px
        if event.y < 30:
            # 点击标题栏区域，开始拖动
            self.start_drag_input(event)
        # 点击非标题栏区域，不开始拖动
        else:
            # 清除可能的拖动状态
            if hasattr(self.input_window, 'x'):
                delattr(self.input_window, 'x')
            if hasattr(self.input_window, 'y'):
                delattr(self.input_window, 'y')
    
    def set_input_placeholder(self, text):
        """设置输入框提示文本"""
        self.input_text.insert(tk.END, text)
        self.input_text.config(fg="gray")
    
    def on_input_focus_in(self, event):
        """输入框获得焦点时"""
        if self.input_text.get("1.0", tk.END).strip() == "请输入你要跟ELR对话内容":
            self.input_text.delete("1.0", tk.END)
            self.input_text.config(fg="black")
    
    def on_input_focus_out(self, event):
        """输入框失去焦点时"""
        if not self.input_text.get("1.0", tk.END).strip():
            self.set_input_placeholder("请输入你要跟ELR对话内容")
    
    def on_enter_pressed(self, event):
        """按下回车键发送消息"""
        if event.state & 0x10:  # 检查是否按下了Shift键
            return  # 如果按下了Shift键，允许换行
        else:
            self.send_message()
            return "break"  # 阻止默认行为
    
    def add_input_field(self):
        """添加输入字段（文件选择）"""
        # 打开文件选择对话框
        # 使用input_window作为父窗口，确保在Windows平台上能正常显示
        file_path = filedialog.askopenfilename(
            parent=self.input_window,
            title="选择文件",
            filetypes=[
                ("所有文件", "*.*"),
                ("Python文件", "*.py"),
                ("模型文件", "*.pt *.pth *.onnx"),
                ("配置文件", "*.json *.yaml *.yml"),
                ("图像文件", "*.png *.jpg *.jpeg *.gif"),
                ("音频文件", "*.wav *.mp3 *.flac")
            ]
        )
        
        if file_path:
            # 显示文件路径到输入框
            self.input_text.delete("1.0", tk.END)
            self.input_text.insert(tk.END, f"[文件] {file_path}")
            self.input_text.config(fg="black")
    
    def send_message(self):
        """发送消息"""
        message = self.input_text.get("1.0", tk.END).strip()
        if message and message != "请输入你要跟ELR对话内容":
            # 如果主窗口已关闭或被隐藏，重新显示主窗口
            if not hasattr(self, 'window') or not self.window.winfo_exists():
                # 重新创建主窗口
                if hasattr(self.parent, 'desktop_widget') and hasattr(self.parent.desktop_widget, 'root'):
                    self.window = tk.Toplevel(self.parent.desktop_widget.root)
                else:
                    self.window = tk.Toplevel(self.parent)
                self.window.title("ELR 容器-构件对话")
                self.window.geometry("600x400")
                # 绑定主窗口关闭事件
                self.window.protocol("WM_DELETE_WINDOW", self.close)
                
                # 重新创建主窗口UI
                # 主框架
                main_frame = ttk.Frame(self.window, padding="10")
                main_frame.pack(fill=tk.BOTH, expand=True)
                
                # 对话区域
                dialog_frame = ttk.Frame(main_frame)
                dialog_frame.pack(fill=tk.BOTH, expand=True, pady=(0, 10))
                
                # 对话内容
                self.dialog_text = tk.Text(dialog_frame, bg="black", fg="white", font=("Microsoft YaHei", 10))
                self.dialog_text.pack(fill=tk.BOTH, expand=True)
                
                # 初始文本
                self.dialog_text.insert(tk.END, "ELR 容器对话窗口\n\n")
                self.dialog_text.insert(tk.END, "欢迎使用 ELR 容器对话功能！\n")
                self.dialog_text.insert(tk.END, "您可以在这里与 ELR 容器进行交互。\n")
                self.dialog_text.insert(tk.END, "ELR 容器提示：点击此窗口的上边框或左右两边框都可以触发对话输入窗口跟随\n\n")
                
                # 移除输出区的鼠标事件禁用，允许选中、复制、点击和右击操作
                # 只为输出区的直接父框架绑定事件，确保输出区本身不会触发拖动
                # 但保留主窗口边框区域的事件绑定，以便能够拖动主窗口
                if self.dialog_text.master:
                    self.dialog_text.master.bind("<Button-1>", lambda e: "break")
                    self.dialog_text.master.bind("<B1-Motion>", lambda e: "break")
                
                # 绑定主窗口的点击事件，实现特定的焦点行为
                def on_main_window_click(e):
                    # 获取窗口尺寸
                    width = self.window.winfo_width()
                    height = self.window.winfo_height()
                    
                    # 检查是否点击了标题栏或边框区域，确保在窗口内部
                    if (e.y >= 0 and e.y < 20) or \
                       (e.x >= 0 and e.x < 10) or \
                       (e.x >= width - 10 and e.x < width) or \
                       (e.y >= height - 10 and e.y < height):
                        # 点击标题栏或边框区域，开始拖动
                        self.start_drag_main(e)
                        # 输入窗口失去焦点，主窗口获得焦点
                        if not self.input_window_closed:
                            self.input_window.grab_release()
                        self.window.focus_set()
                    else:
                        # 点击其他区域，输入窗口保持焦点
                        # 注意：不要使用grab_set，因为这会干扰主窗口的拖动操作
                        # 只设置焦点，不强制抓取
                        if not self.input_window_closed:
                            self.input_window.focus_set()
                
                # 绑定主窗口的点击事件
                self.window.bind("<Button-1>", on_main_window_click)
                
                # 绑定主窗口的其他事件，实现特定的焦点行为
                def on_main_window_focus_in(e):
                    # 检查主窗口是否获得焦点
                    if e.widget == self.window:
                        # 获取鼠标位置
                        x = self.window.winfo_pointerx() - self.window.winfo_rootx()
                        y = self.window.winfo_pointery() - self.window.winfo_rooty()
                        # 获取窗口尺寸
                        width = self.window.winfo_width()
                        height = self.window.winfo_height()
                        # 检查是否点击了标题栏或边框区域，确保在窗口内部
                        if (y >= 0 and y < 20) or \
                           (x >= 0 and x < 5) or \
                           (x >= width - 5 and x < width) or \
                           (y >= height - 5 and y < height):
                            # 点击标题栏或边框区域，输入窗口失去焦点
                            if not self.input_window_closed:
                                self.input_window.grab_release()
                        else:
                            # 点击其他区域，输入窗口保持焦点
                            # 注意：不要使用grab_set，因为这会干扰主窗口的拖动操作
                            # 只设置焦点，不强制抓取
                            if not self.input_window_closed:
                                self.input_window.focus_set()
                
                # 绑定主窗口的FocusIn事件
                self.window.bind("<FocusIn>", on_main_window_focus_in)
                
                # 使主窗口可拖动（仅标题栏）
                # 绑定标题栏点击事件
                self.window.bind("<Button-1>", self.on_window_title_click)
                # 绑定鼠标移动事件，但在drag_main方法中检查是否在标题栏
                self.window.bind("<B1-Motion>", self.drag_main)
                
                # 通知主应用主窗口已重新打开
                if self.parent and hasattr(self.parent, 'dialog_window') and self.parent.dialog_window == self:
                    self.parent.dialog_open = True
                    self.parent.update_tray_menu()
            else:
                # 如果主窗口存在但被隐藏，显示它
                self.window.deiconify()
            
            # 确保主窗口在最前面
            self.window.lift()
            self.window.attributes('-topmost', True)
            self.window.attributes('-topmost', False)
            
            # 检查是否是文件上传
            if message.startswith("[文件]"):
                # 提取文件路径
                file_path = message.replace("[文件]", "").strip()
                # 添加文件上传消息到对话区域
                self.dialog_text.insert(tk.END, f"你: [文件上传] {file_path}\n")
                # 处理文件上传
                response = self.upload_file(file_path)
            else:
                # 添加用户消息到对话区域
                self.dialog_text.insert(tk.END, f"你: {message}\n")
                # 模拟 ELR 回应
                response = self.get_elr_response(message)
            
            self.dialog_text.insert(tk.END, f"ELR: {response}\n\n")
            
            # 清空输入框
            self.input_text.delete("1.0", tk.END)
            self.set_input_placeholder("请输入你要跟ELR对话内容")
            
            # 滚动到底部
            self.dialog_text.see(tk.END)
    
    def upload_file(self, file_path):
        """上传文件到 ELR 容器"""
        import os
        
        # 检查文件是否存在
        if not os.path.exists(file_path):
            return f"错误：文件不存在: {file_path}"
        
        # 检查文件大小
        file_size = os.path.getsize(file_path)
        if file_size > 50 * 1024 * 1024:  # 50MB 限制
            return f"错误：文件过大（{file_size / (1024*1024):.2f}MB），最大支持50MB"
        
        # 尝试上传文件
        try:
            # 构建API URL
            upload_url = f"{self.elr_client.api_base_url}/upload"
            
            # 准备文件数据
            with open(file_path, 'rb') as f:
                files = {'file': (os.path.basename(file_path), f)}
                # 发送请求
                response = requests.post(upload_url, files=files, timeout=30)
            
            if response.status_code == 200:
                data = response.json()
                if data.get('success'):
                    return f"文件上传成功！{data.get('message', '')}"
                else:
                    return f"文件上传失败：{data.get('error', '未知错误')}"
            else:
                return f"文件上传失败：HTTP {response.status_code}"
        except Exception as e:
            return f"文件上传失败：{str(e)}"
    
    def get_elr_response(self, message):
        """获取 ELR 回应"""
        # 简单的回应逻辑
        if "status" in message.lower():
            return f"ELR 状态: {self.elr_client.status}"
        elif "containers" in message.lower():
            containers = self.elr_client.get_elr_containers()
            if containers:
                return f"已加载 {len(containers)} 个容器"
            else:
                return "没有加载容器"
        elif "hello" in message.lower() or "你好" in message:
            return "你好！我是 ELR 容器助手，有什么可以帮助你的吗？"
        else:
            return "我收到了你的消息，正在处理..."
    
    def on_input_window_close(self):
        """输入窗口关闭事件处理"""
        self.input_window_closed = True
        # 隐藏输入窗口而不是销毁
        self.input_window.withdraw()
        # 当输入窗口关闭时，关闭主窗口
        self.close()
        # 通知主应用更新托盘菜单
        if self.parent and hasattr(self.parent, 'dialog_window') and self.parent.dialog_window == self:
            self.parent.dialog_open = False
            self.parent.update_tray_menu()
        
    def show_input_window(self):
        """重新显示输入窗口"""
        if self.input_window_closed:
            if hasattr(self, 'input_window') and self.input_window.winfo_exists():
                # 如果输入窗口存在且未被销毁，显示它
                self.input_window.deiconify()
                # 确保输入窗口在最上层
                self.input_window.attributes('-topmost', True)
                self.input_window.attributes('-topmost', False)
            else:
                # 重新创建输入窗口
                self.input_window = tk.Toplevel()  # 不设置父窗口，避免自动跟随
                self.input_window.title("ELR 对话输入")
                
                # 计算输入窗口位置，使其位于主窗口下方
                main_x = self.window.winfo_x()
                main_y = self.window.winfo_y()
                main_height = self.window.winfo_height()
                main_width = self.window.winfo_width()
                if main_height == 1:  # 处理窗口还未完全初始化的情况
                    main_height = 400
                    main_width = 600
                # 直接位于主窗口下方，无间距
                self.input_window.geometry(f"{main_width}x150+{main_x}+{main_y + main_height}")
                
                # 设置输入窗口属性，确保始终在最上层
                self.input_window.attributes('-topmost', True)  # 始终在最上层
                
                # 限制输入窗口只能左右拉伸
                self.input_window.resizable(True, False)
                
                # 使输入窗口可拖动（仅标题栏）
                self.input_window.bind("<Button-1>", self.on_input_window_title_click)
                self.input_window.bind("<B1-Motion>", self.drag_input)
                
                # 输入窗口框架
                input_frame = ttk.Frame(self.input_window, padding="10")
                input_frame.pack(fill=tk.BOTH, expand=True)
                
                # 多行输入框
                self.input_text = tk.Text(input_frame, height=3, wrap=tk.WORD, font=("Microsoft YaHei", 10))
                self.input_text.pack(fill=tk.X, pady=(0, 15))
                
                # 按钮框架
                buttons_frame = ttk.Frame(input_frame)
                buttons_frame.pack(fill=tk.X)
                
                # 左侧按钮（加号、构件模式、调试构件）
                left_buttons_frame = ttk.Frame(buttons_frame)
                left_buttons_frame.pack(side=tk.LEFT)
                
                # 加号按钮
                plus_button = ttk.Button(left_buttons_frame, text="+", width=3, command=self.add_input_field)
                plus_button.pack(side=tk.LEFT, padx=(0, 10))
                
                # 构件模式按钮
                component_button = ttk.Button(left_buttons_frame, text="构件模式 >")
                component_button.pack(side=tk.LEFT, padx=(0, 10))
                
                # 调试构件按钮
                debug_button = ttk.Button(left_buttons_frame, text="调试构件")
                debug_button.pack(side=tk.LEFT, padx=(0, 10))
                
                # 右侧按钮（发送）
                right_buttons_frame = ttk.Frame(buttons_frame)
                right_buttons_frame.pack(side=tk.RIGHT)
                
                # 发送按钮
                send_button = ttk.Button(right_buttons_frame, text="发送", command=self.send_message)
                send_button.pack(side=tk.RIGHT)
                
                # 在输入窗口标题栏添加取消跟随/跟随文字
                self.update_window_title()
                
                # 设置输入框提示文本
                self.set_input_placeholder("请输入你要跟ELR对话内容")
                
                # 绑定输入框事件
                self.input_text.bind("<FocusIn>", self.on_input_focus_in)
                self.input_text.bind("<FocusOut>", self.on_input_focus_out)
                self.input_text.bind("<Return>", self.on_enter_pressed)
                
                # 绑定输入窗口关闭事件
                self.input_window.protocol("WM_DELETE_WINDOW", self.on_input_window_close)
            
            # 重置关闭标志
            self.input_window_closed = False
    
    def close(self):
        """关闭窗口"""
        # 主窗口关闭时，同时隐藏输入子窗口
        self.window.withdraw()
        # 隐藏输入窗口
        if hasattr(self, 'input_window') and self.input_window.winfo_exists():
            self.input_window.withdraw()
        # 当主窗口关闭时，通知主应用更新托盘菜单
        if self.parent and hasattr(self.parent, 'dialog_window') and self.parent.dialog_window == self:
            self.parent.dialog_open = False
            self.parent.update_tray_menu()

class DesktopWidget:
    """桌面助手 widget 类"""
    
    def __init__(self, elr_client):
        self.elr_client = elr_client
        self.root = tk.Tk()
        self.root.title("ELR 桌面助手")
        self.root.overrideredirect(True)  # 去除窗口边框
        self.root.attributes("-topmost", True)  # 置于顶层
        self.root.attributes("-alpha", 0.8)  # 半透明
        
        # 窗口可见性状态
        self.is_visible = True
        
        # 设置窗口大小
        self.width = 300
        self.height = 200
        
        # 计算屏幕右上角位置
        screen_width = self.root.winfo_screenwidth()
        screen_height = self.root.winfo_screenheight()
        x = screen_width - self.width - 10
        y = 10
        self.root.geometry(f"{self.width}x{self.height}+{x}+{y}")
        
        # 绑定鼠标事件
        self.root.bind("<Button-1>", self.start_drag)
        self.root.bind("<B1-Motion>", self.drag)
        
        # 创建 UI 元素
        self.create_ui()
        
        # 注册回调
        self.elr_client.register_callback(self.update_ui)
    
    def create_ui(self):
        """创建 UI 元素"""
        # 主框架
        main_frame = ttk.Frame(self.root, padding="10")
        main_frame.pack(fill=tk.BOTH, expand=True)
        
        # 标题
        title_label = ttk.Label(main_frame, text="ELR 桌面助手", font=("Microsoft YaHei", 12, "bold"))
        title_label.pack(anchor=tk.W, pady=(0, 10))
        
        # ELR 状态
        status_frame = ttk.Frame(main_frame)
        status_frame.pack(fill=tk.X, pady=(0, 10))
        
        status_label = ttk.Label(status_frame, text="ELR 状态:", width=10)
        status_label.pack(side=tk.LEFT)
        
        self.status_var = tk.StringVar()
        self.status_var.set("未连接")
        status_value = ttk.Label(status_frame, textvariable=self.status_var, width=15)
        status_value.pack(side=tk.LEFT)
        
        # 容器列表
        containers_frame = ttk.LabelFrame(main_frame, text="容器列表")
        containers_frame.pack(fill=tk.BOTH, expand=True, pady=(0, 10))
        
        self.containers_tree = ttk.Treeview(containers_frame, columns=("name", "status"), show="headings")
        self.containers_tree.heading("name", text="名称")
        self.containers_tree.heading("status", text="状态")
        self.containers_tree.column("name", width=150)
        self.containers_tree.column("status", width=100)
        self.containers_tree.pack(fill=tk.BOTH, expand=True)
        
        # 按钮框架
        button_frame = ttk.Frame(main_frame)
        button_frame.pack(fill=tk.X, pady=(0, 0))
        
        refresh_button = ttk.Button(button_frame, text="刷新", command=self.refresh)
        refresh_button.pack(side=tk.LEFT, padx=(0, 10))
        
        settings_button = ttk.Button(button_frame, text="设置", command=self.open_settings)
        settings_button.pack(side=tk.LEFT, padx=(0, 10))
        
        hide_button = ttk.Button(button_frame, text="隐藏", command=self.hide)
        hide_button.pack(side=tk.LEFT)
    
    def start_drag(self, event):
        """开始拖动"""
        self.x = event.x
        self.y = event.y
    
    def drag(self, event):
        """拖动窗口"""
        deltax = event.x - self.x
        deltay = event.y - self.y
        x = self.root.winfo_x() + deltax
        y = self.root.winfo_y() + deltay
        self.root.geometry(f"{self.width}x{self.height}+{x}+{y}")
    
    def update_ui(self):
        """更新 UI"""
        # 更新 ELR 状态
        self.status_var.set(self.elr_client.status)
        
        # 清除容器列表
        for item in self.containers_tree.get_children():
            self.containers_tree.delete(item)
        
        # 添加容器列表
        for container in self.elr_client.containers:
            name = container.get("name", "未知")
            status = container.get("status", "未知")
            self.containers_tree.insert("", tk.END, values=(name, status))
        
        # 刷新界面
        self.root.update()
    
    def refresh(self):
        """手动刷新"""
        self.elr_client.get_elr_status()
        self.elr_client.get_elr_containers()
        self.elr_client.notify_callbacks()
    
    def open_settings(self):
        """打开设置窗口"""
        SettingsWindow(self.root, self.elr_client)
    
    def show(self):
        """显示窗口"""
        self.root.deiconify()
        self.is_visible = True
    
    def hide(self):
        """隐藏窗口"""
        self.root.withdraw()
        self.is_visible = False
    
    def run(self):
        """运行主循环"""
        self.root.mainloop()

class ELRDesktopAssistant:
    """ELR 桌面助手主类"""
    
    def __init__(self):
        import logging
        logging.info("初始化 ELRDesktopAssistant")
        self.elr_client = ELRClient()
        logging.info("创建 ELRClient 实例")
        self.desktop_widget = DesktopWidget(self.elr_client)
        logging.info("创建 DesktopWidget 实例")
        self.monitoring_thread = None
        self.running = False
        self.tray_icon = None
        self.dialog_open = False
        self.dialog_window = None
        logging.info("ELRDesktopAssistant 初始化完成")
    
    def start_monitoring(self):
        """开始监控"""
        self.running = True
        self.monitoring_thread = threading.Thread(target=self.monitoring_loop)
        self.monitoring_thread.daemon = True
        self.monitoring_thread.start()
    
    def stop_monitoring(self):
        """停止监控"""
        self.running = False
        if self.monitoring_thread:
            self.monitoring_thread.join()
    
    def monitoring_loop(self):
        """监控循环"""
        while self.running:
            self.elr_client.get_elr_status()
            self.elr_client.get_elr_containers()
            self.elr_client.notify_callbacks()
            time.sleep(5)  # 每5秒检查一次
    
    def toggle_dialog(self):
        """切换对话窗口状态"""
        if self.dialog_open:
            # 关闭对话
            if self.dialog_window:
                self.dialog_window.close()
            self.dialog_open = False
        else:
            # 打开对话
            if self.dialog_window and hasattr(self.dialog_window, 'window') and self.dialog_window.window.winfo_exists():
                # 如果对话窗口实例存在且窗口未被销毁，显示它
                self.dialog_window.window.deiconify()
                # 确保输入窗口也显示
                if hasattr(self.dialog_window, 'input_window') and self.dialog_window.input_window.winfo_exists():
                    self.dialog_window.input_window.deiconify()
                # 重置输入窗口关闭标志
                self.dialog_window.input_window_closed = False
            else:
                # 创建新的对话窗口实例
                self.dialog_window = DialogWindow(self, self.elr_client)
            self.dialog_open = True
        # 更新托盘图标菜单
        if self.tray_icon:
            self.update_tray_menu()
        else:
            self.create_tray_icon()
    
    def toggle_visibility(self):
        """切换桌面显示状态"""
        # 检查窗口可见性
        if self.desktop_widget.is_visible:
            # 隐藏桌面
            self.desktop_widget.hide()
        else:
            # 显示桌面
            self.desktop_widget.show()
        # 等待窗口状态更新
        time.sleep(0.1)
        # 更新托盘图标菜单
        if self.tray_icon:
            self.update_tray_menu()
        else:
            self.create_tray_icon()
    
    def get_menu(self):
        """获取动态菜单"""
        # 动态生成菜单
        menu_items = []
        
        # 第1项：打开对话/关闭对话
        if self.dialog_open:
            menu_items.append(item('关闭对话', lambda: self.toggle_dialog()))
        else:
            menu_items.append(item('打开对话', lambda: self.toggle_dialog()))
        
        # 第2项：显示/隐藏
        # 使用 DesktopWidget 的 is_visible 属性
        if self.desktop_widget.is_visible:
            menu_items.append(item('隐藏', lambda: self.toggle_visibility()))
        else:
            menu_items.append(item('显示', lambda: self.toggle_visibility()))
        
        # 添加其他菜单项
        menu_items.extend([
            item('设置', lambda: self.desktop_widget.open_settings()),
            item('刷新', lambda: self.desktop_widget.refresh()),
            item('退出', lambda: self.exit_app()),
        ])
        
        # 创建菜单
        return tuple(menu_items)
    
    def update_tray_menu(self):
        """更新托盘图标菜单"""
        if self.tray_icon:
            # 更新托盘图标菜单
            self.tray_icon.menu = self.get_menu()
            self.tray_icon.update_menu()
    
    def create_tray_icon(self):
        """创建系统托盘图标"""
        # 获取图标路径
        import os
        import sys
        
        # 确定图标文件路径
        icon = None
        # 优先使用ico文件
        icon_path = None
        
        # 尝试多种路径查找图标文件，优先使用ico文件
        possible_paths = []
        
        if getattr(sys, 'frozen', False):
            # 打包后，尝试使用内置图标
            base_dir = sys._MEIPASS
            # 优先使用ico文件
            possible_paths.append(os.path.join(base_dir, "icons", "elr_icon.ico"))
            possible_paths.append(os.path.join(base_dir, "elr_icon.ico"))
            # 然后使用png文件作为备选
            possible_paths.append(os.path.join(base_dir, "icons", "elr_icon.png"))
            possible_paths.append(os.path.join(base_dir, "elr_icon.png"))
        else:
            # 开发模式
            base_dir = os.path.dirname(os.path.abspath(__file__))
            # 优先使用ico文件
            possible_paths.append(os.path.join(base_dir, "icons", "elr_icon.ico"))
            possible_paths.append(os.path.join(base_dir, "elr_icon.ico"))
            # 然后使用png文件作为备选
            possible_paths.append(os.path.join(base_dir, "icons", "elr_icon.png"))
            possible_paths.append(os.path.join(base_dir, "elr_icon.png"))
        
        # 尝试加载图标
        for path in possible_paths:
            if os.path.exists(path):
                try:
                    icon = Image.open(path)
                    break
                except:
                    pass
        
        # 如果没有找到图标，创建一个默认图标
        if icon is None:
            icon = Image.new("RGB", (64, 64), color="#4CAF50")
        
        # 获取动态菜单
        menu = self.get_menu()
        
        # 创建托盘图标
        self.tray_icon = pystray.Icon("ELR Desktop Assistant", icon, "ELR 桌面助手", menu)
        
        # 启动托盘图标
        self.tray_icon.run()
    
    def exit_app(self):
        """退出应用"""
        self.stop_monitoring()
        if self.tray_icon:
            self.tray_icon.stop()
        self.desktop_widget.root.quit()
        sys.exit()
    
    def run(self):
        """运行应用"""
        import logging
        logging.info("开始运行应用")
        # 开始监控
        logging.info("启动监控线程")
        self.start_monitoring()
        
        # 创建托盘图标线程
        logging.info("创建托盘图标线程")
        tray_thread = threading.Thread(target=self.create_tray_icon)
        tray_thread.daemon = True
        tray_thread.start()
        logging.info("托盘图标线程已启动")
        
        # 运行桌面 widget
        logging.info("运行桌面 widget")
        self.desktop_widget.run()
        logging.info("应用退出")

if __name__ == "__main__":
    # 添加日志文件
    import logging
    import datetime
    import os
    import mmap
    import sys
    import psutil
    
    # 配置日志
    log_dir = os.path.join(os.path.expanduser("~"), "AppData", "Local", "ELRDesktopAssistant")
    os.makedirs(log_dir, exist_ok=True)
    log_file = os.path.join(log_dir, f"elr_assistant_{datetime.datetime.now().strftime('%Y%m%d_%H%M%S')}.log")
    
    logging.basicConfig(
        level=logging.DEBUG,
        format='%(asctime)s - %(levelname)s - %(message)s',
        handlers=[
            logging.FileHandler(log_file),
            logging.StreamHandler()
        ]
    )
    
    logging.info("ELR Desktop Assistant 启动")
    
    # 检查依赖
    try:
        import pystray
        from PIL import Image
        import requests
        logging.info("依赖检查通过")
    except ImportError as e:
        logging.error(f"缺少依赖: {e}")
        print(f"缺少依赖: {e}")
        print("请安装依赖: pip install pystray Pillow requests psutil")
        sys.exit(1)
    
    # 程序唯一ID
    APP_UNIQUE_ID = "ELRDesktopAssistant_v1.0"
    
    # 使用进程ID文件进行多实例检测
    def is_already_running():
        """检查是否已有实例在运行"""
        logging.info("开始检查多实例")
        app_data_dir = os.path.join(os.path.expanduser("~"), "AppData", "Local", "ELRDesktopAssistant")
        os.makedirs(app_data_dir, exist_ok=True)
        pid_file = os.path.join(app_data_dir, "elr_assistant.pid")
        logging.info(f"PID文件路径: {pid_file}")
        
        current_pid = os.getpid()
        logging.info(f"当前 PID: {current_pid}")
        
        # 1. 检查PID文件是否存在
        if os.path.exists(pid_file):
            logging.info("发现PID文件")
            try:
                with open(pid_file, "r") as f:
                    content = f.read().strip()
                    logging.info(f"PID文件内容: {content}")
                    
                    if content:
                        try:
                            # 尝试读取PID
                            pid = int(content.split("|")[0])
                            logging.info(f"PID文件中的PID: {pid}")
                            
                            # 检查这个PID是否真的在运行
                            import psutil
                            proc = psutil.Process(pid)
                            # 检查进程是否仍然存在
                            proc_status = proc.status()
                            # 检查进程名称
                            proc_name = proc.name()
                            if proc_name == 'ELRDesktopAssistant.exe' or 'python' in proc_name.lower():
                                logging.info(f"发现其他 ELRDesktopAssistant 实例: PID {pid}, 名称: {proc_name}, 状态: {proc_status}")
                                return True
                            else:
                                logging.info(f"进程 PID {pid} 存在，但名称已变为: {proc_name}")
                        except Exception as e:
                            # PID不存在或无法访问
                            logging.info(f"PID {content} 不存在或无法访问: {e}")
                            # PID不存在，删除PID文件
                            try:
                                os.unlink(pid_file)
                                logging.info("已删除无效的PID文件")
                            except Exception as e:
                                logging.warning(f"删除PID文件失败: {e}")
            except Exception as e:
                logging.warning(f"读取PID文件失败: {e}")
                # 读取失败，删除PID文件
                try:
                    os.unlink(pid_file)
                    logging.info("已删除无效的PID文件")
                except Exception as e:
                    logging.warning(f"删除PID文件失败: {e}")
        
        # 2. 创建PID文件
        try:
            # 尝试创建PID文件
            with open(pid_file, "w") as f:
                # 写入 PID、时间戳和唯一ID
                timestamp = datetime.datetime.now().isoformat()
                f.write(f"{current_pid}|{timestamp}|{APP_UNIQUE_ID}")
                logging.info("创建PID文件成功，无其他实例运行")
                return False
        except Exception as e:
            # 如果创建失败，认为没有其他实例
            logging.warning(f"创建PID文件失败: {e}")
            logging.info("创建PID文件失败，继续运行")
            return False
    
    # 检查是否已有实例在运行
    if is_already_running():
        logging.info("ELR Desktop Assistant 已经在运行中！")
        print("ELR Desktop Assistant 已经在运行中！")
        sys.exit(1)
    
    try:
        # 运行应用
        logging.info("创建 ELRDesktopAssistant 实例")
        app = ELRDesktopAssistant()
        logging.info("运行应用")
        app.run()
    except Exception as e:
        logging.error(f"应用运行出错: {e}", exc_info=True)
        print(f"应用运行出错: {e}")
        sys.exit(1)
