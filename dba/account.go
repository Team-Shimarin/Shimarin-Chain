package dba

import (
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/Masterminds/squirrel"
	sqlite3 "github.com/mattn/go-sqlite3"
)

type AccountAccess struct{}

// 初期から持ってるお金
const initalbalance = 100

func (a *AccountAccess) Register(account *model.Account) error {
	_, err := squirrel.Insert(model.AccountTable).
		Columns("id", "publickey", "balance").
		Values(account.ID, account.PublicKey, initalbalance).
		RunWith(db).
		Exec()
	if err != nil {
		if err.(sqlite3.Error).ExtendedCode == 2067 {
			return ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (a *AccountAccess) UpdataBalance(accountid string, addbalance int64) error {
	// DBから現在のbalanceを取得、追加するbalanceを計算し、UPDATE
	sql, args, err := squirrel.Select("balance").
		From(model.AccountTable).
		Where("id == " + accountid).
		ToSql()
	if err != nil {
		return err
	}
	account := model.Account{}
	if err := db.QueryRow(sql, args...).Scan(&account); err != nil {
		return nil
	}
	nowbalance := account.Balance
	_, err = squirrel.Update(model.AccountTable).Set(
			"balance", nowbalance + addbalance).Where("id", accountid).Exec()
	if err != nil {
		if err.(sqlite3.Error).ExtendedCode == 2067 {
			return ErrAlreadyExists
		}
		return err
	}
	return nil

}