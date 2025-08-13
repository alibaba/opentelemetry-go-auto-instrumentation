package main

import (
	"context"
	"fmt"
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

	podAddedCh := make(chan struct{})
	podUpdatedCh := make(chan struct{})
	podDeletedCh := make(chan struct{})

	InjectEventChannels(podAddedCh, podUpdatedCh, podDeletedCh)

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

	select {
	case <-podAddedCh:
		fmt.Println("Pod added event received")
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting for pod Added event")
	}

	select {
	case <-podUpdatedCh:
		fmt.Println("Pod updated event received")
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting for pod Updated event")
	}

	deletePolicy := metav1.DeletePropagationForeground
	err = clientset.CoreV1().Pods("default").Delete(ctx, testPod.Name, metav1.DeleteOptions{
		GracePeriodSeconds: new(int64),
		PropagationPolicy:  &deletePolicy,
	})
	if err != nil {
		log.Fatalf("failed to delete test pod: %v", err)
	}

	select {
	case <-podDeletedCh:
		fmt.Println("Pod deleted event received")
	case <-time.After(10 * time.Second):
		log.Fatal("timeout waiting for pod Deleted event")
	}

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
