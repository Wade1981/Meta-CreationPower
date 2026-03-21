# ELR Desktop Assistant

ELR桌面助手是一个基于Qt开发的桌面应用，用于监控和管理ELR（Enlightenment Lighthouse Runtime）环境。

## 功能特性

- **任务栏驻留**：在系统托盘显示ELR状态
- **桌面助手**：可拖动的桌面widget，显示ELR状态和容器信息
- **实时监控**：自动定期检查ELR状态和容器状态
- **容器管理**：显示ELR容器列表和状态
- **美观界面**：半透明设计，支持鼠标拖动

## 技术栈

- **C++**：核心开发语言
- **Qt 6**：GUI框架
- **CMake**：构建系统
- **网络通信**：使用Qt Network模块与ELR API通信

## 目录结构

```
ELR-Desktop-Assistant/
├── CMakeLists.txt          # CMake构建配置
├── main.cpp               # 主入口文件
├── ELRDesktopAssistant.h  # 主应用类头文件
├── ELRDesktopAssistant.cpp # 主应用类实现
├── TrayIcon.h             # 系统托盘图标类头文件
├── TrayIcon.cpp           # 系统托盘图标类实现
├── DesktopWidget.h        # 桌面助手widget头文件
├── DesktopWidget.cpp      # 桌面助手widget实现
├── ELRClient.h            # ELR客户端类头文件
├── ELRClient.cpp          # ELR客户端类实现
├── elr_icons.qrc          # Qt资源文件
├── icons/                 # 图标目录
└── README.md              # 本文件
```

## 构建步骤

1. **安装依赖**
   - Qt 6.0+
   - CMake 3.16+
   - C++17兼容编译器

2. **添加图标**
   - 在`icons`目录中添加以下图标文件：
     - `elr_icon.png`：默认图标
     - `elr_icon_running.png`：运行中状态图标
     - `elr_icon_stopped.png`：停止状态图标

3. **构建项目**
   ```bash
   mkdir build
   cd build
   cmake ..
   cmake --build .
   ```

4. **运行应用**
   - 构建完成后，在build目录中运行`ELRDesktopAssistant.exe`

## ELR API接口

应用默认连接到`http://localhost:8080/api`，需要ELR提供以下API接口：

- **GET /api/status**：获取ELR状态
- **GET /api/containers**：获取容器列表

## 界面说明

- **系统托盘**：显示ELR状态，右键菜单可操作
- **桌面助手**：
  - 显示ELR状态
  - 显示容器列表和状态
  - 支持鼠标拖动
  - 点击"刷新"按钮手动刷新状态
  - 点击"隐藏"按钮隐藏助手

## 注意事项

- 确保ELR服务正在运行
- 确保ELR API端口正确（默认8080）
- 首次运行时，桌面助手会显示在屏幕中央
- 可以通过系统托盘图标控制桌面助手的显示/隐藏