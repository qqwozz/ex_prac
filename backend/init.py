from dotenv import load_dotenv
import os

load_dotenv()

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