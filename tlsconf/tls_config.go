package tlsconf

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func Build(opts ...TLSConfigOption) (*tls.Config, error) {
	config := &tls.Config{}

	for _, opt := range opts {
		err := opt(config)
		if err != nil {
			return nil, err
		}
	}

	//config.BuildNameToCertificate()
	return config, nil
}

type TLSConfigOption func(*tls.Config) error

func WithBaseConfig(baseConfig *tls.Config) TLSConfigOption {
	return func(config *tls.Config) error {
		config = baseConfig
		return nil
	}
}

func WithX509KeyPair(certFile string, keyFile string) TLSConfigOption {
	return func(config *tls.Config) error {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}
		config.Certificates = append(config.Certificates, cert)
		return nil
	}
}

func WithCAFile(caPEMFile string) TLSConfigOption {
	return func(config *tls.Config) error {
		file, err := os.Open(caPEMFile)
		if err != nil {
			return err
		}
		defer file.Close()
		return WithCACertificates(file)(config)
	}
}

func WithRootCAFile(caPEMFile string) TLSConfigOption {
	return func(config *tls.Config) error {
		caCert, err := ioutil.ReadFile(caPEMFile)
		if err != nil {
			log.Fatalf("failed to load cert: %s", err)
		}

		if config.RootCAs == nil {
			caCertPool := x509.NewCertPool()
			ok := caCertPool.AppendCertsFromPEM(caCert)
			if !ok {
				return fmt.Errorf("failed to append certs from PEM: %q", string(caCert))
			}
			config.RootCAs = caCertPool
		} else {
			config.RootCAs.AppendCertsFromPEM(caCert)
		}

		return nil
	}
}

func WithCertificates(certs []tls.Certificate) TLSConfigOption {
	return func(config *tls.Config) error {
		for _, cert := range certs {
			config.Certificates = append(config.Certificates, cert)
		}
		return nil
	}
}

func WithCACertificates(certs ...io.Reader) TLSConfigOption {
	return func(config *tls.Config) error {
		if config.ClientCAs == nil {
			config.ClientCAs = x509.NewCertPool()
		}

		for _, cert := range certs {
			bytes, err := ioutil.ReadAll(cert)
			if err != nil {
				return err
			}
			ok := config.ClientCAs.AppendCertsFromPEM(bytes)
			if !ok {
				return fmt.Errorf("unable to append cert from cert: \n%s", bytes)
			}
		}
		return nil
	}
}

func WithCurvePreferences(curveIDs []tls.CurveID) TLSConfigOption {
	return func(config *tls.Config) error {
		config.CurvePreferences = curveIDs
		return nil
	}
}

func WithPreferenceForServerCipherSuites() TLSConfigOption {
	return func(config *tls.Config) error {
		config.PreferServerCipherSuites = true
		return nil
	}
}

func WithMinVersion(version uint16) TLSConfigOption {
	return func(config *tls.Config) error {
		config.MinVersion = version
		return nil
	}
}

func WithCipherSuites(cipherSuites []uint16) TLSConfigOption {
	return func(config *tls.Config) error {
		config.CipherSuites = cipherSuites
		return nil
	}
}

func WithMutualAuthentication() TLSConfigOption {
	return func(config *tls.Config) error {
		config.ClientAuth = tls.RequireAndVerifyClientCert
		return nil
	}
}

func WithInsecureVerfication() TLSConfigOption {
	return func(config *tls.Config) error {
		config.InsecureSkipVerify = true
		return nil
	}
}

func BuildDefaultServerTLSConfig(caPemFile, serverCertFile, serverKeyFile string) *tls.Config {
	config, _ := Build(

		// used to verify the client cert is signed by the CA and is therefore valid
		WithCAFile(caPemFile),
		// server certificate which is validated by the client
		WithX509KeyPair(serverCertFile, serverKeyFile),
		// this requires a valid client certificate to be supplied during handshake
		WithMutualAuthentication(),
		WithMinVersion(tls.VersionTLS12),
		WithPreferenceForServerCipherSuites(),
		WithCurvePreferences([]tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}),
		WithCipherSuites([]uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		}),
	)

	return config
}

func BuildDefaultClientTLSConfig(caPemFile, clientCertFile, clientKeyFile string) *tls.Config {
	config, _ := Build(
		WithRootCAFile(caPemFile),
		WithX509KeyPair(clientCertFile, clientKeyFile),
		WithInsecureVerfication(), // TODO fix
	)

	return config
}
