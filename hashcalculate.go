// ハッシュ計算君

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/hash"
	"github.com/InvincibleMan/anzu-chain/tx"
	"github.com/garyburd/redigo/redis"
)

type ValidHashJSON struct {
	Transaction tx.Tx `json:"tx"`
	Timestamp   int64 `json:"ts"`
	HP          int64 `json:"hp"`
}

func makePrev_block(prevhash string, myid string) string {
	return fmt.Sprintf("%s%s%s", prevhash, myid, time.Now().Unix())
}

func getPrevHash() string {
	latest_block_ := &dba.BlockAccess{}
	latest_block, err := latest_block_.GetLatestBlock()
	if err != nil {
		panic(err)
	}
	prevhash := latest_block.MyHash
	return prevhash
}

func HashCalculate(myid string, myhp int64, diff int64) {
	log.Println("HashCalculate: Goroutine Start")
	// redis connection
	conf := config.GetConfig()
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		log.Fatalf("Dead HashCalculate Goroutine because %v", err)
	}
	defer c.Close()
	log.Println("HashCalculate: connected to redis")
	//　1秒毎にハッシュ計算
	for {
		all_tx, err := redis.String(c.Do("GET", prevTxPoolKey))
		if err == redis.ErrNil {
			all_tx = ""
			err = nil
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(all_tx) == 0 {
			continue
		} else { //Tx Poolが空じゃない
			// ハッシュ計算

			timestamp := time.Now().Unix()
			prevhash := getPrevHash()
			prev_block := makePrev_block(prevhash, myid)
			// ハッシュ計算が完了
			// TxPoolをPrev_TxPoolに移す
			if hash.IsOKHash(myhp, diff, prev_block) {
				// txPoolを削除
				c.Do("SET", txPoolKey, "")
				// Prev_tx_poolに移す
				c.Do("SET", prevTxPoolKey, all_tx)
				// 報酬金のトランザクション
				sa_tx := tx.SendAsset{
					myid,
					"system",
					100,
				}
				// json 化
				tx_to_valid_hash := tx.Tx{
					[]tx.Commander{&sa_tx},
					myid,
					timestamp,
				}
				publish_request := ValidHashJSON{
					tx_to_valid_hash,
					timestamp,
					myhp,
				}
				json_byte, err := json.Marshal(
					publish_request,
				)
				if err != nil {
					panic(err)
				}
				// ValidHashにpublish
				c.Do("PUBLISH", validHashChan, json_byte)
			}
		}
		time.Sleep(1 * time.Second)
	}
}
