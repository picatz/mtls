package cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
)

func SaveCertAndKey(prefix string, caCertPEM, caPrivKeyPEM []byte) error {
	err := ioutil.WriteFile(prefix+".cert.pem", caCertPEM, 0644)
	if err != nil {
		return err
	}
	ioutil.WriteFile(prefix+".priv.key.pem", caPrivKeyPEM, 0600)
	if err != nil {
		return err
	}
	return nil
}

func ReadCertAndKey(caCertPEM, caPrivKeyPEM io.Reader) (*x509.Certificate, interface{}, error) {
	// Decode CA cert from PEM encoded io.Reader bytes
	var cert *x509.Certificate
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
		cert = crt
	}
	if cert == nil {
		return nil, nil, fmt.Errorf("no cert found")
	}

	// Decode CA private key from PEM encoded io.Reader bytes
	// TODO(kent): Support decoding RSA PEM encoded file
	var key interface{}
	caPrivKeyPEMBytes, err := ioutil.ReadAll(caPrivKeyPEM)
	if err != nil {
		return nil, nil, err
	}
	pblock, _ = pem.Decode(caPrivKeyPEMBytes)
	if pblock != nil {
		pk, err := x509.ParsePKCS8PrivateKey(pblock.Bytes)
		if err != nil {
			return nil, nil, err
		}
		key = pk
	}
	if key == nil {
		return nil, nil, fmt.Errorf("no priv key found")
	}

	// Return
	return cert, key, nil
}

func NewCA(opts ...CertOption) ([]byte, []byte, error) {
	opts = append(opts, IsCA(), WithNewECDSAKey())
	return New(
		opts...,
	)
}
