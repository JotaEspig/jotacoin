package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

// Wallet stores the Private key and the public key
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

func NewWallet() (*Wallet, error) {
	private, pub, err := NewKeyPair()
	if err != nil {
		return nil, err
	}

	return &Wallet{private, pub}, nil
}

func NewKeyPair() (*ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, []byte{}, err
	}

	pub := append(private.X.Bytes(), private.Y.Bytes()...)
	return private, pub, err
}

func PublicKeyHash(pubKey []byte) ([]byte, error) {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		return []byte{}, err
	}

	publicRipMD := hasher.Sum(nil)
	return publicRipMD, nil
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
