package http

import "testing"

func TestClientHttpStatusCodeConverter(t *testing.T) {
	c := ClientHttpStatusCodeConverter{}
	if c.IsError(200) {
		t.Fatalf("200 should not be an error")
	}
	if !c.IsError(600) || !c.IsError(90) {
		t.Fatalf("600 and 90 should be an error")
	}
}

func TestServerHttpStatusCodeConverter(t *testing.T) {
	c := ServerHttpStatusCodeConverter{}
	if c.IsError(200) {
		t.Fatalf("200 should not be an error")
	}
	if !c.IsError(500) || !c.IsError(90) {
		t.Fatalf("600 and 90 should be an error")
	}
}
