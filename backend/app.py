from flask import Flask
from flask_cors import CORS
from routes.proxy import proxy_bp
from routes.ai import ai_bp
from init import PYTHON_PORT
from waitress import serve

app = Flask(__name__)
CORS(app, origins=[
    "http://localhost:5500",
    "http://127.0.0.1:5500",
    "http://127.0.0.1:5501",
    "http://localhost:5501",
    "http://localhost:5000",
])

app.register_blueprint(proxy_bp)
app.register_blueprint(ai_bp)

if __name__ == "__main__":
    print(f"\n--- Rubium Python Server (waitress) ---")
    print(f"  Port: {PYTHON_PORT}")
    print(f"  CORS: localhost: 5500, 5501, 5000\n")
    serve(app, host="0.0.0.0", port=PYTHON_PORT)