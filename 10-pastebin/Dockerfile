# Use a slim version of Python to keep the image small
FROM python:3.11-slim

# Set the working directory
WORKDIR /app
COPY server.py requirements.txt .
RUN pip install --no-cache-dir --upgrade -r requirements.txt

# Create the upload directory and set permissions
RUN mkdir -p received_files && chmod 777 received_files

# Environment variables (can be overridden at runtime)
ENV UPLOAD_PORT=5000
ENV UPLOAD_DIR=received_files

# Run the app using Gunicorn
# -b 0.0.0.0 binds to all interfaces
# --threads 4 allows handling multiple uploads at once
CMD ["sh", "-c", "gunicorn -w 1 --threads 4 -b 0.0.0.0:${UPLOAD_PORT} server:app"]
