package test

import (
	"testing"
)

const gorm_dependency_name = "gorm.io/gorm"
const gorm_module_name = "gorm"

func init() {
	TestCases = append(TestCases, NewGeneralTestCase("gorm_crud_test", gorm_module_name, "v1.22.0", "v1.24.6", "1.18", "", TestGormCrud),
		NewMuzzleTestCase("gorm_muzzle_test", gorm_dependency_name, gorm_module_name, "v1.22.0", "v1.24.6", "1.18", "", []string{"test_gorm_crud.go"}),
		NewLatestDepthTestCase("gorm_latestdepth_test", gorm_dependency_name, gorm_module_name, "v1.22.0", "v1.24.6", "1.18", "", TestGormCrud))
}

func TestGormCrud(t *testing.T, env ...string) {
	mysqlC, mysqlPort := init8xMySqlContainer()
	defer clearDbSqlContainer(mysqlC)
	UseApp("gorm/v1.22.0")
	RunInstrument(t, "-debuglog", "--", "test_gorm_crud.go")
	env = append(env, "MYSQL_PORT="+mysqlPort.Port())
	RunApp(t, "test_gorm_crud", env...)
}
