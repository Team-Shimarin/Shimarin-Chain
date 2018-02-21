package main

import (
	"io"
	"os"

	"log"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/handler"
	anzuredis "github.com/InvincibleMan/anzu-chain/redis"
	"github.com/gin-gonic/gin"
)

const (
	systemId = "system"
	inithp = 0
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

	myhp := int64(100) // TODO: dbaから取る

	go anzuredis.HashCalculate(conf.MinorAccountID, myhp, conf.Diff)
	go anzuredis.ValidHashSubScribe()
	go anzuredis.SubscribeValidHashEach()

	accountHandler := handler.NewAccountHandler(conf, dba.AccountAccess{})

	r.POST("/api/v1/register", accountHandler.Register)
	// HPの受け入れエンドポイント
	r.POST("/api/v1/account/healthpoint/update", accountHandler.UpdateHP)
	r.Run(":8080")
}
