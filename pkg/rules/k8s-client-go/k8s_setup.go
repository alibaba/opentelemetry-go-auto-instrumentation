// Copyright (c) 2025 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s_client_go

import (
	"context"
	"time"
	_ "unsafe"

	"github.com/alibaba/loongsuite-go-agent/pkg/api"
	"k8s.io/client-go/tools/cache"
)

//go:linkname ProcessDeltasOnEnter k8s.io/client-go/tools/cache.ProcessDeltasOnEnter
func ProcessDeltasOnEnter(call api.CallContext, handler cache.ResourceEventHandler, clientState cache.Store, deltas cache.Deltas, isInInitialList bool) {
	if !k8sEnabler.Enable() {
		return
	}
	eventsInfo := k8sEventsInfo{
		isInInitialList: isInInitialList,
		eventCount:      len(deltas),
	}
	ctx := k8sClientGoEventsInstrumenter.Start(context.TODO(), eventsInfo)

	for _, d := range deltas {
		obj := d.Object

		eventInfo := k8sEventInfo{}
		eventInfo.startTime = time.Now()
		eventInfo.eventType = string(d.Type)
		if m, ok := metaAccessor(obj); ok {
			eventInfo.name = m.GetName()
			eventInfo.namespace = m.GetNamespace()
			eventInfo.eventUID = m.GetUID()
			eventInfo.resourceVersion = m.GetResourceVersion()
		}
		if objWithGVK, ok := objectKindAccessor(obj); ok {
			gvk := objWithGVK.GroupVersionKind()
			eventInfo.apiVersion = gvk.Group + "/" + gvk.Version
			eventInfo.kind = gvk.Kind
		}

		ctx = k8sClientGoEventInstrumenter.Start(ctx, eventInfo)

		switch d.Type {
		case cache.Sync, cache.Replaced, cache.Added, cache.Updated:
			if old, exists, err := clientState.Get(obj); err == nil && exists {
				if err := clientState.Update(obj); err != nil {
					call.SetData(err)
					return
				}
				handler.OnUpdate(old, obj)
			} else {
				if err := clientState.Add(obj); err != nil {
					call.SetData(err)
					return
				}
				handler.OnAdd(obj, isInInitialList)
			}
		case cache.Deleted:
			if err := clientState.Delete(obj); err != nil {
				call.SetData(err)
				return
			}
			handler.OnDelete(obj)
		}

		duration := time.Since(eventInfo.startTime)
		eventInfo.processingTime = duration.Microseconds()
		k8sClientGoEventInstrumenter.End(ctx, eventInfo, eventInfo, nil)
	}

	k8sClientGoEventsInstrumenter.End(ctx, eventsInfo, eventsInfo, nil)

	call.SetData(nil)
	call.SetSkipCall(true)
}

//go:linkname ProcessDeltasOnExit k8s.io/client-go/tools/cache.ProcessDeltasOnExit
func ProcessDeltasOnExit(call api.CallContext, err error) {
	if !k8sEnabler.Enable() {
		return
	}
	e := call.GetData()
	if e != nil {
		err = e.(error)
		call.SetReturnVal(0, err)
	}
}
