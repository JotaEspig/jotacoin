package utils

import (
	"bytes"
	"encoding/binary"
)

// ToHex converts an integer to a hexadecimal as slice of bytes
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}

	return buff.Bytes()
}
