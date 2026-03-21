"""Logging utility module for ELR Service Load."""

import logging
import os


def setup_logger():
    """
    Setup logger with appropriate configuration.
    
    Returns:
        logging.Logger: Configured logger instance
    """
    # Get log level from environment or use default
    log_level = os.environ.get('ELR_LOG_LEVEL', 'INFO')
    
    # Create logger
    logger = logging.getLogger('elr-service-load')
    logger.setLevel(log_level)
    
    # Create console handler
    handler = logging.StreamHandler()
    handler.setLevel(log_level)
    
    # Create formatter
    formatter = logging.Formatter('%(asctime)s - %(name)s - %(levelname)s - %(message)s')
    handler.setFormatter(formatter)
    
    # Add handler to logger
    logger.addHandler(handler)
    
    return logger
