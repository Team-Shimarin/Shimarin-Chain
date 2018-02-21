package redis

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	hashpack "github.com/InvincibleMan/anzu-chain/hash"
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/InvincibleMan/anzu-chain/tx"
	"github.com/garyburd/redigo/redis"
	"github.com/k0kubun/pp"
)

const validHashEachChan = "validHashEach"
const validHashChan = "validHash"
const txPoolKey = "TxPool"

type Cnt struct {
	cnt int
	f   bool
}

func SubscribeValidHashEach() {
	log.Println("subscribeValidHashEach: Goroutine Start")
	approveCnt := make(map[string]Cnt)
	rejectCnt := make(map[string]Cnt)
	conf := config.GetConfig()
	// redis connection
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		log.Printf("Dead subscribe_valid_hash_each Goroutine because %v", err)
		panic(err)
	}
	defer c.Close()
	log.Println("subscribeValidHashEach: conected to redis")

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe(validHashEachChan)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			log.Println("subscribeValidHashEach: catch message!")
			// cnt に加算する
			data := struct {
				IsValid   bool   `json:"isValid"`
				Timestamp int64  `json:"timestamp"`
				Data      string `json:"data"`
			}{}

			if err := json.Unmarshal(v.Data, &data); err != nil {
				panic(err)
			}

			if data.IsValid {
				if _, ok := approveCnt[data.Data]; !ok {
					approveCnt[data.Data] = Cnt{}
				}
				approveCnt[data.Data] = Cnt{
					cnt: approveCnt[data.Data].cnt + 1,
					f:   approveCnt[data.Data].f,
				}
			} else {
				if _, ok := rejectCnt[data.Data]; !ok {
					rejectCnt[data.Data] = Cnt{}
				}
				rejectCnt[data.Data] = Cnt{
					cnt: rejectCnt[data.Data].cnt + 1,
					f:   rejectCnt[data.Data].f,
				}
			}

			// 全ノード数取得(固定値)
			N := config.GetConfig().RedisNodeCount

			// CreatorID, Timestamp,
			jsonData := struct {
				Txs       string `json:txs`
				Timestamp int64  `json:"timestamp"`
				CreatorId string `json:"creator_id"`
				HP        int64  `json:"hp"`
			}{}
			if err := json.Unmarshal([]byte(data.Data), &jsonData); err != nil {
				log.Println(data.Data)
				panic(err)
			}
			// txs -> []*tx.tx
			txs := []*tx.Tx{}
			if err := json.Unmarshal([]byte(jsonData.Txs), &txs); err != nil {
				panic(err)
			}

			if approveCnt[data.Data].cnt > N/2 && !approveCnt[data.Data].f {
				// tx実行
				pp.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
				if err := executeTxs(txs); err != nil {
					panic(err)
				}

				ba := dba.BlockAccess{}
				prevhash, _ := ba.GetLatestBlockHash()

				if err := makeBlock(
					jsonData.Txs,
					jsonData.CreatorId,
					jsonData.Timestamp,
					fmt.Sprintf("%s%s%s", prevhash, jsonData.CreatorId, jsonData.HP),
				); err != nil {
					panic(err)
				}

				approveCnt[data.Data] = Cnt{
					cnt: approveCnt[data.Data].cnt,
					f:   true,
				}

			}
		case error:
			log.Print(v)
		}
	}
}

func makeBlock(txs string, creatorID string, timestamp int64, hash string) error {
	// HASH計算
	blockAccess := &dba.BlockAccess{}
	prevHash, err := blockAccess.GetLatestBlockHash()
	if err != nil {
		log.Println(err)
	}
	return blockAccess.AddBlock(
		&model.Block{
			PrevHash:  prevHash,
			Txs:       []string{txs},
			CreatorID: creatorID,
			Timestamp: timestamp,
			Hash:      string(hashpack.SHA256(hash)),
		},
	)

}

func executeTxs(txs []*tx.Tx) error {
	// TxをそれぞれExecuteする
	for _, tx := range txs {
		if err := executeTx(tx); err != nil {
			return err
		}
	}
	return nil
}

func executeTx(transaction *tx.Tx) error {
	aa := dba.AccountAccess{}
	return aa.ApplyTx(transaction)
}

// RedisからJSONのTxsを取得するFunc
func getTxsJSON(c redis.Conn) ([]byte, error) {
	bs, err := redis.Bytes(c.Do("GET", txPoolKey))
	if err == redis.ErrNil {
		return []byte("[]"), nil
	}
	return bs, err
}

// 構造体にするFunc
func getTxs(c redis.Conn) ([]*tx.Tx, error) {
	txsBytes, err := getTxsJSON(c)
	if err != nil {
		return nil, err
	}

	txs := make([]*tx.Tx, 0, 100)
	if err := json.Unmarshal(txsBytes, &txs); err != nil {
		return nil, err
	}

	return txs, nil
}
