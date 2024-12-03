// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package gorm

type gormRequest struct {
	DbName    string
	Endpoint  string
	Operation string
	User      string
	System    string
}
