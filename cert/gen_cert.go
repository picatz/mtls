package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

func baseCert() *x509.Certificate {
	serialNumber, err := GenerateSerialNumber()
	if err != nil {
		panic(err)
	}

	return &x509.Certificate{
		SerialNumber:          serialNumber,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),  // valid for 10 years
		KeyUsage:              x509.KeyUsageDigitalSignature, // ex: x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{},          // ex: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		IsCA:                  false,
		BasicConstraintsValid: true,
	}
}

// New generates a PEM encoded x509 cert and private key.
func New(opts ...CertOption) ([]byte, []byte, error) {
	cerOpts := &CertOptions{
		key:  nil,
		cert: baseCert(),
	}

	for _, opt := range opts {
		err := opt(cerOpts)
		if err != nil {
			return nil, nil, err
		}
	}

	var (
		pubKey  interface{}
		privKey interface{}
	)
	switch k := cerOpts.key.(type) {
	case *ecdsa.PrivateKey:
		pubKey = &k.PublicKey
		privKey = k
	case *rsa.PrivateKey:
		pubKey = &k.PublicKey
		privKey = k
	default:
		panic(fmt.Sprintf("%T not implemented", k))
	}

	cert := cerOpts.cert

	var (
		certBytes []byte
		err       error
	)
	if cert.IsCA { // self sign
		certBytes, err = x509.CreateCertificate(rand.Reader, cert, cert, pubKey, privKey)
	} else { // sign with parent cert
		certBytes, err = x509.CreateCertificate(rand.Reader, cert, cerOpts.parent.cert, pubKey, cerOpts.parent.key)
	}

	if err != nil {
		return nil, nil, err
	}

	certPEMBuffer := new(bytes.Buffer)
	err = pem.Encode(certPEMBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err != nil {
		return nil, nil, err
	}

	b, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, nil, err
	}

	privKeyPEMBuffer := new(bytes.Buffer)
	err = pem.Encode(privKeyPEMBuffer, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: b,
	})
	if err != nil {
		return nil, nil, err
	}

	return certPEMBuffer.Bytes(), privKeyPEMBuffer.Bytes(), nil
}
