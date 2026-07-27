# routes/ai.py
from flask import Blueprint, request, jsonify
import subprocess
import tempfile
import os
import requests
from init import DEEPSEEK_API_KEY, DEEPSEEK_API_URL

ai_bp = Blueprint("ai", __name__)

CHECKER_SYSTEM_PROMPT = """Ты — строгий проверяющий заданий по информатике (ЕГЭ/ОГЭ).
Оцени решение ученика и ответь ТОЛЬКО числом:
1 — если решение правильное
0 — если решение неправильное

Никаких пояснений, только 1 или 0."""


@ai_bp.route("/ai/v1/run", methods=["POST"])
def run_code():
    data = request.get_json()
    code = data.get("code", "")
    language = data.get("language", "python")

    if not code:
        return jsonify({"error": "no code provided"}), 400

    if language == "python":
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

            if result.returncode == 0:
                return jsonify({"output": result.stdout, "error": None})
            else:
                return jsonify({"output": result.stdout, "error": result.stderr})
        except subprocess.TimeoutExpired:
            return jsonify({"output": "", "error": "Time limit exceeded (5s)"})
        except Exception as e:
            return jsonify({"output": "", "error": str(e)})

    return jsonify({"error": f"unsupported language: {language}"}), 400


@ai_bp.route("/ai/v1/check", methods=["POST"])
def ai_check():
    data = request.get_json()
    task_content = data.get("content", "")
    user_code = data.get("answer", "")
    correct_answer = data.get("correct_answer", "")

    if not task_content or not user_code:
        return jsonify({"correct": False, "error": "missing fields"}), 400

    try:
        with tempfile.NamedTemporaryFile(mode="w", suffix=".py", delete=False) as f:
            f.write(user_code)
            tmpfile = f.name
        result = subprocess.run(["python3", tmpfile], capture_output=True, text=True, timeout=5)
        os.unlink(tmpfile)
        user_output = (result.stdout + result.stderr).strip()
    except:
        user_output = ""

    if user_output and correct_answer:
        normalized_out = "".join(user_output.split())
        normalized_ans = "".join(str(correct_answer).split())
        if normalized_out == normalized_ans:
            return jsonify({"correct": True, "method": "exact"})

    try:
        resp = requests.post(
            DEEPSEEK_API_URL,
            headers={
                "Authorization": f"Bearer {DEEPSEEK_API_KEY}",
                "Content-Type": "application/json"
            },
            json={
                "model": "deepseek-chat",
                "messages": [
                    {"role": "system", "content": CHECKER_SYSTEM_PROMPT},
                    {"role": "user", "content": f"Задание:\n{task_content[:500]}\n\nОжидаемый ответ:\n{correct_answer}\n\nКод/ответ ученика:\n{user_code[:500]}\n\nВывод программы:\n{user_output[:200]}"}
                ],
                "temperature": 0,
                "max_tokens": 2
            },
            timeout=10
        )
        answer = resp.json()["choices"][0]["message"]["content"].strip()
        return jsonify({
            "correct": "1" in answer,
            "method": "ai",
            "ai_raw": answer
        })
    except Exception as e:
        return jsonify({"correct": False, "error": str(e)}), 500