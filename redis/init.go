package redis

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
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
