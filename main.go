package main

import (
	"io"
	"os"
	"time"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/handler"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	systemId = "system"
)

func main() {
	// get config
	log.Println("Anzu Wake Up")
	conf := config.GetConfig()
	// redis connection
	var c redis.Conn
	var err error
	for i := 0; i < 200; i++ {
		c, err = redis.Dial("tcp", conf.RedisHost + ":" + conf.RedisPort)
		if err != nil {
			log.Printf("%s:%s", conf.RedisHost, conf.RedisPort)
			log.Printf(err.Error())
			log.Printf("redis connection: retry cnt %d", i)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	defer c.Close()

	r := gin.Default()
	f, _ := os.Create("anzu-access.log")
	gin.DefaultWriter = io.MultiWriter(f)
	r.Use(gin.Logger())

	myId := "hoge"
	myhp := int64(100)
	diff := int64(100)

	go HashCalculate(c, myId, myhp, diff)
	go ValidHashSubScribe(c)
	// NOTE: SUBSCRIBE VALID_HASH_EACH
	go subscribeValidHashEach(c)

	accountHandler := handler.NewAccountHandler(conf, dba.AccountAccess{})

	r.POST("/api/v1/register", accountHandler.Register)
	r.Run(":8080")
}
