package cert

import (
	"bytes"
	"fmt"
	"testing"
)

func TestNewServerFromCA(t *testing.T) {
	caPEM, caPrivKeyPEM, err := NewCA(
		WithCommonName("server"),
	)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(caPEM))
	fmt.Println(string(caPrivKeyPEM))

	// Server Key + Cert
	serverPEM, serverPrivKeyPEM, err := NewServerFromCA(
		bytes.NewReader(caPrivKeyPEM),
		bytes.NewReader(caPEM),
		WithCommonName("server"),
	)

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(serverPEM))
	fmt.Println(string(serverPrivKeyPEM))
}
