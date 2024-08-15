module logrusTest/v1.5.0

go 1.18

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../otel-go-auto-instrumentation

replace github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.5.0

require github.com/sirupsen/logrus v1.5.0

require (
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	golang.org/x/sys v0.0.0-20190422165155-953cdadca894 // indirect
)
