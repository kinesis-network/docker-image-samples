import subprocess
import os
import uuid
import yaml
from flask import Flask, request, jsonify, send_file
from celery import Celery
from pathlib import Path

# --- Flask app setup ---
app = Flask(__name__)
work_dir = '/work'

# --- Celery setup ---
celery = Celery(
    "server",
    broker="redis://localhost:6379/0",
    backend="redis://localhost:6379/0"
)

# --- Celery Task ---
@celery.task(bind=True)
def run_inference_task(self, job_id, config_data):
    log_file = f"{work_dir}/logs/{job_id}.log"
    input_file = f"{work_dir}/inputs/{job_id}.yaml"
    output_dir = f'{work_dir}/outputs/{job_id}'
    os.makedirs(output_dir, exist_ok=True)

    config_data['ckpt_dir'] = 'ckpts'
    config_data['output_dir'] = output_dir
    config_data['each_example_n_times'] = 1

    with open(input_file, "w") as f:
        yaml.dump(config_data, f)

    cmd = ["python3", "inference.py", "--config-file", input_file]
    with open(log_file, "w") as lf:
        process = subprocess.Popen(
            cmd,
            stdout=lf,
            stderr=subprocess.STDOUT,
            text=True
        )
        self.update_state(state='STARTED', meta={'log_file': log_file})
        process.wait()
    status = "success" if process.returncode == 0 else "failed"

    outdir = Path(output_dir)
    for file_path in outdir.iterdir():
        if file_path.is_file():
            video_path = str(file_path)
            break

    return {
        "status": status,
        "video_path": video_path,
        "log_file": log_file
    }

@app.route("/run", methods=["POST"])
def run_inference():
    data = request.get_json(force=True)
    job_id = str(uuid.uuid4())
    task = run_inference_task.delay(job_id, data)
    return {"job_id": job_id, "task_id": task.id}, 202

@app.route("/status/<task_id>")
def status(task_id):
    task = run_inference_task.AsyncResult(task_id)

    if task.state in ["PENDING", "STARTED"]:
        # Try reading partial log
        result = task.info or {}
        log_file = result.get("log_file")
        output = ""
        if log_file and os.path.exists(log_file):
            with open(log_file, "r") as f:
                output = f.read()[-2000:]  # tail last 2000 chars
        return jsonify({"state": "RUNNING", "output": output})
    elif task.state == "SUCCESS":
        result = task.result
        video_path = result.get("video_path")
        if video_path and os.path.exists(video_path):
            return send_file(
                video_path,
                mimetype="video/mp4",
                as_attachment=False,
                download_name=os.path.basename(video_path)
            )
        else:
            return jsonify({"state": "SUCCESS", "message": "Video not found"}), 404
    elif task.state == "FAILURE":
        return jsonify({"state": "FAILURE", "error": str(task.info)})
    else:
        return jsonify({"state": task.state})

if __name__ == "__main__":
    os.makedirs(f'{work_dir}/inputs', exist_ok=True)
    os.makedirs(f'{work_dir}/outputs', exist_ok=True)
    os.makedirs(f'{work_dir}/logs', exist_ok=True)

    # Only start Flask here; Redis and Celery are started via entrypoint.sh
    app.run(host="0.0.0.0", port=8000)
