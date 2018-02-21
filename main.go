package main

import (
	"io"
	"os"

	"log"

	"github.com/Team-Shimarin/Shimarin-Chain/config"
	"github.com/Team-Shimarin/Shimarin-Chain/dba"
	"github.com/Team-Shimarin/Shimarin-Chain/handler"
	anzuredis "github.com/Team-Shimarin/Shimarin-Chain/redis"
	"github.com/gin-gonic/gin"
)

const (
	systemId    = "system"
	inithp      = 0
	initbalance = 0
)

func main() {
	// get config
	log.Println("Anzu Wake Up")
	conf := config.GetConfig()
	log.Print(conf)

	r := gin.Default()
	f, _ := os.Create("anzu-access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())

	go anzuredis.HashCalculate(conf.MinorAccountID, conf.Diff)
	go anzuredis.ValidHashSubScribe()
	go anzuredis.SubscribeValidHashEach()

	accountHandler := handler.NewAccountHandler(conf, dba.AccountAccess{})

	// HP更新
	// req {hp: healthpoint}
	r.POST("/api/v1/account/healthpoint/update", accountHandler.UpdateHP)

	// 送金
	r.POST("/api/v1/balance/remit", accountHandler.Remit)

	r.Run(":8080")
}
