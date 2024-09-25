// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
