# Use an official CUDA-enabled base image for PyTorch + flash-attn compatibility
# FROM nvidia/cuda:13.0.1-cudnn-devel-ubuntu24.04
FROM nvidia/cuda:12.9.0-cudnn-devel-ubuntu24.04

# Avoid interactive prompts during package installation
ENV DEBIAN_FRONTEND=noninteractive

# Install Python and basic dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    git build-essential \
    python3 python3-pip python3-venv \
    libgl1 libglib2.0-0 \
    redis-server \
    && rm -rf /var/lib/apt/lists/*

RUN python3 -m venv /opt/venv
ENV PATH="/opt/venv/bin:${PATH}"
RUN python3 -m pip install --upgrade \
    pip setuptools wheel psutil \
    huggingface_hub \
    flask celery[redis] pyyaml

RUN git clone https://github.com/character-ai/Ovi.git -b main /app
WORKDIR /app

# Install torch, torchvision, torchaudio from official PyTorch wheel index
RUN pip install torch torchvision torchaudio \
    --index-url https://download.pytorch.org/whl/cu129

# Install flash-attn (from source, requires CUDA headers)
# If you only need prebuilt version, this will automatically pick it up
RUN pip install flash-attn --no-build-isolation --use-pep517

RUN pip install -r requirements.txt

COPY . .
EXPOSE 8000

ENTRYPOINT ["./entrypoint.sh"]
