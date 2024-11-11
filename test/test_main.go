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

package test

import (
	"log"
	"runtime"
	"strings"
	"testing"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/version"
)

type TestCase struct {
	TestName           string
	DependencyName     string
	ModuleName         string
	MinVersion         *version.Version
	MaxVersion         *version.Version
	MinGoVersion       *version.Version
	MaxGoVersion       *version.Version
	TestFunc           func(t *testing.T, env ...string)
	LatestDepthFunc    func(t *testing.T, env ...string)
	MuzzleClasses      []string
	IsMuzzleCheck      bool
	IsLatestDepthCheck bool
}

var TestCases = make([]*TestCase, 0)

func NewGeneralTestCase(testName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion string, testFunc func(t *testing.T, env ...string)) *TestCase {
	minVer, err := version.NewVersion(minVersion)
	if minVersion != "" && err != nil {
		log.Printf("Error parsing min version: %v", err)
	}
	maxVer, err := version.NewVersion(maxVersion)
	if maxVersion != "" && err != nil {
		log.Printf("Error parsing max version: %v", err)
	}
	minGoVer, err := version.NewGoVersion(minGoVersion)
	if minGoVersion != "" && err != nil {
		log.Printf("Error parsing min go version: %v", err)
	}
	maxGoVer, err := version.NewGoVersion(maxGoVersion)
	if maxGoVersion != "" && err != nil {
		log.Printf("Error parsing max go version: %v", err)
	}
	goVersion, _ := version.NewGoVersion(strings.ReplaceAll(runtime.Version(), "go", ""))
	if (minGoVer != nil && goVersion.LessThan(minGoVer)) || (maxGoVer != nil && goVersion.GreaterThan(maxGoVer)) {
		log.Printf("This test does not suppport go " + goVersion.String())
		return nil
	}
	return &TestCase{
		TestName:           testName,
		ModuleName:         moduleName,
		MinVersion:         minVer,
		MaxVersion:         maxVer,
		MinGoVersion:       minGoVer,
		MaxGoVersion:       maxGoVer,
		TestFunc:           testFunc,
		IsMuzzleCheck:      false,
		IsLatestDepthCheck: false,
	}
}

func NewMuzzleTestCase(testName, dependencyName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion string, muzzleClasses []string) *TestCase {
	c := NewGeneralTestCase(testName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion, nil)
	if c == nil {
		return nil
	}
	c.IsMuzzleCheck = true
	c.DependencyName = dependencyName
	c.MuzzleClasses = muzzleClasses
	return c
}

func NewLatestDepthTestCase(testName, dependencyName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion string, latestTestFunc func(t *testing.T, env ...string)) *TestCase {
	c := NewGeneralTestCase(testName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion, nil)
	if c == nil {
		return nil
	}
	c.LatestDepthFunc = latestTestFunc
	c.DependencyName = dependencyName
	c.IsLatestDepthCheck = true
	return c
}
