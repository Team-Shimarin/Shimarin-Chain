package redis

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Team-Shimarin/Shimarin-Chain/config"
	"github.com/Team-Shimarin/Shimarin-Chain/tx"
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
		if i == 199 {
			err = errors.New("cannot connect Redis")
		}
		break
	}

	return c, err
}

func init() {
	conf := config.GetConfig()
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		panic(err)
	}

	c.Do("SET", txPoolKey, "[]")
}

func AddSetToTxPoolKey(transaction tx.Tx) error {
	conf := config.GetConfig()
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		panic(err)
	}

	reply, err := c.Do("GET", txPoolKey)
	if reply == nil {
		return redis.ErrNil
	}

	b, err := redis.Bytes(reply, err)
	if err != nil {
		return err
	}

	txs := []tx.Tx{}
	if err := json.Unmarshal(b, &txs); err != nil {
		return err
	}

	txs = append(txs, transaction)

	txsJSON, err := json.Marshal(txs)
	if err != nil {
		return err
	}

	c.Do("SET", txPoolKey, string(txsJSON))

	return nil
}
