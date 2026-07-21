# app.py
from flask import Flask
from flask_cors import CORS
from routes.tasks import tasks_bp

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
    app.run(debug=True, port=5080)