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
- [ ] Single Page Application embedded
- [ ] File Server
  - [ ] Streaming (Music, Videos)
- [x] Upload Server
- [ ] Directory-synchronization server
- [ ] Chatroom
  - [ ] Clipboard Server
- [ ] Fix MIMEType
  - `.apk` file downloaded as `.zip` 
- [ ] Auto-open default browser
  - https://github.com/chromedp/chromedp/blob/master/allocate.go#L331

### Known Issue(s)
- The mime-type in Windows registry may be modified by other software causing wrong `Content-Type` in response, see issue(s)
  - https://github.com/golang/go/issues/32350
    - How to fix: `regedit` -> go to `Computer\HKEY_CLASSES_ROOT\.js` -> set `Content Type` to `application/javascript`
    - It should be platform-independent anyway
    
- Some server-side files does not have their correct `Content-Type` in response (`text/html; charset=utf-8`), found
  - [ ] `.woff`
  - [ ] `.ttf`

### Reference
- https://github.com/svenstaro/miniserve
- [HFS ~ HTTP File Server](https://www.rejetto.com/hfs/)
- https://github.com/caddyserver/caddy
