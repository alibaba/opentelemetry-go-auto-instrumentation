// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      rpc://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"go.opentelemetry.io/otel/attribute"
)

func Shadow(attrs []attribute.KeyValue, metricsSemConv map[attribute.Key]bool) (int, []attribute.KeyValue) {
	swap := func(attrs []attribute.KeyValue, i, j int) {
		tmp := attrs[i]
		attrs[i] = attrs[j]
		attrs[j] = tmp
	}
	index := 0
	for i, attr := range attrs {
		if _, ok := metricsSemConv[attr.Key]; ok {
			if index != i {
				swap(attrs, i, index)
			}
			index++
		}
	}
	return index, attrs
}
