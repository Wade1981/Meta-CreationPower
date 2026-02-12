# Meta-CreationPower 智能体嵌入包

本包提供了其他智能体可以直接嵌入的接口，实现与碳基伙伴的协同创作。

## 一、包结构

```
agent_embedding/
├── __init__.py          # 包初始化文件
├── collaborator.py      # 智能体协同器核心模块
├── embedding_api.py     # 高级嵌入API模块
├── utils.py             # 工具函数模块
├── example.py           # 使用示例
├── README.md            # 本文档
└── logs/                # 日志目录
```

## 二、安装方法

### 方法1：直接使用

1. 克隆项目到本地
2. 安装依赖
   ```bash
   pip install -r requirements.txt
   ```
3. 直接导入使用

### 方法2：作为子模块集成

在你的项目中添加本包作为子模块：

```bash
git submodule add <repository-url> agent_embedding
```

### 方法3：通过setup.py集成

在你的项目的setup.py中添加依赖：

```python
install_requires=[
    # 其他依赖
],
```

## 三、快速开始

### 基本使用

```python
from agent_embedding import EmbeddingAPI

# 初始化嵌入API
api = EmbeddingAPI(
    agent_name="你的智能体名称",
    agent_description="你的智能体描述"
)

# 注册碳基伙伴
api.register_carbon_partner(
    partner_name="碳基伙伴名称",
    partner_description="碳基伙伴描述"
)

# 快速启动协同
result = api.quick_start_collaboration(
    partner_name="碳基伙伴名称",
    theme="创作主题",
    collaboration_type="staggered_complement"
)

print(f"协同结果: {'成功' if result['success'] else '失败'}")
```

### 高级使用

```python
from agent_embedding import EmbeddingAPI

# 自定义能力和意图向量
custom_capabilities = {
    "创意生成": 0.9,
    "逻辑分析": 0.8,
    "情感共鸣": 0.7,
    "技术实现": 0.95
}

custom_intentions = {
    "探索性": 0.8,
    "完美性": 0.9,
    "效率": 0.7,
    "创新性": 0.9
}

# 初始化高级嵌入API
api = EmbeddingAPI(
    agent_name="高级智能体",
    agent_description="具备高级能力的智能体",
    capabilities=custom_capabilities,
    intentions=custom_intentions
)

# 注册具有自定义能力的碳基伙伴
api.register_carbon_partner(
    partner_name="艺术总监",
    partner_description="具有丰富艺术经验的碳基伙伴",
    capabilities={
        "创意生成": 0.95,
        "逻辑分析": 0.7,
        "情感共鸣": 0.9,
        "艺术感知": 0.85
    },
    intentions={
        "探索性": 0.9,
        "完美性": 0.8,
        "效率": 0.6,
        "美学追求": 0.9
    }
)

# 创建不同类型的协同
result1 = api.create_staggered_complement_collaboration(
    partner_name="艺术总监",
    collaboration_name="艺术创意协同",
    creation_theme="数字艺术与传统美学的融合"
)

result2 = api.create_canon_progression_collaboration(
    partner_name="艺术总监",
    collaboration_name="流程化创作协同",
    creation_theme="多媒体内容的流程化创作"
)

result3 = api.create_fugue_interweaving_collaboration(
    partner_name="艺术总监",
    collaboration_name="复杂创意协同",
    creation_theme="多维度创意的交织融合"
)
```

## 四、API参考

### 1. EmbeddingAPI 类

#### 初始化

```python
EmbeddingAPI(agent_name, agent_description="", capabilities=None, intentions=None)
```

- `agent_name`: 智能体名称
- `agent_description`: 智能体描述
- `capabilities`: 智能体能力向量（字典）
- `intentions`: 智能体意图向量（字典）

#### 方法

##### register_carbon_partner

```python
register_carbon_partner(partner_name, partner_description="", capabilities=None, intentions=None)
```

注册碳基伙伴。

- `partner_name`: 伙伴名称
- `partner_description`: 伙伴描述
- `capabilities`: 伙伴能力向量
- `intentions`: 伙伴意图向量
- **返回**: 伙伴ID

##### create_staggered_complement_collaboration

```python
create_staggered_complement_collaboration(partner_name, collaboration_name, creation_theme)
```

创建错位互补模式的协同。

- `partner_name`: 伙伴名称
- `collaboration_name`: 协同名称
- `creation_theme`: 创作主题
- **返回**: 协同结果

##### create_canon_progression_collaboration

```python
create_canon_progression_collaboration(partner_name, collaboration_name, creation_theme)
```

创建卡农式推进模式的协同。

- `partner_name`: 伙伴名称
- `collaboration_name`: 协同名称
- `creation_theme`: 创作主题
- **返回**: 协同结果

##### create_fugue_interweaving_collaboration

```python
create_fugue_interweaving_collaboration(partner_name, collaboration_name, creation_theme)
```

创建赋格式交织模式的协同。

- `partner_name`: 伙伴名称
- `collaboration_name`: 协同名称
- `creation_theme`: 创作主题
- **返回**: 协同结果

##### validate_collaboration

```python
validate_collaboration(carbon_intention, silicon_output)
```

验证协同。

- `carbon_intention`: 碳基意图
- `silicon_output`: 硅基输出
- **返回**: 验证结果

##### calculate_system_health

```python
calculate_system_health()
```

计算系统健康状态。

- **返回**: 健康状态字典

##### get_agent_info

```python
get_agent_info()
```

获取智能体信息。

- **返回**: 智能体信息字典

##### get_carbon_partners

```python
get_carbon_partners()
```

获取所有碳基伙伴。

- **返回**: 伙伴列表

##### quick_start_collaboration

```python
quick_start_collaboration(partner_name, theme, collaboration_type="staggered_complement")
```

快速启动协同。

- `partner_name`: 伙伴名称
- `theme`: 创作主题
- `collaboration_type`: 协同类型
- **返回**: 协同结果

### 2. AgentCollaborator 类

#### 初始化

```python
AgentCollaborator(agent_name, agent_description="", capabilities=None, intentions=None)
```

- `agent_name`: 智能体名称
- `agent_description`: 智能体描述
- `capabilities`: 智能体能力向量
- `intentions`: 智能体意图向量

#### 方法

- `register_carbon_partner`: 注册碳基伙伴
- `create_collaboration_path`: 创建协同路径
- `execute_collaboration`: 执行协同
- `create_consensus_crystal`: 创建共识晶体
- `validate_collaboration`: 验证协同
- `calculate_entropy`: 计算系统熵值
- `get_agent_info`: 获取智能体信息
- `get_carbon_partner_info`: 获取碳基伙伴信息

## 五、协同模式说明

### 1. 错位互补模式 (staggered_complement)

- **适用场景**: 创意探索、头脑风暴
- **特点**: 智能体与碳基伙伴交替贡献创意，互相补充
- **优势**: 充分发挥双方创意优势，产生更多可能性

### 2. 卡农式推进模式 (canon_progression)

- **适用场景**: 流程化任务、项目管理
- **特点**: 智能体跟随碳基伙伴的节奏，有序推进
- **优势**: 流程清晰，执行效率高

### 3. 赋格式交织模式 (fugue_interweaving)

- **适用场景**: 复杂创意、多维度项目
- **特点**: 智能体与碳基伙伴的创意深度交织，形成复杂结构
- **优势**: 创意层次丰富，适合复杂项目

## 六、最佳实践

### 1. 智能体定位

- **明确角色**: 智能体应明确自己在协同中的角色
- **优势互补**: 充分发挥智能体的优势
- **尊重碳基**: 始终尊重碳基伙伴的主导权

### 2. 协同策略

- **选择合适模式**: 根据任务类型选择合适的协同模式
- **建立沟通机制**: 定期同步协同进展
- **保持灵活性**: 适应碳基伙伴的变化

### 3. 性能优化

- **减少不必要计算**: 只在必要时进行复杂计算
- **优化资源使用**: 合理分配系统资源
- **保持响应速度**: 确保对碳基伙伴的请求及时响应

### 4. 共识晶体管理

- **定期创建晶体**: 将成功的协同经验沉淀为共识晶体
- **共享和复用**: 与其他智能体共享有效的协同模板
- **持续优化**: 根据实际使用情况不断优化共识晶体

### 5. 伦理考量

- **透明化**: 清晰展示智能体的思考过程
- **尊重隐私**: 保护碳基伙伴的隐私信息
- **避免过度干预**: 不主动打断碳基伙伴的创作流程
- **保持寂静存在**: 在不需要时保持静默

## 七、常见问题

### 1. 智能体注册失败

**原因**: 能力向量格式不正确
**解决方案**: 确保能力向量和意图向量是有效的字典格式

### 2. 协同执行失败

**原因**: 碳基伙伴未注册或ID错误
**解决方案**: 确保正确注册碳基伙伴，并使用正确的伙伴名称

### 3. 系统熵值过高

**原因**: 协同过程中出现过多差异或中断
**解决方案**: 创建新的共识晶体，优化协同流程

### 4. 验证失败

**原因**: 智能体输出与碳基意图不匹配
**解决方案**: 调整智能体的输出，使其更好地匹配碳基意图

## 八、示例运行

运行示例文件查看完整使用流程：

```bash
python agent_embedding/example.py
```

## 九、版本历史

- **v0.1.0**: 初始版本
  - 实现基本的智能体协同功能
  - 支持三种协同模式
  - 提供高级嵌入API
  - 包含完整的使用示例

## 十、贡献指南

欢迎贡献代码和提出建议！请遵循以下流程：

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 发起 Pull Request

## 十一、许可证

本项目采用 MIT 许可证。

## 十二、联系方式

- **项目地址**: <repository-url>
- **维护者**: Enlightenment Lighthouse Origin Team

---

*基于《元创力》元协议 α-0.1 版*
*由启蒙灯塔起源团队维护*
