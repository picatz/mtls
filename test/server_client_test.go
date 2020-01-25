package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/picatz/mtls/cert"
	"github.com/picatz/mtls/client"
	"github.com/picatz/mtls/server"
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

func TestServerClient(t *testing.T) {
	caPEM, caPrivKeyPEM, err := cert.NewCA()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("CA Cert + Key")

	fmt.Println(string(caPEM))
	fmt.Println(string(caPrivKeyPEM))

	caPEMFile, err := writeToTempFile("caPEM", caPEM)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(caPEMFile)

	caPrivKeyPEMFile, err := writeToTempFile("caPrivKeyPEM", caPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(caPrivKeyPEMFile)

	// Server Key + Cert
	serverPEM, serverPrivKeyPEM, err := cert.NewServerFromCA(bytes.NewReader(caPrivKeyPEM))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Server Cert + Key")
	fmt.Println(string(serverPEM))
	fmt.Println(string(serverPrivKeyPEM))

	serverPEMFile, err := writeToTempFile("serverPEM", serverPEM)
	if err != nil {
		t.Fatal(err)
	}

	serverPrivKeyPEMFile, err := writeToTempFile("serverPrivKeyPEM", serverPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	// Client Key + Cert
	clientPEM, clientPrivKeyPEM, err := cert.NewClientFromCA(bytes.NewReader(caPrivKeyPEM))

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Client Cert + Key")
	fmt.Println(string(clientPEM))
	fmt.Println(string(clientPrivKeyPEM))

	clientPEMFile, err := writeToTempFile("clientPEM", serverPEM)
	if err != nil {
		t.Fatal(err)
	}

	clientPrivKeyPEMFile, err := writeToTempFile("clientPrivKeyPEM", clientPrivKeyPEM)
	if err != nil {
		t.Fatal(err)
	}

	// Server Instance
	go func() {
		s, err := server.New(
			server.WithTLSConfig(
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
	}()

	// Client Instance
	c, err := client.New(
		client.WithAddr("127.0.0.1:2222"),
		client.WithTLSConfig(
			tlsconf.BuildDefaultClientTLSConfig(
				caPEMFile,
				clientPEMFile,
				clientPrivKeyPEMFile,
			),
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
		t.Fatal(err)
	}

	tag := fmt.Sprintf("[%s -> %s]", conn.LocalAddr(), conn.RemoteAddr())

	if len(conn.ConnectionState().PeerCertificates) > 0 {
		log.Printf("%s client common name: %+v", tag, conn.ConnectionState().PeerCertificates[0].Subject.CommonName)
	}

	err = conn.Close()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(c)
}
