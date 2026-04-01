#!/usr/bin/env python3
import sys
import json

def main():
    try:
        # 1. Read request from stdin
        line = sys.stdin.readline()
        if not line:
            return
        
        req = json.loads(line)
        method = req.get("method")
        req_id = req.get("id")

        # 2. Protocol Dispatch
        if method == "handshake":
            resp = {
                "jsonrpc": "2.0",
                "id": req_id,
                "result": {
                    "name": "example-python",
                    "version": "1.0.0",
                    "api_version": 1,
                    "checks_supported": ["python_check"]
                }
            }
            print(json.dumps(resp))
            sys.stdout.flush()

        elif method == "run_check":
            resp = {
                "jsonrpc": "2.0",
                "id": req_id,
                "result": {
                    "status": "pass",
                    "finding": {
                        "id": "F-PY-001",
                        "check_id": "python_check",
                        "severity": "info",
                        "title": "Python Plugin Running",
                        "description": "Python-based plugin is communicating correctly via JSON-RPC.",
                        "evidence": {"interpreter": sys.version},
                        "remediation": "None"
                    }
                }
            }
            print(json.dumps(resp))
            sys.stdout.flush()

    except Exception as e:
        # All non-protocol output to stderr
        sys.stderr.write(f"Error: {str(e)}\n")
        sys.exit(1)

if __name__ == "__main__":
    main()
