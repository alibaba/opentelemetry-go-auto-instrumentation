// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package databasesql

import "github.com/alibaba/opentelemetry-go-auto-instrumentation/api"

func init() {
	api.NewStructRule("database/sql", "DB", "Endpoint", "string").Register()
	api.NewStructRule("database/sql", "DB", "DriverName", "string").Register()
	api.NewStructRule("database/sql", "DB", "DSN", "string").Register()
	api.NewStructRule("database/sql", "Stmt", "Data", "map[string]string").Register()
	api.NewStructRule("database/sql", "Stmt", "DriverName", "string").Register()
	api.NewStructRule("database/sql", "Stmt", "DSN", "string").Register()
	api.NewStructRule("database/sql", "Tx", "Endpoint", "string").Register()
	api.NewStructRule("database/sql", "Tx", "DriverName", "string").Register()
	api.NewStructRule("database/sql", "Tx", "DSN", "string").Register()
	api.NewStructRule("database/sql", "Conn", "Endpoint", "string").Register()
	api.NewStructRule("database/sql", "Conn", "DriverName", "string").Register()
	api.NewStructRule("database/sql", "Conn", "DSN", "string").Register()

	api.NewRule("database/sql", "Open", "", "beforeOpenInstrumentation", "afterOpenInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "PingContext", "*DB", "beforePingContextInstrumentation", "afterPingContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "PrepareContext", "*DB", "beforePrepareContextInstrumentation", "afterPrepareContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "ExecContext", "*DB", "beforeExecContextInstrumentation", "afterExecContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "QueryContext", "*DB", "beforeQueryContextInstrumentation", "afterQueryContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "BeginTx", "*DB", "beforeTxInstrumentation", "afterTxInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "Conn", "*DB", "beforeConnInstrumentation", "afterConnInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()

	api.NewRule("database/sql", "PingContext", "*Conn", "beforeConnPingContextInstrumentation", "afterConnPingContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "PrepareContext", "*Conn", "beforeConnPrepareContextInstrumentation", "afterConnPrepareContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "ExecContext", "*Conn", "beforeConnExecContextInstrumentation", "afterConnExecContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "QueryContext", "*Conn", "beforeConnQueryContextInstrumentation", "afterConnQueryContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "BeginTx", "*Conn", "beforeConnTxInstrumentation", "afterConnTxInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()

	api.NewRule("database/sql", "StmtContext", "*Tx", "beforeTxStmtContextInstrumentation", "afterTxStmtContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "PrepareContext", "*Tx", "beforeTxPrepareContextInstrumentation", "afterTxPrepareContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "ExecContext", "*Tx", "beforeTxExecContextInstrumentation", "afterTxExecContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "QueryContext", "*Tx", "beforeTxQueryContextInstrumentation", "afterTxQueryContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "Commit", "*Tx", "beforeTxCommitInstrumentation", "afterTxCommitInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "Rollback", "*Tx", "beforeTxRollbackInstrumentation", "afterTxRollbackInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()

	api.NewRule("database/sql", "ExecContext", "*Stmt", "beforeStmtExecContextInstrumentation", "afterStmtExecContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()
	api.NewRule("database/sql", "QueryContext", "*Stmt", "beforeStmtQueryContextInstrumentation", "afterStmtQueryContextInstrumentation").WithFileDeps("databasesql_data_type.go", "databasesql_otel_instrumenter.go", "databasesql_parser.go").Register()

}
