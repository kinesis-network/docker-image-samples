FROM nvidia/cuda:12.8.0-cudnn-devel-ubuntu24.04
WORKDIR /app

COPY . .
RUN ./compile.sh \
  && rm Makefile ./compile.sh

CMD ["./bin/t", "-d"]
