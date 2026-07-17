import json
import os
import requests
from flask import Flask, request, jsonify
from flask_cors import CORS
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)
CORS(
    app,
    origins=[
        "http://localhost:5500",
        "http://127.0.0.1:5500",
        "http://localhost:5000",
    ],
)

# ─── Конфигурация ───
with open("backend/model/AI_API_KEY.txt") as f:
    DEEPSEEK_API_KEY = f.read().strip()

DEEPSEEK_API_URL = "https://api.deepseek.com/v1/chat/completions"
SUPABASE_URL = "https://nrmihghshpteellhmzuh.supabase.co"
SUPABASE_ANON_KEY = os.environ.get("SUPABASE_ANON_KEY", "")
SUPABASE_SERVICE_KEY = os.environ.get("SUPABASE_SERVICE_KEY", "")

def supabase_headers(use_service_role=False):
    key = (
        SUPABASE_SERVICE_KEY
        if use_service_role and SUPABASE_SERVICE_KEY
        else SUPABASE_ANON_KEY
    )
    return {
        "apikey": key,
        "Authorization": f"Bearer {key}",
        "Content-Type": "application/json",
    }
    
def get_users():
    try:
        r = requests.get(
            f"{SUPABASE_URL}/rest/v1/users?"
            f"select=id,first_name,last_name,email,role,photo_url",
            headers=supabase_headers(),
        )
        if r.status_code == 200:
            return r.json()
    except Exception as e:
        print(f"get_users: {e}")
    return []

def get_tasks():
    try:
        r = requests.get(
            f"{SUPABASE_URL}/rest/v1/tasks?"
            f"select=*",
            headers=supabase_headers(),
        )
        if r.status_code == 200:
            return r.json()
    except Exception as e:
        print(f"get tasks error: {e}")

request_users = get_users()
request_tasks = get_tasks()
for i in request_tasks[0]:
    print(f"{i} : {request_tasks[0][i]}")