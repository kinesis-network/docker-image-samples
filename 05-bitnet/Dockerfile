FROM python:3.12 AS builder

# --- Install system deps ---
RUN apt-get update \
  && apt-get install -y --no-install-recommends \
    git build-essential cmake clang \
  && rm -rf /var/lib/apt/lists/*

# --- Clone repos ---
WORKDIR /app
RUN git clone --recursive https://github.com/microsoft/BitNet.git -b main

WORKDIR /app/BitNet

# --- Install Python dependencies ---
RUN pip install --no-cache-dir \
  huggingface_hub \
  -r requirements.txt

# --- Download model and compile ---
RUN hf download microsoft/BitNet-b1.58-2B-4T-gguf \
  --local-dir models/BitNet-b1.58-2B-4T
RUN python setup_env.py -md models/BitNet-b1.58-2B-4T -q i2_s

FROM python:3.12

WORKDIR /app/BitNet

RUN pip install --no-cache-dir open-webui

COPY --from=builder \
  /app/BitNet/build/bin/llama-server \
  /app/BitNet/build/3rdparty/llama.cpp/src/libllama.so \
  /app/BitNet/build/3rdparty/llama.cpp/ggml/src/libggml.so \
  /app/BitNet/
COPY --from=builder \
  /app/BitNet/models/* \
  /app/BitNet/models/
COPY --chmod=755 entrypoint.sh /app/

ENV BACKEND_PORT=10000
ENV PORT=8080
ENV RAG_EMBEDDING_ENGINE=openai
ENV OPENAI_API_BASE_URL=http://localhost:10000/v1
ENV WEBUI_AUTH=False

ENTRYPOINT ["/app/entrypoint.sh"]
