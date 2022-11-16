package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"log"
)

// ToHex converts an integer to a hexadecimal as slice of bytes
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// Serialize returns a []byte representative of data
func Serialize(data any) ([]byte, error) {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(data)
	if err != nil {
		return []byte{}, nil
	}

	return result.Bytes(), nil
}
