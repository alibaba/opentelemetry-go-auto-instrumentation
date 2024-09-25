// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import "testing"

func TestDbNameExtractor(t *testing.T) {
	dbSpanNameExtractor := DBSpanNameExtractor[testRequest]{
		Getter: mongoAttrsGetter{},
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
