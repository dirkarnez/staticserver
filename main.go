package main

import (
  "log"
  "net/http"
  "flag"
  "fmt"
  "os"
)

var (
	port string
	root string
)

func main() {
  flag.StringVar(&port, "port", "", "Port, default is 9999")
  flag.StringVar(&root, "root", "", "Absolute path for root directory")
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

  http.Handle("/", http.FileServer(http.Dir(root)))
  log.Println(fmt.Sprintf("Listening on %s, serving %s", port, root))
  http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}