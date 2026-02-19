import sys
print(f"Python version: {sys.version}")
print(f"Python executable: {sys.executable}")
print(f"Current directory: {sys.path}")

try:
    import websockets
    print("✓ websockets module imported successfully")
except ImportError as e:
    print(f"✗ Failed to import websockets: {e}")

try:
    import asyncio
    print("✓ asyncio module imported successfully")
except ImportError as e:
    print(f"✗ Failed to import asyncio: {e}")

try:
    import json
    print("✓ json module imported successfully")
except ImportError as e:
    print(f"✗ Failed to import json: {e}")

print("Environment test completed.")
