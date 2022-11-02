package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/mr-tron/base58/base58"
	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

// Wallet stores the Private key and the public key. More info at:
// https://blocktrade.com/wallet-addresses-public-and-private-keys-explained/
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte
}

// Address gets the address of the wallet
func (w *Wallet) Address() ([]byte, error) {
	pubHash, err := PublicKeyHash(w.PublicKey)
	if err != nil {
		return []byte{}, err
	}

	versionedHash := append([]byte{version}, pubHash...)
	checksumVal := checksum(versionedHash)

	fullHash := append(versionedHash, checksumVal...)
	return []byte(base58.Encode(fullHash)), nil
}

// NewWallet creates a new wallet with random keys
func NewWallet() (*Wallet, error) {
	private, pub, err := NewKeyPair()
	if err != nil {
		return nil, err
	}

	return &Wallet{private, pub}, nil
}

// NewKeyPair generates the private and the public key randomly
func NewKeyPair() (*ecdsa.PrivateKey, []byte, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, []byte{}, err
	}

	pub := append(private.X.Bytes(), private.Y.Bytes()...)
	return private, pub, err
}

//PublicKeyHash generates the hash of the public key
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

func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
