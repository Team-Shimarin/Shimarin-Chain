package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/hash"
	"github.com/InvincibleMan/anzu-chain/tx"
	"github.com/garyburd/redigo/redis"
)

func ValidHashSubScribe() {
	log.Println("ValidHashSubScribe: Goutine start")
	conf := config.GetConfig()
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	defer c.Close()
	log.Println("ValidHashSubScribe: connected to redis")
	psc := redis.PubSubConn{Conn: c}
	err = psc.Subscribe(validHashChan)
	if err != nil {
		log.Print("error: subscribe ", validHashChan, ":because ", err)
	}
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			log.Println("validhHashChan get Message ", v.Data)
			str, err := redis_get(txPoolKey, c)
			if err != nil {
				log.Printf("error: in validHashSubScribe: %v", err)
			}
			if str != "" {
				valid_data := struct {
					Transaction tx.Tx `json:"tx"`
					Timestamp   int64 `json:"ts"`
					HP          int64 `json:"hp"`
				}{}

				if err := json.Unmarshal(v.Data, &valid_data); err != nil {
					log.Print("json Unmarshall error in valid Hash Subscrive ", err)
					continue
				}
				timestamp := valid_data.Timestamp
				hp := valid_data.HP
				tx := fmt.Sprint(valid_data.Transaction)

				data := struct {
					isValid   bool  `json:"isValid"`
					timestamp int64 `json:"timestamp"`
				}{}
				if hash.IsOKHash(hp, timestamp, tx) {
					// validhash each に true をpub
					data.isValid = true
				} else {
					// validhash each に falseをpub
					data.isValid = false
				}
				data.timestamp = time.Now().Unix()
				c.Do("PUBLISH", validHashEachChan, data)
			} else {
				continue
			}
		default:
			log.Printf("error when psc.Receive switch-case in validHashSubScribe: %v", v)
		}
	}
}

func redis_get(key string, c redis.Conn) (string, error) {
	reply, err := c.Do("GET", key)
	if err != nil || reply == nil {
		return "", err
	}
	s, err := redis.String(reply, err)
	if err != nil {
		return "", err
	}
	return s, nil
}
