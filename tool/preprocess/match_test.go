package preprocess

import (
	"os"
	"testing"

	"github.com/alibaba/loongsuite-go-agent/tool/config"
	"github.com/alibaba/loongsuite-go-agent/tool/rules"
)

func init() {
	// Create a minimal config file for testing
	configFile := "/tmp/test_config.json"
	os.WriteFile(configFile, []byte(`{"Verbose": false}`), 0644)
	os.Setenv("OTELTOOL_CONF", configFile)

	config.InitConfig()
}

// TestMatchDependencies tests the matchDependencies function
func TestMatchDependencies(t *testing.T) {
	tests := []struct {
		name        string
		rule        rules.InstRule
		projectDeps map[string]bool
		want        bool
	}{
		{
			name: "non-InstFuncRule should return true",
			rule: &rules.InstStructRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/struct"},
				StructType:   "TestStruct",
			},
			projectDeps: map[string]bool{},
			want:        true,
		},
		{
			name: "InstFuncRule with no dependencies should return true",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
			},
			projectDeps: map[string]bool{},
			want:        true,
		},
		{
			name: "InstFuncRule with all dependencies present should return true",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
				Dependencies: []string{"github.com/gin-gonic/gin", "github.com/jinzhu/gorm"},
			},
			projectDeps: map[string]bool{
				"github.com/gin-gonic/gin": true,
				"github.com/jinzhu/gorm":   true,
			},
			want: true,
		},
		{
			name: "InstFuncRule with missing dependency should return false",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
				Dependencies: []string{"github.com/gin-gonic/gin", "github.com/missing/dep"},
			},
			projectDeps: map[string]bool{
				"github.com/gin-gonic/gin": true,
			},
			want: false,
		},
		{
			name: "InstFuncRule with single dependency present should return true",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
				Dependencies: []string{"github.com/gin-gonic/gin"},
			},
			projectDeps: map[string]bool{
				"github.com/gin-gonic/gin": true,
			},
			want: true,
		},
		{
			name: "InstFuncRule with single dependency missing should return false",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
				Dependencies: []string{"github.com/missing/dep"},
			},
			projectDeps: map[string]bool{},
			want:        false,
		},
		{
			name: "InstFuncRule with empty dependencies should return true",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
				Dependencies: []string{},
			},
			projectDeps: map[string]bool{},
			want:        true,
		},
		{
			name: "InstFuncRule with nil dependencies should return true",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/func"},
				Function:     "TestFunc",
				Dependencies: nil,
			},
			projectDeps: map[string]bool{},
			want:        true,
		},
		{
			name: "InstFuncRule with complex dependencies should return true when all present",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/complex"},
				Function:     "ComplexFunc",
				Dependencies: []string{
					"github.com/gin-gonic/gin",
					"github.com/jinzhu/gorm",
					"github.com/sirupsen/logrus",
					"github.com/stretchr/testify",
				},
			},
			projectDeps: map[string]bool{
				"github.com/gin-gonic/gin":    true,
				"github.com/jinzhu/gorm":      true,
				"github.com/sirupsen/logrus":  true,
				"github.com/stretchr/testify": true,
				"github.com/extra/package":    true, // extra packages don't affect result
			},
			want: true,
		},
		{
			name: "InstFuncRule with complex dependencies should return false when one missing",
			rule: &rules.InstFuncRule{
				InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/complex"},
				Function:     "ComplexFunc",
				Dependencies: []string{
					"github.com/gin-gonic/gin",
					"github.com/jinzhu/gorm",
					"github.com/missing/package",
					"github.com/sirupsen/logrus",
				},
			},
			projectDeps: map[string]bool{
				"github.com/gin-gonic/gin":   true,
				"github.com/jinzhu/gorm":     true,
				"github.com/sirupsen/logrus": true,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := &ruleMatcher{
				projectDeps: tt.projectDeps,
			}
			if got := rm.matchDependencies(tt.rule); got != tt.want {
				t.Errorf("matchDependencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestRuleMatcherCreation tests that ruleMatcher can be created and used
func TestRuleMatcherCreation(t *testing.T) {
	// This test ensures the ruleMatcher struct is properly defined
	rm := &ruleMatcher{
		availableRules: make(map[string][]rules.InstRule),
		projectDeps:    make(map[string]bool),
	}

	// Test with empty rule
	rule := &rules.InstFuncRule{
		InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/empty"},
		Function:     "EmptyFunc",
	}

	result := rm.matchDependencies(rule)
	if !result {
		t.Errorf("Expected empty rule to match, got %v", result)
	}
}

// TestMatchDependenciesEdgeCases tests edge cases for matchDependencies
func TestMatchDependenciesEdgeCases(t *testing.T) {
	rm := &ruleMatcher{
		projectDeps: map[string]bool{
			"": true, // empty string as dependency
		},
	}

	// Test with empty string dependency
	rule := &rules.InstFuncRule{
		InstBaseRule: rules.InstBaseRule{ImportPath: "test.com/edge"},
		Function:     "EdgeFunc",
		Dependencies: []string{""},
	}

	result := rm.matchDependencies(rule)
	if !result {
		t.Errorf("Expected empty string dependency to match, got %v", result)
	}
}
