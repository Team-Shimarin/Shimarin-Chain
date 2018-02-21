package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Team-Shimarin/Shimarin-Chain/config"
	"github.com/Team-Shimarin/Shimarin-Chain/dba"
	"github.com/Team-Shimarin/Shimarin-Chain/hash"
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
			log.Println("validhHashChan get Message ", string(v.Data))
			publishedData := struct {
				Txs       string `json:"txs"`
				Timestamp int64  `json:"timestamp"`
				CreatorID string `json:"creator_id"`
				HP        int64  `json:"hp"`
			}{}
			if err := json.Unmarshal(v.Data, &publishedData); err != nil {
				panic(err)
			}
			log.Print("validHashSubScribe: publishedData is ", publishedData)
			if publishedData.Txs != "" {
				// 計算!
				ba := dba.BlockAccess{}
				prevHash, _ := ba.GetLatestBlockHash()
				// hash.IsOKHash(publishedData.HP, conf.Diff, fmt.Sprintf("%s%s%s", prevHash, publishedData.CreatorID, publishedData.Timestamp))
				//
				data := struct {
					IsValid   bool   `json:"isValid"`
					Timestamp int64  `json:"timestamp"`
					Data      string `json:"data"`
				}{}
				if hash.IsOKHash(publishedData.HP, conf.Diff, fmt.Sprintf("%s%s%s", prevHash, publishedData.CreatorID, publishedData.Timestamp)) {
					// validhash each に true をpub
					data.IsValid = true
				} else {
					// validhash each に falseをpub
					data.IsValid = false
				}
				data.Timestamp = time.Now().Unix()
				data.Data = string(v.Data)
				bsJSON, err := json.Marshal(data)
				if err != nil {
					panic(err)
				}
				pubc, _ := getRedisConn(conf.RedisHost, conf.RedisPort)
				pubc.Do("PUBLISH", validHashEachChan, string(bsJSON))
			} else {
				continue
			}
		case redis.Subscription:
			fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
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
