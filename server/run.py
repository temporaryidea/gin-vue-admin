import os
import subprocess
import sys

def build_go():
    # Build the Go binary
    build_process = subprocess.run(['go', 'build', '-o', 'server'], 
                                 cwd=os.path.dirname(os.path.abspath(__file__)))
    if build_process.returncode != 0:
        sys.exit(1)

def run_server():
    # Run the compiled binary
    server_process = subprocess.run(['./server'])
    sys.exit(server_process.returncode)

if __name__ == '__main__':
    build_go()
    run_server()
