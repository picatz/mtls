package cert

import (
	"crypto/rand"
	"math/big"
)

// GenerateSerialNumber creates a new serial number for a x509 certificate.
func GenerateSerialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}
