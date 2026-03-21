"""Configuration management module for ELR Service Load."""

import os


def load_config():
    """
    Load configuration from environment variables or use defaults.
    
    Returns:
        dict: Configuration dictionary
    """
    config = {
        'host': os.environ.get('ELR_SERVICE_HOST', '0.0.0.0'),
        'port': int(os.environ.get('ELR_SERVICE_PORT', '8000')),
        'archive_file': os.environ.get('ELR_ARCHIVE_FILE', 'el_cscc_archive.json'),
        'log_level': os.environ.get('ELR_LOG_LEVEL', 'INFO')
    }
    return config
