package resource

import (
	"embed"
	"strings"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/tool/util"

	"github.com/alibaba/opentelemetry-go-auto-instrumentation/api"
)

func listFiles(fs embed.FS, dir string) ([]string, error) {
	list, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, item := range list {
		path := dir + "/" + item.Name()
		if item.IsDir() {
			subFiles, err := listFiles(fs, path)
			if err != nil {
				return nil, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, path)
		}
	}
	return files, nil
}

func CopyAPITo(target string, pkgName string) (string, error) {
	apiSnippet := strings.Replace(api.ExportAPITemplate(), "package api", "package "+pkgName, 1)
	return util.WriteStringToFile(target, apiSnippet)
}
