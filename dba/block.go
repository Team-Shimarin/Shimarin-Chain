package dba

import (
	"encoding/json"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/Team-Shimarin/Shimarin-Chain/model"
)

type BlockAccess struct{}

func (a *BlockAccess) AddBlock(block *model.Block) error {
	txJSON, err := json.Marshal(block.Txs)
	if err != nil {
		return err
	}
	txStrJSON := string(txJSON)

	sql, args, err := squirrel.Insert(model.BlockTable).
		Columns("prevhash", "txs", "creator_id", "timestamp", "hash").
		Values(block.PrevHash, txStrJSON, block.CreatorID, block.Timestamp, block.Hash).
		ToSql()

	_, err = db.Exec(
		sql,
		args...,
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *BlockAccess) GetLatestBlockHash() (string, error) {
	sql, args, err := squirrel.Select("hash").
		From(model.BlockTable).
		Where("timestamp IN (SELECT MAX(timestamp) FROM block)").
		ToSql()
	if err != nil {
		return "", err
	}
	var res string
	if err := db.QueryRow(sql, args...).Scan(&res); err != nil {
		return "", err
	}
	return res, nil
}

func (a *BlockAccess) GetAllBlock() (*model.Block, error) {
	sql, args, err := squirrel.Select("*").
		From(model.BlockTable).
		ToSql()
	if err != nil {
		return nil, err
	}
	log.Println(sql, args)
	block := model.Block{}
	if err := db.QueryRow(sql, args...).Scan(&block.PrevHash, &block.Txs, &block.CreatorID, &block.Hash, &block.Timestamp); err != nil {
		return nil, err
	}

	return &block, nil
}
