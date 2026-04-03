#!/bin/bash
set -e

# GPU Detection Logic
detect_gpus() {
    # Count GPUs using the "noheader" CSV method to get a clean number
    # We take the first line ('head -n 1') because as we discovered, it repeats per GPU.
    local count
    count=$(nvidia-smi --query-gpu=count --format=csv,noheader,nounits 2>/dev/null | head -n 1 || echo "1")
    if [[ ! "$count" =~ ^[0-9]+$ ]]; then
        echo "[!] Error: No GPUs detected."
        exit 1
    fi
    echo "$count"
}

NUM_GPUS=$(detect_gpus)
echo "[*] Detected GPU count: $NUM_GPUS"

# 1. Enable WireGuard if config exists
if [ -f "/etc/wireguard/wg0.conf" ]; then
    echo "[*] Found wg0.conf. Starting WireGuard..."
    wg-quick up /etc/wireguard/wg0.conf || echo "[*] Failed to start WireGuard. Check permissions."
else
    echo "[*] No wg0.conf found. Skipping WireGuard setup."
fi

# 2. Network Interface Selection for NCCL
if [ -z "$NCCL_SOCKET_IFNAME" ]; then
    echo "[*] NCCL_SOCKET_IFNAME not set. Detecting interface..."
    # Check if wg0 exists in the system
    if ip addr show wg0 > /dev/null 2>&1; then
        export NCCL_SOCKET_IFNAME="wg0"
        echo "[*] Selected wg0 for NCCL."
    else
        # Fallback: Get the interface associated with the default route
        PRIMARY_IF=$(ip route | grep default | awk '{print $5}' | head -n1)
        export NCCL_SOCKET_IFNAME=$PRIMARY_IF
        echo "[*] wg0 not found. Selected primary interface: $NCCL_SOCKET_IFNAME"
    fi
else
    echo "[*] Using provided NCCL_SOCKET_IFNAME: $NCCL_SOCKET_IFNAME"
fi

# 3. Dynamic Argument Construction
# Extract the IP address assigned to the chosen interface
LOCAL_IP=$(ip -4 addr show "$NCCL_SOCKET_IFNAME" | grep -oP '(?<=inet\s)\d+(\.\d+){3}')
if [ -z "$LOCAL_IP" ]; then
    echo "[!] Error: Could not find an IP address for $NCCL_SOCKET_IFNAME"
    exit 1
fi
echo "[*] Resolved local IP on $NCCL_SOCKET_IFNAME: $LOCAL_IP"

# Extract the Host/IP from RDZV_ENDPOINT (handle IP:PORT or just IP)
# We use parameter expansion to strip the port if it exists
RDZV_HOST="${RDZV_ENDPOINT%%:*}"
echo "[*] RDZV_HOST: $RDZV_HOST"

# Initialize an array to hold our dynamic torchrun arguments
# CURRENT_TASKID starts at 1, node_rank should be 0, 1, 2...
EXTRA_ARGS=(
    "--nnodes=$TOTAL_TASKS"
    "--nproc_per_node=$NUM_GPUS"
    "--node_rank=$((CURRENT_TASKID - 1))"
    "--rdzv_id=torchrun_job"
    "--rdzv_backend=c10d"
    "--rdzv_endpoint=$RDZV_ENDPOINT"
    "--local_addr=$LOCAL_IP"
)

# Apply is_host=1 if this node's IP matches the rendezvous endpoint IP
# This works for both specific IPs and 'localhost' if testing locally
if [[ "$LOCAL_IP" == "$RDZV_HOST" || "$RDZV_HOST" == "localhost" || "$RDZV_HOST" == "127.0.0.1" ]]; then
    echo "[*] This node is recognized as the master host. Adding is_host=1."
    EXTRA_ARGS+=("--rdzv_conf=is_host=1")
fi

# 4. Execute torchrun
# We slip in our calculated EXTRA_ARGS before passing the rest of the user's arguments ($@)
echo "[*] Launching torchrun with: ${EXTRA_ARGS[@]}"
exec torchrun "${EXTRA_ARGS[@]}" "$@"
