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
//go:build ignore

package databasesql

import (
	"context"
	"database/sql"
	"log"
	"strings"
)

var databaseSqlInstrumenter = BuildDatabaseSqlOtelInstrumenter()

func beforeOpenInstrumentation(call sql.CallContext, driverName, dataSourceName string) {
	addr, err := parseDSN(driverName, dataSourceName)
	if err != nil {
		log.Printf("failed to parse dsn: %v", err)
	}
	call.SetData(map[string]string{
		"endpoint": addr,
		"driver":   driverName,
		"dsn":      dataSourceName,
	})
}

func afterOpenInstrumentation(call sql.CallContext, db *sql.DB, err error) {
	if db == nil {
		return
	}
	data, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	endpoint, ok := data["endpoint"]
	if ok {
		db.Endpoint = endpoint
	}
	driver, ok := data["driver"]
	if ok {
		db.DriverName = driver
	}
	dsn, ok := data["dsn"]
	if ok {
		db.DSN = dsn
	}
}

func beforePingContextInstrumentation(call sql.CallContext, db *sql.DB, ctx context.Context) {
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "ping", "ping", db.Endpoint, db.DriverName, db.DSN)
}

func afterPingContextInstrumentation(call sql.CallContext, err error) {
	instrumentEnd(call, err)
}

func beforePrepareContextInstrumentation(call sql.CallContext, db *sql.DB, ctx context.Context, query string) {
	if db == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": db.Endpoint,
		"sql":      query,
		"driver":   db.DriverName,
		"dsn":      db.DSN,
	})
}

func afterPrepareContextInstrumentation(call sql.CallContext, stmt *sql.Stmt, err error) {
	if stmt == nil {
		return
	}
	callDataMap, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{
		"endpoint": callDataMap["endpoint"],
		"sql":      callDataMap["sql"],
		"driver":   callDataMap["driver"],
	}
	stmt.DSN = callDataMap["dsn"]
}

func beforeExecContextInstrumentation(call sql.CallContext, db *sql.DB, ctx context.Context, query string, args ...any) {
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "exec", query, db.Endpoint, db.DriverName, db.DSN)
}

func afterExecContextInstrumentation(call sql.CallContext, result sql.Result, err error) {
	instrumentEnd(call, err)
}

func beforeQueryContextInstrumentation(call sql.CallContext, db *sql.DB, ctx context.Context, query string, args ...any) {
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "query", query, db.Endpoint, db.DriverName, db.DSN)
}

func afterQueryContextInstrumentation(call sql.CallContext, rows *sql.Rows, err error) {
	instrumentEnd(call, err)
}

func beforeTxInstrumentation(call sql.CallContext, db *sql.DB, ctx context.Context, opts *sql.TxOptions) {
	if db == nil {
		return
	}
	instrumentStart(call, ctx, "begin", "START TRANSACTION", db.Endpoint, db.DriverName, db.DSN)
}

func afterTxInstrumentation(call sql.CallContext, tx *sql.Tx, err error) {
	if tx == nil {
		return
	}
	callData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	dbRequest, ok := callData["dbRequest"].(databaseSqlRequest)
	if !ok {
		return
	}
	tx.Endpoint = dbRequest.endpoint
	tx.DriverName = dbRequest.driverName
	tx.DSN = dbRequest.dsn
	instrumentEnd(call, err)
}

func beforeConnInstrumentation(call sql.CallContext, db *sql.DB, ctx context.Context) {
	if db == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": db.Endpoint,
		"driver":   db.DriverName,
		"dsn":      db.DSN,
	})
}

func afterConnInstrumentation(call sql.CallContext, conn *sql.Conn, err error) {
	if conn == nil {
		return
	}
	data, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	endpoint, ok := data["endpoint"]
	if ok {
		conn.Endpoint = endpoint
	}
	driverName, ok := data["driver"]
	if ok {
		conn.DriverName = driverName
	}
	dsn, ok := data["dsn"]
	if ok {
		conn.DSN = dsn
	}
}

func beforeConnPingContextInstrumentation(call sql.CallContext, conn *sql.Conn, ctx context.Context) {
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "ping", "ping", conn.Endpoint, conn.DriverName, conn.DSN)
}

func afterConnPingContextInstrumentation(call sql.CallContext, err error) {
	instrumentEnd(call, err)
}

func beforeConnPrepareContextInstrumentation(call sql.CallContext, conn *sql.Conn, ctx context.Context, query string) {
	if conn == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": conn.Endpoint,
		"sql":      query,
		"driver":   conn.DriverName,
		"dsn":      conn.DSN,
	})
}

func afterConnPrepareContextInstrumentation(call sql.CallContext, stmt *sql.Stmt, err error) {
	if stmt == nil {
		return
	}
	callDataMap, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{
		"endpoint": callDataMap["endpoint"],
		"sql":      callDataMap["sql"],
		"driver":   callDataMap["driver"],
	}
	stmt.DSN = callDataMap["dsn"]
}

func beforeConnExecContextInstrumentation(call sql.CallContext, conn *sql.Conn, ctx context.Context, query string, args ...any) {
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "exec", query, conn.Endpoint, conn.DriverName, conn.DSN)
}

func afterConnExecContextInstrumentation(call sql.CallContext, result sql.Result, err error) {
	instrumentEnd(call, err)
}

func beforeConnQueryContextInstrumentation(call sql.CallContext, conn *sql.Conn, ctx context.Context, query string, args ...any) {
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "query", query, conn.Endpoint, conn.DriverName, conn.DSN)
}

func afterConnQueryContextInstrumentation(call sql.CallContext, rows *sql.Rows, err error) {
	instrumentEnd(call, err)
}

func beforeConnTxInstrumentation(call sql.CallContext, conn *sql.Conn, ctx context.Context, opts *sql.TxOptions) {
	if conn == nil {
		return
	}
	instrumentStart(call, ctx, "start", "START TRANSACTION", conn.Endpoint, conn.DriverName, conn.DSN)
}

func afterConnTxInstrumentation(call sql.CallContext, tx *sql.Tx, err error) {
	instrumentEnd(call, err)
}

func beforeTxPrepareContextInstrumentation(call sql.CallContext, tx *sql.Tx, ctx context.Context, query string) {
	if tx == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": tx.Endpoint,
		"sql":      query,
		"driver":   tx.DriverName,
		"dsn":      tx.DSN,
	})
}

func afterTxPrepareContextInstrumentation(call sql.CallContext, stmt *sql.Stmt, err error) {
	if stmt == nil {
		return
	}
	callDataMap, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{
		"endpoint": callDataMap["endpoint"],
		"sql":      callDataMap["sql"],
		"driver":   callDataMap["driver"],
	}
	stmt.DSN = callDataMap["dsn"]
}

func beforeTxStmtContextInstrumentation(call sql.CallContext, tx *sql.Tx, ctx context.Context, stmt *sql.Stmt) {
	if stmt == nil {
		return
	}
	call.SetData(map[string]string{
		"endpoint": stmt.Data["endpoint"],
		"driver":   stmt.Data["driver"],
		"dsn":      stmt.DSN,
	})
}

func afterTxStmtContextInstrumentation(call sql.CallContext, stmt *sql.Stmt) {
	if stmt == nil {
		return
	}
	data, ok := call.GetData().(map[string]string)
	if !ok {
		return
	}
	stmt.Data = map[string]string{}
	endpoint, ok := data["endpoint"]
	if ok {
		stmt.Data["endpoint"] = endpoint
	}
	driverName, ok := data["driver"]
	if ok {
		stmt.Data["driver"] = driverName
	}
	dsn, ok := data["dsn"]
	if ok {
		stmt.Data["dsn"] = dsn
	}
}

func beforeTxExecContextInstrumentation(call sql.CallContext, tx *sql.Tx, ctx context.Context, query string, args ...any) {
	if tx == nil {
		return
	}
	instrumentStart(call, ctx, "exec", query, tx.Endpoint, tx.DriverName, tx.DSN)
}

func afterTxExecContextInstrumentation(call sql.CallContext, result sql.Result, err error) {
	instrumentEnd(call, err)
}

func beforeTxQueryContextInstrumentation(call sql.CallContext, tx *sql.Tx, ctx context.Context, query string, args ...any) {
	if tx == nil {
		return
	}
	instrumentStart(call, ctx, "query", query, tx.Endpoint, tx.DriverName, tx.DSN)
}

func afterTxQueryContextInstrumentation(call sql.CallContext, rows *sql.Rows, err error) {
	instrumentEnd(call, err)
}

func beforeTxCommitInstrumentation(call sql.CallContext, tx *sql.Tx) {
	if tx == nil {
		return
	}
	instrumentStart(call, context.Background(), "commit", "COMMIT", tx.Endpoint, tx.DriverName, tx.DSN)
}

func afterTxCommitInstrumentation(call sql.CallContext, err error) {
	instrumentEnd(call, err)
}

func beforeTxRollbackInstrumentation(call sql.CallContext, tx *sql.Tx) {
	if tx == nil {
		return
	}
	instrumentStart(call, context.Background(), "rollback", "ROLLBACK", tx.Endpoint, tx.DriverName, tx.DSN)
}

func afterTxRollbackInstrumentation(call sql.CallContext, err error) {
	instrumentEnd(call, err)
}

func beforeStmtExecContextInstrumentation(call sql.CallContext, stmt *sql.Stmt, ctx context.Context, args ...any) {
	if stmt == nil {
		return
	}
	sql, endpoint, driverName, dsn := "", "", "", ""
	if stmt.Data != nil {
		sql, endpoint, driverName, dsn = stmt.Data["sql"], stmt.Data["endpoint"], stmt.Data["driver"], stmt.DSN
	}
	instrumentStart(call, ctx, "exec", sql, endpoint, driverName, dsn)
}

func afterStmtExecContextInstrumentation(call sql.CallContext, result sql.Result, err error) {
	instrumentEnd(call, err)
}

func beforeStmtQueryContextInstrumentation(call sql.CallContext, stmt *sql.Stmt, ctx context.Context, args ...any) {
	if stmt == nil {
		return
	}
	sql, endpoint, driverName, dsn := "", "", "", ""
	if stmt.Data != nil {
		sql, endpoint, driverName, dsn = stmt.Data["sql"], stmt.Data["endpoint"], stmt.Data["driver"], stmt.DSN
	}
	instrumentStart(call, ctx, "query", sql, endpoint, driverName, dsn)
}

func afterStmtQueryContextInstrumentation(call sql.CallContext, rows *sql.Rows, err error) {
	instrumentEnd(call, err)
}

func instrumentStart(call sql.CallContext, ctx context.Context, spanName, query, endpoint, driverName, dsn string) {
	req := databaseSqlRequest{
		opType:     calOp(query),
		sql:        query,
		endpoint:   endpoint,
		driverName: driverName,
		dsn:        dsn,
	}
	newCtx := databaseSqlInstrumenter.Start(ctx, req)
	call.SetData(map[string]interface{}{
		"dbRequest": req,
		"newCtx":    newCtx,
	})
}

func instrumentEnd(call sql.CallContext, err error) {
	callData, ok := call.GetData().(map[string]interface{})
	if !ok {
		return
	}
	dbRequest, ok := callData["dbRequest"].(databaseSqlRequest)
	if !ok {
		return
	}
	newCtx, ok := callData["newCtx"].(context.Context)
	if !ok {
		return
	}
	databaseSqlInstrumenter.End(newCtx, dbRequest, nil, err)
}

func calOp(sql string) string {
	sqls := strings.Split(sql, " ")
	var op string
	if len(sqls) > 0 {
		op = sqls[0]
	}
	return op
}
