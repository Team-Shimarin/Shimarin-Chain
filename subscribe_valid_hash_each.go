package main

import (
	"encoding/json"
	"log"

	"github.com/InvincibleMan/anzu-chain/config"
	"github.com/InvincibleMan/anzu-chain/dba"
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/InvincibleMan/anzu-chain/tx"
	"github.com/garyburd/redigo/redis"
)

const validHashEachChan = "validHashEach"
const validHashChan = "validHash"
const prevTxPoolKey = "prev_TxPool"
const txPoolKey = "TxPool"

func subscribeValidHashEach() {
	log.Println("subscribeValidHashEach: Goroutine Start")
	approveCnt := 0
	rejectCnt := 0
	conf := config.GetConfig()
	// redis connection
	c, err := getRedisConn(conf.RedisHost, conf.RedisPort)
	if err != nil {
		log.Fatalf("Dead subscribe_valid_hash_each Goroutine because %v", err)
	}
	defer c.Close()
	log.Println("subscribeValidHashEach: conected to redis")

	psc := redis.PubSubConn{Conn: c}
	psc.Subscribe(validHashEachChan)
	for {
		switch v := psc.Receive().(type) {
		case redis.Message:
			// cnt に加算する
			data := struct {
				isValid   bool  `json:"isValid"`
				timestamp int64 `json:"timestamp"`
			}{}

			if err := json.Unmarshal(v.Data, &data); err != nil {
				log.Print(err)
				continue
			}

			if data.isValid {
				approveCnt += 1
			} else {
				rejectCnt += 1
			}

			// 全ノード数取得(固定値)
			N := config.GetConfig().RedisNodeCount

			// each cntをresetするfunc
			resetCnt := func() {
				approveCnt = 0
				rejectCnt = 0
			}

			// 複数回必要なので予めとっておく
			txs, err := getTxs(c)
			if err != nil {
				log.Print(err)
				continue
			}

			if approveCnt > N/2 {
				if err := executeTxs(txs); err != nil {
					log.Print(err)
					continue
				}

				if err := makeBlock(c, txs[len(txs)-1].CreatorID, data.timestamp); err != nil {
					log.Print(err)
					continue
				}

				resetCnt()
			} else if rejectCnt > N/2 {
				if err := backTxPool(c, txs); err != nil {
					log.Print(err)
					continue
				}

				resetCnt()
			}
		case error:
			log.Print(v)
		}
	}
}

func makeBlock(c redis.Conn, creatorID string, timestamp int64) error {
	// HASH計算
	blockAccess := &dba.BlockAccess{}
	latestblock, err := blockAccess.GetLatestBlock()
	if err != nil {
		log.Println(err)
	}
	prevhash := latestblock.MyHash

	// ブロック生成
	j, err := getTxsJSON(c)
	if err != nil {
		return err
	}

	return blockAccess.AddBlock(
		&model.Block{
			PrevHash:  prevhash,
			Txs:       string(j),
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
	for _, cmd := range transaction.Cmd {
		if c, ok := cmd.(*tx.AddAsset); ok {
			/// ADDASSET システムからお金をもらう
			//  DB接続なので、dba/の実装ができ次第かく
			txdata := tx.AddAsset{
				c.ToID,
				c.Value,
			}
			tx_json, err := json.Marshal(txdata)
			if err != nil {
				log.Println(err)
			}
			tx := model.Tx{
				// TODO: statuidを挿入
				[]byte(tx_json),
				1,
			}
			txaccsess := dba.TransactionAccsess{}
			err = txaccsess.AddTransaction(&tx)
			if err != nil {
				log.Println(err)
			}

			// アカウントにお金を追加
			accountaccess := dba.AccountAccess{}
			err = accountaccess.UpdataBalance(c.ToID, c.Value)
			if err != nil {
				log.Println(err)
			}
			log.Println(c)
		}
		if c, ok := cmd.(*tx.SendAsset); ok {
			// SENDASSET 送金を
			// DB接続なので、dba/の実装ができ次第かく
			txdata := tx.SendAsset{
				c.ToID,
				c.FromID,
				c.Value,
			}
			tx_json, err := json.Marshal(txdata)
			if err != nil {
				log.Println(err)
			}
			tx := model.Tx{
				// TODO: statuidを挿入
				[]byte(tx_json),
				2,
			}
			txaccsess := dba.TransactionAccsess{}
			err = txaccsess.AddTransaction(&tx)
			if err != nil {
				log.Println(err)
			}

			// アカウントにお金を追加
			accountaccess := dba.AccountAccess{}
			err = accountaccess.UpdataBalance(c.ToID, c.Value)
			if err != nil {
				log.Println(err)
			}

			// アカウントからお金を徴収
			err = accountaccess.UpdataBalance(c.FromID, -1*c.Value)
			if err != nil {
				log.Println(err)
			}
			log.Println(c)
		}
		if c, ok := cmd.(*tx.CreateAccount); ok {
			// DB接続なので、dba/の実装ができ次第かく
			account := model.Account{
				c.ID,
				c.Pubkey,
				0,
			}
			accountaccsess := dba.AccountAccess{}
			err := accountaccsess.Register(&account)
			if err != nil {
				log.Println(err)
			}
			log.Println(c)
		}
	}
	return err
}

func backTxPool(c redis.Conn, txs []*tx.Tx) error {
	// 一番ケツのTxのCreator(報酬生成Txなので)を見て、それが自分のIDなら、
	// TxPoolの先頭にprev_TxPoolを差し込み、prev_TxPoolを空にする
	if txs[len(txs)-1].CreatorID == config.GetConfig().MinorAccountID {
		// prevもってくる
		prevPool, err := redis.Bytes(c.Do("GET", prevTxPoolKey))
		if err == redis.ErrNil {
			err = nil
			prevPool = []byte("")
		}
		if err != nil {
			return err
		}
		// prevに空文字入れる
		if _, err := c.Do("SET", prevTxPoolKey, ""); err != nil {
			return err
		}
		// txpoolもってくる
		txPool, err := redis.Bytes(c.Do("GET", txPoolKey))
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
		if _, err := c.Do("SET", txPoolKey, string(newTxsJson)); err != nil {
			return err
		}
	}
	return nil
}

// RedisからJSONのTxsを取得するFunc
func getTxsJSON(c redis.Conn) ([]byte, error) {
	bs, err := redis.Bytes(c.Do("GET", prevTxPoolKey))
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
