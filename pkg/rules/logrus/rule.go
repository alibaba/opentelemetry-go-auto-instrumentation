package rule

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func init() {
	api.NewRule("github.com/sirupsen/logrus", "SetFormatter", "", "", "logNewOnExit").
		WithVersion("[1.5.0,1.9.4)").
		Register()

}
