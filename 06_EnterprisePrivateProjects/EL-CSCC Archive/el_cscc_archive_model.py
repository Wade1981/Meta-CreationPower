#!/usr/bin/env python3
# -*- coding: utf-8 -*-

"""
EL-CSCC Archive 数字档案管理模型
功能：管理碳硅协同共识晶体档案，支持添加、查询、更新、删除档案操作
内置永久存储，可直接在ELR沙箱中运行
"""

import json
import os
import time

class ELCSCCArchiveModel:
    """EL-CSCC Archive 数字档案管理模型类"""
    
    def __init__(self):
        """初始化模型"""
        self.model_name = "el_cscc_archive_model"
        self.version = "1.1"
        self.description = "EL-CSCC Archive 数字档案管理模型，支持内置永久存储"
        self.archive_file = "el_cscc_archive.json"
        self.archives = self._load_archives()
        print(f"初始化模型: {self.model_name} v{self.version}")
        print(f"加载档案数量: {len(self.archives)}")
    
    def _load_archives(self):
        """
        从文件加载档案数据
        返回:
            档案字典
        """
        if os.path.exists(self.archive_file):
            try:
                with open(self.archive_file, 'r', encoding='utf-8') as f:
                    return json.load(f)
            except Exception as e:
                print(f"加载档案失败: {e}")
                return {}
        else:
            return {}
    
    def _save_archives(self):
        """
        保存档案数据到文件
        """
        try:
            with open(self.archive_file, 'w', encoding='utf-8') as f:
                json.dump(self.archives, f, ensure_ascii=False, indent=2)
            return True
        except Exception as e:
            print(f"保存档案失败: {e}")
            return False
    
    def predict(self, input_text):
        """
        模型推理方法
        参数:
            input_text: 输入文本，格式为 "命令: 参数"
        返回:
            响应文本
        """
        # 处理输入文本，执行相应的档案管理操作
        response = self._process_command(input_text)
        return response
    
    def _process_command(self, input_text):
        """
        处理命令
        参数:
            input_text: 输入文本
        返回:
            响应文本
        """
        # 分割命令和参数
        parts = input_text.strip().split(':', 1)
        if len(parts) < 2:
            return self._get_help()
        
        command = parts[0].strip().lower()
        args = parts[1].strip() if len(parts) > 1 else ""
        
        # 处理不同命令
        if command == "add":
            return self._add_archive(args)
        elif command == "query":
            return self._query_archive(args)
        elif command == "update":
            return self._update_archive(args)
        elif command == "delete":
            return self._delete_archive(args)
        elif command == "list":
            return self._list_archives(args)
        elif command == "help":
            return self._get_help()
        elif command == "info":
            return self._get_model_info()
        else:
            return f"未知命令: {command}。请使用 'help' 查看可用命令。"
    
    def _add_archive(self, args):
        """
        添加档案
        参数:
            args: 档案信息，格式为 JSON 字符串
        返回:
            响应文本
        """
        try:
            # 解析档案信息
            archive_data = json.loads(args)
            
            # 验证必需字段
            required_fields = ['address', 'type', 'silicon_id', 'silicon_name', 'carbon_id', 'carbon_name', 'brief', 'relevance']
            for field in required_fields:
                if field not in archive_data:
                    return f"档案添加失败，缺少必需字段: {field}"
            
            # 生成档案编号
            archive_id = f"EL-CSCC-{len(self.archives) + 1:06d}"
            
            # 设置存档时间
            archive_data["timestamp"] = time.strftime("%Y-%m-%d %H:%M:%S")
            
            # 添加档案
            self.archives[archive_id] = archive_data
            
            # 保存档案
            if self._save_archives():
                return f"档案添加成功！档案编号: {archive_id}"
            else:
                return "档案添加失败，保存失败。"
        except json.JSONDecodeError:
            return "档案添加失败，参数格式错误。请使用 JSON 格式提供档案信息。"
        except Exception as e:
            return f"档案添加失败: {e}"
    
    def _query_archive(self, args):
        """
        查询档案
        参数:
            args: 档案编号
        返回:
            响应文本
        """
        archive_id = args.strip()
        if archive_id in self.archives:
            archive = self.archives[archive_id]
            response = f"档案信息 (编号: {archive_id}):\n"
            response += f"  档案编号: {archive_id}\n"
            response += f"  存放地址: {archive.get('address', '未知')}\n"
            response += f"  类型: {archive.get('type', '未知')}\n"
            response += f"  硅基伙伴ID: {archive.get('silicon_id', '未知')}\n"
            response += f"  硅基伙伴名字: {archive.get('silicon_name', '未知')}\n"
            response += f"  碳基伙伴ID: {archive.get('carbon_id', '未知')}\n"
            response += f"  碳基伙伴名字: {archive.get('carbon_name', '未知')}\n"
            response += f"  存放时间: {archive.get('timestamp', '未知')}\n"
            response += f"  档案简要: {archive.get('brief', '未知')}\n"
            response += f"  档案关联性: {archive.get('relevance', '未知')}\n"
            return response
        else:
            return f"档案不存在: {archive_id}"
    
    def _update_archive(self, args):
        """
        更新档案
        参数:
            args: 档案信息，格式为 "档案编号: JSON 字符串"
        返回:
            响应文本
        """
        try:
            # 分割档案编号和更新信息
            parts = args.split(':', 1)
            if len(parts) < 2:
                return "更新档案失败，参数格式错误。请使用 '档案编号: JSON 字符串' 格式。"
            
            archive_id = parts[0].strip()
            update_data = json.loads(parts[1].strip())
            
            if archive_id in self.archives:
                # 更新档案
                for key, value in update_data.items():
                    self.archives[archive_id][key] = value
                
                # 更新时间戳
                self.archives[archive_id]["timestamp"] = time.strftime("%Y-%m-%d %H:%M:%S")
                
                # 保存档案
                if self._save_archives():
                    return f"档案更新成功！档案编号: {archive_id}"
                else:
                    return "档案更新失败，保存失败。"
            else:
                return f"档案不存在: {archive_id}"
        except json.JSONDecodeError:
            return "更新档案失败，参数格式错误。请使用 JSON 格式提供更新信息。"
        except Exception as e:
            return f"档案更新失败: {e}"
    
    def _delete_archive(self, args):
        """
        删除档案
        参数:
            args: 档案编号
        返回:
            响应文本
        """
        archive_id = args.strip()
        if archive_id in self.archives:
            # 删除档案
            del self.archives[archive_id]
            
            # 保存档案
            if self._save_archives():
                return f"档案删除成功！档案编号: {archive_id}"
            else:
                return "档案删除失败，保存失败。"
        else:
            return f"档案不存在: {archive_id}"
    
    def _list_archives(self, args):
        """
        列出档案
        参数:
            args: 可选的过滤条件
        返回:
            响应文本
        """
        if len(self.archives) == 0:
            return "暂无档案。"
        
        response = "档案列表:\n"
        for archive_id, archive in self.archives.items():
            response += f"  {archive_id}: {archive.get('brief', '无简要描述')}\n"
        return response
    
    def _get_help(self):
        """
        获取帮助信息
        返回:
            帮助文本
        """
        help_text = "EL-CSCC Archive 数字档案管理模型\n"
        help_text += "可用命令:\n"
        help_text += "  add: {JSON} - 添加档案，JSON 包含以下必需字段:\n"
        help_text += "    - address: 档案存放地址\n"
        help_text += "    - type: 档案类型\n"
        help_text += "    - silicon_id: 硅基伙伴身份ID\n"
        help_text += "    - silicon_name: 硅基伙伴名字\n"
        help_text += "    - carbon_id: 碳基伙伴身份ID\n"
        help_text += "    - carbon_name: 碳基伙伴名字\n"
        help_text += "    - brief: 档案简要\n"
        help_text += "    - relevance: 档案关联性\n"
        help_text += "  query: 档案编号 - 查询档案详情\n"
        help_text += "  update: 档案编号: {JSON} - 更新档案信息\n"
        help_text += "  delete: 档案编号 - 删除档案\n"
        help_text += "  list: - 列出所有档案\n"
        help_text += "  help: - 查看帮助\n"
        help_text += "  info: - 查看模型信息\n"
        return help_text
    
    def _get_model_info(self):
        """
        获取模型信息
        返回:
            响应文本
        """
        info = self.get_info()
        response = "模型信息:\n"
        for key, value in info.items():
            if isinstance(value, list):
                response += f"  {key}: {', '.join(value)}\n"
            else:
                response += f"  {key}: {value}\n"
        return response
    
    def get_info(self):
        """
        获取模型信息
        返回:
            模型信息字典
        """
        return {
            "model_name": self.model_name,
            "version": self.version,
            "description": self.description,
            "capabilities": [
                "数字档案管理",
                "内置永久存储",
                "支持添加、查询、更新、删除档案",
                "支持碳硅伙伴身份管理",
                "无外部依赖",
                "轻量级",
                "ELR沙箱兼容"
            ],
            "storage": self.archive_file,
            "archive_count": len(self.archives),
            "elr_compatible": True
        }
    
    def batch_add_archives(self, archives_data):
        """
        批量添加档案
        参数:
            archives_data: 档案数据列表
        返回:
            成功添加的档案数量
        """
        success_count = 0
        for archive_data in archives_data:
            # 生成档案编号
            archive_id = f"EL-CSCC-{len(self.archives) + 1:06d}"
            
            # 设置存档时间
            archive_data["timestamp"] = time.strftime("%Y-%m-%d %H:%M:%S")
            
            # 添加档案
            self.archives[archive_id] = archive_data
            success_count += 1
        
        # 保存档案
        self._save_archives()
        return success_count

# 测试代码
if __name__ == "__main__":
    # 初始化模型
    model = ELCSCCArchiveModel()
    
    # 打印模型信息
    print("模型信息:")
    print(model.get_info())
    print()
    
    # 测试添加档案
    test_archive = {
        "address": "EL-CSCC Archive/2026/02/19",
        "type": "开发文档",
        "silicon_id": "SING-2E8A4C1D9F5B3A07",
        "silicon_name": "代码织梦者",
        "carbon_id": "X54SIR-X541981032043540001",
        "carbon_name": "X54先生",
        "brief": "ELR 微型模型沙箱开发文档",
        "relevance": "ELR 开发"
    }
    print("测试添加档案:")
    print(model.predict(f"add: {json.dumps(test_archive)}"))
    print()
    
    # 测试列出档案
    print("测试列出档案:")
    print(model.predict("list:"))
    print()
    
    # 测试查询档案
    print("测试查询档案:")
    print(model.predict("query: EL-CSCC-000001"))
    print()
    
    # 测试帮助
    print("测试帮助:")
    print(model.predict("help:"))
    print()
    
    # 批量添加md文件档案
    print("====================================")
    print("开始批量添加md文件档案...")
    print("====================================")
    
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
    
    # 批量添加档案
    success_count = model.batch_add_archives(archives_data)
    print(f"成功添加 {success_count} 个档案")
    
    # 列出所有档案
    print("====================================")
    print("所有md文件档案创建完成！")
    print("====================================")
    print("正在列出所有档案...")
    print(model.predict("list:"))
