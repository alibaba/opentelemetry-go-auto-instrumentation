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
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"
)

type K8sClientGoInnerEnabler struct {
	enabled bool
}

var k8sEnabler = K8sClientGoInnerEnabler{
	enabled: os.Getenv("OTEL_INSTRUMENTATION_K8S_CLIENT_GO_ENABLED") != "false",
}

func (k K8sClientGoInnerEnabler) Enable() bool {
	return k.enabled
}

type k8sEventInfo struct {
	eventType       string
	eventUID        string
	namespace       string
	name            string
	resourceVersion string
	apiVersion      string
	kind            string
	startTime       time.Time
	processingTime  int64
}

type k8sEventsInfo struct {
	isInInitialList bool
	eventCount      int
}

type K8sObjectMeta interface {
	GetName() string
	GetNamespace() string
	GetUID() string
	GetResourceVersion() string
}

type GVK interface {
	GroupVersionKind() (group, version, kind string)
}

type ObjectMetaAccessor interface {
	GetName() string
	GetNamespace() string
	GetUID() string
	GetResourceVersion() string
}

type ObjectKind interface {
	GroupVersionKind() GroupVersionKind
}

type GroupVersionKind struct {
	Group   string
	Version string
	Kind    string
}

// metaAccessor provides the same functionality as k8s.io/apimachinery/pkg/api/meta.Accessor
func metaAccessor(obj interface{}) (ObjectMetaAccessor, bool) {
	if obj == nil {
		return nil, false
	}

	if accessor, ok := obj.(ObjectMetaAccessor); ok {
		return accessor, true
	}

	return getObjectMetaViaReflection(obj)
}

// objectKindAccessor provides the same functionality as schema.ObjectKind
func objectKindAccessor(obj interface{}) (ObjectKind, bool) {
	if obj == nil {
		return nil, false
	}

	if kind, ok := obj.(ObjectKind); ok {
		return kind, true
	}

	return getObjectKindViaReflection(obj)
}

// Reflection-based helpers for accessing Kubernetes object fields
func getObjectMetaViaReflection(obj interface{}) (ObjectMetaAccessor, bool) {
	if obj == nil {
		return nil, false
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.IsValid() {
		return nil, false
	}

	// Look for ObjectMeta field
	if metaField := val.FieldByName("ObjectMeta"); metaField.IsValid() {
		return &objectMetaWrapper{value: metaField}, true
	}

	// Look for nested ObjectMeta in TypeMeta/ObjectMeta pattern
	if typeMetaField := val.FieldByName("TypeMeta"); typeMetaField.IsValid() {
		if metaField := val.FieldByName("ObjectMeta"); metaField.IsValid() {
			return &objectMetaWrapper{value: metaField}, true
		}
	}

	return nil, false
}

func getObjectKindViaReflection(obj interface{}) (ObjectKind, bool) {
	if obj == nil {
		return nil, false
	}

	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if !val.IsValid() {
		return nil, false
	}

	// Look for TypeMeta field which contains GVK information
	if typeMetaField := val.FieldByName("TypeMeta"); typeMetaField.IsValid() {
		return &objectKindWrapper{value: typeMetaField}, true
	}

	return nil, false
}

// Wrapper types to provide interface implementations

type objectMetaWrapper struct {
	value reflect.Value
}

func (w *objectMetaWrapper) GetName() string {
	if field := w.value.FieldByName("Name"); field.IsValid() && field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

func (w *objectMetaWrapper) GetNamespace() string {
	if field := w.value.FieldByName("Namespace"); field.IsValid() && field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

func (w *objectMetaWrapper) GetUID() string {
	if field := w.value.FieldByName("UID"); field.IsValid() {
		if field.Kind() == reflect.String {
			return field.String()
		}
		// Handle types.k8s.io/apimachinery/pkg/types.UID
		return fmt.Sprintf("%v", field.Interface())
	}
	return ""
}

func (w *objectMetaWrapper) GetResourceVersion() string {
	if field := w.value.FieldByName("ResourceVersion"); field.IsValid() && field.Kind() == reflect.String {
		return field.String()
	}
	return ""
}

type objectKindWrapper struct {
	value reflect.Value
}

func (w *objectKindWrapper) GroupVersionKind() GroupVersionKind {
	gvk := GroupVersionKind{}

	if apiVersion := w.value.FieldByName("APIVersion"); apiVersion.IsValid() && apiVersion.Kind() == reflect.String {
		version := apiVersion.String()
		if parts := strings.Split(version, "/"); len(parts) == 2 {
			gvk.Group = parts[0]
			gvk.Version = parts[1]
		} else {
			gvk.Version = version
		}
	}

	if kind := w.value.FieldByName("Kind"); kind.IsValid() && kind.Kind() == reflect.String {
		gvk.Kind = kind.String()
	}

	return gvk
}
