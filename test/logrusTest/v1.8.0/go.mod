module logrusTest/v1.8.0

go 1.18

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../../otel-go-auto-instrumentation

replace github.com/sirupsen/logrus => github.com/sirupsen/logrus v1.8.0

require github.com/sirupsen/logrus v1.8.0

require (
	github.com/magefile/mage v1.10.0 // indirect
	golang.org/x/sys v0.0.0-20191026070338-33540a1f6037 // indirect
)
