package gopg

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api-semconv/instrumenter/db"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/instrumenter"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/utils"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/inst-api/version"
	"go.opentelemetry.io/otel/sdk/instrumentation"
)

type gogpAttrsGetter struct{}

func (g gogpAttrsGetter) GetSystem(gopgRequest gopgRequest) string {
	return gopgRequest.System
}

func (g gogpAttrsGetter) GetServerAddress(gopgRequest gopgRequest) string {
	return gopgRequest.Addr
}

func (g gogpAttrsGetter) GetStatement(gopgRequest gopgRequest) string {
	return gopgRequest.Statement
}

func (g gogpAttrsGetter) GetCollection(gopgRequest gopgRequest) string {
	// TBD: We need to implement retrieving the collection later.
	return ""
}

func (g gogpAttrsGetter) GetOperation(gopgRequest gopgRequest) string {
	return string(gopgRequest.QueryOp)
}

func (g gogpAttrsGetter) GetParameters(gopgRequest gopgRequest) []any {
	return nil
}

func (g gogpAttrsGetter) GetDbNamespace(gopgRequest gopgRequest) string {
	return gopgRequest.DbName
}

func (g gogpAttrsGetter) GetBatchSize(gopgRequest gopgRequest) int {
	return 0
}

func BuildGopgInstrumenter() instrumenter.Instrumenter[gopgRequest, interface{}] {
	builder := instrumenter.Builder[gopgRequest, interface{}]{}
	getter := gogpAttrsGetter{}
	return builder.Init().SetSpanNameExtractor(&db.DBSpanNameExtractor[gopgRequest]{Getter: getter}).SetSpanKindExtractor(&instrumenter.AlwaysClientExtractor[gopgRequest]{}).
		AddAttributesExtractor(&db.DbClientAttrsExtractor[gopgRequest, any, gogpAttrsGetter]{Base: db.DbClientCommonAttrsExtractor[gopgRequest, any, gogpAttrsGetter]{Getter: getter}}).
		SetInstrumentationScope(instrumentation.Scope{
			Name:    utils.GOPG_SCOPE_NAME,
			Version: version.Tag,
		}).
		BuildInstrumenter()
}
