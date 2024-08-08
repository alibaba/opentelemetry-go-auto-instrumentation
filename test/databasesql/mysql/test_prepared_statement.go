package main

import (
	"context"
	"database/sql"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	db, err := sql.Open("mysql",
		"test:test@tcp(127.0.0.1:"+os.Getenv("MYSQL_PORT")+")/test")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS users`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `INSERT INTO users (id, name, age) VALUE ( ?, ?, ?)`, "0", "foo", 10); err != nil {
		log.Fatal(err)
	}
	stmt, err := db.Prepare("select id, name from users where id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {

	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "DROP", "", "mysql", "", "127.0.0.1", "DROP TABLE IF EXISTS users", "DROP")
		verifier.VerifyDbAttributes(stubs[1][0], "CREATE", "", "mysql", "", "127.0.0.1", "CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)", "CREATE")
		verifier.VerifyDbAttributes(stubs[2][0], "INSERT", "", "mysql", "", "127.0.0.1", "INSERT INTO users (id, name, age) VALUE ( ?, ?, ?)", "INSERT")
		verifier.VerifyDbAttributes(stubs[3][0], "select", "", "mysql", "", "127.0.0.1", "select id, name from users where id = ?", "select")
	})
}
