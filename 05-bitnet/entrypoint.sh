#!/bin/bash
set -e
open-webui serve &
./llama-server \
  --model models/ggml-model-i2_s.gguf \
  --host 0.0.0.0 \
  --port ${BACKEND_PORT} \
  --ctx-size 2048
