package model

import (
	uuid "github.com/satori/go.uuid"
)

type Account struct {
	ID      string
	Balance int64
}

func NewAccount() (*Account, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:      u.String(),
		Balance: 1000000000,
	}, nil
}

const AccountTable = "account"
