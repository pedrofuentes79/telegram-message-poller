#!/usr/bin/env python3
import subprocess

def list_audio_sinks():
    result = subprocess.run(['pactl', 'list', 'sinks', 'short'], 
                          capture_output=True, text=True, check=True)
    
    print("Available audio sinks:")
    print("=" * 60)
    for line in result.stdout.strip().split('\n'):
        parts = line.split('\t')
        if len(parts) >= 2:
            sink_id = parts[0]
            sink_name = parts[1]
            status = parts[3] if len(parts) > 3 else "UNKNOWN"
            print(f"ID: {sink_id}")
            print(f"Name: {sink_name}")
            print(f"Status: {status}")
            print("-" * 60)
    
    # Show current default
    current = subprocess.run(['pactl', 'get-default-sink'], 
                           capture_output=True, text=True, check=True)
    print(f"\nCurrently active: {current.stdout.strip()}")

if __name__ == '__main__':
    list_audio_sinks()

