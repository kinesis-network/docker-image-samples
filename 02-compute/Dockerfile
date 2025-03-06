# FROM nvidia/cuda:12.8.0-cudnn-devel-ubuntu24.04
FROM ubuntu:24.04
WORKDIR /app

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
      g++ make \
  && rm -rf /var/lib/apt/lists/*

COPY . .
RUN ./compile.sh \
  && rm Makefile ./compile.sh

CMD ["./bin/t"]
