package main

import (
	"io"
	"os"
	"time"

	"log"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/handler"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

const (
	systemId = "system"
)

func getRedisConn(host, port string) (redis.Conn, error) {
	var c redis.Conn
	var err error
	for i := 0; i < 200; i++ {
		c, err = redis.Dial("tcp", host+":"+port)
		if err != nil {
			log.Print(host + ":" + port)
			log.Printf(err.Error())
			log.Printf("redis connection: retry cnt %d", i)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}

	return c, err
}

func main() {
	// get config
	log.Println("Anzu Wake Up")
	conf := config.GetConfig()

	r := gin.Default()
	f, _ := os.Create("anzu-access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())

	myId := "hoge"
	myhp := int64(100)
	diff := int64(100)

	go HashCalculate(myId, myhp, diff)
	go ValidHashSubScribe()
	go subscribeValidHashEach()

	accountHandler := handler.NewAccountHandler(conf, dba.AccountAccess{})

	r.POST("/api/v1/register", accountHandler.Register)
	r.Run(":8080")
}
