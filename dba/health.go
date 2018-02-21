package dba

import (
	"github.com/InvincibleMan/anzu-chain/model"
	"github.com/Masterminds/squirrel"
	sqlite3 "github.com/mattn/go-sqlite3"
	"log"
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

func (a *AccountAccess) GetHealth(accountid string) (*model.Health, error) {
	// HPを取得
	sql, args, err := squirrel.Select("*").
		From(model.HealthTable).Where("account_id == '" + accountid + "'").
		ToSql()
	if err != nil {
		return nil, err
	}
	account := model.Health{}
	log.Println(sql, args)
	if err := db.QueryRow(sql, args...).Scan(&account.Id,&account.Accout_id, &account.Hp); err != nil {
		return nil, err
	}

	return &account, nil
}

func (a *AccountAccess) InsertHealth(account *model.Account, healthpoint int64)(error){
	sql, args, err := squirrel.Insert(model.HealthTable).Columns("account_id", "hp").Values(account.ID, healthpoint).ToSql()
	if err != nil {
		log.Println(err)
		return err
	}
	if err := db.QueryRow(sql, args...).Scan(); err != nil {
		return err
	}
	return nil
}
