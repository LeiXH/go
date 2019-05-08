package main

import (
	"flag"
	"github.com/labstack/echo"
	"net/http"
	"os"
	"pkg/logger"
	server2 "pkg/server"
	"pkg/server/controller"
)

var c string

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func init()  {
	isExist, _ :=PathExists("./bin/foo.db")

	if !isExist {
		_, err :=os.Create("./bin/foo.db")
		if err != nil {
			logger.Fatal("no foo.db file set")
		}
	}
	flag.StringVar(&c, "config", "./bin/development.yml", "config file path")
	flag.Parse()
}

func main() {

	if len(c) == 0 {
		logger.Fatal("no config file set")
	}
	server2.New(c)

	e := echo.New()

	e.Static("/static", "html")
	e.File("/", "html/index.html")
	e.File("/favicon.ico", "html/images/favicon.ico")
	e.File("/list", "html/list.html")

	e.POST("/api/sign", controller.GinDoUserSignManually)

	e.GET("/ws/wait", controller.PushFaceDetectResultToFront)

	e.GET("/api/reprint", controller.GinJustPrintUserLabel)

	e.GET("/api/face", controller.GinFaceSign)

	e.GET("/api/all", controller.All)

	e.POST("/api/import", controller.Import)

	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":1323"))
}
