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
	defer db.Close()
	if _, err := db.ExecContext(ctx, `DROP TABLE IF EXISTS users`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)`); err != nil {
		log.Fatal(err)
	}
	stmt, err := db.PrepareContext(ctx, `INSERT INTO users (id, name, age) VALUE ( ?, ?, ?)`)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		_ = stmt.Close()
	}()

	_, err = stmt.ExecContext(ctx, "1", "bar", 11)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tx.Exec(`INSERT INTO users (id, name, age) VALUE ( ?, ?, ? )`, "2", "foobar", 24); err != nil {
		log.Fatal(err)
	}

	if _, err := tx.ExecContext(ctx, `UPDATE users SET name = ? WHERE id = ?`, "foobar", "0"); err != nil {
		log.Fatal(err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "DROP", "", "mysql", "", "127.0.0.1", "DROP TABLE IF EXISTS users", "DROP")
		verifier.VerifyDbAttributes(stubs[1][0], "CREATE", "", "mysql", "", "127.0.0.1", "CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)", "CREATE")
		verifier.VerifyDbAttributes(stubs[2][0], "INSERT", "", "mysql", "", "127.0.0.1", "INSERT INTO users (id, name, age) VALUE ( ?, ?, ?)", "INSERT")
		verifier.VerifyDbAttributes(stubs[3][0], "START", "", "mysql", "", "127.0.0.1", "START TRANSACTION", "START")
		verifier.VerifyDbAttributes(stubs[4][0], "INSERT", "", "mysql", "", "127.0.0.1", "INSERT INTO users (id, name, age) VALUE ( ?, ?, ? )", "INSERT")
		verifier.VerifyDbAttributes(stubs[5][0], "UPDATE", "", "mysql", "", "127.0.0.1", "UPDATE users SET name = ? WHERE id = ?", "UPDATE")
		verifier.VerifyDbAttributes(stubs[6][0], "COMMIT", "", "mysql", "", "127.0.0.1", "COMMIT", "COMMIT")
	})
}
