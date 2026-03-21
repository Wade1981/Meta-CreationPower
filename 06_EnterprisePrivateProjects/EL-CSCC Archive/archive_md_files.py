#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
脚本功能：遍历Meta-CreationPower文件夹及其子文件夹，为所有md文件在EL-CSCC Archive中创建档案记录
"""

import os
import json

# 工作目录
working_dir = "e:\\X54\\github\\Meta-CreationPower\\06_EnterprisePrivateProjects\\EnlightenmentLighthouseOriginTeam\\EL-CSCC Archive"
os.chdir(working_dir)

# 遍历的根目录
root_dir = "e:\\X54\\github\\Meta-CreationPower"

# 碳硅伙伴信息
silicon_id = "SING-2E8A4C1D9F5B3A07"
silicon_name = "代码织梦者"
carbon_id = "X54SIR-X541981032043540001"
carbon_name = "X54先生"

# 查找所有md文件
print("正在查找所有md文件...")
md_files = []
for root, dirs, files in os.walk(root_dir):
    for file in files:
        if file.endswith('.md'):
            full_path = os.path.join(root, file)
            md_files.append(full_path)

print(f"找到 {len(md_files)} 个md文件")

# 构建档案数据列表
archives_data = []
for file_path in md_files:
    # 提取文件名
    file_name = os.path.basename(file_path)
    # 提取目录名
    dir_name = os.path.dirname(file_path)
    
    # 确定档案关联性
    relevance = ""
    if "01_EnlightenmentLighthouseOriginTeam" in dir_name or "07_ExclusiveFolder_CodeWeaver" in dir_name:
        relevance = "团队文档"
    elif "02_PublicCollaboration" in dir_name or "CarbonSiliconSynergyTh" in dir_name:
        relevance = "公共协作"
    elif "03_CodeAndAchievements" in dir_name or "CodeRepository" in dir_name:
        relevance = "代码与成就"
    elif "04_RootConfigFiles" in dir_name:
        relevance = "配置文件"
    elif "05_Open_source_ProjectRepository" in dir_name or "AIAgentFramework" in dir_name or "SiliconPartn" in dir_name:
        relevance = "开源项目"
    elif "06_EnterprisePrivateProjects" in dir_name or "EnlightenmentLighthouseOriginTeam" in dir_name:
        relevance = "企业私有项目"
    elif "docs" in dir_name:
        relevance = "文档"
    else:
        relevance = "其他"
    
    # 构建档案数据
    archive_data = {
        "address": file_path,
        "type": "文档",
        "silicon_id": silicon_id,
        "silicon_name": silicon_name,
        "carbon_id": carbon_id,
        "carbon_name": carbon_name,
        "brief": f"Markdown文档: {file_name}",
        "relevance": relevance
    }
    
    archives_data.append(archive_data)

# 导入模型
from el_cscc_archive_model import ELCSCCArchiveModel

# 初始化模型
model = ELCSCCArchiveModel()

# 批量添加档案
print("正在批量添加档案...")
success_count = model.batch_add_archives(archives_data)
print(f"成功添加 {success_count} 个档案")

# 列出所有档案
print("====================================")
print("所有md文件档案创建完成！")
print("====================================")
print("正在列出所有档案...")
list_result = model.predict("list:")
print(list_result)
