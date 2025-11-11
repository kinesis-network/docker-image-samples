FROM nvidia/cuda:12.9.0-cudnn-devel-ubuntu24.04 AS build
WORKDIR /build
RUN apt update && apt install git cmake -y
RUN git clone https://github.com/NVIDIA/cuda-samples -b v12.9

COPY src/matrixMul.cu cuda-samples/Samples/0_Introduction/matrixMul/
RUN mkdir cuda-samples/build \
  && cd cuda-samples/build \
  && cmake .. \
  && make -j$(nproc) -C Samples/0_Introduction

FROM ubuntu:24.04
WORKDIR /app
COPY --from=build /build/cuda-samples/build/Samples /app/

# Default environment variables (optional)
ENV DAEMON_MODE=1
ENV MATRIX_SIZE=

COPY entrypoint.sh /app/
CMD ["sh", "./entrypoint.sh"]
