package blockchain

// BlockChain Represents a chain of blocks
type BlockChain []*Block

// AddBlock adds a block into the chain of blocks
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := (*chain)[len(*chain)-1] // get last element
	newBlock := NewBlock(data, prevBlock.Hash)
	*chain = append(*chain, newBlock)
}

// NewBlockChain creates a new blockchain, starting with 'genesis'
func NewBlockChain() BlockChain {
	genesis := NewBlock("genesis", []byte{})
	return []*Block{genesis}
}
