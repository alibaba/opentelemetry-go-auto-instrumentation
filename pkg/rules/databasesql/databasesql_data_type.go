// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package databasesql

type databaseSqlRequest struct {
	opType     string
	sql        string
	endpoint   string
	driverName string
	dsn        string
	params     []any
}
