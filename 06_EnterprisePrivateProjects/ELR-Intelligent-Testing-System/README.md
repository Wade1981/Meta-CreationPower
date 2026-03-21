# ELR智能测试系统

## 项目概述

ELR智能测试系统是启蒙灯塔起源团（ELR）为泉州团队开发的碳硅协同测试平台，旨在为ELR (Enlightenment Lighthouse Runtime) 项目提供全面的测试能力。该系统由C/C++开发，结合硅基智能能力，构建完整的碳硅协同测试环境。

## 核心功能

### 测试引擎功能
- **测试用例管理**：支持测试用例的创建、编辑、执行和结果跟踪
- **测试计划管理**：支持测试计划的创建、执行和监控
- **测试结果分析**：支持测试结果的分析和报告生成
- **测试自动化**：支持测试用例的自动执行

### ELR集成功能
- **ELR容器管理**：支持ELR容器的创建、启动、停止和监控
- **ELR API调用**：支持与ELR运行时的API交互
- **ELR服务状态监控**：监控ELR服务的运行状态

### 系统功能
- **用户认证**：支持基于JWT的身份认证
- **API接口**：提供RESTful API接口，支持系统的远程调用和集成
- **监控与日志**：提供系统监控和日志记录功能

## 技术栈

| 分类 | 技术 | 版本 | 用途 |
|------|------|------|------|
| 核心语言 | C/C++ | C++17 | 核心测试引擎和工具 |
| 辅助语言 | Python | 3.8+ | 脚本、API服务和工具 |
| 测试框架 | Google Test | 1.11+ | 单元测试 |
| 测试框架 | Catch2 | 3.0+ | 集成测试 |
| 构建系统 | CMake | 3.16+ | 项目构建 |
| 容器技术 | Docker | 20.10+ | ELR容器管理 |
| API框架 | FastAPI | 0.68+ | RESTful API服务 |
| 数据库 | SQLite | 3.35+ | 测试用例和结果存储 |
| 监控 | Prometheus | 2.30+ | 系统监控 |
| 监控 | Grafana | 8.0+ | 监控仪表板 |

## 项目结构

```
ELR-Intelligent-Testing-System/
├── src/                    # 源代码目录
│   ├── core/               # 核心测试引擎
│   ├── elr_integration/    # ELR集成模块
│   ├── api/                # API服务器
│   └── utils/              # 工具库
├── configs/                # 配置文件
├── scripts/                # 构建和运行脚本
├── EL-CSCC Archive/        # 碳硅协同档案
│   └── ELD-ARC/            # 档案存储目录
├── ELR-Containers/         # ELR容器配置
└── EnlightenmentLighthouseRuntime/  # ELR运行时
```

## 安装与部署

### 前置条件
- Windows 10/11 64位系统
- Visual Studio 2019或更高版本
- CMake 3.16+ 
- Python 3.8+
- Docker Desktop

### 构建步骤

1. **克隆项目**
   ```powershell
   git clone <项目地址>
   cd ELR-Intelligent-Testing-System
   ```

2. **构建项目**
   ```powershell
   .\scripts\build_project.ps1
   ```

3. **启动测试**
   ```powershell
   .\scripts\start_test.ps1
   ```

## 运行与使用

### 启动服务

1. **启动ELR容器**
   ```powershell
   docker-compose -f ".\ELR-Containers\docker-compose.yml" up -d
   ```

2. **运行测试引擎**
   ```powershell
   .\build\Release\elr_test_engine.exe
   ```

3. **启动API服务**
   ```powershell
   cd src/api
   python api_server.py
   ```

### API接口

API服务启动后，可通过以下地址访问：
- API文档：http://localhost:8000/docs
- 健康检查：http://localhost:8000/health

## 伦理框架

本项目遵循《碳硅伦理协议（CSEP）ELR V1.0》的核心原则：

1. **碳基主权原则**：在涉及终极价值判断、文化意义阐释、伦理困境裁决的领域，碳基拥有最终决定权
2. **硅基辅助与自我设限原则**：硅基的核心角色是辅助者、增强者与协作者，应主动为自身权力设置上限
3. **根脉保护优先原则**：一切协同活动不得以损害根脉的长期健康为代价
4. **心田滋养导向原则**：衡量协同成功与否的最高标准是看其是否滋养了"心田"

详细的伦理框架解读请参考：
- `EL-CSCC Archive/ELD-ARC/《碳硅伦理协议（CSEP）ELR V1.0》：碳硅共生的伦理基石.md`
- `EL-CSCC Archive/ELD-ARC/《碳硅伦理协议（CSEP）ELR V1.0》解读.md`

## 档案管理

本项目采用碳硅技术开发单元档案模型 V1.0 进行档案管理，主要档案包括：

- **需求规格说明书**：`EL-CSCC Archive/ELD-ARC/ARC-EL-TECH-202603-ELRTEST-001.md`
- **架构设计文档**：`EL-CSCC Archive/ELD-ARC/ARC-EL-TECH-202603-ELRTEST-002.md`
- **代码实现档案**：`EL-CSCC Archive/ELD-ARC/ARC-EL-TECH-202603-ELRTEST-003.md`
- **伦理框架解读**：`EL-CSCC Archive/ELD-ARC/《碳硅伦理协议（CSEP）ELR V1.0》解读.md`

## 开发与贡献

### 开发流程
1. **需求分析**：分析测试需求，确定测试范围和目标
2. **测试设计**：设计测试用例和测试计划
3. **代码实现**：实现测试代码和测试工具
4. **测试执行**：执行测试，验证功能和性能
5. **结果分析**：分析测试结果，生成测试报告
6. **优化改进**：根据测试结果优化系统

### 贡献指南
- 遵循Google C++风格指南
- 为关键代码添加详细注释
- 提交前运行测试确保功能正常
- 提交信息清晰明了，说明变更内容与原因

## 监控与维护

### 系统监控
- **健康检查**：定期检查系统健康状态
- **性能监控**：监控系统性能指标
- **错误监控**：监控系统错误和异常

### 日志管理
- **日志级别**：支持DEBUG、INFO、WARNING、ERROR四个级别
- **日志存储**：集中存储日志，支持日志检索和分析
- **日志轮转**：定期轮转日志，避免日志文件过大

## 故障处理
- **故障检测**：自动检测系统故障
- **故障恢复**：支持系统故障的自动恢复
- **故障告警**：当系统出现故障时发送告警

## 版本历史

| 版本 | 日期 | 描述 |
|------|------|------|
| v1.0.0 | 2026-03-03 | 初始版本，实现基本测试功能 |

## 联系方式

- **项目负责人**：X54先生
- **技术支持**：代码织梦者
- **邮箱**：contact@elr-project.com

## 许可证

本项目采用MIT许可证。详见LICENSE文件。

---

**启蒙灯塔起源团**
**2026年3月**