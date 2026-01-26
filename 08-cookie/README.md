## Repro

Run this.

```
export ENDPOINT=http://localhost:5000
rm c.txt; for i in {1..5}; do curl -s -b c.txt -c c.txt ${ENDPOINT};sleep 1; done
```

If everything works fine, you'll see something like this.

```
NEW SESSION: Created ID [02c1ca958f5f4a169a6ff38ff2c9e355]
OK: Persistent session. I recognize ID [02c1ca958f5f4a169a6ff38ff2c9e355].
OK: Persistent session. I recognize ID [02c1ca958f5f4a169a6ff38ff2c9e355].
OK: Persistent session. I recognize ID [02c1ca958f5f4a169a6ff38ff2c9e355].
OK: Persistent session. I recognize ID [02c1ca958f5f4a169a6ff38ff2c9e355].
```

If the request goes to a different node, you'see something like this.

```
NEW SESSION: Created ID [0b578c33a5ed4dd3ac12118dc90ecc8e]
DISCREPANCY DETECTED!
Client sent ID: [0b578c33a5ed4dd3ac12118dc90ecc8e]
Server Status: I have no record of issuing this ID.
DISCREPANCY DETECTED!
Client sent ID: [0b578c33a5ed4dd3ac12118dc90ecc8e]
Server Status: I have no record of issuing this ID.
OK: Persistent session. I recognize ID [0b578c33a5ed4dd3ac12118dc90ecc8e].
DISCREPANCY DETECTED!
Client sent ID: [0b578c33a5ed4dd3ac12118dc90ecc8e]
Server Status: I have no record of issuing this ID.
```
