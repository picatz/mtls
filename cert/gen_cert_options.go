package cert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"time"
)

type CertOptions struct {
	parent struct {
		key  interface{}
		cert *x509.Certificate
	}
	key  interface{}
	cert *x509.Certificate
}

type CertOption func(*CertOptions) error

func WithParent(cert *x509.Certificate, key interface{}) CertOption {
	return func(o *CertOptions) error {
		o.parent.cert = cert
		o.parent.key = key
		return nil
	}
}

func WithKey(key interface{}) CertOption {
	return func(o *CertOptions) error {
		o.key = key
		return nil
	}
}

func WithNewECDSAKey() CertOption {
	return func(o *CertOptions) error {
		privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return err
		}
		o.key = privKey
		return nil
	}
}

func IsCA() CertOption {
	return func(o *CertOptions) error {
		o.cert.IsCA = true
		o.cert.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign
		o.cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}
		return nil
	}
}

func IsServer() CertOption {
	return func(o *CertOptions) error {
		o.cert.IsCA = false
		o.cert.KeyUsage = x509.KeyUsageDigitalSignature
		o.cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}
		return nil
	}
}

func IsClient() CertOption {
	return func(o *CertOptions) error {
		o.cert.IsCA = false
		o.cert.KeyUsage = x509.KeyUsageDigitalSignature
		o.cert.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}
		return nil
	}
}

func IsValidFor(d time.Duration) CertOption {
	return func(o *CertOptions) error {
		o.cert.NotBefore = time.Now()
		o.cert.NotAfter = time.Now().Add(d)
		return nil
	}
}

func WithCommonName(name string) CertOption {
	return func(o *CertOptions) error {
		o.cert.Subject = pkix.Name{
			CommonName: name,
		}
		return nil
	}
}
