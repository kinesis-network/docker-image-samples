from flask import Flask, request, \
  send_from_directory, render_template_string, redirect, url_for
import os
import time

app = Flask(__name__)

# Configuration from Environment Variables
PORT = int(os.environ.get("UPLOAD_PORT", 5000))
UPLOAD_DIR = os.environ.get("UPLOAD_DIR", "received_files")

if not os.path.exists(UPLOAD_DIR):
  os.makedirs(UPLOAD_DIR)

# HTML Template with Delete Buttons
HTML_TEMPLATE = """
<!DOCTYPE html>
<html>
<head>
  <title>File Manager</title>
  <style>
    body { font-family: sans-serif; margin: 40px; line-height: 1.6; }
    .file-list { list-style: none; padding: 0; }
    .file-item { 
      background: #f4f4f4; 
      margin: 10px 0; 
      padding: 10px; 
      display: flex; 
      justify-content: space-between;
      border-radius: 4px;
    }
    .delete-btn { color: white; background: #ff4444; border: none; padding: 5px 10px; border-radius: 3px; cursor: pointer; }
    .empty-msg { color: #888; font-style: italic; }
  </style>
</head>
<body>
  <h2>📁 Uploaded Files</h2>
  {% if files %}
    <ul class="file-list">
      {% for file in files %}
      <li class="file-item">
        <a href="{{ url_for('download_file', filename=file) }}"><strong>{{ file }}</strong></a>
        <form action="{{ url_for('delete_file', filename=file) }}" method="POST" style="display:inline;">
          <button type="submit" class="delete-btn" onclick="return confirm('Delete this file?')">Delete</button>
        </form>
      </li>
      {% endfor %}
    </ul>
  {% else %}
    <p class="empty-msg">No files uploaded yet.</p>
  {% endif %}
  <p><small>Storage: {{ UPLOAD_DIR }}/</small></p>
</body>
</html>
"""

@app.route('/upload', methods=['POST'])
def upload_file():
  # 1. Check Content-Length header first (fastest check)
  content_length = request.content_length
  if content_length is not None and content_length == 0:
    return "Error: No data received (Content-Length is 0).\n", 400

  filename = f"data_{int(time.time())}"
  filepath = os.path.join(UPLOAD_DIR, filename)
  
  bytes_written = 0
  with open(filepath, 'wb') as f:
    chunk_size = 4096
    while True:
      chunk = request.stream.read(chunk_size)
      if not chunk:
        break
      f.write(chunk)
      bytes_written += len(chunk)

  # 2. Final check: if we somehow wrote 0 bytes despite the header
  if bytes_written == 0:
    if os.path.exists(filepath):
      os.remove(filepath) # Don't leave empty files
    return "Error: No data received in POST body.\n", 400

  return f"Successfully saved {bytes_written} bytes to {filename}\n", 200

@app.route('/browse')
def browse_files():
  # Get list of files, excluding hidden ones
  files = [f for f in os.listdir(UPLOAD_DIR) if os.path.isfile(os.path.join(UPLOAD_DIR, f))]
  files.sort(reverse=True) # Newest first
  return render_template_string(HTML_TEMPLATE, files=files)

@app.route('/download/<filename>')
def download_file(filename):
  # as_attachment=True forces the browser to download rather than display
  return send_from_directory(UPLOAD_DIR, filename, as_attachment=True)

@app.route('/delete/<filename>', methods=['POST'])
def delete_file(filename):
  filepath = os.path.join(UPLOAD_DIR, filename)
  # Basic security check to prevent directory traversal
  if os.path.exists(filepath) and ".." not in filename:
    os.remove(filepath)
  return redirect(url_for('browse_files'))

if __name__ == '__main__':
  # Using the PORT variable from the environment
  app.run(host='0.0.0.0', port=PORT)
