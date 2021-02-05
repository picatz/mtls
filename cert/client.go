package cert

import (
	"io"
)

func NewClientFromCA(caPrivKeyPEM, caCertPEM io.Reader, opts ...CertOption) ([]byte, []byte, error) {
	allOpts := []CertOption{}
	allOpts = append(allOpts, WithNewECDSAKey())
	allOpts = append(allOpts, opts...)
	allOpts = append(allOpts, IsClient())

	// Decode CA cert and private key from PEM encoded io.Reader bytes
	caCert, caPrivKey, err := ReadCertAndKey(caCertPEM, caPrivKeyPEM)
	if err != nil {
		return nil, nil, err
	}
	allOpts = append(allOpts, WithParent(caCert, caPrivKey))

	return New(allOpts...)
}
