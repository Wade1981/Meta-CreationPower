#!/bin/bash

# Startup script for Enlightenment Lighthouse Runtime (ELR)

set -e

# Script version
SCRIPT_VERSION="1.0.0"

# Working directory
WORKDIR="$(cd "$(dirname "$0")/.." && pwd)"

# Log directory
LOG_DIR="$WORKDIR/logs"
mkdir -p "$LOG_DIR"

# Log file
LOG_FILE="$LOG_DIR/elr-start.log"

# Function: log message
log() {
    local timestamp=$(date "%Y-%m-%d %H:%M:%S")
    local level="$1"
    local message="$2"
    echo "[$timestamp] [$level] $message" | tee -a "$LOG_FILE"
}

# Function: check if command exists
check_command() {
    if ! command -v "$1" &> /dev/null; then
        log "ERROR" "Command '$1' not found"
        return 1
    fi
    return 0
}

# Function: build ELR
build_elr() {
    log "INFO" "Building Enlightenment Lighthouse Runtime..."
    
    cd "$WORKDIR"
    
    # Check if Go is installed
    if ! check_command "go"; then
        log "ERROR" "Go is not installed, cannot build ELR"
        return 1
    fi
    
    # Build ELR
    log "INFO" "Building ELR binary..."
    go build -o "$WORKDIR/bin/elr" "$WORKDIR/cli"
    
    if [ $? -eq 0 ]; then
        log "INFO" "ELR built successfully"
        return 0
    else
        log "ERROR" "Failed to build ELR"
        return 1
    fi
}

# Function: install ELR
install_elr() {
    log "INFO" "Installing Enlightenment Lighthouse Runtime..."
    
    # Check if bin directory exists
    if [ ! -d "$WORKDIR/bin" ]; then
        mkdir -p "$WORKDIR/bin"
    fi
    
    # Build ELR if not already built
    if [ ! -f "$WORKDIR/bin/elr" ]; then
        if ! build_elr; then
            return 1
        fi
    fi
    
    # Add ELR to PATH (temporarily)
    export PATH="$WORKDIR/bin:$PATH"
    
    log "INFO" "ELR installed successfully"
    return 0
}

# Function: start ELR
start_elr() {
    log "INFO" "Starting Enlightenment Lighthouse Runtime..."
    
    # Install ELR if not already installed
    if ! install_elr; then
        return 1
    fi
    
    # Start ELR runtime
    log "INFO" "Starting ELR runtime..."
    "$WORKDIR/bin/elr" start &
    
    # Save PID
    ELR_PID=$!
    echo "$ELR_PID" > "$WORKDIR/elr.pid"
    
    log "INFO" "ELR runtime started with PID $ELR_PID"
    
    # Wait a bit for ELR to start
    sleep 2
    
    # Check if ELR is running
    if ps -p "$ELR_PID" > /dev/null; then
        log "INFO" "ELR runtime started successfully"
        return 0
    else
        log "ERROR" "ELR runtime failed to start"
        return 1
    fi
}

# Function: stop ELR
stop_elr() {
    log "INFO" "Stopping Enlightenment Lighthouse Runtime..."
    
    # Check if ELR PID file exists
    if [ ! -f "$WORKDIR/elr.pid" ]; then
        log "ERROR" "ELR PID file not found, ELR may not be running"
        return 1
    fi
    
    # Read PID
    ELR_PID=$(cat "$WORKDIR/elr.pid")
    
    # Check if ELR is running
    if ! ps -p "$ELR_PID" > /dev/null; then
        log "ERROR" "ELR is not running"
        return 1
    fi
    
    # Stop ELR
    log "INFO" "Stopping ELR runtime with PID $ELR_PID..."
    kill "$ELR_PID"
    
    # Wait for ELR to stop
    sleep 2
    
    # Check if ELR is stopped
    if ! ps -p "$ELR_PID" > /dev/null; then
        log "INFO" "ELR runtime stopped successfully"
        rm "$WORKDIR/elr.pid"
        return 0
    else
        log "ERROR" "ELR runtime failed to stop"
        return 1
    fi
}

# Main function
main() {
    log "INFO" "Enlightenment Lighthouse Runtime Startup Script v$SCRIPT_VERSION"
    log "INFO" "Working directory: $WORKDIR"
    
    # Parse command line arguments
    if [ $# -eq 0 ]; then
        log "INFO" "No command specified, starting ELR"
        start_elr
        exit $?
    fi
    
    command="$1"
    
    case "$command" in
        "start")
            start_elr
            ;;
        "stop")
            stop_elr
            ;;
        "build")
            build_elr
            ;;
        "install")
            install_elr
            ;;
        "help")
            echo "Enlightenment Lighthouse Runtime Startup Script"
            echo "Usage: $0 [command]"
            echo ""
            echo "Commands:"
            echo "  start    Start ELR runtime"
            echo "  stop     Stop ELR runtime"
            echo "  build    Build ELR binary"
            echo "  install  Install ELR"
            echo "  help     Show this help message"
            ;;
        *)
            log "ERROR" "Unknown command: $command"
            echo "Unknown command: $command"
            echo "Use '$0 help' for available commands"
            exit 1
            ;;
    esac
    
    exit $?
}

# Run main function
main "$@"
