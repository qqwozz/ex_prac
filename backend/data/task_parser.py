# data/task_parser.py
import requests
from init import SUPABASE_URL, supabase_headers

def parse(subject=None):
    url = f"{SUPABASE_URL}/rest/v1/tasks?select=*"
    if subject:
        url += f"&subject=eq.{subject}"

    try:
        r = requests.get(url, headers=supabase_headers())
        if r.status_code == 200:
            raw_tasks = r.json()
            return [_map_task(t) for t in raw_tasks]
    except Exception as e:
        print(f"parse error: {e}")

    return []


def _map_task(t):
    return {
        "id": t["id"],
        "subject": t.get("subject", "math"),
        "exam": t.get("exam_type", "ege"),
        "type": _infer_type(t.get("task_type")),
        "number": t.get("task_number"),
        "difficulty": _map_level(t.get("level")),
        "topic": t.get("topic"),
        "text": t.get("content", ""),
        "options": None,
        "answer": t.get("answer", ""),
        "explanation": t.get("solution") or "",
        "image_url": None,
    }


def _infer_type(task_type):
    mapping = {
        "fipi": "number",
        "choice": "choice",
        "multi": "multi",
        "string": "string",
    }
    return mapping.get(task_type, "number")


def _map_level(level):
    mapping = {"base": 1, "medium": 3, "hard": 5}
    return mapping.get(level, 2)

def get_by_id(task_id):
    url = f"{SUPABASE_URL}/rest/v1/tasks?id=eq.{task_id}&limit=1"
    try:
        r = requests.get(url, headers=supabase_headers())
        if r.status_code == 200:
            data = r.json()
            return _map_task(data[0]) if data else None
    except Exception as e:
        print(f"get_by_id error: {e}")
    return None