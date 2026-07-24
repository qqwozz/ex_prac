# app.py
from flask import Flask
from flask_cors import CORS
from routes.tasks import tasks_bp
import requests
from init import SUPABASE_URL, supabase_headers, DEEPSEEK_API_URL, DEEPSEEK_API_KEY
from time import sleep

app = Flask(__name__)
CORS(app, origins=[
    "http://localhost:5500",
    "http://127.0.0.1:5500",
    "http://127.0.0.1:5501",
    "http://localhost:5501",
    "http://localhost:5000",
])

app.register_blueprint(tasks_bp)

if __name__ == "__main__":
    print("Supabase:", requests.get(
        f"{SUPABASE_URL}/rest/v1/users?limit=1",
        headers=supabase_headers()
    ).status_code)

    headers = {"Authorization": f"Bearer {DEEPSEEK_API_KEY}"}
    print("DeepSeek:", requests.get("https://api.deepseek.com/v1/models", headers=headers).status_code)

    print("\nStarting server...")
    app.run(debug=True, port=5080)
    