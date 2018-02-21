package dba

import (
	"github.com/Masterminds/squirrel"
	"github.com/Team-Shimarin/Shimarin-Chain/model"
)

type HealthAccess struct{}

func (a *HealthAccess) UpdateHP(accountid string, hp int64) error {
	// HPをアップデート
	_, err := squirrel.Update(model.HealthTable).Set("hp", hp).Where(squirrel.Eq{"account_id": accountid}).Exec()
	if err != nil {
		return err
	}
	return nil
}

func (a *HealthAccess) GetHealth(accountid string) (int64, error) {
	// HPを取得
	sql, args, err := squirrel.Select("hp").
		From(model.HealthTable).Where(squirrel.Eq{"account_id": accountid}).
		ToSql()
	if err != nil {
		return 0, err
	}
	var v int64
	if err := db.QueryRow(sql, args...).Scan(&v); err != nil {
		return 0, err
	}

	return v, nil
}
