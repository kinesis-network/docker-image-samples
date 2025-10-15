## 01-python-web

This image is a kinesified version of [Ovi](https://github.com/character-ai/Ovi).

```
docker rmi kinesisorg/03-ovi; \
  docker build -t kinesisorg/03-ovi .

docker run --rm -it --gpus=all \
  -v ${PWD}:/work \
  docker.io/kinesisorg/03-ovi \
  bash
```
