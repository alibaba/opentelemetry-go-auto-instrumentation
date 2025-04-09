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

package test

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

func init() {
	TestCases = append(TestCases,
		NewGeneralTestCase("databasesql-mysql-8-access-database-test", "databasesql", "", "", "1.18", "", TestMySql8xAccessDatabase),
		NewGeneralTestCase("databasesql-mysql-8-fetching-database-test", "databasesql", "", "", "1.18", "", TestMySql8xFetchingDatabase),
		NewGeneralTestCase("databasesql-mysql-8-modify-data-test", "databasesql", "", "", "1.18", "", TestMySql8xModifyData),
		NewGeneralTestCase("databasesql-mysql-8-prepared-statement-test", "databasesql", "", "", "1.18", "", TestPreparedStatement),
		NewGeneralTestCase("databasesql-mysql-8-single-row-query-test", "databasesql", "", "", "1.18", "", TestSingleRowQuery),
		NewGeneralTestCase("databasesql-mysql-8-single-transaction-test", "databasesql", "", "", "1.18", "", TestTransaction),

		NewGeneralTestCase("databasesql-mysql-5-access-database-test", "databasesql", "", "", "1.18", "", TestMySql5xAccessDatabase),
		NewGeneralTestCase("databasesql-mysql-5-fetching-database-test", "databasesql", "", "", "1.18", "", TestMySql5xFetchingDatabase),
		NewGeneralTestCase("databasesql-mysql-5-modify-data-test", "databasesql", "", "", "1.18", "", TestMySql5xModifyData),
		NewGeneralTestCase("databasesql-mysql-5-prepared-statement-test", "databasesql", "", "", "1.18", "", TestMySql5xPreparedStatement),
		NewGeneralTestCase("databasesql-mysql-5-single-row-query-test", "databasesql", "", "", "1.18", "", TestMySql5xSingleRowQuery),
		NewGeneralTestCase("databasesql-mysql-5-single-transaction-test", "databasesql", "", "", "1.18", "", TestMySql5xTransaction),
	)
}

func TestMySql5xAccessDatabase(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init5xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_access_database.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_access_database", env...)
}

func TestMySql5xFetchingDatabase(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init5xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_fetching_database.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_fetching_database", env...)
}

func TestMySql5xModifyData(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init5xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_modify_data.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_modify_data", env...)
}

func TestMySql5xSingleRowQuery(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init5xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_single_row_query.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_single_row_query", env...)
}

func TestMySql5xTransaction(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init5xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_transaction.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_transaction", env...)
}

func TestMySql5xPreparedStatement(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init5xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_prepared_statement.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_prepared_statement", env...)
}

func TestMySql8xAccessDatabase(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_access_database.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_access_database", env...)
}

func TestMySql8xFetchingDatabase(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_fetching_database.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_fetching_database", env...)
}

func TestMySql8xModifyData(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_modify_data.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_modify_data", env...)
}

func TestSingleRowQuery(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_single_row_query.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_single_row_query", env...)
}

func TestTransaction(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_transaction.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_transaction", env...)
}

func TestPreparedStatement(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer testcontainers.CleanupContainer(t, mysqlC)
	UseApp("databasesql/mysql")
	RunGoBuild(t, "go", "build", "test_prepared_statement.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	env = append(env, "OTEL_INSTRUMENTATION_DB_EXPERIMENTAL_ENABLE=true")
	RunApp(t, "test_prepared_statement", env...)
}

func init5xMySqlContainer() (testcontainers.Container, nat.Port) {
	ctx := context.Background()
	mysqlContainer, err := mysql.Run(ctx, "mysql:5.6")
	if err != nil {
		panic(err)
	}
	if err := mysqlContainer.Start(ctx); err != nil {
		panic(err)
	}
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}
	return mysqlContainer, port
}

func init8xMySqlContainer() (testcontainers.Container, nat.Port) {
	ctx := context.Background()
	mysqlContainer, err := mysql.Run(ctx, "mysql:8.0.36")
	if err != nil {
		panic(err)
	}
	if err := mysqlContainer.Start(ctx); err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		panic(err)
	}
	return mysqlContainer, port
}
