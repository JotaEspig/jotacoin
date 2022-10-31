package blockchain

// TxInput represents an input of a transaction. For more information:
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch05.html
type TxInput struct {
	PrevTxHash []byte // previous transacion ID, where the output ("balance") is stored
	OutIdx     int    // idx of output in the transaction struct
	Sig        string
}

// TxOutput represents an output of a transaction. For more information:
// https://www.oreilly.com/library/view/mastering-bitcoin/9781491902639/ch05.html
type TxOutput struct {
	Value  int
	PubKey string
}

// IsMadeBy checks if the address has made the input
func (txin *TxInput) IsMadeBy(address string) bool {
	return txin.Sig == address
}

// IsFor checks if the address is the receiver of the output
func (txout *TxOutput) IsFor(address string) bool {
	return txout.PubKey == address
}
