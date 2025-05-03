from http.server import ThreadingHTTPServer, BaseHTTPRequestHandler
from urllib.parse import urlparse, parse_qs
from pathlib import Path
from random import random
import logging, os, requests, sys, tempfile, time

logging.basicConfig(
  format='%(asctime)s %(message)s',
  datefmt='%Y-%m-%d %H:%M:%S',
)
logger = logging.getLogger()
logger.setLevel(logging.INFO)

def stress_cpu(duration_sec=5):
  logger.info(f'perform square root for {duration_sec} seconds')
  end = time.time() + duration_sec
  while time.time() < end:
    _ = int(random() * 1000000) ** 0.5

def stress_ram(size_mb=100):
  logger.info(f'generating {size_mb} MB of random data')
  _ = [int(random() * 256) for _ in range(size_mb * 1024 * 1024)]

def stress_disk(size_mb=100):
  container_root = "/work"
  if not Path(container_root).is_dir():
    container_root = None

  with tempfile.NamedTemporaryFile(
    dir=container_root, delete_on_close=True, mode="w") as f:
    logger.info(f'writing {size_mb}MB to {f.name}')
    f.write("x" * 1024 * 1024 * size_mb)
    f.close()

def stress_nw_in(url="http://example.com", count=10):
  for i in range(count):
    try:
      logger.info(f'[{i + 1}/{count}] loading {url}')
      requests.get(url, timeout=1)
    except requests.RequestException:
      pass

def stress_nw_out(url="http://example.com", size_mb=10):
  try:
    logger.info(f'posting {size_mb} MB to {url}')
    requests.post(url, data=b"x" * 1024 * 1024 * size_mb)
  except requests.RequestException:
    pass

class handler(BaseHTTPRequestHandler):
  def do_GET(self):
    parsed_url = urlparse(self.path)
    path = parsed_url.path
    query = parse_qs(parsed_url.query)

    if path == "/cpu":
      duration = int(query.get("p", [5])[0])
      stress_cpu(duration)
    elif path == "/ram":
      size = int(query.get("p", [100])[0])
      stress_ram(size)
    elif path == "/disk":
      size = int(query.get("p", [100])[0])
      stress_disk(size)
    elif path == "/nwin":
      count = int(query.get("p", [10])[0])
      stress_nw_in("https://kinesiscloud.com", count)
    elif path == "/nwout":
      size = int(query.get("p", [10])[0])
      stress_nw_out("https://kinesiscloud.com", size)

    self.send_response(200)
    self.send_header("Content-type", "text/plain")
    self.end_headers()

    message = f"Handled request: {self.path}\n"
    self.wfile.write(message.encode("utf-8"))
    logger.info('done.')

def main(server_class=ThreadingHTTPServer, handler_class=handler):
  if not 'PORT' in os.environ:
    logger.info("PORT is not specified")
    sys.exit(1)

  port = int(os.environ["PORT"])
  server_address = ('', port)
  httpd = server_class(server_address, handler_class)
  logger.info(f"Starting server on port {port}")
  httpd.serve_forever()

if __name__ == '__main__':
  main()
