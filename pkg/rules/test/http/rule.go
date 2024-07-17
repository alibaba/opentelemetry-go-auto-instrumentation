package http

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {

	api.NewRule("net/http", "NewRequestWithContext", "", "onEnterNewRequestWithContext2", "").
		WithRuleName("testrule").
		Register()

	api.NewRule("net/http", "Do", "*Client", "onEnterClientDo2", "").
		WithRuleName("testrule").
		Register()
}
