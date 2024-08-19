package rule

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	// record dbinfo
	api.NewStructRule("gorm.io/driver/mysql", "Dialector", "DbInfo", "interface{}").
		Register()
	// add callback
	api.NewRule("gorm.io/gorm", "Open", "", "", "afterGormOpen").
		WithVersion("[1.22.0,1.25.10)").
		WithFileDeps("gorm_data_type.go", "gorm_otel_instrumenter.go").
		Register()

}
