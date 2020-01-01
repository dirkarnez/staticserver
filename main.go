package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	textTemplate "text/template"
)

type Config struct {
	Path      string   `json:"path"`
	Source    string   `json:"source"`
	Template  string   `json:"template"`
	Arguments []string `json:"arguments"`
}

var (
	port           string
	root           string
	mode           string
	configFilePath string
	configMap      map[string]Config
	router         *gin.Engine
)

func init() {
	configMap = make(map[string]Config)

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

func main() {
	flag.StringVar(&root, "root", "", "Absolute path for root directory")
	flag.StringVar(&port, "port", "", "Port, default is 80")
	flag.StringVar(&mode, "mode", "", "Mode")
	flag.StringVar(&configFilePath, "config", "", "Config file path")
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

	if len(configFilePath) > 0 {
		raw, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			fmt.Println("Config file not found. No configuration is loaded.")
		}

		var configArr []Config
		if err := json.Unmarshal(raw, &configArr); err != nil {
			fmt.Println("Cannot parse config file. No configuration is loaded.")
		}

		for _, config := range configArr {
			configPath := config.Path
			configMap[configPath] = config
			router.POST(configPath, handlerCreator(configPath))
		}
	}

	switch mode {
	case "spa":
		router.GET("/*page", func(c *gin.Context) {
			urlPath := c.Request.URL.Path
			fullPath := filepath.FromSlash(path.Join(root, urlPath))
			rel, _ := filepath.Rel(root, fullPath)

			if len(rel) > 0 && rel != "." {
				fileInfo, err := os.Stat(fullPath)
				if !os.IsNotExist(err) && fileInfo.Mode().IsRegular() {
					c.File(fullPath)
					return
				}
			}
			c.File(path.Join(root, "index.html"))
		})
	case "fs":
	default:
		http.Handle("/", http.FileServer(http.Dir(root)))
	}

	log.Println(fmt.Sprintf("Listening on %s, serving %s, in %s mode", port, root, mode))
	router.Run(fmt.Sprintf(":%s", port))
}

func handlerCreator(key string) func(c *gin.Context) {
	return func(c *gin.Context) {
		value, ok := configMap[key]
		if !ok {
			c.AbortWithStatus(404)
		} else {
			source := value.Source
			template := value.Template
			arguments := value.Arguments

			if len(template) > 1 {
				doc, err := jsonquery.LoadURL(source)
				if err != nil {
					c.AbortWithStatus(404)
				}

				tmpl, err := textTemplate.New("template").Parse(template)
				if err != nil {
					c.AbortWithStatus(404)
					return
				}

				buf := new(bytes.Buffer)

				queries := make([]string, len(arguments))
				for i, argument := range arguments {
					nodeNameNode := jsonquery.FindOne(doc, argument)
					if nodeNameNode != nil {
						queries[i] = nodeNameNode.InnerText()
					} else {
						queries[i] = ""
					}
				}

				err = tmpl.Execute(buf, queries)
				if err != nil {
					c.AbortWithStatus(404)
					return
				}

				c.JSON(200, buf.String())
			} else {
				c.JSON(200, source)
			}
		}
	}
}