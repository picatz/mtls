package cert

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewClientFromCA(t *testing.T) {
	// CA Key + Cert
	caPEM, caPrivKeyPEM, err := NewCA(
		WithNewECDSAKey(),
		WithCommonName("ca"),
	)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(caPEM))
	fmt.Println(string(caPrivKeyPEM))

	// Client Key + Cert
	clientPEM, clientPrivKeyPEM, err := NewClientFromCA(
		bytes.NewReader(caPrivKeyPEM),
		bytes.NewReader(caPEM),
		WithNewECDSAKey(),
		WithCommonName("client"),
	)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(clientPEM))
	fmt.Println(string(clientPrivKeyPEM))
}
