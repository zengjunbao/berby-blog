package main

import (
	"fmt"
	"github.com/cjyzwg/forestblog/config"
	// "github.com/cjyzwg/forestblog/helper"
	"github.com/cjyzwg/forestblog/routes"
	"net/http"
	"strconv"
)

func main() {
	routes.InitRoute()
	if err := http.ListenAndServe( ":" + strconv.Itoa(config.Cfg.Port) , nil); err != nil{
		fmt.Println("ServeErr:",err)
	}
}
