package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/hash"
	"github.com/garyburd/redigo/redis"
)

func ValidHashSubScribe() {
	log.Println("ValidHashSubScribe: Goutine start")
	conf := config.GetConfig()
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	defer c.Close()
	log.Println("ValidHashSubScribe: connected to redis")
	psc := redis.PubSubConn{Conn: c}
	err = psc.Subscribe(ValidHashChan)
	if err != nil {
		log.Print("error: subscribe ", ValidHashChan, ":because ", err)
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
				log.Fatal(err)
			}
			if publishedData.Txs != "" {
				// 計算!
				prevHash := "" // TODO: dbaから取得
				// hash.IsOKHash(publishedData.HP, conf.Diff, fmt.Sprintf("%s%s%s", prevHash, publishedData.CreatorID, publishedData.Timestamp))
				//
				data := struct {
					isValid   bool   `json:"isValid"`
					timestamp int64  `json:"timestamp"`
					data      string `json:"data"`
				}{}
				if hash.IsOKHash(publishedData.HP, conf.Diff, fmt.Sprintf("%s%s%s", prevHash, publishedData.CreatorID, publishedData.Timestamp)) {
					// validhash each に true をpub
					data.isValid = true
				} else {
					// validhash each に falseをpub
					data.isValid = false
				}
				data.timestamp = time.Now().Unix()
				data.data = string(v.Data)
				c.Do("PUBLISH", ValidHashEachChan, data)
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
