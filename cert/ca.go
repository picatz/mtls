package cert

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"time"
)

func readCACertAndKey(caCertPEM, caPrivKeyPEM io.Reader) (*x509.Certificate, *ecdsa.PrivateKey, error) {
	// Decode CA cert from PEM encoded io.Reader bytes
	var caCert *x509.Certificate
	caCertPEMBytes, err := ioutil.ReadAll(caCertPEM)
	if err != nil {
		return nil, nil, err
	}
	pblock, _ := pem.Decode(caCertPEMBytes)
	if pblock != nil {
		crt, err := x509.ParseCertificate(pblock.Bytes)
		if err != nil {
			return nil, nil, err
		}
		caCert = crt
	}
	if caCert == nil {
		return nil, nil, fmt.Errorf("no CA cert found")
	}

	// Decode CA private key from PEM encoded io.Reader bytes
	var caPrivKey *ecdsa.PrivateKey
	caPrivKeyPEMBytes, err := ioutil.ReadAll(caPrivKeyPEM)
	if err != nil {
		return nil, nil, err
	}
	pblock, _ = pem.Decode(caPrivKeyPEMBytes)
	if pblock != nil {
		pk, err := x509.ParseECPrivateKey(pblock.Bytes)
		if err != nil {
			return nil, nil, err
		}
		caPrivKey = pk
	}
	if caPrivKey == nil {
		return nil, nil, fmt.Errorf("no CA priv key found")
	}

	// Return
	return caCert, caPrivKey, nil
}

func NewCA() ([]byte, []byte, error) {
	// TODO add rsa support/option
	// privKey, err := rsa.GenerateKey(rand.Reader, 4096)
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			CommonName: "ssh.ca.name",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // valid for 10 years
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privKey.PublicKey, privKey)
	if err != nil {
		return nil, nil, err
	}

	caPEMBuffer := new(bytes.Buffer)
	err = pem.Encode(caPEMBuffer, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})
	if err != nil {
		return nil, nil, err
	}

	b, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		return nil, nil, err
	}

	caPrivKeyPEMBuffer := new(bytes.Buffer)
	err = pem.Encode(caPrivKeyPEMBuffer, &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: b,
	})
	if err != nil {
		return nil, nil, err
	}

	return caPEMBuffer.Bytes(), caPrivKeyPEMBuffer.Bytes(), nil
}
