import os
import torch
import torch.nn as nn
import torch.optim as optim
import torch.distributed as dist
from torch.nn.parallel import DistributedDataParallel as DDP

def setup():
    # Initialize the process group
    # NCCL is the standard backend for NVIDIA GPUs
    dist.init_process_group(backend="nccl")
    torch.cuda.set_device(int(os.environ["LOCAL_RANK"]))

def cleanup():
    dist.destroy_process_group()

def run_training():
    setup()

    local_rank = int(os.environ["LOCAL_RANK"])
    steps = int(os.environ["STEPS"])
    rank = int(os.environ["RANK"])
    device = torch.device(f"cuda:{local_rank}")

    # 1. Define a tiny model
    model = nn.Linear(10, 10).to(device)
    model = DDP(model, device_ids=[local_rank])

    # 2. Setup Loss and Optimizer
    loss_fn = nn.MSELoss()
    optimizer = optim.SGD(model.parameters(), lr=0.001)

    # 3. Simple training loop (Synthetic data)
    print(f"[Rank {rank}] Starting training...")

    for step in range(steps):
        # Create random data on the fly
        inputs = torch.randn(20, 10).to(device)
        labels = torch.randn(20, 10).to(device)

        optimizer.zero_grad()
        outputs = model(inputs)
        loss = loss_fn(outputs, labels)
        loss.backward()
        optimizer.step()

        if step % 10 == 0 and rank == 0:
            print(f"Step {step} | Loss: {loss.item():.4f}")

    print(f"[Rank {rank}] Training complete.")
    cleanup()

if __name__ == "__main__":
    run_training()