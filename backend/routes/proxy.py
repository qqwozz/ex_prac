# routes/proxy.py
from flask import Blueprint, request, Response
import requests
import json
import time
from init import GO_PORT

GO_URL = f"http://localhost:{GO_PORT}"

proxy_bp = Blueprint("proxy", __name__)

COLORS = {
    "GREEN": "\033[32m",
    "YELLOW": "\033[33m",
    "RED": "\033[31m",
    "CYAN": "\033[36m",
    "GRAY": "\033[90m",
    "RESET": "\033[0m",
}

@proxy_bp.route("/api/v1/<path:endpoint>", methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
def proxy(endpoint):
    url = f"{GO_URL}/api/v1/{endpoint}"
    headers = {k: v for k, v in request.headers if k.lower() != "host"}
    method = request.method

    # Request log
    print(f"\n{COLORS['CYAN']}→ PROXY {method} /api/v1/{endpoint}{COLORS['RESET']}")
    if request.args:
        print(f"  {COLORS['GRAY']}Params: {dict(request.args)}{COLORS['RESET']}")
    body = request.get_json(silent=True)
    if body:
        print(f"  {COLORS['GRAY']}Body: {json.dumps(body, ensure_ascii=False)[:150]}{COLORS['RESET']}")

    start = time.time()
    try:
        if method == "GET":
            resp = requests.get(url, params=request.args, headers=headers, timeout=10)
        else:
            resp = requests.request(method, url, json=body, headers=headers, timeout=10)

        elapsed = (time.time() - start) * 1000
        color = COLORS["GREEN"] if resp.ok else COLORS["RED"]

        print(f"← {color}{resp.status_code}{COLORS['RESET']} {COLORS['GRAY']}({elapsed:.0f}ms){COLORS['RESET']}")
        if not resp.ok:
            print(f"  {COLORS['RED']}{resp.text[:200]}{COLORS['RESET']}")

    except requests.exceptions.ConnectionError:
        print(f"  {COLORS['RED']}✗ GO server not reachable at {GO_URL}{COLORS['RESET']}")
        return Response("GO backend unavailable", status=502)
    except requests.exceptions.Timeout:
        print(f"  {COLORS['RED']}✗ GO server timeout{COLORS['RESET']}")
        return Response("GO backend timeout", status=504)

    return Response(resp.content, resp.status_code, resp.headers.items())