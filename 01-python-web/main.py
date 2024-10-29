import os, sys
from http.server import BaseHTTPRequestHandler, HTTPServer

class handler(BaseHTTPRequestHandler):
  def do_GET(self):
    self.send_response(200)
    self.send_header('Content-type','text/html')
    self.end_headers()

    message = "Hello via GET!"
    self.wfile.write(bytes(message, "utf8"))

  def do_POST(self):
    if self.path == "/stdout":
      out_file = sys.stdout
    elif self.path == "/stderr":
      out_file = sys.stderr
    else:
      self.send_response(404)
      self.end_headers()
      return

    self.send_response(200)
    self.end_headers()

    if 'Content-Length' in self.headers:
      content_length = int(self.headers['Content-Length'])
      post_body = self.rfile.read(content_length)
      print(post_body.decode('utf-8'), file=out_file)

    self.wfile.write(bytes("ok", "utf8"))

if not 'PORT' in os.environ:
  print("PORT is not specified")
  sys.exit(1)

with HTTPServer(('', int(os.environ["PORT"])), handler) as server:
  print("Listening to port", server.server_port)
  server.serve_forever()
