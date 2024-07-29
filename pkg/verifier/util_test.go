package verifier

import "testing"

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Fatal("port is 0")
	}
}
