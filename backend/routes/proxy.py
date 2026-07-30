from flask import Blueprint, request, Response
import requests
import json
import time
from init import GO_PORT

GO_URL = f"http://localhost:{GO_PORT}"
proxy_bp = Blueprint("proxy", __name__)

@proxy_bp.route("/api/v1/<path:endpoint>", methods=["GET", "POST", "PUT", "DELETE", "PATCH"])
def proxy(endpoint):
    url = f"{GO_URL}/api/v1/{endpoint}"
    headers = {k: v for k, v in request.headers if k.lower() != "host"}
    method = request.method

    body = request.get_json(silent=True)
    
    start = time.time()
    try:
        resp = requests.request(method, url, params=request.args, json=body, headers=headers, timeout=10)
        elapsed = (time.time() - start) * 1000
        print(f"→ {method} /api/v1/{endpoint} → {resp.status_code} ({elapsed:.0f}ms)")
    except requests.exceptions.ConnectionError:
        return Response("GO backend unavailable", status=502)
    except requests.exceptions.Timeout:
        return Response("GO backend timeout", status=504)

    return Response(resp.content, status=resp.status_code, mimetype='application/json')