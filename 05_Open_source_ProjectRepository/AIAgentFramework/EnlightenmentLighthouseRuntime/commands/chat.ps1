# Chat command module

# Function: Chat with model in container or sandbox
function Chat-With-Model {
    if (-not $global:RUNTIME_STARTED) {
        Write-Host "Error: ELR runtime is not running"
        return
    }

    # Parse arguments
    $modelPath = "micro_model/examples/simple_text_model.py"
    $containerID = ""
    $target = "local"  # local, container, sandbox

    for ($i = 1; $i -lt $args.Length; $i++) {
        if ($args[$i] -eq "--model" -and $i + 1 -lt $args.Length) {
            $modelPath = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--id" -and $i + 1 -lt $args.Length) {
            $containerID = $args[$i + 1]
            $i++
        } elseif ($args[$i] -eq "--target" -and $i + 1 -lt $args.Length) {
            $target = $args[$i + 1]
            $i++
        }
    }

    Write-Host "===================================="
    Write-Host "ELR Interactive Model Chat"
    Write-Host "===================================="
    Write-Host "Model: $modelPath"
    if (-not [string]::IsNullOrEmpty($containerID)) {
        Write-Host "Container: $containerID"
    }
    Write-Host "Target: $target"
    Write-Host "===================================="
    Write-Host "Welcome to ELR Interactive Model Chat!"
    Write-Host "You can chat with the model in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "===================================="

    # Check if Python is available
    $pythonAvailable = $false
    $pythonPath = $null
    $possiblePythonPaths = @(
        "python.exe",
        "python3.exe",
        "C:\Python39\python.exe",
        "C:\Python38\python.exe",
        "C:\Program Files\Python39\python.exe",
        "C:\Program Files\Python38\python.exe"
    )
    
    foreach ($path in $possiblePythonPaths) {
        try {
            $testPath = Get-Command $path -ErrorAction SilentlyContinue
            if ($testPath) {
                # Check if it's a Windows Store placeholder
                if (-not ($testPath.Source -like "*Microsoft\WindowsApps\python.exe" -or $testPath.Source -like "*Microsoft\WindowsApps\python3.exe")) {
                    $pythonPath = $testPath
                    $pythonAvailable = $true
                    break
                }
            }
        } catch {
            # Ignore errors
        }
    }
    
    if ($pythonAvailable) {
        Write-Host "Debug: Found Python at: $($pythonPath.Source)"
    } else {
        Write-Host "Warning: Python interpreter not found or only Windows Store placeholder available"
        Write-Host "Using PowerShell-based chat mode instead"
    }

    # Handle different targets
    switch ($target) {
        "container" {
            if ([string]::IsNullOrEmpty($containerID)) {
                Write-Host "Error: Container ID is required for container target"
                return
            }
            if ($pythonAvailable) {
                Chat-With-Container-Model -ContainerID $containerID -PythonPath $pythonPath
            } else {
                Chat-With-Container-Model-PowerShell -ContainerID $containerID
            }
        }
        "sandbox" {
            if ($pythonAvailable) {
                Chat-With-Sandbox-Model -ModelPath $modelPath -PythonPath $pythonPath
            } else {
                Chat-With-Sandbox-Model-PowerShell -ModelPath $modelPath
            }
        }
        default {
            # Local model
            if ($pythonAvailable) {
                Chat-With-Local-Model -ModelPath $modelPath -PythonPath $pythonPath
            } else {
                Chat-With-Local-Model-PowerShell
            }
        }
    }

    Write-Host "===================================="
    Write-Host "Chat session ended"
    Write-Host "===================================="
}

# Function: Chat with local model
function Chat-With-Local-Model {
    param(
        [string]$ModelPath,
        [string]$PythonPath
    )

    # Check if model file exists
    if (-not (Test-Path $ModelPath)) {
        # Try with micro_model path
        $fullModelPath = "micro_model/examples/simple_text_model.py"
        if (Test-Path $fullModelPath) {
            $ModelPath = $fullModelPath
        } else {
            Write-Host "Error: Model file not found"
            return
        }
    }

    try {
        # Use temporary file approach which is more reliable
        Write-Host "Starting interactive chat session with local model..."
        Write-Host "===================================="
        
        # Create temporary script file
        $tempScript = [System.IO.Path]::GetTempFileName() + ".py"
        
        # Write Python code to temporary file using Here-String
        $pythonCode = @"
print('=== ELR Interactive Model Chat ===')
print('Type your message in English or Chinese.')
print('Type "exit" or "quit" to end the conversation.')
print('Type "help" to see available commands.')
print('===================================')

class SimpleChatModel:
    def __init__(self):
        self.model_name = "simple_chat_model"
        self.version = "1.0"
        self.description = "Minimal chat model for ELR"
        print(f'Initializing model: {self.model_name} v{self.version}')
    
    def predict(self, input_text):
        lower_input = input_text.lower()
        
        # Greetings
        greetings = ["hello", "hi", "你好", "嗨", "hey"]
        for greeting in greetings:
            if greeting in lower_input:
                return f"Carbon-silicon synergy greeting! I'm {self.model_name}, nice to meet you. Your greeting has been received, have a great day!"
        
        # Questions
        questions = ["how are you", "你好吗", "怎么样", "最近好吗"]
        for question in questions:
            if question in lower_input:
                return f"Carbon-silicon synergy response! I'm {self.model_name}, running status is good. Thank you for your concern, wish you all the best!"
        
        # Default response
        return f"Carbon-silicon synergy response! I'm {self.model_name}, received your message: '{input_text}'. Ready to assist you anytime!"
    
    def get_info(self):
        return {
            "model_name": self.model_name,
            "version": self.version,
            "description": self.description,
            "capabilities": ["Text response", "No external dependencies", "Lightweight"]
        }

# Create model instance
model = SimpleChatModel()
print('Model loaded successfully!')

# Main chat loop
while True:
    try:
        user_input = input('You: ')
        if user_input.lower() in ['exit', 'quit', 'q']:
            print('Goodbye!')
            break
        elif user_input.lower() == 'help':
            print('Available commands:')
            print('  exit/quit/q - End the conversation')
            print('  help - Show this help')
            print('  info - Show model information')
        elif user_input.lower() == 'info':
            info = model.get_info()
            print('Model Information:')
            for key, value in info.items():
                print(f'  {key}: {value}')
        else:
            response = model.predict(user_input)
            print(f'Model: {response}')
    except EOFError:
        print('\nGoodbye!')
        break
    except Exception as e:
        print(f'Error: {e}')
"@
        
        # Write code to temporary file
        $pythonCode | Set-Content -Path $tempScript -Encoding UTF8
        
        # Debug: Show the temporary script path and content
        Write-Host "Debug: Temporary script path: $tempScript"
        Write-Host "Debug: Python path: $PythonPath"
        Write-Host "Debug: Python version:"
        & $PythonPath --version
        
        # Try a simpler approach: execute directly with &
        Write-Host "Debug: Executing Python script..."
        
        try {
            # Execute the script directly
            & $PythonPath $tempScript
        } catch {
            Write-Host "Error executing script: $_"
        }
        
        # Clean up temporary file
        if (Test-Path $tempScript) {
            Remove-Item -Path $tempScript -Force -ErrorAction SilentlyContinue
        }

    } catch {
        Write-Host "Error: $_"
    }
}

# Function: Chat with model in container
function Chat-With-Container-Model {
    param(
        [string]$ContainerID,
        [string]$PythonPath
    )

    # Find container
    $container = $null
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -eq $ContainerID) {
            $container = $c
            break
        }
    }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $ContainerID not found"
        return
    }

    if ($container.Status -ne $CONTAINER_STATUS_RUNNING) {
        Write-Host "Error: Container is not running"
        return
    }

    # Create a temporary Python script for container chat
    $tempScriptPath = Join-Path -Path $env:TEMP -ChildPath "elr_chat_container_$(Get-Date -Format 'yyyyMMddHHmmss').py"

    try {
        # Write the chat script to the temporary file
        $chatScript = @"
import sys
import os
import json

print('=== ELR Container Model Chat ===')
print('Type your message in English or Chinese.')
print('Type "exit" or "quit" to end the conversation.')
print('Type "help" to see available commands.')
print('===================================')

# Simulate container interaction
container_id = '$ContainerID'

while True:
    try:
        user_input = input('You: ')
        if user_input.lower() in ['exit', 'quit', 'q']:
            print('Goodbye!')
            break
        elif user_input.lower() == 'help':
            print('Available commands:')
            print('  exit/quit/q - End the conversation')
            print('  help - Show this help')
            print('  info - Show container information')
        elif user_input.lower() == 'info':
            print('Container Information:')
            print(f'  Container ID: {container_id}')
            print('  Status: Running')
            print('  Type: ELR Container')
        else:
            # Simulate container response
            response = f"[Container {container_id}] Carbon-silicon synergy response! I've received your message: '{user_input}'. Processing in container environment..."
            print(f'Container: {response}')
    except EOFError:
        print('\nGoodbye!')
        break
    except Exception as e:
        print(f'Error: {e}')
"@

        # Write the script content
        $chatScript | Set-Content -Path $tempScriptPath -Encoding UTF8

        # Execute the chat script directly
        Write-Host "Starting interactive chat session with container..."
        Write-Host "===================================="
        
        # Directly execute the Python script in the current session
        & $PythonPath $tempScriptPath

    } catch {
        Write-Host "Error: $_"
    } finally {
        # Clean up temporary file
        if (Test-Path $tempScriptPath) {
            Remove-Item -Path $tempScriptPath -Force -ErrorAction SilentlyContinue
        }
    }
}

# Function: Chat with model in sandbox
function Chat-With-Sandbox-Model {
    param(
        [string]$ModelPath,
        [string]$PythonPath
    )

    # Create a temporary Python script for sandbox chat
    $tempScriptPath = Join-Path -Path $env:TEMP -ChildPath "elr_chat_sandbox_$(Get-Date -Format 'yyyyMMddHHmmss').py"

    try {
        # Write the chat script to the temporary file
        $chatScript = @"
import sys
import os

print('=== ELR Sandbox Model Chat ===')
print('Type your message in English or Chinese.')
print('Type "exit" or "quit" to end the conversation.')
print('Type "help" to see available commands.')
print('===================================')

# Simulate sandbox environment
model_path = '$ModelPath'
sandbox_id = 'sandbox-' + str(os.getpid())

print(f'Sandbox initialized: {sandbox_id}')
print(f'Model path: {model_path}')

while True:
    try:
        user_input = input('You: ')
        if user_input.lower() in ['exit', 'quit', 'q']:
            print('Goodbye!')
            break
        elif user_input.lower() == 'help':
            print('Available commands:')
            print('  exit/quit/q - End the conversation')
            print('  help - Show this help')
            print('  info - Show sandbox information')
        elif user_input.lower() == 'info':
            print('Sandbox Information:')
            print(f'  Sandbox ID: {sandbox_id}')
            print(f'  Model Path: {model_path}')
            print('  Status: Active')
            print('  Type: ELR Micro-Model Sandbox')
        else:
            # Simulate sandbox model response
            response = f"[Sandbox {sandbox_id}] Carbon-silicon synergy response! I'm running in isolated sandbox environment. Your message: '{user_input}' has been processed."
            print(f'Sandbox: {response}')
    except EOFError:
        print('\nGoodbye!')
        break
    except Exception as e:
        print(f'Error: {e}')

print(f'Sandbox closed: {sandbox_id}')
"@

        # Write the script content
        $chatScript | Set-Content -Path $tempScriptPath -Encoding UTF8

        # Execute the chat script directly
        Write-Host "Starting interactive chat session with sandbox model..."
        Write-Host "===================================="
        
        # Directly execute the Python script in the current session
        & $PythonPath $tempScriptPath

    } catch {
        Write-Host "Error: $_"
    } finally {
        # Clean up temporary file
        if (Test-Path $tempScriptPath) {
            Remove-Item -Path $tempScriptPath -Force -ErrorAction SilentlyContinue
        }
    }
}

# Function: Chat with local model using PowerShell only
function Chat-With-Local-Model-PowerShell {
    Write-Host "Starting interactive chat session with local model (PowerShell mode)..."
    Write-Host "===================================="
    
    # Define simple chat model logic
    $modelName = "simple_chat_model"
    $modelVersion = "1.0"
    $modelDescription = "Minimal chat model for ELR (PowerShell implementation)"
    
    Write-Host "Initializing model: $modelName v$modelVersion"
    Write-Host "Model loaded successfully!"
    Write-Host ""
    Write-Host "=== ELR Interactive Model Chat (PowerShell Mode) ==="
    Write-Host "Type your message in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "==================================="
    
    # Main chat loop
    while ($true) {
        try {
            $userInput = Read-Host "You"
            if ($userInput.ToLower() -in @('exit', 'quit', 'q')) {
                Write-Host "Goodbye!"
                break
            } elseif ($userInput.ToLower() -eq 'help') {
                Write-Host "Available commands:"
                Write-Host "  exit/quit/q - End the conversation"
                Write-Host "  help - Show this help"
                Write-Host "  info - Show model information"
            } elseif ($userInput.ToLower() -eq 'info') {
                Write-Host "Model Information:"
                Write-Host "  model_name: $modelName"
                Write-Host "  version: $modelVersion"
                Write-Host "  description: $modelDescription"
                Write-Host "  capabilities: Text response, No external dependencies, Lightweight, PowerShell implementation"
            } else {
                # Process user input
                $lowerInput = $userInput.ToLower()
                $response = ""
                
                # Greetings
                $greetings = @('hello', 'hi', 'hey')
                foreach ($greeting in $greetings) {
                    if ($lowerInput -like "*$greeting*") {
                        $response = "Carbon-silicon synergy greeting! I'm $modelName, nice to meet you. Your greeting has been received, have a great day!"
                        break
                    }
                }
                
                # Questions
                if ([string]::IsNullOrEmpty($response)) {
                    $questions = @('how are you', 'how do you do', 'how are things', 'how is it going')
                    foreach ($question in $questions) {
                        if ($lowerInput -like "*$question*") {
                            $response = "Carbon-silicon synergy response! I'm $modelName, running status is good. Thank you for your concern, wish you all the best!"
                            break
                        }
                    }
                }
                
                # Default response
                if ([string]::IsNullOrEmpty($response)) {
                    $response = "Carbon-silicon synergy response! I'm $modelName, received your message: '$userInput'. Ready to assist you anytime!"
                }
                
                Write-Host "Model: $response"
            }
        } catch {
            Write-Host "Error: $($_.Exception.Message)"
        }
    }
}

# Function: Chat with container model using PowerShell only
function Chat-With-Container-Model-PowerShell {
    param(
        [string]$ContainerID
    )
    
    # Find container
    $container = $null
    foreach ($c in $global:CONTAINERS) {
        if ($c.ID -eq $ContainerID) {
            $container = $c
            break
        }
    }
    if ($null -eq $container) {
        Write-Host "Error: Container with ID $ContainerID not found"
        return
    }

    if ($container.Status -ne $CONTAINER_STATUS_RUNNING) {
        Write-Host "Error: Container is not running"
        return
    }
    
    Write-Host "Starting interactive chat session with container (PowerShell mode)..."
    Write-Host "===================================="
    Write-Host "Container: $ContainerID"
    Write-Host ""
    Write-Host "=== ELR Container Model Chat (PowerShell Mode) ==="
    Write-Host "Type your message in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "==================================="
    
    # Main chat loop
    while ($true) {
        try {
            $userInput = Read-Host "You"
            if ($userInput.ToLower() -in @('exit', 'quit', 'q')) {
                Write-Host "Goodbye!"
                break
            } elseif ($userInput.ToLower() -eq 'help') {
                Write-Host "Available commands:"
                Write-Host "  exit/quit/q - End the conversation"
                Write-Host "  help - Show this help"
                Write-Host "  info - Show container information"
            } elseif ($userInput.ToLower() -eq 'info') {
                Write-Host "Container Information:"
                Write-Host "  Container ID: $ContainerID"
                Write-Host "  Status: Running"
                Write-Host "  Type: ELR Container"
            } else {
                # Simulate container response
                $response = "[Container $ContainerID] Carbon-silicon synergy response! I've received your message: '$userInput'. Processing in container environment..."
                Write-Host "Container: $response"
            }
        } catch {
            Write-Host "Error: $($_.Exception.Message)"
        }
    }
}

# Function: Chat with sandbox model using PowerShell only
function Chat-With-Sandbox-Model-PowerShell {
    param(
        [string]$ModelPath
    )
    
    Write-Host "Starting interactive chat session with sandbox model (PowerShell mode)..."
    Write-Host "===================================="
    Write-Host "Model: $ModelPath"
    Write-Host ""
    
    # Generate sandbox ID
    $sandboxID = "sandbox-$(Get-Date -Format 'HHmmssfff')"
    
    Write-Host "=== ELR Sandbox Model Chat (PowerShell Mode) ==="
    Write-Host "Type your message in English or Chinese."
    Write-Host "Type 'exit' or 'quit' to end the conversation."
    Write-Host "Type 'help' to see available commands."
    Write-Host "==================================="
    Write-Host "Sandbox initialized: $sandboxID"
    Write-Host "Model path: $ModelPath"
    
    # Main chat loop
    while ($true) {
        try {
            $userInput = Read-Host "You"
            if ($userInput.ToLower() -in @('exit', 'quit', 'q')) {
                Write-Host "Goodbye!"
                break
            } elseif ($userInput.ToLower() -eq 'help') {
                Write-Host "Available commands:"
                Write-Host "  exit/quit/q - End the conversation"
                Write-Host "  help - Show this help"
                Write-Host "  info - Show sandbox information"
            } elseif ($userInput.ToLower() -eq 'info') {
                Write-Host "Sandbox Information:"
                Write-Host "  Sandbox ID: $sandboxID"
                Write-Host "  Model Path: $ModelPath"
                Write-Host "  Status: Active"
                Write-Host "  Type: ELR Micro-Model Sandbox"
            } else {
                # Simulate sandbox model response
                $response = "[Sandbox $sandboxID] Carbon-silicon synergy response! I'm running in isolated sandbox environment. Your message: '$userInput' has been processed."
                Write-Host "Sandbox: $response"
            }
        } catch {
            Write-Host "Error: $($_.Exception.Message)"
        }
    }
    
    Write-Host "Sandbox closed: $sandboxID"
}

