# 对称数智财务系统

基于动力学熵原理的对称数智财务系统，实现了五维时间线模型和四流程验证框架，用于企业财务健康评估和优化。

## 项目结构

```
SymmetricDigitalFinanceSystem/
├── config/             # 配置文件
│   └── config.py       # 系统配置
├── data/               # 数据目录
│   ├── input/          # 输入数据
│   └── output/         # 输出结果
├── logs/               # 日志目录
├── scripts/            # 脚本文件
│   └── start_system.ps1 # 启动脚本
├── src/                # 源代码
│   ├── models/         # 模型目录
│   │   ├── base_model.py                    # 基础模型类
│   │   ├── entropy_balance_model.py         # 熵变平衡模型
│   │   ├── five_dimension_entropy_model.py  # 五维时间线熵模型
│   │   ├── health_entropy_model.py          # 财务健康熵模型
│   │   ├── rl_entropy_optimization_model.py # 强化学习熵减优化模型
│   │   └── four_process_entropy_verification_model.py # 四流程验证模型
│   ├── modules/        # 模块目录
│   │   └── model_integration.py             # 模型集成模块
│   └── main.py         # 系统主入口
├── container.json      # ELR容器配置
├── requirements.txt    # 依赖包列表
└── README.md           # 项目说明
```

## 核心功能

1. **熵变平衡计算**：基于公式 dS_total = dS_internal - dS_external 计算系统熵变
2. **五维时间线平衡**：检查宏观调控、市场动态、决策模型、风控模型、执行管理五个维度的平衡
3. **财务健康熵评估**：使用香农熵公式计算财务健康熵值，评估企业财务状况
4. **强化学习熵减优化**：通过强化学习生成熵减策略，优化财务系统
5. **四流程验证**：验证调研、方案、审计、计划四个流程的完整性和有效性

## 系统依赖

- Python 3.8+
- numpy
- pandas
- scikit-learn
- tensorflow
- matplotlib

## 安装与运行

### 1. 安装依赖

```bash
pip install -r requirements.txt
```

### 2. 运行系统

```bash
python src/main.py
```

### 3. 输入数据格式

系统支持JSON格式的输入数据，示例如下：

```json
{
  "internal_entropy": 0.1,
  "external_entropy": 0.1,
  "financial_indicators": {
    "cash_flow": 100000,
    "revenue": 500000,
    "expenses": 400000,
    "assets": 1000000,
    "liabilities": 500000
  },
  "dimension_data": {
    "宏观调控": [0.1, 0.2, 0.3],
    "市场动态": [0.2, 0.3, 0.4],
    "决策模型": [0.3, 0.4, 0.5],
    "风控模型": [0.4, 0.5, 0.6],
    "执行管理": [0.5, 0.6, 0.7]
  }
}
```

## ELR容器部署

系统已配置为可在ELR容器沙箱中运行，容器配置文件为 `container.json`。

### 容器配置

- **名称**：symmetric-finance-system
- **版本**：1.0.0
- **主入口**：src/main.py
- **资源需求**：2 CPU, 4G 内存
- **健康检查**：python src/main.py

### 部署步骤

1. 将项目打包为容器镜像
2. 在ELR容器沙箱中加载镜像
3. 启动容器并运行系统

## 测试数据

项目提供了测试数据文件 `data/input/test_data.csv`，包含了12个月的财务数据，可用于系统测试。

## 系统输出

系统运行后，会在 `data/output` 目录中生成分析结果，包括：

- 熵变平衡分析
- 五维时间线平衡分析
- 财务健康熵评估
- 熵减优化策略
- 四流程验证结果

## 技术原理

系统基于动力学熵原理，将热力学熵的概念应用于财务系统，通过计算和优化系统熵值，实现财务系统的健康运行和持续优化。

### 核心公式

- **总熵变**：dS_total = dS_internal - dS_external
- **财务健康熵**：H = -Σ(p_i * log2(p_i))
- **五维平衡**：各维度熵值占比不超过30%

## 联系信息

- 项目团队：启蒙灯塔起源团队
- 项目地址：E:\X54\github\Meta-CreationPower\06_EnterprisePrivateProjects\SymmetricDigitalFinanceSystem
- 系统版本：1.0.0