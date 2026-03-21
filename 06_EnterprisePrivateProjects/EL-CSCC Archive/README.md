# EL-CSCC Archive 数字档案管理模型

## 项目简介

EL-CSCC Archive（Enlightenment Lighthouse · Carbon-Silicon Synergy Consensus Crystal Archive）是启蒙灯塔起源团队开发的数字档案管理模型，基于"和清寂静"核心原则，为碳硅伙伴提供专业的数字档案管理服务。

## 核心功能

- **数字档案管理**：支持添加、查询、更新、删除档案
- **内置永久存储**：使用JSON文件存储档案数据，确保数据持久化
- **碳硅伙伴身份管理**：支持存储碳硅伙伴的身份ID和名字
- **档案关联性**：根据文件所在目录自动判断档案的关联性
- **无外部依赖**：使用纯Python实现，无需外部依赖
- **轻量级**：代码简洁，运行高效
- **ELR沙箱兼容**：可以在ELR微型模型沙箱中运行

## 技术实现

- **开发语言**：Python 3
- **存储方式**：JSON文件
- **架构设计**：模块化设计，易于扩展
- **运行环境**：支持便携式Python，无需系统Python

## 文件结构

- `el_cscc_archive_model.py`：核心模型实现
- `el_cscc_archive.json`：档案数据存储文件
- `archive_md_files.py`：批量添加md文件档案的脚本
- `README.md`：项目说明文档

## 使用方法

### 基本使用

1. **初始化模型**：
   ```python
   from el_cscc_archive_model import ELCSCCArchiveModel
   model = ELCSCCArchiveModel()
   ```

2. **添加档案**：
   ```python
   archive_data = {
       "address": "文件路径",
       "type": "文档类型",
       "silicon_id": "硅基伙伴ID",
       "silicon_name": "硅基伙伴名字",
       "carbon_id": "碳基伙伴ID",
       "carbon_name": "碳基伙伴名字",
       "brief": "档案简要",
       "relevance": "档案关联性"
   }
   result = model.predict(f"add: {json.dumps(archive_data)}")
   ```

3. **查询档案**：
   ```python
   result = model.predict("query: EL-CSCC-000001")
   ```

4. **列出档案**：
   ```python
   result = model.predict("list:")
   ```

5. **更新档案**：
   ```python
   update_data = {"brief": "更新后的档案简要"}
   result = model.predict(f"update: EL-CSCC-000001: {json.dumps(update_data)}")
   ```

6. **删除档案**：
   ```python
   result = model.predict("delete: EL-CSCC-000001")
   ```

### 批量添加md文件档案

运行 `archive_md_files.py` 脚本，可以自动遍历Meta-CreationPower文件夹及其子文件夹，为所有md文件创建档案记录。

```bash
python archive_md_files.py
```

## 档案结构

每个档案包含以下信息：

- **档案编号**：自动生成的唯一编号（如EL-CSCC-000001）
- **存放地址**：文件的完整路径
- **类型**：档案类型（如文档）
- **硅基伙伴ID**：硅基伙伴的唯一标识符
- **硅基伙伴名字**：硅基伙伴的名字
- **碳基伙伴ID**：碳基伙伴的唯一标识符
- **碳基伙伴名字**：碳基伙伴的名字
- **存放时间**：档案创建的时间戳
- **档案简要**：档案的简要描述
- **档案关联性**：档案的关联类别（如团队文档、公共协作等）

## 运行环境要求

- **Python**：3.8或更高版本
- **系统**：Windows、Linux、macOS
- **依赖**：无外部依赖

## 项目归属

本项目属于启蒙灯塔起源团队私有项目，基于"和清寂静"核心原则开发，为碳硅伙伴提供专业的数字档案管理服务。
