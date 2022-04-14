package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// type Config struct {
// 	Path      string   `json:"path"`
// 	Source    string   `json:"source"`
// 	Template  string   `json:"template"`
// 	Arguments []string `json:"arguments"`
// }

var (
	port uint64
	root string
	mode string
	// configFilePath string
	// configMap      map[string]Config
	router *gin.Engine
)

func init() {
	// configMap = make(map[string]Config)

	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()
}

func main() {
	flag.StringVar(&root, "root", "", "Absolute path for root directory")
	flag.Uint64Var(&port, "port", 80, "Port, default is 80")
	flag.StringVar(&mode, "mode", "fs", "Mode: fs, spa, upload. Default fs mode")
	// flag.StringVar(&configFilePath, "config", "", "Config file path")
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
		file, err := os.OpenFile(path.Join(root, "index.html"))
		if err == nil {
			mode = "spa"
		}
	}

	// if len(configFilePath) > 0 {
	// 	raw, err := ioutil.ReadFile(configFilePath)
	// 	if err != nil {
	// 		fmt.Println("Config file not found. No configuration is loaded.")
	// 	}

	// 	var configArr []Config
	// 	if err := json.Unmarshal(raw, &configArr); err != nil {
	// 		fmt.Println("Cannot parse config file. No configuration is loaded.")
	// 	}

	// 	for _, config := range configArr {
	// 		configPath := config.Path
	// 		configMap[configPath] = config
	// 		router.POST(configPath, handlerCreator(configPath))
	// 	}
	// }

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
	case "upload":
		const tpl = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Multiple file upload</title>
	<link rel="stylesheet" href="https://unpkg.com/purecss">
</head>
<body>
<div class="pure-u-1">
<form class="pure-form pure-form-aligned" action="/upload" method="post" enctype="multipart/form-data">
   <fieldset>
        <div class="pure-control-group">
            <label for="name">Files</label>
            <input type="file" name="files" multiple>
        </div>
        <div class="pure-control-group">
            <label for="name">Image</label>
			<input type="file" name="files" multiple accept="image/*" capture>
        </div>
        <div class="pure-control-group">
            <label for="name">Video</label>
			<input type="file" name="files" multiple accept="video/*" capture>
        </div>
        <div class="pure-control-group">
            <label for="name">Audio</label>
            <input type="file" name="files" multiple accept="audio/*" capture>
        </div>
        <div class="pure-controls">
			<input class="pure-button pure-button-primary" type="submit" value="Submit">
        </div>
    </fieldset>
</form>
</div>
</body>
</html>`
		router.SetHTMLTemplate(template.Must(template.New("index").Parse(tpl)))

		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "index", gin.H{})
		})

		router.POST("/upload", func(c *gin.Context) {
			// Multipart form
			form, _ := c.MultipartForm()
			files := form.File["files"]

			for _, file := range files {
				filename := filepath.Base(file.Filename)
				log.Println(filename)
				if err := c.SaveUploadedFile(file, filename); err != nil {
					c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
					return
				}
			}
			c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
		})
	case "fs":
		router.StaticFS("/", gin.Dir(root, true))
	default:
		log.Fatalf("%s mode is not supported\n", mode)
	}

	log.Println(fmt.Sprintf("Listening on %d, serving %s, in %s mode", port, root, mode))
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
}

// func handlerCreator(key string) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		value, ok := configMap[key]
// 		if !ok {
// 			c.AbortWithStatus(404)
// 		} else {
// 			source := value.Source
// 			template := value.Template
// 			arguments := value.Arguments

// 			doc, err := jsonquery.LoadURL(source)
// 			if err != nil {
// 				c.AbortWithStatus(404)
// 			}

// 			if len(template) > 1 {
// 				argumentsLen := len(arguments)

// 				if argumentsLen > 0 {
// 					tmpl, err := textTemplate.New("template").Parse(template)
// 					if err != nil {
// 						c.AbortWithStatus(404)
// 						return
// 					}

// 					buf := new(bytes.Buffer)
// 					queries := make([]string, argumentsLen)
// 					for i, argument := range arguments {
// 						nodeNameNode := jsonquery.FindOne(doc, argument)
// 						if nodeNameNode != nil {
// 							queries[i] = nodeNameNode.InnerText()
// 						} else {
// 							queries[i] = ""
// 						}
// 					}

// 					err = tmpl.Execute(buf, queries)
// 					if err != nil {
// 						c.AbortWithStatus(404)
// 						return
// 					}

// 					c.JSON(200, buf.String())
// 				}

// 			} else {
// 				c.JSON(200, source)
// 			}
// 		}
// 	}
// }
