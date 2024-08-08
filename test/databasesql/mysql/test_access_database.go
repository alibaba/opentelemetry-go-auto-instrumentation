package main

import (
	"database/sql"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	_ "github.com/go-sql-driver/mysql"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
)

func main() {
	db, err := sql.Open("mysql",
		"test:test@tcp(127.0.0.1:"+os.Getenv("MYSQL_PORT")+")/test")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "ping", "", "mysql", "", "127.0.0.1", "ping", "ping")
	})
}
