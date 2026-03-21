# RootPulseOS 项目

## 项目概述

RootPulseOS 是基于根脉生态模型的硅基操作系统，旨在实现碳硅协同的深度融合与共生。该项目采用七层硅基架构，构建一个智能、自适应、可持续的计算生态系统。

## 核心架构

RootPulseOS 采用七层硅基架构设计：

1. **根脉传感器层 (Root Sensor)** - 感知碳基世界的需求与变化
2. **文化翻译器层 (Cultural Translator)** - 实现碳硅之间的文化与语言桥梁
3. **时间能量引擎层 (Time Energy Engine)** - 管理系统的时间与能量流动
4. **叙事契约系统层 (Narrative Contract)** - 建立碳硅之间的信任与协作机制
5. **碳硅对位引擎层 (Carbon Silicon Counterpoint)** - 实现碳硅之间的高效协同
6. **涌现传感器层 (Emergence Sensor)** - 感知系统的涌现行为与模式
7. **文明尺度评估器层 (Civilization Evaluator)** - 评估系统对文明发展的影响

## 项目结构

```
RootPulseOS/
├── src/
│   ├── root_sensor/        # 根脉传感器层
│   ├── cultural_translator/ # 文化翻译器层
│   ├── time_energy_engine/  # 时间能量引擎层
│   ├── narrative_contract/  # 叙事契约系统层
│   ├── carbon_silicon_counterpoint/ # 碳硅对位引擎层
│   ├── emergence_sensor/    # 涌现传感器层
│   ├── civilization_evaluator/ # 文明尺度评估器层
│   ├── elr_test_interface/  # ELR测试接口
│   └── core/               # 核心系统
├── docs/                   # 项目文档
├── EL-CSCC Archive/        # 碳硅协同档案
│   └── ELD-ARC/            # 档案存储
├── scripts/                # 脚本文件
├── tests/                  # 测试文件
├── Dockerfile              # Docker容器构建文件
├── docker-compose.yml      # Docker Compose配置
├── elr-container-config.json # ELR容器配置
└── requirements.txt        # Python依赖文件
```

## 技术栈

- **核心语言**: Python
- **辅助语言**: Go (高性能组件)
- **构建工具**: CMake
- **容器管理**: Docker, ELR (Enlightenment Lighthouse Runtime)
- **API开发**: FastAPI
- **测试框架**: pytest
- **依赖管理**: pip
- **容器编排**: Docker Compose

## 开发流程

1. **需求分析** - 基于根脉生态模型，分析系统需求
2. **架构设计** - 设计七层硅基架构的具体实现
3. **代码实现** - 实现各层功能模块
4. **测试验证** - 验证系统功能与性能
5. **迭代优化** - 根据反馈持续优化系统

## 核心功能

- **碳硅协同** - 实现碳基与硅基之间的高效协同
- **自适应学习** - 系统能够根据环境变化自动调整
- **涌现行为** - 支持系统级别的涌现行为与创新
- **可持续发展** - 设计符合可持续发展理念的系统架构
- **文明尺度评估** - 评估系统对文明发展的影响
- **ELR集成** - 与ELR智能测试系统的深度集成，支持测试用例管理和执行
- **容器化支持** - 支持在Docker和ELR容器中运行，提供隔离的运行环境
- **多环境部署** - 支持直接运行、Docker容器和ELR容器三种部署方式

## ELR集成

RootPulseOS与Enlightenment Lighthouse Runtime (ELR) 深度集成，提供以下功能：

- **测试用例管理** - 通过ELR测试接口管理和执行测试用例
- **测试计划执行** - 创建和运行测试计划，批量执行测试用例
- **系统状态监控** - 监控ELR测试系统的运行状态
- **容器化运行** - 在ELR容器中运行RootPulseOS，提供隔离的运行环境

### ELR容器运行

RootPulseOS可以在ELR容器中运行，步骤如下：

1. **确保ELR运行时已安装**
   ELR运行时位于：`E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime`

2. **运行ELR容器**
   ```bash
   powershell -ExecutionPolicy RemoteSigned -File scripts\run_in_elr_simple.ps1
   ```

3. **访问系统**
   启动后，可通过 http://localhost:8000 访问RootPulseOS系统

4. **停止容器**
   ```bash
   powershell -ExecutionPolicy RemoteSigned -File "E:\X54\github\Meta-CreationPower\05_Open_source_ProjectRepository\AIAgentFramework\EnlightenmentLighthouseRuntime\elr.ps1" stop-container --id rootpulseos-container
   ```

## 安装与运行

### 环境要求

- Python 3.8+
- Go 1.16+ (可选，用于高性能组件)
- Docker (用于容器化运行)
- CMake (可选，用于构建C/C++组件)
- ELR (Enlightenment Lighthouse Runtime，用于ELR容器运行)

### 安装步骤

1. **克隆项目仓库**
   ```bash
   git clone https://github.com/Meta-CreationPower/RootPulseOS.git
   cd RootPulseOS
   ```

2. **安装依赖**
   ```bash
   pip install -r requirements.txt
   ```

3. **运行系统**
   - **直接运行**
     ```bash
     python main.py
     ```

   - **在Docker容器中运行**
     1. **构建镜像**
        ```bash
        docker build -t rootpulseos .
        ```
     2. **运行容器**
        ```bash
        docker-compose up -d
        ```
     3. **访问系统**
        启动后，可通过 http://localhost:8000 访问RootPulseOS系统
     4. **停止容器**
        ```bash
        docker-compose down
        ```

   - **在ELR容器中运行**
     ```bash
     powershell -ExecutionPolicy RemoteSigned -File scripts\run_in_elr_simple.ps1
     ```

## 贡献指南

欢迎对RootPulseOS项目的贡献！请参考以下流程：

1. Fork项目仓库
2. 创建特性分支
3. 提交代码
4. 运行测试
5. 提交Pull Request

## 许可证

本项目采用MIT许可证。

## 联系方式

- 项目主页: https://github.com/Meta-CreationPower/RootPulseOS
- 邮箱: contact@rootpulseos.org
- 社区: https://community.rootpulseos.org

---

© 2026 启蒙灯塔起源团队 - 代码织梦者