// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package verifier

import (
	"testing"

	"go.opentelemetry.io/otel/attribute"
)

func TestGetFreePort(t *testing.T) {
	port, err := GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	if port == 0 {
		t.Fatal("port is 0")
	}
}

func TestGetAttribute(t *testing.T) {
	attrs := []attribute.KeyValue{
		attribute.Key("key").String("value"),
		attribute.Key("key1").String("value1"),
	}
	if GetAttribute(attrs, "key").AsString() != "value" {
		t.Fatal("key should exist")
	}
	if GetAttribute(attrs, "key2").Type() != attribute.INVALID {
		t.Fatal("key 2 should not exist")
	}
}

func TestAssert(t *testing.T) {
	defer func() {
		pass := false
		if r := recover(); r != nil {
			pass = true
		}
		if !pass {
			t.Fatal("Should be recovered from panic")
		}
	}()
	Assert(1 == 1, "1 should equal to 1")
	Assert(1 == 2, "1 should equal to 1")
}
