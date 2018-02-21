package tx

type Tx struct {
	To    string `json:"to"`
	From  string `json:"from"`
	Value int64  `json:"value"`
}
