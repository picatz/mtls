package cert

import (
	"fmt"
	"testing"
)

func TestNewCA(t *testing.T) {
	caPEM, caPrivKeyPEM, err := NewCA()

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(caPEM))
	fmt.Println(string(caPrivKeyPEM))
}

func TestNewCAWithNewRSAKey(t *testing.T) {
	caPEM, caPrivKeyPEM, err := NewCA(WithNewRSAKey())

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(string(caPEM))
	fmt.Println(string(caPrivKeyPEM))
}
