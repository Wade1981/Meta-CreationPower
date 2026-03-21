#!/usr/bin/env python3
"""
测试ELR-Desktop-Assistant的连接逻辑
"""

import requests
import time

def test_elr_connection():
    """测试ELR连接"""
    print("测试ELR-Desktop-Assistant连接逻辑...")
    
    # 模拟ELRClient的连接逻辑
    api_base_url = "http://localhost:8080/api"
    
    print(f"测试地址: {api_base_url}")
    
    # 测试get_elr_status
    print("\n1. 测试get_elr_status:")
    try:
        response = requests.get(f"{api_base_url}/status", timeout=2)
        print(f"状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
        if response.status_code == 200:
            data = response.json()
            if "status" in data:
                print(f"解析到status: {data['status']}")
            else:
                print("响应中没有status字段")
        else:
            print("响应状态码不是200")
    except Exception as e:
        print(f"异常: {e}")
    
    # 测试get_elr_containers
    print("\n2. 测试get_elr_containers:")
    try:
        response = requests.get(f"{api_base_url}/containers", timeout=2)
        print(f"状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
        if response.status_code == 200:
            data = response.json()
            if isinstance(data, list):
                print(f"解析到容器列表，长度: {len(data)}")
            else:
                print("响应不是列表")
        else:
            print("响应状态码不是200")
    except Exception as e:
        print(f"异常: {e}")
    
    # 测试健康检查
    print("\n3. 测试健康检查:")
    try:
        response = requests.get("http://localhost:8080/health", timeout=2)
        print(f"状态码: {response.status_code}")
        print(f"响应内容: {response.text}")
    except Exception as e:
        print(f"异常: {e}")

if __name__ == "__main__":
    test_elr_connection()
