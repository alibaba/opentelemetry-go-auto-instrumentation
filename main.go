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

package main

import (
	"log"
	"os"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool"
)

func main() {
	err := tool.Init()
	if err != nil {
		log.Printf("failed to init options: %v", err)
		os.Exit(1)

	}
	err = tool.Run()
	if err != nil {
		log.Printf("failed to run the tool: %v", err)
		os.Exit(1)
	}
}
