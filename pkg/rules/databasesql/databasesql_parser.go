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

package databasesql

import (
	"errors"
	"fmt"
	nurl "net/url"
	"github.com/xwb1989/sqlparser"
)

func parseDSN(driverName, dsn string) (addr string, err error) {
	// TODO: need a more delegate DFA
	switch driverName {
	case "mysql":
		return parseMySQL(dsn)
	case "postgres":
		fallthrough
	case "postgresql":
		return parsePostgres(dsn)
	}

	return "", errors.New("invalid DSN")
}

func parsePostgres(url string) (addr string, err error) {
	u, err := nurl.Parse(url)
	if err != nil {
		return "", err
	}

	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return "", fmt.Errorf("invalid connection protocol: %s", u.Scheme)
	}

	return u.Host + ":" + u.Port(), nil
}

func parseMySQL(dsn string) (addr string, err error) {
	n := len(dsn)
	i, j := -1, -1
	for k := 0; k < n; k++ {
		if dsn[k] == '(' {
			i = k
		}
		if dsn[k] == ')' {
			j = k
			break
		}
	}
	if i >= 0 && j > i {
		return dsn[i+1 : j], nil
	}
	return "", errors.New("invalid MySQL DSN")
}

func collectionExtractor(query string) string {
	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return ""
	}
	// Only support DML currently
	switch stmt := stmt.(type) {
	case *sqlparser.Select:
		return getTableName(stmt.From)
	case *sqlparser.Update:
		return getTableName(stmt.TableExprs)
	case *sqlparser.Insert:
		return stmt.Table.Name.String()
	case *sqlparser.Delete:
		return getTableName(stmt.TableExprs)
	default:
		return ""
	}
}

func getTableName(node sqlparser.SQLNode) string {
    switch n := node.(type) {
    case sqlparser.TableName:
        return n.Name.String()
    case sqlparser.TableExprs:
        for _, expr := range n {
            aliasedExpr, ok := expr.(*sqlparser.AliasedTableExpr)
            if !ok {
                continue
            }
            tableName, ok := aliasedExpr.Expr.(sqlparser.TableName)
            if ok {
                return tableName.Name.String()
            }
        }
    } 
    return ""
}