package model

import uuid "github.com/satori/go.uuid"

type Account struct {
	ID        string
	PublicKey string
	Balance int64
	HP int64
}

func NewAccount(pubkey string) (*Account, error) {
	u, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	return &Account{
		ID:        u.String(),
		PublicKey: pubkey,
		Balance: 0,
	}, nil
}

const AccountTable = "account"
