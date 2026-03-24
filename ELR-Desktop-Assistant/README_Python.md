# ELR Desktop Assistant (Python 版本)

ELR 桌面助手是一个基于 Python 开发的桌面应用，用于监控和管理 ELR（Enlightenment Lighthouse Runtime）环境。

## 功能特性

- **任务栏驻留**：在系统托盘显示 ELR 状态
- **桌面助手**：可拖动的桌面 widget，显示 ELR 状态和容器信息
- **实时监控**：自动定期检查 ELR 状态和容器状态
- **容器管理**：显示 ELR 容器列表和状态
- **美观界面**：半透明设计，支持鼠标拖动

## 技术栈

- **Python 3.6+**：核心开发语言
- **Tkinter**：Python 标准库，用于创建 GUI
- **pystray**：用于创建系统托盘图标
- **Pillow**：用于图像处理
- **requests**：用于与 ELR API 通信

## 依赖项

- **pystray**：用于系统托盘功能
- **Pillow**：用于图像处理
- **requests**：用于 HTTP 请求

## 安装依赖

```bash
pip install pystray Pillow requests
```

## 使用方法

1. **确保 ELR 服务正在运行**
2. **运行 Python 脚本**

```bash
python elr_desktop_assistant.py
```

## 界面说明

- **系统托盘**：显示 ELR 状态，右键菜单可操作
- **桌面助手**：
  - 显示 ELR 状态
  - 显示容器列表和状态
  - 支持鼠标拖动
  - 点击"刷新"按钮手动刷新状态
  - 点击"隐藏"按钮隐藏助手

## ELR API 接口

应用默认连接到 `http://localhost:8080/api`，需要 ELR 提供以下 API 接口：

- **GET /api/status**：获取 ELR 状态
- **GET /api/containers**：获取容器列表

## 注意事项

- 确保 ELR 服务正在运行
- 确保 ELR API 端口正确（默认 8080）
- 首次运行时，桌面助手会显示在屏幕中央
- 可以通过系统托盘图标控制桌面助手的显示/隐藏

## 与 Qt 版本的对比

| 特性 | Python 版本 | Qt 版本 |
|------|------------|--------|
| 依赖项 | 轻量级（pystray, Pillow, requests） | 重量级（Qt 6） |
| 安装难度 | 简单（pip 安装） | 复杂（需要安装 Qt 开发环境） |
| 跨平台性 | 跨平台（Windows, Linux, macOS） | 跨平台（需要对应平台的 Qt 版本） |
| 功能完整性 | 完整（包含所有核心功能） | 完整 |
| 性能 | 适中 | 较高 |
| 可定制性 | 高（Python 代码易于修改） | 中（需要 C++ 知识） |

## 项目结构

```
ELR-Desktop-Assistant/
├── elr_desktop_assistant.py  # Python 版本的 ELR 桌面助手
├── README_Python.md          # 本文件
├── CMakeLists.txt            # Qt 版本的 CMake 配置
├── DesktopWidget.cpp         # Qt 版本的桌面 widget 实现
├── DesktopWidget.h           # Qt 版本的桌面 widget 头文件
├── ELRClient.cpp             # Qt 版本的 ELR 客户端实现
├── ELRClient.h               # Qt 版本的 ELR 客户端头文件
├── ELRDesktopAssistant.cpp   # Qt 版本的主应用实现
├── ELRDesktopAssistant.h     # Qt 版本的主应用头文件
├── README.md                 # Qt 版本的说明文件
├── TrayIcon.cpp              # Qt 版本的托盘图标实现
├── TrayIcon.h                # Qt 版本的托盘图标头文件
├── elr_icons.qrc             # Qt 版本的资源文件
└── main.cpp                  # Qt 版本的主入口文件
```

## 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。
