staticserver
======================
```
staticserver 
[--port={port | 80 default}] \
[--root={root (default to working directory)}] \
[--mode={spa | fs default}] \
[--config={config JSON file path (optional)}]
```

### How config JSON file looks like
```
[
  {
    "path": "endpoint",
    "source": "{ url for json }",
    "template": "{ jsonquery on {source}, in golang array-templated string }",
    "arguments": [
      "//nodeName",
      "/internalEndpoint/ports/*[port=8080]/nodePort"
    ]
  }
]
```

### jsonquery
- [github.com/antchfx/jsonquery](github.com/antchfx/jsonquery)

### Golang array-templated string format
- ```{{index . 0}}{{index . 1}}{{index . 2}}....```
