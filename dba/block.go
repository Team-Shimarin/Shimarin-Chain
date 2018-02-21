package dba

import (
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/Masterminds/squirrel"
	"fmt"
	"log"
)

type BlockAccess struct{}

func (a *BlockAccess) AddBlock(block *model.Block) error {
	_, err := squirrel.Insert(model.BlockTable).
		Columns("prevhash", "txs", "creator_id", "timestamp", "hash").
		Values(block.PrevHash, fmt.Sprint(block.Txs), block.CreatorID, block.Timestamp, block.Hash).
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *BlockAccess) GetLatestBlock() (*model.Block, error) {
	sql, args, err := squirrel.Select("*").
		From(model.BlockTable).
		Where("timestamp = (SELECT MAX(timestamp) FROM block)").
		ToSql()
	if err != nil {
		return nil, err
	}
	log.Println(sql, args)
	block := model.Block{}
	//row, err := db.Query(sql)
	//if err != nil{
	//	log.Print(err)
	//	return nil, err
	//}
	if err := db.QueryRow(sql, args...).Scan(&block.PrevHash, &block.Txs, &block.CreatorID, &block.Hash, &block.Timestamp); err != nil {
		return nil, err
	}
	return &block, nil
}

func (a *BlockAccess) GetAllBlock()(*model.Block, error){
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