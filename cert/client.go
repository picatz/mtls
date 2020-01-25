package cert

import (
	"io"
)

func NewClientFromCA(caPrivKeyPEM, caCertPEM io.Reader) ([]byte, []byte, error) {
	// Decode CA cert and private key from PEM encoded io.Reader bytes
	caCert, caPrivKey, err := ReadCertAndKey(caCertPEM, caPrivKeyPEM)
	if err != nil {
		return nil, nil, err
	}

	return New(
		WithParent(caCert, caPrivKey),
		WithNewECDSAKey(),
		IsClient(),
		WithCommonName("ssh.client.name"),
	)
}
