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

package main

import (
	"context"
	"log"
	"time"

	"github.com/alibaba/loongsuite-go-agent/test/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	stopCh := make(chan struct{})
	defer close(stopCh)

	ctx := context.Background()
	clientset := CreateK8sClient(stopCh)

	testPod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "otel-demo-pod",
			Namespace: "default",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "60"},
				},
			},
		},
	}

	_, err := clientset.CoreV1().Pods("default").Create(ctx, testPod, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("failed to create pod: %v", err)
	}

	time.Sleep(10 * time.Second)

	deletePolicy := metav1.DeletePropagationForeground
	err = clientset.CoreV1().Pods("default").Delete(ctx, testPod.Name, metav1.DeleteOptions{
		GracePeriodSeconds: new(int64),
		PropagationPolicy:  &deletePolicy,
	})
	if err != nil {
		log.Fatalf("failed to delete test pod: %v", err)
	}

	time.Sleep(10 * time.Second)

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		if len(stubs) == 0 {
			log.Fatal("No traces found")
		}

		var hasPodAdded, hasPodUpdated, hasPodDelete bool

		for _, stub := range stubs {
			for _, span := range stub {
				if span.Name == "k8s.informer.event.process" {
					var eventType, objectName string
					for _, attr := range span.Attributes {
						switch attr.Key {
						case "k8s.event.type":
							eventType = attr.Value.AsString()
						case "k8s.object.name":
							objectName = attr.Value.AsString()
						}
					}

					if objectName == "otel-demo-pod" {
						switch eventType {
						case "Added":
							hasPodAdded = true
						case "Updated":
							hasPodUpdated = true
						case "Deleted":
							hasPodDelete = true
						}
					}
				}
			}
		}

		if !hasPodAdded {
			log.Fatal("Expected 'Added' event for pod 'otel-demo-pod' not found")
		}
		if !hasPodUpdated {
			log.Fatal("Expected 'Updated' event for pod 'otel-demo-pod' not found")
		}
		if !hasPodDelete {
			log.Fatal("Expected 'Deleted' event for pod 'otel-demo-pod' not found")
		}

		log.Println("All expected events found: Added, Updated, Deleted")
	}, 5)
}
