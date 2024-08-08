//go:build ignore

package databasesql

import (
	"errors"
	"fmt"
	nurl "net/url"
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
