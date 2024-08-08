//go:build ignore

package databasesql

type databaseSqlRequest struct {
	opType     string
	sql        string
	endpoint   string
	driverName string
	dsn        string
}
