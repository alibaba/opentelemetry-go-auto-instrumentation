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

package service_holder

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/experimental"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client/naming_cache"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"log"
	"os"
	"strconv"
)

type nacosEnabler struct{}

func (n nacosEnabler) Enable() bool {
	return os.Getenv("OTEL_INSTRUMENTATION_NACOS_EXPERIMENTAL_ENABLE") == "true"
}

var enabler nacosEnabler

func beforeNewServiceInfoHolder(call api.CallContext, namespace, cacheDir string, updateCacheWhenEmpty, notLoadCacheAtStart bool) {
	if !enabler.Enable() {
		return
	}
	call.SetKeyData("namespace", namespace)
	call.SetKeyData("cacheDir", cacheDir)
	call.SetKeyData("updateCacheWhenEmpty", strconv.FormatBool(updateCacheWhenEmpty))
	call.SetKeyData("notLoadCacheAtStart", strconv.FormatBool(notLoadCacheAtStart))
}

func afterNewServiceInfoHolder(call api.CallContext, holder *naming_cache.ServiceInfoHolder) {
	if !enabler.Enable() {
		return
	}
	reg, err := experimental.GlobalMeter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		attrSet := attribute.NewSet(attribute.KeyValue{
			Key:   "namespace",
			Value: attribute.StringValue(call.GetKeyData("namespace").(string)),
		}, attribute.KeyValue{
			Key:   "cache.dir",
			Value: attribute.StringValue(call.GetKeyData("cacheDir").(string)),
		}, attribute.KeyValue{
			Key:   "update.cache.when.empty",
			Value: attribute.StringValue(call.GetKeyData("updateCacheWhenEmpty").(string)),
		}, attribute.KeyValue{
			Key:   "not.load.cache.at.start",
			Value: attribute.StringValue(call.GetKeyData("notLoadCacheAtStart").(string)),
		})
		observer.ObserveInt64(experimental.ClientServiceInfoMapSize, int64(holder.ServiceInfoMap.Count()), metric.WithAttributeSet(attrSet))
		return nil
	}, experimental.ClientServiceInfoMapSize)
	if err != nil {
		log.Printf("[otel nacos] failed to register metrics for service info holder")
	} else {
		holder.OtelReg = reg
	}
}
