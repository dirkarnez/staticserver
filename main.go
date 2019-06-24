package main

import (
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
	"strings"
)

type Config struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

var (
	port           string
	root           string
	mode           string
	configFilePath string
	configMap      map[string]string
	router         *gin.Engine
)

func init() {
	configMap = make(map[string]string)
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
			configMap[config.Path] = config.Value
			if config.Path == "config" {
				log.Fatal("/config is reserved")
			} else {
				router.POST(config.Path, handlerCreator(config.Path))
			}
		}
	}

	router.POST("/config", func(c *gin.Context) {
		c.Request.ParseForm()
		form := c.Request.PostForm
		for key := range form {
			if _, ok := configMap[key]; !ok {
				router.POST(key, handlerCreator(key))
			}
			configMap[key] = form.Get(key)
		}

		c.JSON(200, configMap)
	})

	switch mode {
	case "spa":
		router.GET("/*page", func(c *gin.Context) {
			urlPath := c.Request.URL.Path
			if strings.Contains(urlPath, ".") {
				c.File(path.Join(root, urlPath))
			} else {
				c.File(path.Join(root, "index.html"))
			}
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
			tokens := strings.Split(value, "->")
			if len(tokens) > 1 {
				fmt.Println(tokens[0], tokens[1])

				doc, err := jsonquery.LoadURL(tokens[0])
				if err != nil {
					c.AbortWithStatus(404)
				}
				var nodeName string
				nodeNameNode := jsonquery.FindOne(doc, tokens[1])
				if nodeNameNode != nil {
					nodeName = nodeNameNode.InnerText()
				}

				c.JSON(200, nodeName)
			} else {
				c.JSON(200, value)
			}
		}
	}
}