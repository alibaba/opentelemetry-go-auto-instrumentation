// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package service

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService)
