# app.py
from flask import Flask
from flask_cors import CORS
from routes.proxy import proxy_bp
from init import PYTHON_PORT, SUPABASE_URL, GO_PORT
from time import sleep

COLORS = {
    "GREEN": "\033[32m",
    "YELLOW": "\033[33m",
    "RED": "\033[31m",
    "CYAN": "\033[36m",
    "GRAY": "\033[90m",
    "BOLD": "\033[1m",
    "RESET": "\033[0m",
}

app = Flask(__name__)
CORS(app, origins=[
    "http://localhost:5500",
    "http://127.0.0.1:5500",
    "http://127.0.0.1:5501",
    "http://localhost:5501",
    "http://localhost:5000",
])

app.register_blueprint(proxy_bp)

if __name__ == "__main__":
    print(f"\n{COLORS['BOLD']}{COLORS['CYAN']}╔══════════════════════════════════════╗{COLORS['RESET']}")
    print(f"{COLORS['BOLD']}{COLORS['CYAN']}║{COLORS['RESET']}   {COLORS['BOLD']}Rubium Python Server{COLORS['RESET']}              {COLORS['BOLD']}{COLORS['CYAN']} ║{COLORS['RESET']}")
    print(f"{COLORS['BOLD']}{COLORS['CYAN']}╚══════════════════════════════════════╝{COLORS['RESET']}")
    print(f"\n  {COLORS['GRAY']}Proxy:{COLORS['RESET']} :{PYTHON_PORT} → Go : {GO_PORT} → Supabase")
    print(f"  {COLORS['GRAY']}AI:{COLORS['RESET']}    : {PYTHON_PORT} → DeepSeek")
    print(f"  {COLORS['GRAY']}CORS:{COLORS['RESET']}  localhost: 5500, 5501, 5000")
    print(f"\n  {COLORS['YELLOW']}Starting server...{COLORS['RESET']}\n")
    sleep(1)
    app.run(debug=False, port=PYTHON_PORT)