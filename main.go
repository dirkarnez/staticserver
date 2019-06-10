package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	port string
	root string
	mode string
)

func main() {
	flag.StringVar(&port, "port", "", "Port, default is 9999")
	flag.StringVar(&root, "root", "", "Absolute path for root directory")
	flag.StringVar(&mode, "mode", "", "Mode")
	flag.Parse()

	if len(root) < 1 {
		dir, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		root = dir
	}

	if len(port) < 1 {
		port = "80"
	}

	if len(mode) < 1 {
		mode = "fs"
	}

	switch mode {
	case "spa":
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			p := r.URL.Path[1:]
			if strings.Contains(p, ".") {
				http.ServeFile(w, r, p)
			} else {
				http.ServeFile(w, r, path.Join(root, "index.html"))
			}
		})
	case "fs":
	default:
		http.Handle("/", http.FileServer(http.Dir(root)))
	}

	log.Println(fmt.Sprintf("Listening on %s, serving %s", port, root))
	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}