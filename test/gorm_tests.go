package test

import (
	"testing"
)

const gorm_dependency_name = "gorm.io/gorm"
const gorm_module_name = "gorm"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("gorm_crud_test", gorm_module_name, "v1.23.0", "v1.24.6", "1.18", "", TestGormCrud1231),
		NewLatestDepthTestCase("gorm_latestdepth_test", gorm_dependency_name, gorm_module_name, "v1.23.0", "v1.24.6", "1.18", "", TestGormCrud1231),
		NewGeneralTestCase("gorm_crud_test", gorm_module_name, "v1.22.0", "v1.23.0", "1.18", "", TestGormCrud1220))
}

func TestGormCrud1231(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer clearDbSqlContainer(mysqlC)
	UseApp("gorm/v1.23.1")
	RunInstrument(t, "-debuglog", "--", "test_gorm_crud.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	RunApp(t, "test_gorm_crud", env...)
}

func TestGormCrud1220(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer clearDbSqlContainer(mysqlC)
	UseApp("gorm/v1.22.0")
	RunInstrument(t, "-debuglog", "--", "test_gorm_crud.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	RunApp(t, "test_gorm_crud", env...)
}
