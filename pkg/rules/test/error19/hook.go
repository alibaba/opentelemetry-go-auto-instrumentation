// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package error19

import (
	_ "unsafe"

	"github.com/alibaba/loongsuite-go-agent/pkg/api"
)

// onEnterGeneric is a generic hook function that can be used to handle all
// entry points of a function. The only parameter is a CallContext, which
// contains all the information about the function call.
//
//go:linkname onEnterGeneric errorstest/all.onEnterGeneric
func onEnterGeneric(call api.CallContext) {
	println("xian")
}

//go:linkname onEnterGeneric2 errorstest/all.onEnterGeneric2
func onEnterGeneric2(call api.CallContext) {
	println("shanxi")
}

//go:linkname onEnterGeneric3 errorstest/all.onEnterGeneric3
func onEnterGeneric3(call api.CallContext) {
	println("zhejiang")
}

//go:linkname onEnterGeneric4 errorstest/all.onEnterGeneric4
func onEnterGeneric4(call api.CallContext) {
	println("beijing")
}

//go:linkname onEnterGeneric5 errorstest/all.onEnterGeneric5
func onEnterGeneric5(call api.CallContext) {
	println("entering" + call.GetFuncName())
	println("within" + call.GetPackageName())
}
