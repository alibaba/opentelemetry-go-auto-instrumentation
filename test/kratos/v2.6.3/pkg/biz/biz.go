// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(NewGreeterUsecase)
