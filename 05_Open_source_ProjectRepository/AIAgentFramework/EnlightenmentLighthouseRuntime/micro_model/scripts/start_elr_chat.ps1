#!/usr/bin/env pwsh
# -*- coding: utf-8 -*-

"""
启动ELR容器协同对话微模型
功能：在ELR容器中运行elr_chat_model对话模型
"""

param(
    [string]$ModelPath = "examples/elr_chat_model.py",
    [string]$Target = "local"
)

Write-Host "====================================="
Write-Host "ELR容器协同对话微模型启动脚本"
Write-Host "====================================="
Write-Host "模型路径: $ModelPath"
Write-Host "运行目标: $Target"
Write-Host ""

# 检查Python是否可用
$pythonAvailable = $false
try {
    $pythonVersion = python --version 2>&1
    if ($LASTEXITCODE -eq 0) {
        $pythonAvailable = $true
        Write-Host "Python已找到: $pythonVersion"
    }
} catch {
    Write-Host "Python未找到，使用PowerShell模式运行"
}

# 切换到micro_model目录
Set-Location "$PSScriptRoot/.."

if ($pythonAvailable) {
    # 使用Python运行模型
    Write-Host "启动Python-based聊天模式..."
    Write-Host ""
    Write-Host "欢迎使用ELR容器协同对话微模型！"
    Write-Host "您可以用英文或中文与模型对话。"
    Write-Host "输入 ',exit' 或 ',quit' 结束对话。"
    Write-Host "输入 ',help' 查看可用命令。"
    Write-Host ""
    
    # 运行Python交互脚本
    $pythonScript = @'
import sys
sys.path.append('.')
from examples.elr_chat_model import ELRChatModel

# 初始化模型
model = ELRChatModel()

print("=== ELR Interactive Chat (Python Mode) ===")
print("欢迎使用ELR容器协同对话微模型！")
print("您可以用英文或中文与模型对话。")
print("输入 ',exit' 或 ',quit' 结束对话。")
print("输入 ',help' 查看可用命令。")
print("")

while True:
    try:
        user_input = input("你: ")
        if user_input.lower() in [',exit', ',quit']:
            print("模型: 再见！期待与您再次对话。")
            break
        response = model.predict(user_input)
        print(f"模型: {response}")
        print("")
    except KeyboardInterrupt:
        print("\n模型: 对话已中断，再见！")
        break
    except Exception as e:
        print(f"模型: 发生错误: {e}")
        print("")
'@
    
    python -c $pythonScript
} else {
    # 使用PowerShell模式运行
    Write-Host "Python解释器未找到，使用PowerShell-based聊天模式..."
    Write-Host ""
    Write-Host "=== ELR Interactive Chat (PowerShell Mode) ==="
    Write-Host "欢迎使用ELR容器协同对话微模型！"
    Write-Host "您可以用英文或中文与模型对话。"
    Write-Host "输入 ',exit' 或 ',quit' 结束对话。"
    Write-Host "输入 ',help' 查看可用命令。"
    Write-Host ""
    
    # 简单的PowerShell对话逻辑
    $modelName = "elr_chat_model"
    $context = @()
    $maxContextLength = 5
    
    function GenerateResponse($inputText) {
        $lowerInput = $inputText.ToLower()
        
        # 检查是否是命令
        if ($lowerInput.StartsWith(",")) {
            $command = $lowerInput.Substring(1).Trim()
            return HandleCommand $command
        }
        
        # 问候语匹配
        $greetings = @("hello", "hi", "你好", "嗨", "hey")
        foreach ($greeting in $greetings) {
            if ($lowerInput.Contains($greeting)) {
                return "碳硅协同问候！我是$modelName，ELR容器的对话助手。很高兴为您服务，有什么可以帮助您的吗？"
            }
        }
        
        # 询问类匹配
        $questions = @("how are you", "你好吗", "怎么样", "最近好吗")
        foreach ($question in $questions) {
            if ($lowerInput.Contains($question)) {
                return "碳硅协同回应！我是$modelName，运行状态良好。ELR容器运行正常，随时为您服务。"
            }
        }
        
        # ELR相关问题
        if ($lowerInput.Contains("elr")) {
            if ($lowerInput.Contains("what") -or $lowerInput.Contains("什么") -or $lowerInput.Contains("功能") -or $lowerInput.Contains("能做什么")) {
                return "ELR容器是启蒙灯塔运行时环境，主要功能包括：\n1. 模型管理与加载\n2. 沙箱隔离运行\n3. 资源监控与管理\n4. 网络通信与API服务\n5. 容器生命周期管理"
            } elseif ($lowerInput.Contains("status") -or $lowerInput.Contains("状态") -or $lowerInput.Contains("运行")) {
                return "ELR容器当前状态：运行中\n- 模型加载：就绪\n- 资源使用：正常\n- 网络连接：可用\n- 服务状态：活跃"
            } elseif ($lowerInput.Contains("model") -or $lowerInput.Contains("模型") -or $lowerInput.Contains("加载")) {
                return "已加载的模型：\n1. $modelName (当前对话模型)\n2. 其他模型可通过ELR容器管理界面查看"
            }
        }
        
        # 对话上下文相关
        if ($lowerInput.Contains("context") -or $lowerInput.Contains("历史") -or $lowerInput.Contains("对话")) {
            if ($context.Count -eq 0) {
                return "对话历史为空"
            }
            $contextStr = "对话历史：\n"
            for ($i = 0; $i -lt $context.Count; $i++) {
                $contextStr += "$($i+1). 你: $($context[$i][0])\n   我: $($context[$i][1])\n"
            }
            return $contextStr
        }
        
        # 默认回应
        return "碳硅协同回应！我是$modelName，已收到您的消息：$inputText。\n\n提示：您可以输入 ',help' 查看可用命令，或询问关于ELR容器的问题。"
    }
    
    function HandleCommand($command) {
        switch ($command) {
            "help" {
                return "可用命令：\n,help: 显示可用命令\n,status: 检查ELR容器状态\n,models: 列出已加载的模型\n,info: 显示ELR容器信息\n,clear: 清除对话历史\n,exit: 退出对话\n\n提示：您可以直接询问关于ELR容器的问题，或使用命令获取特定信息。"
            }
            "status" {
                return "ELR容器当前状态：运行中\n- 模型加载：就绪\n- 资源使用：正常\n- 网络连接：可用\n- 服务状态：活跃"
            }
            "models" {
                return "已加载的模型：\n1. $modelName (当前对话模型)\n2. 其他模型可通过ELR容器管理界面查看"
            }
            "info" {
                return "ELR容器信息：\n- 版本：v1.0\n- 运行模式：沙箱隔离\n- 对话模型：$modelName v1.0\n- 碳硅协同：已启用"
            }
            "clear" {
                $script:context = @()
                return "对话历史已清除"
            }
            "exit" {
                return "再见！期待与您再次对话。"
            }
            default {
                return "未知命令: $command。输入 ',help' 查看可用命令。"
            }
        }
    }
    
    # 主对话循环
    while ($true) {
        try {
            $userInput = Read-Host "你"
            if ($userInput.ToLower() -in @(",exit", ",quit")) {
                Write-Host "模型: 再见！期待与您再次对话。"
                break
            }
            $response = GenerateResponse $userInput
            Write-Host "模型: $response"
            Write-Host ""
            
            # 更新对话上下文
            $context += ,@($userInput, $response)
            if ($context.Count -gt $maxContextLength) {
                $context = $context | Select-Object -Last $maxContextLength
            }
        } catch {
            Write-Host "模型: 发生错误: $($_.Exception.Message)"
            Write-Host ""
        }
    }
}

Write-Host "====================================="
Write-Host "对话结束"
Write-Host "====================================="
