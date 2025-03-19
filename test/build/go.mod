module build

go 1.22

replace google.golang.org/genproto => google.golang.org/genproto v0.0.0-20240822170219-fc7c04adadcd

replace github.com/alibaba/opentelemetry-go-auto-instrumentation => ../../../opentelemetry-go-auto-instrumentation

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg => ../../../opentelemetry-go-auto-instrumentation/pkg

replace github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier => ../../../opentelemetry-go-auto-instrumentation/test/verifier
