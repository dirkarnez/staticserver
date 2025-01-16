staticserver
======================
Minimalistic HTTP web-server for file-sharing, single-page application hosting and many more. Intended for development use / as an utility, not for production use

```
staticserver 
[--port={port | default to 80 }] \
[--root={root (default to working directory)}] \
[--mode={spa | upload | fs default}]
```

### TODOs
- [x] **HTTPS**
  - `%USERPROFILE%\Downloads\PortableGit\usr\bin\openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -subj '/CN=localhost' -keyout server.key -out server.crt`
- [ ] **pdf mode** **(requires https)**
  - a LAN version of [dirkarnez/webcam-playground](https://github.com/dirkarnez/webcam-playground) which the `.pdf` file is uploaded to the host computer
- [ ] **encoder mode**
- [ ] **proxy mode**
- [ ] **media streaming**
- [ ] prompt to auto assign port if default 80 is used
- [ ] directory upload
  - ```html
    <input type="file" id="ctrl" webkitdirectory directory multiple/>
    ``` 
- [ ] force no-caching in client-side
- [ ] File Server
  - [ ] Streaming (Music, Videos, text files)
- [x] Upload Server
- [ ] Directory-synchronization server
  - https://github.com/elgs/filesync
  - [Building a Go-based file synchronization tool | Reintech media](https://reintech.io/blog/building-go-based-file-synchronization-tool)
  - [no-src/gofs: A cross-platform real-time file synchronization tool out of the box based on Golang](https://github.com/no-src/gofs)
- [ ] Chatroom
  - [ ] Clipboard Server
- [ ] Fix MIMEType, add customization support
  - `application/vnd.android.package-archive` for `.apk`
  - `application/wasm` for `.wasm`
  - `text/css` for `.css`
- [ ] Auto-open default browser
  - https://github.com/chromedp/chromedp/blob/master/allocate.go#L331

### TODOs (`2.0.0`)
- spa mode bugs

### Known Issue(s)
- The mime-type in Windows registry may be modified by other software causing wrong `Content-Type` in response, see issue(s)
  - https://github.com/golang/go/issues/32350
    - How to fix: `regedit` -> go to `Computer\HKEY_CLASSES_ROOT\.js` -> set `Content Type` to `application/javascript`
    - It should be platform-independent anyway
    
- Some server-side files does not have their correct `Content-Type` in response (`text/html; charset=utf-8`), found
  - [ ] `.woff`
  - [ ] `.ttf`

### MIME
- [Common MIME types - HTTP | MDN](https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types)
- Sample `mime.yaml` file (in case user-defined mime is needed, create one under the working directory)
- https://github.com/kataras/iris/blob/main/core/router/mime.go
- [MIME | yhirose/cpp-httplib: A C++ header-only HTTP/HTTPS server and client library](https://github.com/yhirose/cpp-httplib/?tab=readme-ov-file#static-file-server)
```yaml
.otf: font/otf
.woff: font/woff
.woff2: font/woff2
.ttf: font/ttf
.ico: image/vnd.microsoft.icon
.apk: application/vnd.android.package-archive
.wasm: application/wasm
.js: application/javascript
```

### Reference
- [**GoogleChromeLabs/simplehttp2server: A simple HTTP/2 server for development**](https://github.com/GoogleChromeLabs/simplehttp2server)
- [**kangc666/MyPHPServer: A easy PHP Server written in Golang**](https://github.com/kangc666/MyPHPServer)
- [http-server - npm](https://www.npmjs.com/package/http-server)
- [svrxjs/svrx: Server-X: A pluggable frontend server built for efficient front-end development](https://github.com/svrxjs/svrx)
- [svenstaro/miniserve: 🌟 For when you really just want to serve some files over HTTP right now!](https://github.com/svenstaro/miniserve)
- [HFS ~ HTTP File Server](https://www.rejetto.com/hfs/)
- [caddyserver/caddy: Fast, multi-platform web server with automatic HTTPS](https://github.com/caddyserver/caddy)
  - [iris/_examples/caddy at master · kataras/iris](https://github.com/kataras/iris/tree/master/_examples/caddy)
