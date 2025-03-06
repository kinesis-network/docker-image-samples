#!/bin/bash
DOCKER=docker
IMAGE_NAME=02-compute

${DOCKER} rmi -f ${IMAGE_NAME}
${DOCKER} build -t ${IMAGE_NAME} .

${DOCKER} run --rm -it --gpus=1 ${IMAGE_NAME}
