package model

// Block
// { prev_hash, txs, CreateID, hash, timestamp}
// txs [{To, From, Value}]

type Block struct {
	PrevHash string
	Txs []string
	CreatorID string
	Hash    string
	Timestamp int64
}

const BlockTable = "block"
