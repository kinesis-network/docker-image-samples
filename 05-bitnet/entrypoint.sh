#!/bin/bash
set -e

if [ ! -f "models/BitNet-b1.58-2B-4T/ggml-model-i2_s.gguf" ]; then
  hf download microsoft/BitNet-b1.58-2B-4T-gguf \
    --local-dir models/BitNet-b1.58-2B-4T
fi
if [ ! -f "build/bin/llama-server" ]; then
  setup_cmd="python setup_env.py -md models/BitNet-b1.58-2B-4T -q i2_s"
  echo ${setup_cmd}
  ${setup_cmd}
fi

open-webui serve &

build/bin/llama-server \
  --model models/BitNet-b1.58-2B-4T/ggml-model-i2_s.gguf \
  --port ${BACKEND_PORT} \
  --ctx-size 2048
