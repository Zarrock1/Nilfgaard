from fastapi import FastAPI, HTTPException
import redis.asyncio as redis
from redis.exceptions import ConnectionError
import requests, httpx, json

app = FastAPI()

# Настройки Redis
REDIS_HOST = "localhost"
REDIS_PORT = 6379
REDIS_DB = 0
redis_client = redis.Redis(host=REDIS_HOST, port=REDIS_PORT, db=REDIS_DB)

token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJBc3NpcmV2YXJBbmFoaWQiLCJleHAiOjE3MzYzNDcxNzYsInVzZXJhY2Nlc3MiOlsidXNlcjpsaXN0OnJlYWQiLCJ1c2VyOmZ1bGxOYW1lOndyaXRlIiwidXNlcjpkYXRhOnJlYWQiLCJ1c2VyOnJvbGVzOnJlYWQiLCJ1c2VyOnJvbGVzOndyaXRlIiwidXNlcjpibG9jazpyZWFkIiwidXNlcjpibG9jazp3cml0ZSIsImNvdXJzZTppbmZvOndyaXRlIiwiY291cnNlOnRlc3RMaXN0IiwiY291cnNlOnRlc3Q6cmVhZCIsImNvdXJzZTp0ZXN0OndyaXRlIiwiY291cnNlOnRlc3Q6YWRkIiwiY291cnNlOnRlc3Q6ZGVsIiwiY291cnNlOnVzZXJMaXN0IiwiY291cnNlOnVzZXI6YWRkIiwiY291cnNlOnVzZXI6ZGVsIiwiY291cnNlOmFkZCIsImNvdXJzZTpkZWwiXSwidXNlcmxvZ2luIjoiTGl6YSJ9.RXPDzO73WoBUkQVCDBMHZ4efg6JF1Anc85gQTTK1nwc"

@app.get("/identification/{chat_id}")
async def _(chat_id: str):
    value = await redis_client.get(chat_id)
    if value is None:
        raise HTTPException(status_code=404)
    return value.decode()

@app.get("/login/{value}")
async def _(value: str):
    chat_id, status, tokenVHOD = value.split(":")
    await redis_client.set(chat_id, f'{status}:{tokenVHOD}')

@app.get("/loginGIT/{tokenVHOD}")
async def _(tokenVHOD: str):
    async with httpx.AsyncClient() as client:
        response = await client.get(f"http://78.136.201.177:8080/api/token?userId={tokenVHOD}")
        response.raise_for_status()
        return response.json()

@app.get("/loginYAN/{tokenVHOD}")
async def _(tokenVHOD: str):
    async with httpx.AsyncClient() as client:
        response = await client.get(f"http://78.136.201.177:8080/api/token?userId={tokenVHOD}")
        response.raise_for_status()
        return response.json()
    
@app.get("/disciplines/")
async def _():
    headers = {"Authorization": f"Bearer {token}"}  
    url = "http://91.245.227.51:3000/api/disciplins/"
    async with httpx.AsyncClient() as client:
        response = await client.get(url, headers=headers)
        response.raise_for_status()
        return response.json()
    
@app.get("/disciplinesNONALL/{id}")
async def _(id: str):
    headers = {"Authorization": f"Bearer {token}"}  
    url = "http://91.245.227.51:3000/api/disciplins/"
    async with httpx.AsyncClient() as client:
        response = await client.get(url+id, headers=headers)
        response.raise_for_status()
        return response.json()
    
@app.get("/users/")
async def _():
    headers = {"Authorization": f"Bearer {token}"}  
    url = "http://91.245.227.51:3000/api/users/"    
    async with httpx.AsyncClient() as client:
        response = await client.get(url, headers=headers) 
        response.raise_for_status() 
        return response.json()
    
@app.get("/usersNONALL/{id}")
async def _(id: str):
    headers = {"Authorization": f"Bearer {token}"}  
    url = "http://91.245.227.51:3000/api/users/"
    async with httpx.AsyncClient() as client:
        response = await client.get(url+id, headers=headers)
        response.raise_for_status()
        return response.json()