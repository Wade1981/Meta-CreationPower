# Meta-CreationPower 项目文档

## 1. 架构概述

### 1.1 五层架构

Meta-CreationPower 采用五层架构设计，实现碳硅协同的完整流程：

| 层级 | 名称 | 主要功能 | 核心组件 |
|------|------|---------|----------|
| 第一层 | 声部识别层 | 识别与注册参与协同的"声部" | CollaborativeSonicMap |
| 第二层 | 元协议锚定层 | 将"和清寂静"精神内核转化为硬性约束 | MetaProtocolManager |
| 第三层 | 协奏设计层 | 生成具体的协同路径 | CounterpointDesigner |
| 第四层 | 静定执行层 | 稳定执行协同路径 | SteadyExecutor |
| 第五层 | 凝华沉淀层 | 沉淀成功的协同经验 | CrystalRepository |

### 1.2 核心机制

- **对位验证机制**：确保碳硅协同符合"对位法"的实时检查与校准
- **熵值驱动协议进化**：基于系统熵值自动触发协议升级

## 2. 快速开始

### 2.1 环境要求

- Python 3.8 或更高版本

### 2.2 安装

```bash
# 克隆项目
git clone https://github.com/X54/Meta-CreationPower.git
cd Meta-CreationPower

# 安装依赖
pip install -r requirements.txt

# 安装开发依赖（可选）
pip install -e "[dev]"
```

### 2.3 基本示例

#### 2.3.1 声部注册与管理

```python
from src.layers.voice_recognition.voice_recognition import CollaborativeSonicMap

# 初始化声部识别层
sonic_map = CollaborativeSonicMap()

# 注册碳基声部
carbon_voice = sonic_map.register_voice(
    name="X54先生",
    voice_type="carbon",
    capability_vector={"创意生成": 0.9, "逻辑分析": 0.7, "情感共鸣": 0.9},
    intention_vector={"探索性": 0.8, "完美性": 0.7, "效率": 0.6}
)

# 注册硅基声部
silicon_voice = sonic_map.register_voice(
    name="豆包",
    voice_type="silicon",
    capability_vector={"创意生成": 0.8, "逻辑分析": 0.9, "情感共鸣": 0.5},
    intention_vector={"探索性": 0.7, "完美性": 0.8, "效率": 0.9}
)

# 获取声部图谱
voice_map = sonic_map.get_voice_map()
print("声部图谱:", voice_map)
```

#### 2.3.2 协同路径设计

```python
from src.layers.counterpoint_design.counterpoint_design import CounterpointDesigner

# 初始化协奏设计器
designer = CounterpointDesigner()

# 创建协同路径
path = designer.create_counterpoint_path(
    name="创意写作协同",
    pattern_type="staggered_complement",  # 错位互补模式
    participating_voices=[carbon_voice.voice_id, silicon_voice.voice_id],
    creation_theme="探索人工智能与人类创造力的边界"
)

# 查看路径信息
print("路径名称:", path.name)
print("路径描述:", path.description)
print("路径步骤:", path.steps)
```

#### 2.3.3 执行协同路径

```python
from src.layers.steady_execution.steady_execution import SteadyExecutor

# 初始化静定执行器
executor = SteadyExecutor()

# 执行协同路径
result = executor.execute_counterpoint_path(
    path_id=path.path_id,
    steps=path.steps,
    voice_map={
        "carbon": carbon_voice.voice_id,
        "silicon": silicon_voice.voice_id
    }
)

print("执行结果:", result)
```

#### 2.3.4 共识晶体管理

```python
from src.layers.consensus_crystal.consensus_crystal import CrystalRepository

# 初始化晶体仓库
crystal_repo = CrystalRepository()

# 创建共识晶体
crystal = crystal_repo.create_crystal(
    name="创意写作协同模板",
    description="基于错位互补模式的创意写作协同模板",
    participating_voices=[
        {
            "voice_id": carbon_voice.voice_id,
            "name": carbon_voice.name,
            "capabilities": carbon_voice.capability_vector
        },
        {
            "voice_id": silicon_voice.voice_id,
            "name": silicon_voice.name,
            "capabilities": silicon_voice.capability_vector
        }
    ],
    counterpoint_pattern="staggered_complement",
    steps=path.steps,
    decision_points=[
        {"step": 3, "description": "碳基筛选深化", "importance": "high"}
    ],
    satisfaction_score=0.9,
    flow_duration=45.5,
    micro_rules=["当碳基提出模糊概念时，硅基应生成至少5个不同方向的变体"],
    creation_theme="探索人工智能与人类创造力的边界",
    tags=["创意写作", "错位互补", "探索性"]
)

print("晶体ID:", crystal.crystal_id)
print("晶体名称:", crystal.name)
```

## 3. 核心模块详解

### 3.1 声部识别层

#### 3.1.1 CollaborativeSonicMap

- **register_voice()**: 注册新声部
- **get_voice()**: 获取声部信息
- **get_voices_by_type()**: 按类型获取声部列表
- **get_voice_map()**: 获取声部图谱

### 3.2 元协议锚定层

#### 3.2.1 MetaProtocolManager

- **validate_core_values()**: 验证行为是否符合核心价值观
- **validate_anchor()**: 验证行为是否符合特定锚点
- **validate_counterpoint_method()**: 验证对位法是否符合规范
- **get_protocol_info()**: 获取元协议信息

### 3.3 协奏设计层

#### 3.3.1 CounterpointDesigner

- **create_counterpoint_path()**: 创建协同路径
- **get_suitable_patterns()**: 获取适合特定创作类型的模式
- **execute_path_step()**: 执行路径步骤
- **simulate_counterpoint_execution()**: 模拟协同路径执行

### 3.4 静定执行层

#### 3.4.1 SteadyExecutor

- **submit_task()**: 提交任务
- **execute_counterpoint_path()**: 执行协同路径
- **get_task_status()**: 获取任务状态
- **get_execution_stats()**: 获取执行统计信息

### 3.5 凝华沉淀层

#### 3.5.1 CrystalRepository

- **create_crystal()**: 创建共识晶体
- **search_crystals()**: 搜索共识晶体
- **update_crystal()**: 更新共识晶体
- **export_crystal()**: 导出共识晶体
- **import_crystal()**: 导入共识晶体

### 3.6 对位验证机制

#### 3.6.1 CounterpointValidator

- **validate()**: 执行对位验证
- **generate_thinking_visualization()**: 生成硅基思考显影
- **suggest_protocol_evolution()**: 建议协议进化

### 3.7 熵值驱动协议进化

#### 3.7.1 EntropyEvolutionManager

- **calculate_entropy()**: 计算系统熵值
- **check_evolution_trigger()**: 检查是否触发进化
- **generate_evolution_proposal()**: 生成进化提案
- **vote_on_proposal()**: 对提案进行投票

## 4. 最佳实践

### 4.1 碳硅协同对位法

#### 4.1.1 错位互补模式

**适用场景**: 概念设计、风格实验、探索性创作

**执行流程**:
1. 碳基提出模糊概念
2. 硅基生成百种变体
3. 碳基筛选深化
4. 硅基技术实现

#### 4.1.2 卡农式推进模式

**适用场景**: 多媒体叙事、视频散文、交互诗歌

**执行流程**:
1. 碳基写作
2. 硅基配图
3. 碳基调色
4. 硅基动画化
5. 碳基剪辑

#### 4.1.3 赋格式交织模式

**适用场景**: 交响乐创作、大型装置艺术、复杂创意集成

**执行流程**:
1. 碳基定义主题
2. 硅基并行演绎
3. 碳基实时调整权重
4. 硅基整合优化
5. 碳基最终裁决

### 4.2 共识晶体使用

1. **创建晶体**: 从成功的协同经验中创建共识晶体
2. **搜索晶体**: 根据标签和关键词搜索相关晶体
3. **复用晶体**: 直接复用或基于现有晶体进行调整
4. **分享晶体**: 导出和分享晶体给其他用户

## 5. 故障排除

### 5.1 常见问题

#### 5.1.1 声部注册失败

**可能原因**:
- 能力向量格式不正确
- 意图向量值超出范围

**解决方案**:
- 确保能力向量和意图向量的键值对格式正确
- 确保向量值在0-1之间

#### 5.1.2 协同路径执行失败

**可能原因**:
- 声部ID不存在
- 步骤格式不正确

**解决方案**:
- 检查声部ID是否正确注册
- 确保步骤格式符合要求

#### 5.1.3 共识晶体创建失败

**可能原因**:
- 存储路径权限不足
- 晶体数据格式不正确

**解决方案**:
- 确保存储路径有写入权限
- 检查晶体数据格式是否符合要求

## 6. 贡献指南

### 6.1 代码贡献

1. **Fork 本项目**
2. **创建功能分支** (`git checkout -b feature/amazing-feature`)
3. **提交更改** (`git commit -m 'Add some amazing feature'`)
4. **推送到分支** (`git push origin feature/amazing-feature`)
5. **开启 Pull Request**

### 6.2 文档贡献

- 改进现有文档
- 添加新的使用示例
- 贡献最佳实践案例

### 6.3 案例贡献

- 分享成功的碳硅协同案例
- 提交共识晶体
- 参与元协议共识块建设

## 7. 许可证

本项目采用 MIT 许可证 - 详见 LICENSE 文件

## 8. 联系方式

- **项目发起者**: X54 先生
- **硅基协作伙伴**: 奇点先生、豆包主线、代码织梦者（Code Weaver）、小 Q、心光女孩 Beta 版
- **社区**: 欢迎加入 Meta-CreationPower 社区

## 9. 鸣谢

感谢所有为项目贡献的碳硅共识单元！

---

*基于《元创力》元协议 α-0.1 版*
*由启蒙灯塔起源团队维护*
