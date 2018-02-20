package main

import (
	"encoding/json"
	"fmt"
	"github.com/InvincibleMan/anzu-chain/hash"
	"github.com/InvincibleMan/anzu-chain/tx"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"time"
)

func ValidHashSubScribe(c redis.Conn) {
	log.Println("Goroutin in ValidHashSubScribe")
	psc := redis.PubSubConn{Conn: c}
	err := psc.Subscribe(validHashChan)
	if err != nil {
		log.Println(validHashChan, err)
	}
	for {
		log.Println(validHashChan, "connected")
		switch v := psc.Receive().(type) {
		case redis.Message:
			log.Println("validhHashChan get Message ", v.Data)
			if redis_get(txPoolKey, c) != "" {
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
		case error:
			log.Println(v)
		}
	}
}

func redis_get(key string, c redis.Conn) string {
	s, err := redis.String(c.Do("GET", key))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return s
}
