# ELR Container Installation Guide

This guide explains how to install and run the Enlightenment Lighthouse Runtime (ELR) container on Windows.

## Prerequisites

- Windows 10 or later
- PowerShell 5.1 or later
- Python 3.8+ (optional, for running Python scripts)
- GCC (optional, for compiling C programs)

## Installation Methods

### Method 1: One-line Installation (Recommended)

Similar to uv installation, you can install ELR with a single PowerShell command:

```powershell
powershell -ExecutionPolicy ByPass -c "irm https://example.com/elr/install.ps1 | iex"
```

This command will:
1. Download the latest ELR files from GitHub
2. Install ELR to your user profile directory (or a custom location you specify)
3. Add ELR to your PATH environment variable
4. Create a wrapper script for easy access

### Method 2: Local Installation

If you already have the ELR source code, you can install it locally:

1. Navigate to the ELR directory
2. Run the local installation script:

```powershell
.nstall.ps1
```

## Installation Directory

By default, ELR is installed to:
```
%USERPROFILE%\ELR
```

You can specify a custom installation directory during the installation process.

## Directory Structure

After installation, ELR will have the following directory structure:

```
ELR/
├── bin/            # Executable files and scripts
├── lib/            # Library files and micro_model
├── config/         # Configuration files
├── models/         # Model files
├── containers/     # Container data
└── README.md       # Documentation
```

## Usage

Once installed, you can use the `elr` command from anywhere in your terminal:

### Basic Commands

- `elr start` - Start the ELR runtime
- `elr stop` - Stop the ELR runtime
- `elr status` - Check runtime status
- `elr list` - List all containers
- `elr help` - Show help information

### Advanced Commands

- `elr create --name <name> --image <image>` - Create a new container
- `elr run --name <name> --image <image>` - Create and start a new container
- `elr start-container --id <container-id>` - Start a container
- `elr stop-container --id <container-id>` - Stop a container
- `elr delete --id <container-id>` - Delete a container
- `elr exec --id <container-id> --command <command>` - Execute a command in a container

### Model Commands

- `elr run-python --source <script.py>` - Run a Python script
- `elr run-python --code '<python code>'` - Run Python code directly
- `elr chat` - Start interactive chat with default local model
- `elr chat --model <model.py>` - Start chat with custom model

## Examples

### Start ELR Runtime

```powershell
elr start
```

### Create and Run a Container

```powershell
elr run --name my-container --image ubuntu:latest
```

### Run a Python Script

```powershell
elr run-python --source script.py
```

### Start Interactive Chat

```powershell
elr chat
```

## Troubleshooting

### Python Not Found

If you encounter Python-related errors, ensure Python 3.8+ is installed and in your PATH.

### GCC Not Found

If you need to compile C programs, install GCC through one of these methods:
- MinGW-w64: https://www.mingw-w64.org/
- MSYS2: https://www.msys2.org/
- Cygwin: https://www.cygwin.com/

### Network Issues

If you have network connectivity issues, check your firewall settings and ensure you can access GitHub.

## Updates

To update ELR to the latest version, simply run the installation command again:

```powershell
powershell -ExecutionPolicy ByPass -c "irm https://example.com/elr/install.ps1 | iex"
```

## Uninstallation

To uninstall ELR:

1. Remove the installation directory (default: `%USERPROFILE%\ELR`)
2. Remove the ELR bin directory from your PATH environment variable

## Support

For issues and feature requests, please visit the GitHub repository or contact the development team.
