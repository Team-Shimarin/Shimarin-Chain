package model

type Tx struct {
	Txjson []byte
	StatusID int64
}

const TxTable = "tx"
