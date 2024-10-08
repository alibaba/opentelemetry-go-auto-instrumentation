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

package version

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewVersion(t *testing.T) {
	cases := []struct {
		version string
		err     bool
	}{
		{"", true},
		{"1.2.3", false},
		{"1.0", false},
		{"1", false},
		{"1.2.beta", true},
		{"1.21.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"1.2.3-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"1.2.0.4-x.Y.0+metadata", false},
		{"1.2.0.4-x.Y.0+metadata-width-hyphen", false},
		{"1.2.0-X-1.2.0+metadata~dist", false},
		{"1.2.3.4-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"v1.2.3", false},
		{"foo1.2.3", true},
		{"1.7rc2", false},
		{"v1.7rc2", false},
		{"1.0-", false},
	}

	for _, tc := range cases {
		_, err := NewVersion(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %q", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %q: %s", tc.version, err)
		}
	}
}

func TestVersionCompare(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.4.5", -1},
		{"1.2-beta", "1.2-beta", 0},
		{"1.2", "1.1.4", 1},
		{"1.2", "1.2-beta", 1},
		{"1.2+foo", "1.2+beta", 0},
		{"v1.2", "v1.2-beta", 1},
		{"v1.2+foo", "v1.2+beta", 0},
		{"v1.2.3.4", "v1.2.3.4", 0},
		{"v1.2.0.0", "v1.2", 0},
		{"v1.2.0.0.1", "v1.2", 1},
		{"v1.2", "v1.2.0.0", 0},
		{"v1.2", "v1.2.0.0.1", -1},
		{"v1.2.0.0", "v1.2.0.0.1", -1},
		{"v1.2.3.0", "v1.2.3.4", -1},
		{"1.7rc2", "1.7rc1", 1},
		{"1.7rc2", "1.7", -1},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", 1},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Compare(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %d\nactual: %d",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestVersionEqual_nil(t *testing.T) {
	mustVersion := func(v string) *Version {
		ver, err := NewVersion(v)
		if err != nil {
			t.Fatal(err)
		}
		return ver
	}
	cases := []struct {
		leftVersion  *Version
		rightVersion *Version
		expected     bool
	}{
		{mustVersion("1.0.0"), nil, false},
		{nil, mustVersion("1.0.0"), false},
		{nil, nil, true},
	}

	for _, tc := range cases {
		given := tc.leftVersion.Equal(tc.rightVersion)
		if given != tc.expected {
			t.Fatalf("expected Equal to nil to be %t", tc.expected)
		}
	}
}

func TestComparePreReleases(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2-beta.2", "1.2-beta.2", 0},
		{"1.2-beta.1", "1.2-beta.2", -1},
		{"1.2-beta.2", "1.2-beta.11", -1},
		{"3.2-alpha.1", "3.2-alpha", 1},
		{"1.2-beta.2", "1.2-beta.1", 1},
		{"1.2-beta.11", "1.2-beta.2", 1},
		{"1.2-beta", "1.2-beta.3", -1},
		{"1.2-alpha", "1.2-beta.3", -1},
		{"1.2-beta", "1.2-alpha.3", 1},
		{"3.0-alpha.3", "3.0-rc.1", -1},
		{"3.0-alpha3", "3.0-rc1", -1},
		{"3.0-alpha.1", "3.0-alpha.beta", -1},
		{"5.4-alpha", "5.4-alpha.beta", 1},
		{"v1.2-beta.2", "v1.2-beta.2", 0},
		{"v1.2-beta.1", "v1.2-beta.2", -1},
		{"v3.2-alpha.1", "v3.2-alpha", 1},
		{"v3.2-rc.1-1-g123", "v3.2-rc.2", 1},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Compare(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %d\nactual: %d",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestVersionMetadata(t *testing.T) {
	cases := []struct {
		version  string
		expected string
	}{
		{"1.2.3", ""},
		{"1.2-beta", ""},
		{"1.2.0-x.Y.0", ""},
		{"1.2.0-x.Y.0+metadata", "metadata"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "metadata~dist"},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v.Metadata()
		expected := tc.expected
		if actual != expected {
			t.Fatalf("expected: %s\nactual: %s", expected, actual)
		}
	}
}

func TestVersionPrerelease(t *testing.T) {
	cases := []struct {
		version  string
		expected string
	}{
		{"1.2.3", ""},
		{"1.2-beta", "beta"},
		{"1.2.0-x.Y.0", "x.Y.0"},
		{"1.2.0-7.Y.0", "7.Y.0"},
		{"1.2.0-x.Y.0+metadata", "x.Y.0"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "metadata-1.2.0"},
		{"17.03.0-ce", "ce"}, // zero-padded fields
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v.Prerelease()
		expected := tc.expected
		if actual != expected {
			t.Fatalf("expected: %s\nactual: %s", expected, actual)
		}
	}
}

func TestVersionSegments(t *testing.T) {
	cases := []struct {
		version  string
		expected []int
	}{
		{"1.2.3", []int{1, 2, 3}},
		{"1.2-beta", []int{1, 2, 0}},
		{"1-x.Y.0", []int{1, 0, 0}},
		{"1.2.0-x.Y.0+metadata", []int{1, 2, 0}},
		{"1.2.0-metadata-1.2.0+metadata~dist", []int{1, 2, 0}},
		{"17.03.0-ce", []int{17, 3, 0}}, // zero-padded fields
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v.Segments()
		expected := tc.expected
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %#v\nactual: %#v", expected, actual)
		}
	}
}

func TestVersionSegments64(t *testing.T) {
	cases := []struct {
		version  string
		expected []int64
	}{
		{"1.2.3", []int64{1, 2, 3}},
		{"1.2-beta", []int64{1, 2, 0}},
		{"1-x.Y.0", []int64{1, 0, 0}},
		{"1.2.0-x.Y.0+metadata", []int64{1, 2, 0}},
		{"1.4.9223372036854775807", []int64{1, 4, 9223372036854775807}},
	}

	for _, tc := range cases {
		v, err := NewVersion(tc.version)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v.Segments64()
		expected := tc.expected
		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %#v\nactual: %#v", expected, actual)
		}

		{
			expected := actual[0]
			actual[0]++
			actual = v.Segments64()
			if actual[0] != expected {
				t.Fatalf("Segments64 is mutable")
			}
		}
	}
}

func TestVersionString(t *testing.T) {
	cases := [][]string{
		{"1.2.3", "1.2.3"},
		{"1.2-beta", "1.2.0-beta"},
		{"1.2.0-x.Y.0", "1.2.0-x.Y.0"},
		{"1.2.0-x.Y.0+metadata", "1.2.0-x.Y.0+metadata"},
		{"1.2.0-metadata-1.2.0+metadata~dist", "1.2.0-metadata-1.2.0+metadata~dist"},
		{"17.03.0-ce", "17.3.0-ce"}, // zero-padded fields
	}

	for _, tc := range cases {
		v, err := NewVersion(tc[0])
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v.String()
		expected := tc[1]
		if actual != expected {
			t.Fatalf("expected: %s\nactual: %s", expected, actual)
		}
		if actual := v.Original(); actual != tc[0] {
			t.Fatalf("expected original: %q\nactual: %q", tc[0], actual)
		}
	}
}

func TestEqual(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", false},
		{"1.2-beta", "1.2-beta", true},
		{"1.2", "1.1.4", false},
		{"1.2", "1.2-beta", false},
		{"1.2+foo", "1.2+beta", true},
		{"v1.2", "v1.2-beta", false},
		{"v1.2+foo", "v1.2+beta", true},
		{"v1.2.3.4", "v1.2.3.4", true},
		{"v1.2.0.0", "v1.2", true},
		{"v1.2.0.0.1", "v1.2", false},
		{"v1.2", "v1.2.0.0", true},
		{"v1.2", "v1.2.0.0.1", false},
		{"v1.2.0.0", "v1.2.0.0.1", false},
		{"v1.2.3.0", "v1.2.3.4", false},
		{"1.7rc2", "1.7rc1", false},
		{"1.7rc2", "1.7", false},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", false},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Equal(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", false},
		{"1.2-beta", "1.2-beta", false},
		{"1.2", "1.1.4", true},
		{"1.2", "1.2-beta", true},
		{"1.2+foo", "1.2+beta", false},
		{"v1.2", "v1.2-beta", true},
		{"v1.2+foo", "v1.2+beta", false},
		{"v1.2.3.4", "v1.2.3.4", false},
		{"v1.2.0.0", "v1.2", false},
		{"v1.2.0.0.1", "v1.2", true},
		{"v1.2", "v1.2.0.0", false},
		{"v1.2", "v1.2.0.0.1", false},
		{"v1.2.0.0", "v1.2.0.0.1", false},
		{"v1.2.3.0", "v1.2.3.4", false},
		{"1.7rc2", "1.7rc1", true},
		{"1.7rc2", "1.7", false},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", true},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.GreaterThan(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s > %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestLessThan(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", true},
		{"1.2-beta", "1.2-beta", false},
		{"1.2", "1.1.4", false},
		{"1.2", "1.2-beta", false},
		{"1.2+foo", "1.2+beta", false},
		{"v1.2", "v1.2-beta", false},
		{"v1.2+foo", "v1.2+beta", false},
		{"v1.2.3.4", "v1.2.3.4", false},
		{"v1.2.0.0", "v1.2", false},
		{"v1.2.0.0.1", "v1.2", false},
		{"v1.2", "v1.2.0.0", false},
		{"v1.2", "v1.2.0.0.1", true},
		{"v1.2.0.0", "v1.2.0.0.1", true},
		{"v1.2.3.0", "v1.2.3.4", true},
		{"1.7rc2", "1.7rc1", false},
		{"1.7rc2", "1.7", true},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", false},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.LessThan(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s < %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestGreaterThanOrEqual(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", false},
		{"1.2-beta", "1.2-beta", true},
		{"1.2", "1.1.4", true},
		{"1.2", "1.2-beta", true},
		{"1.2+foo", "1.2+beta", true},
		{"v1.2", "v1.2-beta", true},
		{"v1.2+foo", "v1.2+beta", true},
		{"v1.2.3.4", "v1.2.3.4", true},
		{"v1.2.0.0", "v1.2", true},
		{"v1.2.0.0.1", "v1.2", true},
		{"v1.2", "v1.2.0.0", true},
		{"v1.2", "v1.2.0.0.1", false},
		{"v1.2.0.0", "v1.2.0.0.1", false},
		{"v1.2.3.0", "v1.2.3.4", false},
		{"1.7rc2", "1.7rc1", true},
		{"1.7rc2", "1.7", false},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", true},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.GreaterThanOrEqual(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s >= %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestLessThanOrEqual(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected bool
	}{
		{"1.2.3", "1.4.5", true},
		{"1.2-beta", "1.2-beta", true},
		{"1.2", "1.1.4", false},
		{"1.2", "1.2-beta", false},
		{"1.2+foo", "1.2+beta", true},
		{"v1.2", "v1.2-beta", false},
		{"v1.2+foo", "v1.2+beta", true},
		{"v1.2.3.4", "v1.2.3.4", true},
		{"v1.2.0.0", "v1.2", true},
		{"v1.2.0.0.1", "v1.2", false},
		{"v1.2", "v1.2.0.0", true},
		{"v1.2", "v1.2.0.0.1", true},
		{"v1.2.0.0", "v1.2.0.0.1", true},
		{"v1.2.3.0", "v1.2.3.4", true},
		{"1.7rc2", "1.7rc1", false},
		{"1.7rc2", "1.7", true},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", false},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.LessThanOrEqual(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <= %s\nexpected: %t\nactual: %t",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestComplexVersion(t *testing.T) {
	goVersion := "devel go1.21-36f564c189"
	v, err := NewGoVersion(goVersion)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	v22, err := NewGoVersion("1.22")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	fmt.Printf("%v\n", v)
	if !v.LessThan(v22) {
		panic("should be less than v22")
	}
}

func TestNewGoVersion(t *testing.T) {
	cases := []struct {
		version string
		err     bool
	}{
		{"", true},
		{"go1.2.3", false},
		{"go1.0", false},
		{"go1", false},
		{"go1.2.beta", true},
		{"go1.21.beta", true},
		{"gofoo", true},
		{"go1.2-5", false},
		{"1.2-beta.5", false},
		{"go\n1.2", true},
		{"\ngo1.2", false},
		{"go1.2.0-x.Y.0+metadata", false},
		{"go1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"go1.2.3-rc1-with-hyphen", false},
		{"go1.2.3.4", false},
		{"go1.2.0.4-x.Y.0+metadata", false},
		{"go1.2.0.4-x.Y.0+metadata-width-hyphen", false},
		{"go1.2.0-X-1.2.0+metadata~dist", false},
		{"go1.2.3.4-rc1-with-hyphen", false},
		{"go1.2.3.4", false},
		{"gov1.2.3", false},
		{"gofoo1.2.3", true},
		{"go1.7rc2", false},
		{"gov1.7rc2", false},
		{"go1.0-", false},
	}

	for _, tc := range cases {
		_, err := NewGoVersion(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %q", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %q: %s", tc.version, err)
		}
	}
}
