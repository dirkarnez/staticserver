package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/dirkarnez/staticserver/views"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/v12"
	"gopkg.in/yaml.v3"
)

var (
	port uint64
	root string
	mode string
)

func main() {
	flag.StringVar(&root, "root", "", "Absolute path for root directory")
	flag.Uint64Var(&port, "port", 80, "Port, default is 80")
	flag.StringVar(&mode, "mode", "", "Mode: fs, spa, upload, clipboard. Default fs mode")
	flag.Parse()

	if len(root) < 1 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		root = dir
	}

	if port > 1<<16-1 {
		log.Fatal("Port number too large")
	}

	if len(mode) < 1 {
		_, err := os.Open(filepath.Join(root, "index.html"))
		if err != nil {
			mode = "fs"
		} else {
			mode = "spa"
		}
	}

	iris.RegisterOnInterrupt(func() {
		// TODO
	})

	app := iris.New()
	app.Use(iris.Compression)
	app.Use(iris.NoCache)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})
	app.Use(crs)

	app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context) {
		ctx.Writef("404 not found here")
	})

	switch mode {
	case "spa":
		app.Use(func(ctx iris.Context) {
			mimeOverrides := loadMIMEOverrides()
			extension := filepath.Ext(ctx.RequestPath(true))
			mimeForExtension, ok := mimeOverrides[extension]
			if ok {
				log.Println(extension, mimeForExtension, "user-defined")
			} else {
				mimeForExtension := mime.TypeByExtension(extension)
				log.Println(extension, mimeForExtension, "iris built-in")
			}

			ctx.ContentType(mimeForExtension)
			ctx.Next()
		})
		app.HandleDir("/", iris.Dir(root), iris.DirOptions{
			IndexName: "index.html",
			SPA:       true,
		})
	case "clipboard":
		app.Get("/", func(ctx iris.Context) {
			ctx.HTML(views.Clipboard)
		})

		app.Post("/clipboard", func(ctx iris.Context) {
			input := ctx.FormValue("input")
			fileName := fmt.Sprintf("%s.txt", strings.ReplaceAll(time.Now().Format(time.RFC3339), ":", "-"))
			file, err := os.Create(fileName)
			if err != nil {
				ctx.StopWithError(iris.StatusInternalServerError, err)
				return
			}
			defer file.Close()

			n, err := file.WriteString(input)
			if err != nil {
				ctx.StopWithError(iris.StatusInternalServerError, err)
				return
			}

			if n != len(input) {
				ctx.StopWithError(iris.StatusInternalServerError, fmt.Errorf("incomplete data is written"))
				return
			}
			log.Printf("%s is created with %d!", fileName, n)
			ctx.StatusCode(http.StatusOK)
		})
	default:
		log.Fatalf("%s mode is not supported\n", mode)
	}
	
	log.Printf("Listening on %d, serving %s, in %s mode\n", port, root, mode)

	
	err := app.Run(
		iris.TLS(fmt.Sprintf(":%d", port), "mycert.crt", "mykey.key"),
		// skip err server closed when CTRL/CMD+C pressed:
		iris.WithoutServerError(iris.ErrServerClosed),
		// enables faster json serialization and more:
		iris.WithOptimizations,
	)
	if err != nil {
		log.Fatal(err)
	}
}

func loadMIMEOverrides() map[string]string {
	m := make(map[string]string)
	buf := new(bytes.Buffer)
	file, err := os.Open(filepath.Join(root, "mime.yaml"))
	if err != nil {
		return m
	}
	defer file.Close()
	_, err = buf.ReadFrom(file)
	if err != nil {
		return m
	}
	err = yaml.Unmarshal(buf.Bytes(), &m)
	if err != nil {
		return m
	}
	return m
}

// openURL opens a browser window to the specified location.
// This code originally appeared at:
//   http://stackoverflow.com/questions/10377243/how-can-i-launch-a-process-that-is-not-a-file-in-go
func openURL(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", "http://localhost:4001/").Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Cannot open URL %s on this platform", url)
	}
	return err
}
