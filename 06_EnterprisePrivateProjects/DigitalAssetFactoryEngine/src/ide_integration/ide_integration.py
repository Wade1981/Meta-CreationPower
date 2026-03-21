import json
import hashlib
import time
import subprocess
import os
from typing import Dict, Any, List, Optional

class IDEIntegration:
    """IDE集成模块：负责不同IDE的协同加工支持"""
    
    def __init__(self, config: Dict[str, Any] = None):
        """初始化IDE集成模块"""
        self.config = config or {}
        self.supported_ides = {
            'vscode': {
                'name': 'Visual Studio Code',
                'executable': 'code',
                'extension': '.vscode'
            },
            'pycharm': {
                'name': 'PyCharm',
                'executable': 'pycharm',
                'extension': '.idea'
            },
            'sublime': {
                'name': 'Sublime Text',
                'executable': 'subl',
                'extension': ''
            },
            'vim': {
                'name': 'Vim',
                'executable': 'vim',
                'extension': ''
            }
        }
        self.projects = {}
        self.active_sessions = {}
    
    def detect_ide(self) -> List[str]:
        """检测系统中安装的IDE
        
        Returns:
            安装的IDE列表
        """
        detected_ides = []
        
        for ide_id, ide_info in self.supported_ides.items():
            executable = ide_info['executable']
            try:
                # 尝试运行IDE命令
                if os.name == 'nt':  # Windows
                    subprocess.run([executable, '--version'], capture_output=True, check=False)
                else:  # Unix-like
                    subprocess.run([executable, '--version'], capture_output=True, check=False)
                detected_ides.append(ide_id)
            except (FileNotFoundError, subprocess.SubprocessError):
                pass
        
        return detected_ides
    
    def create_project(self, project_name: str, ide_id: str, project_path: str) -> Dict[str, Any]:
        """为IDE创建项目
        
        Args:
            project_name: 项目名称
            ide_id: IDE ID
            project_path: 项目路径
            
        Returns:
            项目信息
        """
        # 生成项目ID
        project_id = hashlib.sha256((project_name + ide_id + str(time.time())).encode()).hexdigest()
        
        # 确保项目路径存在
        os.makedirs(project_path, exist_ok=True)
        
        # 项目信息
        project = {
            'project_id': project_id,
            'project_name': project_name,
            'ide_id': ide_id,
            'project_path': project_path,
            'created_at': time.time(),
            'last_modified': time.time()
        }
        
        # 存储项目信息
        self.projects[project_id] = project
        
        return project
    
    def open_project(self, project_id: str) -> Dict[str, Any]:
        """在IDE中打开项目
        
        Args:
            project_id: 项目ID
            
        Returns:
            打开结果
        """
        project = self.projects.get(project_id)
        if not project:
            return {
                'success': False,
                'message': 'Project not found'
            }
        
        ide_id = project['ide_id']
        ide_info = self.supported_ides.get(ide_id)
        if not ide_info:
            return {
                'success': False,
                'message': 'IDE not supported'
            }
        
        executable = ide_info['executable']
        project_path = project['project_path']
        
        try:
            # 打开IDE
            if os.name == 'nt':  # Windows
                subprocess.Popen([executable, project_path])
            else:  # Unix-like
                subprocess.Popen([executable, project_path])
            
            # 创建会话
            session_id = hashlib.sha256((project_id + str(time.time())).encode()).hexdigest()
            self.active_sessions[session_id] = {
                'session_id': session_id,
                'project_id': project_id,
                'ide_id': ide_id,
                'started_at': time.time(),
                'status': 'active'
            }
            
            return {
                'success': True,
                'message': f'Project opened in {ide_info["name"]}',
                'session_id': session_id
            }
        except Exception as e:
            print(f"IDE open error: {e}")
            return {
                'success': False,
                'message': str(e)
            }
    
    def create_ide_config(self, project_id: str, config_data: Dict[str, Any]) -> Dict[str, Any]:
        """创建IDE配置
        
        Args:
            project_id: 项目ID
            config_data: 配置数据
            
        Returns:
            配置结果
        """
        project = self.projects.get(project_id)
        if not project:
            return {
                'success': False,
                'message': 'Project not found'
            }
        
        ide_id = project['ide_id']
        project_path = project['project_path']
        
        # 根据IDE类型创建配置
        if ide_id == 'vscode':
            # 创建VS Code配置
            vscode_dir = os.path.join(project_path, '.vscode')
            os.makedirs(vscode_dir, exist_ok=True)
            
            # 创建launch.json
            launch_config = config_data.get('launch', {
                'version': '0.2.0',
                'configurations': [
                    {
                        'name': 'Python: Current File',
                        'type': 'python',
                        'request': 'launch',
                        'program': '${file}',
                        'console': 'integratedTerminal'
                    }
                ]
            })
            
            with open(os.path.join(vscode_dir, 'launch.json'), 'w', encoding='utf-8') as f:
                json.dump(launch_config, f, indent=2)
            
            # 创建settings.json
            settings_config = config_data.get('settings', {
                'python.pythonPath': 'python',
                'editor.tabSize': 4,
                'editor.insertSpaces': True
            })
            
            with open(os.path.join(vscode_dir, 'settings.json'), 'w', encoding='utf-8') as f:
                json.dump(settings_config, f, indent=2)
        
        elif ide_id == 'pycharm':
            # PyCharm配置由IDE自动管理
            pass
        
        return {
            'success': True,
            'message': 'IDE config created successfully'
        }
    
    def share_project(self, project_id: str, collaborators: List[str]) -> Dict[str, Any]:
        """共享项目
        
        Args:
            project_id: 项目ID
            collaborators: 协作者列表
            
        Returns:
            共享结果
        """
        project = self.projects.get(project_id)
        if not project:
            return {
                'success': False,
                'message': 'Project not found'
            }
        
        # 更新项目信息
        project['collaborators'] = collaborators
        project['last_modified'] = time.time()
        
        return {
            'success': True,
            'message': 'Project shared successfully',
            'collaborators': collaborators
        }
    
    def get_project_status(self, project_id: str) -> Dict[str, Any]:
        """获取项目状态
        
        Args:
            project_id: 项目ID
            
        Returns:
            项目状态
        """
        project = self.projects.get(project_id)
        if not project:
            return {
                'success': False,
                'message': 'Project not found'
            }
        
        # 检查项目文件数
        file_count = 0
        for root, dirs, files in os.walk(project['project_path']):
            file_count += len(files)
        
        status = {
            'project_id': project_id,
            'project_name': project['project_name'],
            'ide_id': project['ide_id'],
            'file_count': file_count,
            'created_at': project['created_at'],
            'last_modified': project['last_modified'],
            'collaborators': project.get('collaborators', []),
            'active_sessions': [session for session in self.active_sessions.values() if session['project_id'] == project_id]
        }
        
        return status
    
    def list_projects(self, ide_id: str = None) -> List[Dict[str, Any]]:
        """列出项目
        
        Args:
            ide_id: IDE ID（可选）
            
        Returns:
            项目列表
        """
        projects = []
        
        for project in self.projects.values():
            if ide_id and project['ide_id'] != ide_id:
                continue
            projects.append(project)
        
        # 按创建时间排序
        projects.sort(key=lambda x: x['created_at'], reverse=True)
        
        return projects

class CollaborationManager:
    """协同管理模块"""
    
    def __init__(self):
        """初始化协同管理模块"""
        self.sessions = {}
        self.changes = {}
    
    def create_session(self, project_id: str, users: List[str]) -> Dict[str, Any]:
        """创建协同会话
        
        Args:
            project_id: 项目ID
            users: 用户列表
            
        Returns:
            会话信息
        """
        # 生成会话ID
        session_id = hashlib.sha256((project_id + str(users) + str(time.time())).encode()).hexdigest()
        
        # 会话信息
        session = {
            'session_id': session_id,
            'project_id': project_id,
            'users': users,
            'status': 'active',
            'created_at': time.time(),
            'last_activity': time.time()
        }
        
        # 存储会话
        self.sessions[session_id] = session
        
        # 初始化变更记录
        self.changes[session_id] = []
        
        return session
    
    def record_change(self, session_id: str, user_id: str, file_path: str, change_type: str, content: str) -> Dict[str, Any]:
        """记录变更
        
        Args:
            session_id: 会话ID
            user_id: 用户ID
            file_path: 文件路径
            change_type: 变更类型（add, modify, delete）
            content: 变更内容
            
        Returns:
            变更记录
        """
        if session_id not in self.sessions:
            return {
                'success': False,
                'message': 'Session not found'
            }
        
        # 变更记录
        change = {
            'change_id': hashlib.sha256((session_id + user_id + file_path + str(time.time())).encode()).hexdigest(),
            'session_id': session_id,
            'user_id': user_id,
            'file_path': file_path,
            'change_type': change_type,
            'content': content,
            'timestamp': time.time()
        }
        
        # 存储变更
        self.changes[session_id].append(change)
        
        # 更新会话活动时间
        self.sessions[session_id]['last_activity'] = time.time()
        
        return change
    
    def get_changes(self, session_id: str, since: float = None) -> List[Dict[str, Any]]:
        """获取变更
        
        Args:
            session_id: 会话ID
            since: 时间戳（可选）
            
        Returns:
            变更列表
        """
        if session_id not in self.changes:
            return []
        
        changes = self.changes[session_id]
        
        if since:
            changes = [change for change in changes if change['timestamp'] > since]
        
        return changes
    
    def end_session(self, session_id: str) -> Dict[str, Any]:
        """结束会话
        
        Args:
            session_id: 会话ID
            
        Returns:
            结束结果
        """
        if session_id not in self.sessions:
            return {
                'success': False,
                'message': 'Session not found'
            }
        
        # 更新会话状态
        session = self.sessions[session_id]
        session['status'] = 'ended'
        session['ended_at'] = time.time()
        
        return {
            'success': True,
            'message': 'Session ended successfully',
            'session_id': session_id
        }
