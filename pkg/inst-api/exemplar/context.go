// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package exemplar

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"go.opentelemetry.io/otel/trace"
)

// ContextManager manages trace context for exemplar association
type ContextManager struct {
	mu       sync.RWMutex
	contexts map[uint64]context.Context
}

var (
	manager = &ContextManager{
		contexts: make(map[uint64]context.Context),
	}
)

// GetManager returns the global exemplar context manager
func GetManager() *ContextManager {
	return manager
}

// StoreContext associates a goroutine ID with a context containing trace info
func (m *ContextManager) StoreContext(gid uint64, ctx context.Context) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	// Only store if there's an active, sampled trace
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() && span.SpanContext().IsSampled() {
		m.contexts[gid] = ctx
	}
}

// GetContext retrieves the context for the current goroutine
func (m *ContextManager) GetContext(gid uint64) context.Context {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	if ctx, ok := m.contexts[gid]; ok {
		return ctx
	}
	return context.Background()
}

// CleanupContext removes the context for a goroutine
func (m *ContextManager) CleanupContext(gid uint64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.contexts, gid)
}

// GetGoroutineID returns current goroutine ID
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	var gid uint64
	fmt.Sscanf(string(b), "goroutine %d ", &gid)
	return gid
}
