package test

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/version"
	"log"
	"runtime"
	"strings"
	"testing"
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
	c.IsMuzzleCheck = true
	c.DependencyName = dependencyName
	c.MuzzleClasses = muzzleClasses
	return c
}

func NewLatestDepthTestCase(testName, dependencyName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion string, latestTestFunc func(t *testing.T, env ...string)) *TestCase {
	c := NewGeneralTestCase(testName, moduleName, minVersion, maxVersion, minGoVersion, maxGoVersion, nil)
	c.LatestDepthFunc = latestTestFunc
	c.DependencyName = dependencyName
	c.IsLatestDepthCheck = true
	return c
}
