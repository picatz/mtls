package cert

import "testing"

func TestGenerateSerialNumber(t *testing.T) {
	n, err := GenerateSerialNumber()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(n)
}
