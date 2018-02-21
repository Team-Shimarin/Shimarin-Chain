package dba

import (
	"github.com/Masterminds/squirrel"
	"github.com/Team-Shimarin/Shimarin-Chain/model"
	"github.com/Team-Shimarin/Shimarin-Chain/tx"
)

type AccountAccess struct{}

// 初期から持ってるお金
const initalbalance = 10000

func (a *AccountAccess) UpdateBalance(accountid string, addbalance int64) error {
	sql, args, err := squirrel.Select("balance").
		From(model.AccountTable).
		Where(squirrel.Eq{"id": accountid}).
		ToSql()
	if err != nil {
		return err
	}
	var nowBalance int64
	if err := db.QueryRow(sql, args...).Scan(&nowBalance); err != nil {
		return nil
	}
	_, err = squirrel.Update(model.AccountTable).Set(
		"balance", nowBalance+addbalance).Where(squirrel.Eq{"id": accountid}).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountAccess) GetBalance(accountid string) (int64, error) {
	// accountテーブルからBalanceを取得
	sql, args, err := squirrel.Select("balance").
		From(model.AccountTable).
		Where(squirrel.Eq{"id": accountid}).
		ToSql()
	if err != nil {
		return 0, err
	}
	var balance int64
	if err := db.QueryRow(sql, args...).Scan(&balance); err != nil {
		return 0, err
	}

	return balance, nil
}

func (a *AccountAccess) ApplyTx(tx *tx.Tx) error {
	// ふやす
	sql, args, err := squirrel.Select("balance").
		From(model.AccountTable).
		Where(squirrel.Eq{"id": tx.To}).
		ToSql()
	if err != nil {
		return err
	}
	var nowBalance int64
	if err := db.QueryRow(sql, args...).Scan(&nowBalance); err != nil {
		return nil
	}
	_, err = squirrel.Update(model.AccountTable).Set(
		"balance", nowBalance+tx.Value).Where(squirrel.Eq{"id": tx.To}).RunWith(db).Exec()
	if err != nil {
		return err
	}

	if tx.From == "system" {
		return nil
	}
	// へらす
	sql, args, err = squirrel.Select("balance").
		From(model.AccountTable).
		Where(squirrel.Eq{"id": tx.From}).
		ToSql()
	if err != nil {
		return err
	}
	nowBalance = 0
	if err := db.QueryRow(sql, args...).Scan(&nowBalance); err != nil {
		return nil
	}
	_, err = squirrel.Update(model.AccountTable).Set(
		"balance", nowBalance-tx.Value).Where(squirrel.Eq{"id": tx.From}).RunWith(db).Exec()
	if err != nil {
		return err
	}
	return nil
}
