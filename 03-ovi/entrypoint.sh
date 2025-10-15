#!/bin/bash
set -e

# Start Redis in background (if you’re running Redis inside the container)
redis-server --daemonize yes

# Start Celery worker in background
celery -A server.celery worker --loglevel=INFO --concurrency=1 &

# Start Flask API in foreground
python3 server.py
