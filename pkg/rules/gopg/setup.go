package gopg

import (
	"context"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/api"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"os"
)

var requestKey = "otel-request"

type gopgInnerEnabler struct {
	enabled bool
}

func (g gopgInnerEnabler) Enable() bool {
	return g.enabled
}

var gopgEnabler = gopgInnerEnabler{
	enabled: os.Getenv("OTEL_INSTRUMENTATION_GOPG_ENABLED") != "false",
}

var gopgInstrumenter = BuildGopgInstrumenter()

type otelQueryHooker struct {
	db *pg.DB
}

type queryOperation interface {
	Operation() orm.QueryOp
}

func (o otelQueryHooker) BeforeQuery(ctx context.Context, event *pg.QueryEvent) (context.Context, error) {
	var query string
	var operation orm.QueryOp

	if v, ok := event.Query.(queryOperation); ok {
		operation = v.Operation()
	}
	if operation == orm.InsertOp {
		if b, err := event.UnformattedQuery(); err == nil {
			query = string(b)
		}
	} else {

		if b, err := event.FormattedQuery(); err == nil {
			// fixme event.UnformattedQuery()
			query = string(b)
		}
	}
	request := gopgRequest{
		QueryOp:   event.Query.(queryOperation).Operation(),
		System:    "postgresql",
		Statement: query,
		Addr:      o.db.Options().Addr,
		User:      o.db.Options().User,
		DbName:    o.db.Options().Database,
	}
	return context.WithValue(gopgInstrumenter.Start(ctx, request), requestKey, request), nil
}

func (o otelQueryHooker) AfterQuery(ctx context.Context, event *pg.QueryEvent) error {
	request := ctx.Value(requestKey).(gopgRequest)
	gopgInstrumenter.End(ctx, request, nil, event.Err)
	return nil
}

//go:linkname afterGopgConnect github.com/go-pg/pg/v10/pg.afterGopgConnect
func afterGopgConnect(_ api.CallContext, db *pg.DB) {
	if !gopgEnabler.Enable() {
		return
	}
	if db == nil {
		return
	}
	db.AddQueryHook(&otelQueryHooker{db: db})
}
