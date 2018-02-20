package tx

type Commander interface {
	Command()
}

type SendAsset struct {
	ToID   string
	FromID string
	Value  int64
}

func (s *SendAsset) Command() {}

type AddAsset struct {
	ToID  string
	Value int64
}

func (a *AddAsset) Command() {}

type CreateAccount struct {
	ID     string
	Pubkey string
}

func (c *CreateAccount) Command() {}

type Tx struct {
	Cmd       []Commander
	CreatorID string
	Timestamp int64
}
