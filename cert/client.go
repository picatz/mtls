package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"time"
)

func NewClientFromCA(caPrivKeyPEM, caCertPEM io.Reader) ([]byte, []byte, error) {
	// Decode CA cert and private key from PEM encoded io.Reader bytes
	caCert, caPrivKey, err := readCACertAndKey(caCertPEM, caPrivKeyPEM)

	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "ssh.client.name",
		},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // valid for 10 years
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, caCert, &privKey.PublicKey, caPrivKey)
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

	b, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return nil, nil, err
	}

	privKeyPEMBuffer := new(bytes.Buffer)
	err = pem.Encode(privKeyPEMBuffer, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: b,
	})
	if err != nil {
		return nil, nil, err
	}

	return certPEMBuffer.Bytes(), privKeyPEMBuffer.Bytes(), nil
}
