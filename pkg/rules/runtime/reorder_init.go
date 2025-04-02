//go:build ignore

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

package runtime

import (
	"internal/abi"
	"unsafe"
)

// -----------------------------------------------------------------------------
// Initialization Reordering
//
// We hope that our hook initialization (including otel setup, etc.) occurs
// before the user's code (including init functions). This is almost impossible
// under a regular Golang runtime because the Go linker accurately calculates
// the init initialization order and writes it into the binary file, and the
// runtime simply executes them in order, without any magic. To achieve our goal,
// we need to hack the runtime a bit, which is precisely what tools are good at,
// so naturally, we can rely on tools to do this. Next, we will modify the runtime
// code and hijack the inittask initialization function, reordering it before the
// initialization takes place. We will move the inittasks of the user project to
// the end (but before main) and preserve their relative order. Our rationale is
// that other inittasks B and C, which inittask A depends on, will be executed
// first. Thus, as long as the current inittask A is not depended on by others,
// it should be free to move further back.

var counter = 0

// It must contain illegal characters that are not allowed in the package path,
// otherwise we might accidentally target some unlucky one.
const localPrefix = "<REORDER>"

const OtelPkgDir = "otel_pkg"

func isPrefixOf(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func reorderInitTasks(ts []*initTask) {
	// Skip first runtime package initialization. I can't think of any reason
	// we should jump the gun on this.
	if counter != 1 {
		counter++
		return
	}
	// Reorder the init tasks so that the local package is initialized last.
	r := 0
	for i := 0; i < len(ts)-r; i++ {
		t := ts[i]
		// 0 = uninitialized, 1 = in progress, 2 = done
		if t.state != 0 {
			continue
		}
		firstFunc := add(unsafe.Pointer(t), 8)
		f := *(*func())(unsafe.Pointer(&firstFunc))
		pkg := funcpkgpath(findfunc(abi.FuncPCABIInternal(f)))
		if (isPrefixOf(pkg, localPrefix) &&
			!isPrefixOf(pkg, localPrefix+"/"+OtelPkgDir)) ||
			isPrefixOf(pkg, "main") {
			// Move the task to the end but keep the relative order.
			// main should be moved as well, because we want to it to be the
			// former one of the local package.
			ts = append(ts[:i], ts[i+1:]...)
			ts = append(ts, t)
			i--
			r++
		}
	}
}
