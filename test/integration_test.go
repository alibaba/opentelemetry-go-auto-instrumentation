package test

import (
	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/test/version"
	"testing"
)

func TestPlugins(t *testing.T) {
	for _, c := range TestCases {
		if c == nil {
			t.Skip()
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
			t.Skip()
		}
		if !c.IsMuzzleCheck {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			ExecMuzzle(t, c.DependencyName, c.ModuleName, c.MinVersion, c.MaxVersion)
		})
	}
}

func TestLatest(t *testing.T) {
	for _, c := range TestCases {
		if c == nil {
			t.Skip()
		}
		if !c.IsLatestDepthCheck {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			ExecLatestTest(t, c.DependencyName, c.ModuleName, c.MinVersion, c.MaxVersion, c.LatestDepthFunc)
		})
	}
}
