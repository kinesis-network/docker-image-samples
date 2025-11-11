FROM nvidia/cuda:12.9.0-cudnn-devel-ubuntu24.04 AS builder
WORKDIR /build
RUN apt update && apt install git cmake -y
COPY main.cu .
RUN nvcc -O0 -o vram_stress main.cu

FROM ubuntu:24.04
WORKDIR /app
COPY --from=builder /build/vram_stress /usr/local/bin/vram_stress
CMD ["vram_stress", "0.5", "0.9", "10"]
