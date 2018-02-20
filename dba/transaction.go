package dba

import (
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/Masterminds/squirrel"
)

type TransactionAccsess struct{}

func (a *TransactionAccsess) AddTransaction(tx *model.Tx) error {
	_, err := squirrel.Insert(model.TxTable).
		Columns("tx", "StatusID").
			// TODO: tx のjsonとStatusIDを入力
		Values().
		RunWith(db).
		Exec()
	if err != nil {
		return err
	}
	return nil
}