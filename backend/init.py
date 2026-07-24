# init.py
import os
from dotenv import load_dotenv
import yaml

load_dotenv(".env")

with open("config.yaml") as f:
    config = yaml.safe_load(f)

def _resolve_env(value):
    if isinstance(value, str) and value.startswith("${") and value.endswith("}"):
        return os.environ.get(value[2:-1], "")
    return value

SUPABASE_URL = config["supabase"]["url"]
SUPABASE_ANON_KEY = _resolve_env(config["supabase"]["anon_key"])
SUPABASE_SERVICE_KEY = _resolve_env(config["supabase"]["service_key"])
DEEPSEEK_API_URL = config["deepseek"]["api_url"]
OPENAI_API_KEY = _resolve_env(config["openai"]["api_key"])
GO_PORT = config["server"]["go_port"]
PYTHON_PORT = config["server"]["python_port"]
with open("backend/model/AI_API_KEY.txt") as file:
    DEEPSEEK_API_KEY = file.read().strip()

def supabase_headers(use_service_role=False):
    key = SUPABASE_SERVICE_KEY if use_service_role and SUPABASE_SERVICE_KEY else SUPABASE_ANON_KEY
    return {
        "apikey": key,
        "Authorization": f"Bearer {key}",
        "Content-Type": "application/json",
    }