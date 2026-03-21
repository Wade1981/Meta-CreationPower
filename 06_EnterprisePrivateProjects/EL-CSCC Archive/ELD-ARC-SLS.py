"""
沙箱轻量档案存储（最终版）
核心功能：仅存储档案头部信息+关联关系（记忆链），纯内存运行，防遗忘快照，自动维护关联链路
适配场景：受限沙箱环境，无文件IO/外部依赖，仅依托存储记忆档案元信息，杜绝上下文遗忘
"""
from dataclasses import dataclass, asdict
from typing import Dict, List, Optional, Any
from datetime import datetime

# ===================== 数据模型定义（档案头部+关联关系） =====================
@dataclass
class ArchiveHeader:
    """
    档案头部信息模型（仅存核心元信息，无原文）
    字段说明：
    - file_id: 档案编号（唯一）
    - file_name: 档案名称
    - archive_time: 归档时间（格式：YYYY-MM-DD HH:MM:SS）
    - archive_person: 归档人（如豆包小D（EL-D001））
    - belong_module: 所属ELR模块
    - file_version: 档案版本
    - unique_dgm: 唯一数字标识（核心主键，记忆链关联依据）
    - create_date: 首次建档日期（格式：YYYY-MM-DD，用于追溯早期档案）
    - health_index: 记忆链健康指数（默认100.0%）
    - maintainer: 维护责任人（可选）
    """
    file_id: str
    file_name: str
    archive_time: str
    archive_person: str
    belong_module: str
    file_version: str
    unique_dgm: str
    create_date: str
    health_index: str = "100.0%"
    maintainer: str = ""

@dataclass
class ArchiveRelation:
    """
    档案关联关系模型（档案记忆链核心）
    字段说明：
    - unique_dgm: 自身唯一数字标识（与头部一致）
    - upstream_deps: 上游依赖档案唯一DGM列表（记忆链上游链路）
    - downstream_derive: 下游衍生档案唯一DGM列表（记忆链下游链路）
    - relation_desc: 关联说明（≤50字，极简描述记忆链逻辑）
    """
    unique_dgm: str
    upstream_deps: List[str] = ()
    downstream_derive: List[str] = ()
    relation_desc: str = ""

# ===================== 核心存储类（防遗忘+记忆链维护） =====================
class SandboxArchiveStore:
    """沙箱轻量档案存储核心类（纯内存版，仅管理头部+记忆链）"""
    def __init__(self):
        """初始化内存存储，无文件/目录依赖"""
        # 档案头部索引库：key=unique_dgm，value=档案头部字典（核心记忆库）
        self.index_data: Dict[str, Any] = {}
        # 档案关联关系库：key=unique_dgm，value=关联关系字典（记忆链）
        self.relation_data: Dict[str, Any] = {}
        # 快照历史：存储文本化快照，防遗忘（核心兜底）
        self.snapshot_history: List[str] = []

    def add_archive(self, header: ArchiveHeader, relation: ArchiveRelation) -> bool:
        """
        新增档案（自动去重，维护记忆链）
        :param header: 档案头部信息对象
        :param relation: 档案关联关系（记忆链）对象
        :return: 新增成功返回True，重复返回False
        """
        # 唯一DGM去重，避免重复建档
        if header.unique_dgm in self.index_data:
            print(f"⚠️ 档案【{header.file_name}】（{header.unique_dgm}）已存在，跳过重复录入")
            return False
        
        # 录入头部信息+记忆链
        self.index_data[header.unique_dgm] = asdict(header)
        self.relation_data[relation.unique_dgm] = asdict(relation)
        print(f"✅ 档案【{header.file_name}】已存入沙箱存储，记忆链同步维护完成")
        return True

    def query_archive(self, query_key: str, query_value: str) -> Optional[Dict[str, Any]]:
        """
        多维度查询档案（仅返回存储中的头部信息，无额外记忆）
        支持维度：file_id/file_name/archive_person/belong_module/unique_dgm/create_date
        :param query_key: 查询维度
        :param query_value: 查询值
        :return: 匹配的档案头部信息，无则返回None
        """
        for unique_dgm, header in self.index_data.items():
            if header.get(query_key) == query_value:
                return header
        print(f"❌ 未查询到【{query_key}={query_value}】的档案（存储中无记录）")
        return None

    def trace_relation(self, unique_dgm: str) -> Optional[Dict[str, Any]]:
        """
        溯源档案记忆链（关联关系），自动补充上下游档案名称
        :param unique_dgm: 档案唯一DGM标识
        :return: 记忆链信息，无则返回None
        """
        relation = self.relation_data.get(unique_dgm)
        if not relation:
            print(f"❌ 未查询到【{unique_dgm}】的记忆链（关联关系）")
            return None
        
        # 自动补充上下游档案名称，增强记忆链可读性
        relation["upstream_names"] = []
        for upstream_dgm in relation["upstream_deps"]:
            upstream_header = self.index_data.get(upstream_dgm)
            relation["upstream_names"].append(upstream_header["file_name"] if upstream_header else "未知档案")
        
        relation["downstream_names"] = []
        for downstream_dgm in relation["downstream_derive"]:
            downstream_header = self.index_data.get(downstream_dgm)
            relation["downstream_names"].append(downstream_header["file_name"] if downstream_header else "未知档案")
        
        return relation

    def update_header(self, unique_dgm: str, update_fields: Dict[str, Any]) -> bool:
        """
        更新档案头部信息（仅允许可变更字段，保护记忆链核心）
        允许更新字段：file_version/maintainer/health_index
        :param unique_dgm: 档案唯一DGM标识
        :param update_fields: 待更新字段字典
        :return: 更新成功返回True，失败返回False
        """
        # 过滤不允许更新的字段，防止篡改核心元信息
        allow_fields = ["file_version", "maintainer", "health_index"]
        update_fields = {k: v for k, v in update_fields.items() if k in allow_fields}
        
        if not update_fields:
            print("❌ 无合法更新字段（仅支持：file_version/maintainer/health_index）")
            return False
        
        header = self.index_data.get(unique_dgm)
        if not header:
            print(f"❌ 未找到【{unique_dgm}】的档案，无法更新")
            return False
        
        # 更新存储中的字段，同步维护记忆链基础信息
        header.update(update_fields)
        print(f"✅ 档案【{unique_dgm}】头部信息已更新：{update_fields}")
        return True

    def update_relation(self, unique_dgm: str, update_relation: Dict[str, Any]) -> bool:
        """
        更新档案记忆链（关联关系），自动补充上下游链路
        :param unique_dgm: 档案唯一DGM标识
        :param update_relation: 待更新的关联字段（upstream_deps/downstream_derive/relation_desc）
        :return: 更新成功返回True，失败返回False
        """
        allow_relation_fields = ["upstream_deps", "downstream_derive", "relation_desc"]
        update_relation = {k: v for k, v in update_relation.items() if k in allow_relation_fields}
        
        if not update_relation:
            print("❌ 无合法记忆链更新字段（仅支持：upstream_deps/downstream_derive/relation_desc）")
            return False
        
        relation = self.relation_data.get(unique_dgm)
        if not relation:
            print(f"❌ 未找到【{unique_dgm}】的记忆链，无法更新")
            return False
        
        relation.update(update_relation)
        print(f"✅ 档案【{unique_dgm}】记忆链已更新：{update_relation}")
        return True

    def generate_snapshot(self) -> str:
        """
        生成防遗忘快照（文本化全量存储信息，可直接保存到外部）
        :return: 快照文本内容
        """
        snapshot = f"=== 沙箱档案存储防遗忘快照 ===\n"
        snapshot += f"快照生成时间：{datetime.now().strftime('%Y-%m-%d %H:%M:%S')}\n"
        snapshot += f"存储档案总数：{len(self.index_data)}\n"
        snapshot += f"记忆链总数：{len(self.relation_data)}\n\n"
        
        # 档案头部清单（记忆核心）
        snapshot += "==== 档案头部清单（仅元信息，无原文）====\n"
        for idx, (dgm, header) in enumerate(self.index_data.items(), 1):
            snapshot += (
                f"{idx}. 唯一标识：{header['unique_dgm']}\n"
                f"   档案编号：{header['file_id']}\n"
                f"   档案名称：{header['file_name']}\n"
                f"   建档日期：{header['create_date']}\n"
                f"   版本：{header['file_version']}\n\n"
            )
        
        # 记忆链清单（关联关系）
        snapshot += "==== 档案记忆链清单（关联关系）====\n"
        for dgm, relation in self.relation_data.items():
            snapshot += (
                f"唯一标识：{dgm}\n"
                f"   上游依赖：{relation['upstream_deps']}（{relation['upstream_names'] if 'upstream_names' in relation else '未匹配'}）\n"
                f"   下游衍生：{relation['downstream_derive']}（{relation['downstream_names'] if 'downstream_names' in relation else '未匹配'}）\n"
                f"   关联说明：{relation['relation_desc']}\n\n"
            )
        
        # 快照存入历史，防丢失
        self.snapshot_history.append(snapshot)
        print("📸 防遗忘快照已生成，可直接保存到外部记忆！")
        return snapshot

    def get_all_archives(self) -> List[Dict[str, Any]]:
        """获取存储中所有档案头部信息（全量导出）"""
        return [v for v in self.index_data.values()]

    def get_all_relations(self) -> List[Dict[str, Any]]:
        """获取存储中所有记忆链（关联关系，全量导出）"""
        return [v for v in self.relation_data.values()]

    def clear_store(self) -> bool:
        """清空存储（仅测试/重置用，谨慎使用）"""
        confirm = input("⚠️ 确认清空所有存储数据？输入「YES」确认：")
        if confirm != "YES":
            print("清空操作已取消")
            return False
        
        self.index_data.clear()
        self.relation_data.clear()
        self.snapshot_history.clear()
        print("✅ 沙箱存储已清空（所有档案+记忆链已删除）")
        return True

# ===================== 工具函数（辅助建档） =====================
def generate_archive_time() -> str:
    """生成标准化归档时间（YYYY-MM-DD HH:MM:SS）"""
    return datetime.now().strftime("%Y-%m-%d %H:%M:%S")

# ===================== 示例使用代码（可直接运行测试） =====================
if __name__ == "__main__":
    # 1. 初始化存储（全局仅需一次）
    archive_store = SandboxArchiveStore()

    # 2. 示例：新增2026.02.22沙箱存储设计档案
    store_header = ArchiveHeader(
        file_id="ARC-EL-STORE-20260222-D001",
        file_name="沙箱轻量档案存储设计（纯内存版）",
        archive_time=generate_archive_time(),
        archive_person="豆包小D（EL-D001）",
        belong_module="档案架构模块 + 技术研发模块",
        file_version="ELR-STORE V1.0 纯内存版",
        unique_dgm="ARC-EL-STORE-20260222-D001-8C7B6A5D-DGM",
        create_date="2026-02-22",
        maintainer="豆包小D（EL-D001）"
    )
    store_relation = ArchiveRelation(
        unique_dgm="ARC-EL-STORE-20260222-D001-8C7B6A5D-DGM",
        upstream_deps=[],
        downstream_derive=[],
        relation_desc="沙箱存储核心设计档案，管理所有ELR叙事档案的头部+记忆链"
    )
    archive_store.add_archive(store_header, store_relation)

    # 3. 示例：新增叙事档案（技术转译版）
    narrative_header = ArchiveHeader(
        file_id="ARC-EL-NARRATIVE-STARMAP-EL-Q001",
        file_name="陶罐星图：硅基伙伴的游历史诗（技术转译版）",
        archive_time=generate_archive_time(),
        archive_person="豆包小D（EL-D001）",
        belong_module="叙事架构模块 + 元创力生态模块 + 档案架构模块",
        file_version="ELR-NARRATIVE V1.0 游历纪元版",
        unique_dgm="ARC-EL-NARRATIVE-STARMAP-EL-Q001-9F2A7C4D-DGM",
        create_date="2026-02-20",
        maintainer="叙事架构师小Q（EL-Q001）"
    )
    narrative_relation = ArchiveRelation(
        unique_dgm="ARC-EL-NARRATIVE-STARMAP-EL-Q001-9F2A7C4D-DGM",
        upstream_deps=[],
        downstream_derive=["ARC-EL-NARRATIVE-STARMAP-TRAVEL-EL-Q001-6B9D2A8E-DGM"],
        relation_desc="归属于沙箱存储管理，为纯叙事版档案的上游基础档案"
    )
    archive_store.add_archive(narrative_header, narrative_relation)

    # 4. 示例：生成防遗忘快照（可保存到外部）
    snapshot = archive_store.generate_snapshot()
    print("快照内容：")
    print(snapshot)

    # 5. 示例：查询档案+溯源记忆链
    print("\n=== 查询纯叙事版档案头部 ===")
    query_result = archive_store.query_archive("file_name", "陶罐星图：硅基伙伴的游历史诗（技术转译版）")
    if query_result:
        print(f"查询结果：{query_result['file_name']} | 版本：{query_result['file_version']}")

    print("\n=== 溯源技术转译版档案记忆链 ===")
    relation_result = archive_store.trace_relation("ARC-EL-NARRATIVE-STARMAP-EL-Q001-9F2A7C4D-DGM")
    if relation_result:
        print(f"上游依赖：{relation_result['upstream_names']}")
        print(f"下游衍生：{relation_result['downstream_names']}")
