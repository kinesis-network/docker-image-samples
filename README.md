# Docker Image Samples

## Example commands

```sh
# To rebuild an image in each directory
docker rmi $(basename $PWD); docker build -t $(basename $PWD) .

# To create an one-time container
docker run --rm -p 80:8888 -e PORT=8888 -v .:/work $(basename $PWD)
```
