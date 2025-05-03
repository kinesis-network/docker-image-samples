## 01-python-web

This image is to stress system resources via HTTP

```
# CPU
curl localhost:8888/cpu?p=10

# RAM
curl localhost:8888/ram?p=100

# NetworkIn
curl localhost:8888/nwin?p=10

# NetworkOut
curl localhost:8888/nwout?p=100

# Storage
curl localhost:8888/disk?p=10
```
