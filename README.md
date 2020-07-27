staticserver
======================
```
staticserver 
[--port={port | 80 default}] \
[--root={root (default to working directory)}] \
[--mode={spa | upload | fs default}]
```

### Warning(s)
- The mime-type in Windows registry may be modified by other software causing wrong `Content-Type` in response, see issue(s)
  - https://github.com/golang/go/issues/32350
    - How to fix: `regedit` -> go to `Computer\HKEY_CLASSES_ROOT\.js` -> set `Content Type` to `application/javascript`