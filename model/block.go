package model

type Block struct {
	PrevHash  string
	Txs       string
	CreatorID string
	Timestamp int64
	MyHash string
}

const BlockTable = "block"
