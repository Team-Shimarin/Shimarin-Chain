package redis

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
		time.Sleep(1 * time.Second)
		all_tx, err := redis.String(c.Do("GET", txPoolKey))
		if err == redis.ErrNil {
			all_tx = ""
			err = nil
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(all_tx) == 0 {
			log.Printf("HashCalculate: nothing Tx in TxPool at %v", time.Now().Unix())
			continue
		} else { //Tx Poolが空じゃない
			log.Printf("HashCalculate: Tx is exist in TxPool at %v", time.Now().Unix())
			// ハッシュ計算
			timestamp := time.Now().Unix()
			prevhash := getPrevHash()
			sha256raw := makePrev_block(prevhash, myid)
			// ハッシュ計算が完了
			// TxPoolをPrev_TxPoolに移す
			log.Print(hash.IsOKHash(myhp, diff, sha256raw))
			if hash.IsOKHash(myhp, diff, sha256raw) {
				log.Printf("HashCalculate: success! hashcalc at %v", time.Now().Unix())
				// txPoolを削除
				c.Do("SET", txPoolKey, "")
				// Prev_tx_poolに移す
				c.Do("SET", prevTxPoolKey, all_tx)
				// all_txをあんまーしゃるする
				jsonAllTx := make([]tx.Tx, 0)
				if err := json.Unmarshal([]byte(all_tx), &jsonAllTx); err != nil {
					log.Fatal(err)
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
					log.Fatal("cannnot marshal")
				}
				// CreatorID, Timestamp,　を付与して送信
				js := struct {
					CreatorID string `json:creator_id`
					Timestamp int64  `json:timestamp`
					Txs       string `json:txs`
					HP        int64  `json:"hp"`
				}{
					CreatorID: myid,
					Timestamp: timestamp,
					Txs:       string(jBytes),
					HP:        myhp,
				}
				// ValidHashにpublish
				c.Do("PUBLISH", validHashChan, js)
				c.Do("SET", txPoolKey, "")
			} else {
				log.Printf("HashCalculate: missing PoH! at %v", time.Now().Unix())
			}
		}
	}
}
