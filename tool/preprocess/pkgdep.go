package preprocess

import (
	"fmt"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/resource"
)

const (
	OtelPkgDepsDir = "otel_pkgdep"
)

func (dp *DepProcessor) copyPkgDep() error {
	dir := OtelPkgDepsDir
	err := resource.CopyPkgTo(dir)
	dp.addGeneratedDep(dir)
	if err != nil {
		return fmt.Errorf("failed to copy pkg deps: %w", err)
	}
	return nil
}
