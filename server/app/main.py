from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
import httpx
import uvicorn
import os
import subprocess
import threading
import time

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Start Go backend in a separate thread
def run_go_backend():
    subprocess.run(["go", "build", "-o", "server"], 
                  cwd=os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    subprocess.Popen(["./server"],
                    cwd=os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
    
threading.Thread(target=run_go_backend, daemon=True).start()
time.sleep(5)  # Wait for Go backend to start

# Forward all requests to Go backend
@app.api_route("/{path:path}", methods=["GET", "POST", "PUT", "DELETE"])
async def proxy(path: str, request: Request):
    async with httpx.AsyncClient() as client:
        url = f"http://localhost:8888/{path}"
        
        # Get the request body
        body = await request.body()
        
        # Forward the request with the same method and headers
        response = await client.request(
            method=request.method,
            url=url,
            content=body,
            headers=dict(request.headers),
            follow_redirects=True
        )
        
        return response.json()

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
