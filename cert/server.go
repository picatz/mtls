package cert

import (
	"io"
)

func NewServerFromCA(caPrivKeyPEM, caCertPEM io.Reader) ([]byte, []byte, error) {
	// Decode CA cert and private key from PEM encoded io.Reader bytes
	caCert, caPrivKey, err := readCertAndKey(caCertPEM, caPrivKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	return New(
		WithParent(caCert, caPrivKey),
		WithNewECDSAKey(),
		IsServer(),
		WithCommonName("ssh.server.name"),
	)
}
