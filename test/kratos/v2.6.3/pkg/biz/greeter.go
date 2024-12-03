// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package biz

import (
	"context"

	v1 "kratos/v2.5.2/pkg/api/helloworld/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model.
type Greeter struct {
	Hello string
}

// GreeterUsecase is a Greeter usecase.
type GreeterUsecase struct {
	log *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewGreeterUsecase(logger log.Logger) *GreeterUsecase {
	return &GreeterUsecase{log: log.NewHelper(logger)}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return &Greeter{"Hello"}, nil
}
