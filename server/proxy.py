from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
import uvicorn
import httpx
import subprocess
import os
import signal
import time
import threading

app = FastAPI()

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["https://frontend-web-app-ca74ufzk.devinapps.com"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Start Go server in a separate thread
def run_go_server():
    global go_process
    go_process = subprocess.Popen(["./server"], cwd=os.getcwd())

threading.Thread(target=run_go_server, daemon=True).start()
time.sleep(2)  # Wait for Go server to start

@app.middleware("http")
async def proxy_middleware(request: Request, call_next):
    url = httpx.URL(path=request.url.path, query=request.url.query.encode("utf-8"))
    go_url = f"http://localhost:8888{url}"
    
    async with httpx.AsyncClient() as client:
        try:
            response = await client.request(
                method=request.method,
                url=go_url,
                headers=dict(request.headers),
                content=await request.body()
            )
            return response
        except Exception as e:
            return {"error": str(e)}

if __name__ == "__main__":
    port = int(os.getenv("PORT", 8000))
    uvicorn.run(app, host="0.0.0.0", port=port)
