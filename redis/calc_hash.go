package redis

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Team-Shimarin/Shimarin-Chain/config"
	"github.com/Team-Shimarin/Shimarin-Chain/dba"
	"github.com/Team-Shimarin/Shimarin-Chain/hash"
	"github.com/Team-Shimarin/Shimarin-Chain/tx"
	"github.com/garyburd/redigo/redis"
)

func makePrev_block(prevhash string, myid string) string {
	return fmt.Sprintf("%s%s%s", prevhash, myid, time.Now().Unix())
}

func getPrevHash() (string, error) {
	latest_block_ := &dba.BlockAccess{}
	hash, err := latest_block_.GetLatestBlockHash()
	if err != nil {
		return "", err
	}
	return hash, nil
}

func HashCalculate(myid string, diff int64) {
	ha := dba.HealthAccess{}
	myhp, _ := ha.GetHealth(config.GetConfig().MinorAccountID)
	log.Println("HashCalculate: Goroutine Start")
	// redis connection
	conf := config.GetConfig()
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		log.Print("Dead HashCalculate Goroutine because %v", err)
		panic(err)
	}
	defer c.Close()
	log.Println("HashCalculate: connected to redis")
	//　1秒毎にハッシュ計算
	for {
		time.Sleep(1 * time.Second)
		all_tx, err := redis.String(c.Do("GET", txPoolKey))
		if err == redis.ErrNil {
			all_tx = ""
			err = nil
		}
		if err != nil {
			log.Print(err)
			panic(err)
		}
		if len(all_tx) <= 3 {
			log.Printf("HashCalculate: nothing Tx in TxPool at %v", time.Now().Unix())
			continue
		} else { //Tx Poolが空じゃない
			log.Printf("HashCalculate: Tx is exist in TxPool at %v", time.Now().Unix())
			// ハッシュ計算
			timestamp := time.Now().Unix()
			prevhash, _ := getPrevHash()
			sha256raw := makePrev_block(prevhash, myid)
			// ハッシュ計算が完了
			// TxPoolをPrev_TxPoolに移す
			log.Print(hash.IsOKHash(myhp, diff, sha256raw))
			if hash.IsOKHash(myhp, diff, sha256raw) && len(all_tx) >= 3 {
				log.Printf("HashCalculate: success! hashcalc at %v", time.Now().Unix())
				// txPoolを削除
				c.Do("SET", txPoolKey, "[]")
				// all_txをあんまーしゃるする
				jsonAllTx := make([]tx.Tx, 0)
				log.Print(all_tx)
				if err := json.Unmarshal([]byte(all_tx), &jsonAllTx); err != nil {
					log.Print(err)
					panic(err)
				}
				// 報酬金のトランザクション
				sa_tx := tx.Tx{
					To:    myid,
					From:  "system",
					Value: 100,
				}
				jsonAllTx = append(jsonAllTx, sa_tx)
				// json 化
				jBytes, err := json.Marshal(jsonAllTx)
				if err != nil {
					log.Print("cannnot marshal")
					panic(err)
				}
				// CreatorID, Timestamp,　を付与して送信
				js := struct {
					CreatorID string `json:"creator_id"`
					Timestamp int64  `json:"timestamp"`
					Txs       string `json:"txs"`
					HP        int64  `json:"hp"`
				}{
					CreatorID: myid,
					Timestamp: timestamp,
					Txs:       string(jBytes),
					HP:        myhp,
				}
				// json 化
				jjBytes, err := json.Marshal(js)
				if err != nil {
					log.Print("cannnot marshal")
					panic(err)
				}
				// ValidHashにpublish
				c.Do("PUBLISH", validHashChan, string(jjBytes))
				c.Do("SET", txPoolKey, "[]")
			} else {
				log.Printf("HashCalculate: missing PoH! at %v", time.Now().Unix())
			}
		}
	}
}
