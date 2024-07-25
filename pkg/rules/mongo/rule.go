package mongo

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	//client
	api.NewRule("go.mongodb.org/mongo-driver/mongo", "NewClient", "", "mongoOnEnter", "").
		WithVersion("[1.11.1,1.15.2)").
		WithFileDeps("mongo_otel_instrumenter.go").
		WithFileDeps("mongo_data_type.go").
		Register()
}
