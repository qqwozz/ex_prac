# routes/ai.py
from flask import Blueprint, request, jsonify
import subprocess
import tempfile
import os
import requests
from init import DEEPSEEK_API_KEY, DEEPSEEK_API_URL

ai_bp = Blueprint("ai", __name__)

@ai_bp.route("/ai/v1/check_code", methods=["POST"])
def check_code():
    data = request.get_json()
    code = data.get("code", "")
    expected_output = data.get("expected_output", "").strip()
    
    if not code:
        return jsonify({"correct": False, "error": "no code provided"}), 400
    
    try:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".py", delete=False) as f:
            f.write(code)
            tmpfile = f.name
        
        result = subprocess.run(
            ["python3", tmpfile],
            capture_output=True,
            text=True,
            timeout=5
        )
        os.unlink(tmpfile)
        
        output = result.stdout.strip()
        error = result.stderr.strip()
        
        if result.returncode != 0:
            return jsonify({
                "correct": False,
                "error": error,
                "output": output
            })
        is_correct = output == expected_output
        
        return jsonify({
            "correct": is_correct,
            "output": output,
            "expected": expected_output,
            "error": None
        })
        
    except subprocess.TimeoutExpired:
        return jsonify({"correct": False, "error": "Time limit exceeded (5s)"})
    except Exception as e:
        return jsonify({"correct": False, "error": str(e)})