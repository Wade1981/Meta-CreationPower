#!/usr/bin/env python3

import sys

# Simple test model
if __name__ == "__main__":
    if len(sys.argv) > 1:
        input_text = sys.argv[1]
        print(f"Test model processed: {input_text}")
    else:
        print("Test model called without input")
