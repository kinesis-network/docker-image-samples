# Use the official PyTorch image with CUDA support
FROM pytorch/pytorch:2.11.0-cuda12.8-cudnn9-devel

WORKDIR /app

# Install packages to enable Wireguard
RUN apt-get update && apt-get install -y \
    iputils-ping \
    net-tools \
    openssh-server \
    wireguard-tools iptables iproute2 \
    && rm -rf /var/lib/apt/lists/*

# Allow SSH root login for troubleshooting
RUN mkdir /run/sshd \
    && sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' \
        /etc/ssh/sshd_config

COPY train.py /app/train.py
COPY --chmod=755 entrypoint.sh /app/entrypoint.sh

ENV PYTHONUNBUFFERED=1 \
  NCCL_DEBUG=INFO \
  NCCL_IGNORE_DISABLED_P2P=1 \
  STEPS=50

ENTRYPOINT ["/app/entrypoint.sh"]
CMD ["train.py"]
