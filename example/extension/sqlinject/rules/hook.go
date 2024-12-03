// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package rules

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
)

func checkSqlInjection(query string) error {
	patterns := []string{"go", "build", ";", "/*", " or ", " and ", "'"}
	for _, pattern := range patterns {
		if strings.Contains(strings.ToLower(query), pattern) {
			return errors.New("potential SQL injection detected")
		}
	}
	return nil
}

func sqlQueryOnEnter(call api.CallContext, db *sql.DB, query string, args ...interface{}) {
	if err := checkSqlInjection(query); err != nil {
		log.Fatalf("sqlQueryOnEnter %v", err)
	}
}
