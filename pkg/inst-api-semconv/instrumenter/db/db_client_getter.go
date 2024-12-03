// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package db

type DbClientCommonAttrsGetter[REQUEST any] interface {
	GetSystem(REQUEST) string
	GetServerAddress(REQUEST) string
}

type DbClientAttrsGetter[REQUEST any] interface {
	DbClientCommonAttrsGetter[REQUEST]
	GetStatement(REQUEST) string
	GetOperation(REQUEST) string
	GetParameters(REQUEST) []any
}

type SqlClientAttributesGetter[REQUEST any] interface {
	DbClientCommonAttrsGetter[REQUEST]
	GetRawStatement(REQUEST) string
}
