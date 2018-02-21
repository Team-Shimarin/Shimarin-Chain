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
	"github.com/InvincibleMan/anzu-chain/model"
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

	myhp := int64(100) // TODO: dbaから取る


	// accout Create
	accountaccess := dba.AccountAccess{}
	account, err := model.NewAccount()
	if err != nil{
		log.Println("failed to create newaccount", err)
	}
	err = accountaccess.Register(account)
	if err != nil {
		log.Println("failed to create new account", err)
	}else{
		log.Println("success to create new account", account.ID)
	}
	err = accountaccess.InsertHealth(account, myhp)
	if err != nil{
		log.Println(err)
	}
	r := gin.Default()
	f, _ := os.Create("anzu-access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())


	go anzuredis.HashCalculate(conf.MinorAccountID, myhp, conf.Diff)
	go anzuredis.ValidHashSubScribe()
	go anzuredis.SubscribeValidHashEach()

	accountHandler := handler.NewAccountHandler(conf, dba.AccountAccess{})

	// Redister
	r.POST("/api/v1/register", accountHandler.Register)

	// HP更新
	// req {id: accountid, hp: healthpoint}
	r.POST("/api/v1/account/healthpoint/update", accountHandler.UpdateHP)

	// HPの取得
	// req {id: accountid}
	// res {hp: healthpoint}
	r.POST("/api/v1/account/healthpoint/get", accountHandler.GetHP)

	// Balanceの取得
	// req {id: accountid}
	// res {balance: balance}
	r.POST("/api/v1/account/balance/get", accountHandler.GetBalance)

	// Blickを全て取得
	r.GET("/api/v1/block/getall", accountHandler.GetBlock)
	r.Run(":8081")
}
