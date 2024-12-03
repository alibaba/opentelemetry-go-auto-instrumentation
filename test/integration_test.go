// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"os"
	"testing"

	_ "github.com/alibaba/opentelemetry-go-auto-instrumentation/test/version"
)

const test_plugin_name_key = "TEST_PLUGIN_NAME"

func TestPlugins(t *testing.T) {
	testPluginName := os.Getenv(test_plugin_name_key)
	for _, c := range TestCases {
		if c == nil {
			continue
		}
		if c.IsMuzzleCheck || c.IsLatestDepthCheck || (testPluginName != "" && c.TestName != testPluginName) {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			c.TestFunc(t)
		})
	}
}

func TestMuzzle(t *testing.T) {
	testPluginName := os.Getenv(test_plugin_name_key)
	for _, c := range TestCases {
		if c == nil {
			continue
		}
		if !c.IsMuzzleCheck || (testPluginName != "" && c.TestName != testPluginName) {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			ExecMuzzle(t, c.DependencyName, c.ModuleName, c.MinVersion, c.MaxVersion, c.MuzzleClasses)
		})
	}
}

func TestLatest(t *testing.T) {
	for _, c := range TestCases {
		testPluginName := os.Getenv(test_plugin_name_key)
		if c == nil {
			continue
		}
		if !c.IsLatestDepthCheck || (testPluginName != "" && c.TestName != testPluginName) {
			continue
		}
		t.Run(c.TestName, func(t *testing.T) {
			ExecLatestTest(t, c.DependencyName, c.ModuleName, c.MinVersion, c.MaxVersion, c.LatestDepthFunc)
		})
	}
}
