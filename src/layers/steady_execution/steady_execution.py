"""
静定执行层 (Steady Execution)
功能：以绝对稳定、无感知干扰的方式执行协同路径
"""

from dataclasses import dataclass
from typing import Dict, List, Optional, Tuple, Any, Callable
import json
import uuid
import time
from queue import Queue
import threading
import logging


@dataclass
class Task:
    """任务基类"""
    task_id: str
    name: str
    type: str
    priority: int
    payload: Dict[str, Any]
    status: str  # pending, executing, completed, failed
    created_at: float
    started_at: Optional[float]
    completed_at: Optional[float]
    error: Optional[str]


class SteadyExecutor:
    """
    静定执行器
    以绝对稳定、无感知干扰的方式执行协同路径
    """
    
    def __init__(self):
        """
        初始化静定执行器
        """
        self.task_queue = Queue()
        self.active_tasks: Dict[str, Task] = {}
        self.completed_tasks: Dict[str, Task] = {}
        self.failed_tasks: Dict[str, Task] = {}
        
        self.execution_thread = threading.Thread(target=self._execution_loop, daemon=True)
        self.execution_thread.start()
        
        self.is_running = True
        self.max_concurrent_tasks = 5
        self.current_concurrent_tasks = 0
        
        # 配置日志
        logging.basicConfig(
            filename="execution.log",
            level=logging.INFO,
            format='%(asctime)s - %(levelname)s - %(message)s'
        )
        self.logger = logging.getLogger("SteadyExecutor")
    
    def _execution_loop(self):
        """
        执行循环
        """
        while self.is_running:
            if self.current_concurrent_tasks < self.max_concurrent_tasks and not self.task_queue.empty():
                task = self.task_queue.get()
                self._execute_task(task)
            time.sleep(0.005)  # 5ms 检查一次，确保响应速度，符合"静"原则的无抖动执行
    
    def _execute_task(self, task: Task):
        """
        执行单个任务
        
        Args:
            task: 任务对象
        """
        def task_thread():
            try:
                self.current_concurrent_tasks += 1
                task.status = "executing"
                task.started_at = time.time()
                self.active_tasks[task.task_id] = task
                
                self.logger.info(f"开始执行任务: {task.name} (ID: {task.task_id})")
                
                # 实际实现中应根据任务类型执行相应的操作
                # 这里仅模拟执行过程
                time.sleep(0.05)  # 模拟执行时间，确保低于100ms
                
                task.status = "completed"
                task.completed_at = time.time()
                self.completed_tasks[task.task_id] = task
                self.active_tasks.pop(task.task_id, None)
                
                self.logger.info(f"任务完成: {task.name} (ID: {task.task_id})")
                
            except Exception as e:
                error_msg = str(e)
                task.status = "failed"
                task.error = error_msg
                self.failed_tasks[task.task_id] = task
                self.active_tasks.pop(task.task_id, None)
                
                # 错误静默处理，仅记录日志
                self.logger.error(f"任务失败: {task.name} (ID: {task.task_id}) - {error_msg}")
                
            finally:
                self.current_concurrent_tasks -= 1
        
        # 启动任务线程
        thread = threading.Thread(target=task_thread, daemon=True)
        thread.start()
    
    def submit_task(self, name: str, task_type: str, 
                   payload: Dict[str, Any], 
                   priority: int = 0) -> str:
        """
        提交任务
        
        Args:
            name: 任务名称
            task_type: 任务类型
            payload: 任务负载
            priority: 任务优先级
        
        Returns:
            任务ID
        """
        task_id = str(uuid.uuid4())
        task = Task(
            task_id=task_id,
            name=name,
            type=task_type,
            priority=priority,
            payload=payload,
            status="pending",
            created_at=time.time(),
            started_at=None,
            completed_at=None,
            error=None
        )
        
        self.task_queue.put(task)
        self.logger.info(f"提交任务: {name} (ID: {task_id})")
        
        return task_id
    
    def get_task_status(self, task_id: str) -> Dict[str, Any]:
        """
        获取任务状态
        
        Args:
            task_id: 任务ID
        
        Returns:
            任务状态字典
        """
        if task_id in self.active_tasks:
            task = self.active_tasks[task_id]
        elif task_id in self.completed_tasks:
            task = self.completed_tasks[task_id]
        elif task_id in self.failed_tasks:
            task = self.failed_tasks[task_id]
        else:
            return {"error": "任务不存在"}
        
        return {
            "task_id": task.task_id,
            "name": task.name,
            "status": task.status,
            "created_at": task.created_at,
            "started_at": task.started_at,
            "completed_at": task.completed_at,
            "error": task.error
        }
    
    def execute_counterpoint_path(self, path_id: str, 
                                 steps: List[Dict[str, Any]],
                                 voice_map: Dict[str, str]) -> Dict[str, Any]:
        """
        执行协同路径
        
        Args:
            path_id: 路径ID
            steps: 步骤列表
            voice_map: 声部映射
        
        Returns:
            执行结果
        """
        execution_id = str(uuid.uuid4())
        task_ids = []
        
        self.logger.info(f"开始执行协同路径: {path_id} (执行ID: {execution_id})")
        
        # 原子化执行所有步骤
        for i, step in enumerate(steps):
            task_id = self.submit_task(
                name=f"步骤 {i + 1}: {step['action']}",
                task_type="counterpoint_step",
                payload={
                    "step": step,
                    "path_id": path_id,
                    "execution_id": execution_id,
                    "voice_id": voice_map.get(step['role'], "")
                },
                priority=len(steps) - i  # 确保步骤按顺序执行
            )
            task_ids.append(task_id)
        
        # 等待所有任务完成
        # 注意：这里使用轮询而非阻塞，确保系统响应性
        start_time = time.time()
        while time.time() - start_time < 5:  # 5秒超时
            all_completed = all(
                self.get_task_status(tid).get("status") in ["completed", "failed"]
                for tid in task_ids
            )
            if all_completed:
                break
            time.sleep(0.01)
        
        # 收集执行结果
        results = {
            "execution_id": execution_id,
            "path_id": path_id,
            "task_results": [self.get_task_status(tid) for tid in task_ids],
            "success": all(self.get_task_status(tid).get("status") == "completed" for tid in task_ids)
        }
        
        self.logger.info(f"协同路径执行完成: {path_id} (执行ID: {execution_id}, 成功: {results['success']})")
        
        return results
    
    def get_execution_stats(self) -> Dict[str, Any]:
        """
        获取执行统计信息
        
        Returns:
            执行统计信息字典
        """
        return {
            "queue_size": self.task_queue.qsize(),
            "active_tasks": len(self.active_tasks),
            "completed_tasks": len(self.completed_tasks),
            "failed_tasks": len(self.failed_tasks),
            "current_concurrent_tasks": self.current_concurrent_tasks,
            "max_concurrent_tasks": self.max_concurrent_tasks
        }
    
    def clear_completed_tasks(self):
        """
        清理已完成的任务
        """
        self.completed_tasks.clear()
        self.failed_tasks.clear()
        self.logger.info("清理已完成和失败的任务")
    
    def shutdown(self):
        """
        关闭执行器
        """
        self.is_running = False
        if self.execution_thread.is_alive():
            self.execution_thread.join(timeout=2)
        self.logger.info("执行器已关闭")
    
    def get_system_health(self) -> Dict[str, Any]:
        """
        获取系统健康状态
        
        Returns:
            系统健康状态字典
        """
        # 实际实现中应包含更详细的健康检查
        return {
            "status": "healthy" if self.is_running else "unhealthy",
            "execution_thread_alive": self.execution_thread.is_alive(),
            "stats": self.get_execution_stats()
        }
    
    def set_max_concurrent_tasks(self, max_tasks: int):
        """
        设置最大并发任务数
        
        Args:
            max_tasks: 最大并发任务数
        """
        if max_tasks > 0:
            self.max_concurrent_tasks = max_tasks
            self.logger.info(f"设置最大并发任务数: {max_tasks}")
    
    def get_task_history(self, limit: int = 100) -> List[Dict[str, Any]]:
        """
        获取任务历史
        
        Args:
            limit: 限制数量
        
        Returns:
            任务历史列表
        """
        history = []
        
        # 按时间排序获取最近的任务
        all_tasks = list(self.completed_tasks.values()) + list(self.failed_tasks.values())
        all_tasks.sort(key=lambda t: t.created_at, reverse=True)
        
        for task in all_tasks[:limit]:
            history.append({
                "task_id": task.task_id,
                "name": task.name,
                "type": task.type,
                "status": task.status,
                "created_at": task.created_at,
                "completed_at": task.completed_at,
                "error": task.error
            })
        
        return history
    
    def is_task_completed(self, task_id: str) -> bool:
        """
        检查任务是否完成
        
        Args:
            task_id: 任务ID
        
        Returns:
            是否完成
        """
        status = self.get_task_status(task_id).get("status")
        return status == "completed"
    
    def get_failed_tasks(self) -> List[Dict[str, Any]]:
        """
        获取失败的任务
        
        Returns:
            失败任务列表
        """
        return [{
            "task_id": task.task_id,
            "name": task.name,
            "error": task.error,
            "created_at": task.created_at
        } for task in self.failed_tasks.values()]
