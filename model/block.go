package model

import  (
	"github.com/InvincibleMan/anzu-chain/tx"
)

// Block
// { prev_hash, txs, CreateID, hash, timestamp}
// txs [{To, From, Value}]

type Block struct {
	PrevHash string
	Txs []tx.Tx
	CreatorID string
	Timestamp int64
	Hash    string
}

const BlockTable = "block"
