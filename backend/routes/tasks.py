# routes/tasks.py
from flask import Blueprint, request, jsonify
from data.task_parser import parse, get_by_id
from services.checker import check_answer
from data.task_parser import parse, get_by_id

tasks_bp = Blueprint("tasks", __name__)


@tasks_bp.route("/api/v1/tasks", methods=["GET"])
def get_tasks():
    subject = request.args.get("subject", "math")
    topic = request.args.get("topic")
    difficulty = request.args.get("difficulty", type=int)
    limit = request.args.get("limit", 1, type=int)

    tasks = parse(subject)

    if topic:
        tasks = [t for t in tasks if t["topic"] == topic]
    if difficulty:
        tasks = [t for t in tasks if t["difficulty"] <= difficulty]

    return jsonify(tasks[:limit])


@tasks_bp.route("/api/v1/check", methods=["POST"])
def check():
    data = request.get_json()
    task_id = data.get("task_id")
    user_answer = data.get("answer", "")

    task = get_by_id(task_id)

    if not task:
        return jsonify({"error": "task not found"}), 404

    correct = check_answer(task["type"], task["answer"], user_answer)

    return jsonify({
        "correct": correct,
        "correct_answer": task["answer"],
        "explanation": task["explanation"],
    })