FROM python:3.12-slim

# --- Install system deps ---
RUN apt-get update \
  && apt-get install -y \
    git wget curl build-essential cmake clang \
  && rm -rf /var/lib/apt/lists/*

# --- Install Python dependencies ---
RUN pip install --no-cache-dir --upgrade pip \
  && pip install --no-cache-dir \
    huggingface_hub \
    open-webui

# --- Clone repos ---
WORKDIR /app
COPY entrypoint.sh .
RUN chmod 755 entrypoint.sh
RUN git clone --recursive https://github.com/microsoft/BitNet.git -b main

WORKDIR /app/BitNet
RUN pip install --no-cache-dir -r requirements.txt

ENV BACKEND_PORT=10000
ENV PORT=8080
ENV RAG_EMBEDDING_ENGINE=openai
ENV OPENAI_API_BASE_URL=http://localhost:10000/v1
ENV WEBUI_AUTH=False

ENTRYPOINT ["/app/entrypoint.sh"]
