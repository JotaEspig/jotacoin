package blockchain

import (
	"bytes"
	"encoding/gob"
)

// SerializedBlock represents a block that has passes through the process
// of serializing
type SerializedBlock []byte

// bytesToSerializedBlock cast the data to SerializedBlock. Then you can use the Deserialize method
func bytesToSerializedBlock(data []byte) SerializedBlock {
	return data
}

// Deserialize transforms a SerializedBlock into a Block
func (data SerializedBlock) Deserialize() (*Block, error) {
	block := &Block{}
	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(block)
	return block, err
}
