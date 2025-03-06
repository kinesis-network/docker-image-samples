FROM nvidia/cuda:12.8.0-cudnn-devel-ubuntu24.04
WORKDIR /app

COPY . .
RUN ./compile.sh \
  && rm Makefile ./compile.sh

# Default environment variables (optional)
ENV DAEMON_MODE=1
ENV MATRIX_SIZE=

CMD ["sh", "./entrypoint.sh"]
