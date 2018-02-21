package model

// Block
// { prev_hash, txs, CreateID, hash, timestamp}
// txs [{To, From, Value}]

type Block struct {
	PrevHash string
	Txs string
	CreatorID string
	Timestamp int64
	Hash    string
}

const BlockTable = "block"
