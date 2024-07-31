package db

import "testing"

func TestDbNameExtractor(t *testing.T) {
	dbSpanNameExtractor := DBSpanNameExtractor[testRequest]{
		getter: mongoAttrsGetter{},
	}
	if dbSpanNameExtractor.Extract(testRequest{}) != "DB Query" {
		t.Fatalf("Should have returned DB_QUERY")
	}
	if dbSpanNameExtractor.Extract(testRequest{Name: "test"}) != "test" {
		t.Fatalf("Should have returned test")
	}
	if dbSpanNameExtractor.Extract(testRequest{Operation: "op_test"}) != "op_test" {
		t.Fatalf("Should have returned op_test")
	}
}
