package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo"
	"os"
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
			fmt.Println(err)
		}
	}
	flag.StringVar(&c, "config", "./bin", "config file path")
	flag.Parse()
}

func main() {
	e := echo.New()

	fmt.Println("hello word!")
}
