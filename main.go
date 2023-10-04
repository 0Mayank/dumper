package main

import (
	"fmt"

	"github.com/0Mayank/dumper/configs"
	"github.com/gin-gonic/gin"
)

func main() {
	configs.GetDB()

	r := gin.Default()

	r.Run(fmt.Sprintf("0.0.0.0:%v", configs.GetConfig().Port))
}
