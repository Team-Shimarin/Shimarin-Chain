package redis

import (
	"encoding/json"
	"log"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/InvincibleMan/anzu-chain/tx"
	"github.com/garyburd/redigo/redis"
)

const ValidHashEachChan = "validHashEach"
const ValidHashChan = "validHash"
const PrevTxPoolKey = "prev_TxPool"
const TxPoolKey = "TxPool"

func SubscribeValidHashEach() {
	log.Println("subscribeValidHashEach: Goroutine Start")
	approveCnt := make(map[string]int)
	rejectCnt := make(map[string]int)
	conf := config.GetConfig()
	// redis connection
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		log.Fatalf("Dead subscribe_valid_hash_each Goroutine because %v", err)
	}
	defer c.Close()
	log.Println("subscribeValidHashEach: conected to redis")

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe(ValidHashEachChan)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			// cnt に加算する
			data := struct {
				isValid   bool   `json:"isValid"`
				timestamp int64  `json:"timestamp"`
				data      string `json:"data"`
			}{}

			if err := json.Unmarshal(v.Data, &data); err != nil {
				log.Print(err)
				continue
			}

			if data.isValid {
				if _, ok := approveCnt[data.data]; !ok {
					approveCnt[data.data] = 0
				}
				approveCnt[data.data] += 1
			} else {
				if _, ok := rejectCnt[data.data]; !ok {
					rejectCnt[data.data] = 0
				}
				rejectCnt[data.data] += 1
			}

			// 全ノード数取得(固定値)
			N := config.GetConfig().RedisNodeCount

			// CreatorID, Timestamp,
			jsonData := struct {
				txs       string `json:txs`
				timestamp int64  `json:"timestamp"`
				creatorId string `json:"creator_id"`
			}{}
			if err := json.Unmarshal([]byte(data.data), &jsonData); err != nil {
				log.Print(err)
				continue
			}
			// txs -> []*tx.tx
			txs := []*tx.Tx{}
			if err := json.Unmarshal([]byte(jsonData.txs), &txs); err != nil {
				log.Print(err)
				continue
			}

			if approveCnt[data.data] > N/2 {
				// tx実行
				if err := executeTxs(txs); err != nil {
					log.Print(err)
					continue
				}

				if err := makeBlock(jsonData.txs, jsonData.creatorId, jsonData.timestamp); err != nil {
					log.Print(err)
					continue
				}

			} else if rejectCnt[data.data] > N/2 && conf.MinorAccountID == jsonData.creatorId {
				if err := backTxPool(c, txs); err != nil {
					log.Print(err)
					continue
				}

			}
		case error:
			log.Print(v)
		}
	}
}

func makeBlock(txs string, creatorID string, timestamp int64) error {
	// HASH計算
	blockAccess := &dba.BlockAccess{}
	latestblock, err := blockAccess.GetLatestBlock()
	if err != nil {
		log.Println(err)
	}
	prevhash := latestblock.Hash
	return blockAccess.AddBlock(
		&model.Block{
			PrevHash:  prevhash,
			Txs:       []string{txs},
			CreatorID: creatorID,
			Timestamp: timestamp,
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
	var err error
	// TODO: おかねをふやす、へらすDBAそうさ
	return err
}

func backTxPool(c redis.Conn, txs []*tx.Tx) error {
	// 一番ケツのTxのCreator(報酬生成Txなので)を見て、それが自分のIDなら、
	// TxPoolの先頭にprev_TxPoolを差し込み、prev_TxPoolを空にする
	// prevもってくる
	prevPool, err := redis.Bytes(c.Do("GET", PrevTxPoolKey))
	if err == redis.ErrNil {
		err = nil
		prevPool = []byte("")
	}
	if err != nil {
		return err
	}
	// prevに空文字入れる
	if _, err := c.Do("SET", PrevTxPoolKey, ""); err != nil {
		return err
	}
	// txpoolもってくる
	txPool, err := redis.Bytes(c.Do("GET", TxPoolKey))
	if err == redis.ErrNil {
		err = nil
		txPool = []byte("")
	}
	if err != nil {
		return err
	}
	// アンマーシャルする
	prevTx := make([]*tx.Tx, 0, 100)
	if err := json.Unmarshal(prevPool, &prevTx); err != nil {
		return err
	}
	poolTx := make([]*tx.Tx, 0, 100)
	if err := json.Unmarshal(txPool, &poolTx); err != nil {
		return err
	}
	// 頭につける
	newTxs := append(poolTx, prevTx...)
	// マーシャルする
	newTxsJson, err := json.Marshal(newTxs)
	if err != nil {
		return err
	}
	// txpoolに入れる
	if _, err := c.Do("SET", TxPoolKey, string(newTxsJson)); err != nil {
		return err
	}
	return nil
}

// RedisからJSONのTxsを取得するFunc
func getTxsJSON(c redis.Conn) ([]byte, error) {
	bs, err := redis.Bytes(c.Do("GET", PrevTxPoolKey))
	if err == redis.ErrNil {
		return []byte(""), nil
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
