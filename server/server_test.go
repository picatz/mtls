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
)

func writeToTempFile(prefix string, data []byte) (string, error) {
	file, err := ioutil.TempFile("/tmp", "prefix")
	if err != nil {
		return "", err
	}
	written, err := file.Write(data)
	if written != len(data) {
		return "", fmt.Errorf("%q expected %d got %d", prefix, written, len(data))
	}
	err = file.Sync()
	if err != nil {
		return "", err
	}
	return file.Name(), nil
}

func TestNewServer(t *testing.T) {
	caPEM, caPrivKeyPEM, err := cert.NewCA()

	if err != nil {
		t.Fatal(err)
	}

	caPEMFile, err := writeToTempFile("caPEM", caPEM)
	if err != nil {
		t.Fatal(err)
	}

	_, err = writeToTempFile("caPrivKeyPEM", caPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	serverPEM, serverPrivKeyPEM, err := cert.NewServerFromCA(bytes.NewReader(caPrivKeyPEM), bytes.NewReader(caPEM))

	if err != nil {
		t.Fatal(err)
	}

	serverPEMFile, err := writeToTempFile("serverPEM", serverPEM)
	if err != nil {
		t.Fatal(err)
	}

	serverPrivKeyPEMFile, err := writeToTempFile("serverPrivKeyPEM", serverPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	s, err := New(
		WithTLSConfig(
			tlsconf.BuildDefaultServerTLSConfig(
				caPEMFile,
				serverPEMFile,
				serverPrivKeyPEMFile,
			),
		),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Listener().Close()
	fmt.Println(s)
}

func TestServerClient(t *testing.T) {
	// 1. Create new CA Key and Cert
	caPEM, caPrivKeyPEM, err := cert.NewCA()
	if err != nil {
		t.Fatal(err)
	}

	caPEMFile, err := writeToTempFile("caPEM", caPEM)
	if err != nil {
		t.Fatal(err)
	}

	_, err = writeToTempFile("caPrivKeyPEM", caPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	// 2. Create new Server Key and Cert
	serverPEM, serverPrivKeyPEM, err := cert.NewServerFromCA(bytes.NewReader(caPrivKeyPEM), bytes.NewReader(caPEM))
	if err != nil {
		t.Fatal(err)
	}

	serverPEMFile, err := writeToTempFile("serverPEM", serverPEM)
	if err != nil {
		t.Fatal(err)
	}

	serverPrivKeyPEMFile, err := writeToTempFile("serverPrivKeyPEM", serverPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Create new Client Key and Cert
	clientPEM, clientPrivKeyPEM, err := cert.NewClientFromCA(bytes.NewReader(caPrivKeyPEM), bytes.NewReader(caPEM))
	if err != nil {
		t.Fatal(err)
	}

	clientPEMFile, err := writeToTempFile("clientPEM", clientPEM)
	if err != nil {
		t.Fatal(err)
	}

	clientPrivKeyPEMFile, err := writeToTempFile("clientPrivKeyPEM", clientPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	serverTLSConf := tlsconf.BuildDefaultServerTLSConfig(
		caPEMFile,
		serverPEMFile,
		serverPrivKeyPEMFile,
	)

	serverTLSConf.BuildNameToCertificate()

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
	if err != nil {
		t.Fatal("new server:", err)
	}
	defer s.Shutdown()
	s.Start()

	clientTLSConfig := tlsconf.BuildDefaultClientTLSConfig(
		caPEMFile,
		clientPEMFile,
		clientPrivKeyPEMFile,
	)

	clientTLSConfig.BuildNameToCertificate()

	// Client Instance
	c, err := client.New(
		client.WithAddr("127.0.0.1:2222"),
		client.WithTLSConfig(
			clientTLSConfig,
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	conn, err := c.Dial()
	if err != nil {
		t.Fatal(err)
	}

	// this is required to complete the handshake and populate the connection state
	// we are doing this so we can print the peer certificates prior to reading / writing to the connection
	err = conn.Handshake()
	if err != nil {
		t.Fatal("client handshake", err)
	}

	tag := fmt.Sprintf("[%s -> %s]", conn.LocalAddr(), conn.RemoteAddr())

	if len(conn.ConnectionState().PeerCertificates) > 0 {
		log.Printf("client: %s server common name: %+v", tag, conn.ConnectionState().PeerCertificates[0].Subject.CommonName)
	}

	// allow for server/client to complete jobs
	<-serverHandlerDoneForTest

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
}
