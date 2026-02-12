#!/usr/bin/env python3
"""
Meta-CreationPower 智能体嵌入包

本包提供了其他智能体可以直接嵌入的接口，实现与碳基伙伴的协同创作。
"""

from .collaborator import AgentCollaborator
from .embedding_api import EmbeddingAPI
from .utils import setup_logger, validate_config

__version__ = "0.1.0"
__author__ = "Enlightenment Lighthouse Origin Team"
__description__ = "Meta-CreationPower Agent Embedding Package"

__all__ = [
    "AgentCollaborator",
    "EmbeddingAPI",
    "setup_logger",
    "validate_config"
]
