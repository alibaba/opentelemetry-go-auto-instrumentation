//go:build ignore

package mongo

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoInstrumenter = BuildMongoOtelInstrumenter()

func mongoOnEnter(call mongo.CallContext, opts ...*options.ClientOptions) {
	syncMap := sync.Map{}
	for _, opt := range opts {
		hosts := opt.Hosts
		hostLength := len(hosts)
		if hostLength == 0 {
			continue
		}
		configuredMonitor := opt.Monitor
		opt.Monitor = &event.CommandMonitor{
			Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
				if configuredMonitor != nil {
					configuredMonitor.Started(ctx, startedEvent)
				}
				mongoRequest := mongoRequest{
					CommandName:  startedEvent.CommandName,
					ConnectionID: startedEvent.ConnectionID,
					DatabaseName: startedEvent.DatabaseName,
				}
				newCtx := mongoInstrumenter.Start(ctx, mongoRequest)
				syncMap.Store(fmt.Sprintf("%d", startedEvent.RequestID), newCtx)
			},
			Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {
				if configuredMonitor != nil {
					configuredMonitor.Succeeded(ctx, succeededEvent)
				}
				if newCtx, ok := syncMap.LoadAndDelete(fmt.Sprintf("%d", succeededEvent.RequestID)); ok && newCtx != nil {
					newContext, ok := newCtx.(context.Context)
					if ok {
						mongoInstrumenter.End(newContext, mongoRequest{
							CommandName:  succeededEvent.CommandName,
							ConnectionID: succeededEvent.ConnectionID,
						}, nil, nil)
					}
				}
			},
			Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {
				if configuredMonitor != nil {
					configuredMonitor.Failed(ctx, failedEvent)
				}
				if newCtx, ok := syncMap.LoadAndDelete(fmt.Sprintf("%d", failedEvent.RequestID)); ok && newCtx != nil {
					newContext, ok := newCtx.(context.Context)
					if ok {
						mongoInstrumenter.End(newContext, mongoRequest{
							CommandName:  failedEvent.CommandName,
							ConnectionID: failedEvent.ConnectionID,
						}, nil, errors.New(failedEvent.Failure))
					}
				}
			},
		}
	}
}
