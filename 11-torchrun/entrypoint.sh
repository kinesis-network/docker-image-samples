#!/bin/bash
set -e

# 1. Enable WireGuard if config exists
if [ -f "/etc/wireguard/wg0.conf" ]; then
    echo "[*] Found wg0.conf. Starting WireGuard..."
    # Ensure the container has net-admin capabilities
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

# 3. Execute torchrun
# We use 'exec' so that signals (like CTRL+C) are passed correctly to the python process
echo "[*] Launching torchrun..."
exec torchrun "$@"
