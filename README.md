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

### Warning(s)
- The mime-type in Windows registry may be modified by other software causing wrong `Content-Type` in response, see issue(s)
  - https://github.com/golang/go/issues/32350
    - How to fix: `regedit` -> go to `Computer\HKEY_CLASSES_ROOT\.js` -> set `Content Type` to `application/javascript`
