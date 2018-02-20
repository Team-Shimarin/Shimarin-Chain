package dba

import (
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/Masterminds/squirrel"
	sqlite3 "github.com/mattn/go-sqlite3"
	)

func (a *AccountAccess) UpdataHP(accountid string, hp int64) error {
	// HPをアップデート
	_, err := squirrel.Update(model.HealthTable).Set(
		"hp", hp).Where("account_id", accountid).Exec()
	if err != nil {
		if err.(sqlite3.Error).ExtendedCode == 2067 {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}
