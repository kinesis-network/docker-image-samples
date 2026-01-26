from flask import Flask, request, make_response
import uuid

app = Flask(__name__)

# This list stays in the RAM of the specific server instance
issued_cookies = set()

@app.route('/')
def index():
    user_cookie = request.cookies.get('repro_id')

    # CASE 1: No cookie at all
    if not user_cookie:
        new_id = uuid.uuid4().hex
        issued_cookies.add(new_id)

        resp = make_response(f"NEW SESSION: Created ID [{new_id}]\n")
        resp.set_cookie('repro_id', new_id, max_age=10)
        return resp

    # CASE 2: Cookie exists, check if THIS server created it
    if user_cookie in issued_cookies:
        return f"OK: Persistent session. I recognize ID [{user_cookie}].\n"
    else:
        # CASE 3: Cookie exists, but this server doesn't know it!
        return (
            f"DISCREPANCY DETECTED!\n"
            f"Client sent ID: [{user_cookie}]\n"
            f"Server Status: I have no record of issuing this ID.\n"
        )

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
