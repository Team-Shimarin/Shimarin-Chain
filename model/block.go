package model

type Block struct {
	PrevHash string
	// TODO : Txs の定義が終わり次第修正
	CreatorID string
	Timestamp int64
	MyHash    string
}

const BlockTable = "block"
