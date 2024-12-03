// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package db

import "testing"

func TestDbNameExtractor(t *testing.T) {
	dbSpanNameExtractor := DBSpanNameExtractor[testRequest]{
		Getter: mongoAttrsGetter{},
	}
	if dbSpanNameExtractor.Extract(testRequest{}) != "DB Query" {
		t.Fatalf("Should have returned DB Query")
	}
	if dbSpanNameExtractor.Extract(testRequest{Name: "test", Operation: "test"}) != "test" {
		t.Fatalf("Should have returned test")
	}
	if dbSpanNameExtractor.Extract(testRequest{Operation: "op_test"}) != "op_test" {
		t.Fatalf("Should have returned op_test")
	}
}
