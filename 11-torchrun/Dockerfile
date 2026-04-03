# Use the official PyTorch image with CUDA support
FROM pytorch/pytorch:2.11.0-cuda12.8-cudnn9-devel

# Set the working directory
WORKDIR /app

# Install any basic networking tools for troubleshooting
RUN apt-get update && apt-get install -y \
    iputils-ping \
    net-tools \
    wireguard-tools iptables iproute2 \
    && rm -rf /var/lib/apt/lists/*

COPY train.py /app/train.py

ENV PYTHONUNBUFFERED=1 \
  NCCL_SOCKET_IFNAME=wg0 \
  NCCL_DEBUG=INFO \
  NCCL_IGNORE_DISABLED_P2P=1 \
  STEPS=50

# The entrypoint will be torchrun, but we leave it flexible
# so we can pass arguments via 'docker run'
ENTRYPOINT ["torchrun"]
