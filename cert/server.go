package cert

import (
	"io"
)

func NewServerFromCA(caPrivKeyPEM, caCertPEM io.Reader, opts ...CertOption) ([]byte, []byte, error) {
	opts = append(opts, IsServer(), WithNewECDSAKey())
	// Decode CA cert and private key from PEM encoded io.Reader bytes
	caCert, caPrivKey, err := ReadCertAndKey(caCertPEM, caPrivKeyPEM)
	if err != nil {
		return nil, nil, err
	}
	opts = append(opts, WithParent(caCert, caPrivKey))

	return New(opts...)
}
