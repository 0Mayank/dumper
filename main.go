package main

import (
	"fmt"

	"github.com/0Mayank/dumper/configs"
	"github.com/0Mayank/dumper/ws"
	"github.com/gin-gonic/gin"
)

func main() {
	configs.GetDB()

	hub := ws.NewHub()
	go hub.Run()

	r := gin.Default()

	r.GET("/connect", ws.ConnectWs(hub))

	r.Run(fmt.Sprintf("0.0.0.0:%v", configs.GetConfig().Port))
}
