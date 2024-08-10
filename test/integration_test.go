package test

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/test/version"
	"testing"
)

func TestPlugins(t *testing.T) {
	for _, c := range TestCases {
		if c == nil {
			continue
		}
		if c.IsMuzzleCheck || c.IsLatestDepthCheck {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			c.TestFunc(t)
		})
	}
}

func TestMuzzle(t *testing.T) {
	for _, c := range TestCases {
		if c == nil {
			continue
		}
		if !c.IsMuzzleCheck {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			ExecMuzzle(t, c.DependencyName, c.ModuleName, c.MinVersion, c.MaxVersion, c.MuzzleMainClass)
		})
	}
}

func TestLatest(t *testing.T) {
	for _, c := range TestCases {
		if c == nil {
			continue
		}
		if !c.IsLatestDepthCheck {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			ExecLatestTest(t, c.DependencyName, c.ModuleName, c.MinVersion, c.MaxVersion, c.LatestDepthFunc)
		})
	}
}
