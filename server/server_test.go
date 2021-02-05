package server

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/picatz/mtls/cert"
	"github.com/picatz/mtls/client"
	"github.com/picatz/mtls/tlsconf"
	"github.com/stretchr/testify/require"
)

func writeToTempFile(t *testing.T, prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile("/tmp", prefix)
	require.NoError(t, err)
	require.NotNil(t, file)

	written, err := file.Write(data)
	require.NoError(t, err)
	require.Equal(t, len(data), written)

	err = file.Sync()
	require.NoError(t, err)

	return file.Name(), nil
}

func TestNewServer(t *testing.T) {
	caPEM, caPrivKeyPEM, err := cert.NewCA(
		cert.WithCommonName("ca.name"),
	)
	require.NoError(t, err)

	caPEMFile, err := writeToTempFile(t, "caPEM", caPEM)
	require.NoError(t, err)

	_, err = writeToTempFile(t, "caPrivKeyPEM", caPrivKeyPEM)
	require.NoError(t, err)

	serverPEM, serverPrivKeyPEM, err := cert.NewServerFromCA(
		bytes.NewReader(caPrivKeyPEM),
		bytes.NewReader(caPEM),
		cert.WithCommonName("server.name"),
	)
	require.NoError(t, err)

	serverPEMFile, err := writeToTempFile(t, "serverPEM", serverPEM)
	require.NoError(t, err)

	serverPrivKeyPEMFile, err := writeToTempFile(t, "serverPrivKeyPEM", serverPrivKeyPEM)
	require.NoError(t, err)

	s, err := New(
		WithTLSConfig(
			tlsconf.BuildDefaultServerTLSConfig(
				caPEMFile,
				serverPEMFile,
				serverPrivKeyPEMFile,
			),
		),
	)
	require.NoError(t, err)

	defer s.Listener().Close()
	fmt.Println(s)
}

func TestServerClient(t *testing.T) {
	// 1. Create new CA Key and Cert
	caPEM, caPrivKeyPEM, err := cert.NewCA(
		cert.WithCommonName("ca"),
	)
	require.NoError(t, err)

	caPEMFile, err := writeToTempFile(t, "caPEM", caPEM)
	require.NoError(t, err)

	// 2. Create new Server Key and Cert
	serverPEM, serverPrivKeyPEM, err := cert.NewServerFromCA(
		bytes.NewReader(caPrivKeyPEM),
		bytes.NewReader(caPEM),
		cert.WithCommonName("server.name"),
	)
	require.NoError(t, err)

	serverPEMFile, err := writeToTempFile(t, "serverPEM", serverPEM)
	require.NoError(t, err)

	serverPrivKeyPEMFile, err := writeToTempFile(t, "serverPrivKeyPEM", serverPrivKeyPEM)
	require.NoError(t, err)

	// 3. Create new Client Key and Cert
	clientPEM, clientPrivKeyPEM, err := cert.NewClientFromCA(
		bytes.NewReader(caPrivKeyPEM),
		bytes.NewReader(caPEM),
		cert.WithCommonName("client.name"),
	)
	require.NoError(t, err)

	clientPEMFile, err := writeToTempFile(t, "clientPEM", clientPEM)
	require.NoError(t, err)

	clientPrivKeyPEMFile, err := writeToTempFile(t, "clientPrivKeyPEM", clientPrivKeyPEM)
	require.NoError(t, err)

	serverTLSConf := tlsconf.BuildDefaultServerTLSConfig(
		caPEMFile,
		serverPEMFile,
		serverPrivKeyPEMFile,
	)

	// Deprecated:
	// serverTLSConf.BuildNameToCertificate()

	serverHandlerDoneForTest := make(chan struct{})

	// Server hanlder
	serverHandler := func(conn *tls.Conn) {
		defer func() {
			serverHandlerDoneForTest <- struct{}{}
		}()
		defer conn.Close()
		defer log.Println("server: closing tls conn")

		tag := fmt.Sprintf("[%s -> %s]", conn.LocalAddr(), conn.RemoteAddr())

		err := conn.Handshake()
		if err != nil {
			log.Println("server: handshake-error:", err)
		}
		if len(conn.ConnectionState().PeerCertificates) > 0 {
			log.Printf("server: %s client common name: %+v", tag, conn.ConnectionState().PeerCertificates[0].Subject.CommonName)
		}
		log.Println("server: done handling conn")
	}

	// Server Instance
	s, err := New(
		WithTLSConfig(
			serverTLSConf,
		),
		WithHandler(serverHandler),
	)
	require.NoError(t, err)

	defer s.Shutdown()
	s.Start()

	clientTLSConfig := tlsconf.BuildClientTLSConfigWithCustomVerification(
		caPEMFile,
		clientPEMFile,
		clientPrivKeyPEMFile,
		tlsconf.VerifyPeerCertificateInsecureAny,
		// tlsconf.VerifyFirstPeerCert(x509.VerifyOptions{
		// 	DNSName:     "127.0.0.1",
		// 	CurrentTime: time.Now(),
		// }),
	)

	// Deprecated:
	// clientTLSConfig.BuildNameToCertificate()

	// Client Instance
	c, err := client.New(
		client.WithAddr(DefaultAddr),
		client.WithTLSConfig(
			clientTLSConfig,
		),
	)
	require.NoError(t, err)

	conn, err := c.Dial()
	require.NoError(t, err)

	// this is required to complete the handshake and populate the connection state
	// we are doing this so we can print the peer certificates prior to reading / writing to the connection
	err = conn.Handshake()
	require.NoError(t, err)

	tag := fmt.Sprintf("[%s -> %s]", conn.LocalAddr(), conn.RemoteAddr())

	if len(conn.ConnectionState().PeerCertificates) > 0 {
		log.Printf("client: %s server common name: %+v", tag, conn.ConnectionState().PeerCertificates[0].Subject.CommonName)
	}

	// allow for server/client to complete jobs
	<-serverHandlerDoneForTest

	err = conn.Close()
	require.NoError(t, err)
}
