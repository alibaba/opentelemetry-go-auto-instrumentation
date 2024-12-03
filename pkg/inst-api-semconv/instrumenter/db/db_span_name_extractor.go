// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package db

type DBSpanNameExtractor[REQUEST any] struct {
	Getter DbClientAttrsGetter[REQUEST]
}

func (d *DBSpanNameExtractor[REQUEST]) Extract(request REQUEST) string {
	operation := d.Getter.GetOperation(request)
	if operation == "" {
		return "DB Query"
	}
	return operation
}
