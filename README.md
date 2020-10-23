staticserver
======================
Minimalistic HTTP web-server for file-sharing, single-page application hosting

```
staticserver 
[--port={port | default to 80 }] \
[--root={root (default to working directory)}] \
[--mode={spa | upload | fs default}]
```

### TODOs
- [x] Single Page Application
- [x] File Server
- [x] Upload Server
- [ ] Chatroom
- [ ] Directory-synchronization server

### Known Issue(s)
- The mime-type in Windows registry may be modified by other software causing wrong `Content-Type` in response, see issue(s)
  - https://github.com/golang/go/issues/32350
    - How to fix: `regedit` -> go to `Computer\HKEY_CLASSES_ROOT\.js` -> set `Content Type` to `application/javascript`
    
### Reference
- https://github.com/svenstaro/miniserve
- [HFS ~ HTTP File Server](https://www.rejetto.com/hfs/)
