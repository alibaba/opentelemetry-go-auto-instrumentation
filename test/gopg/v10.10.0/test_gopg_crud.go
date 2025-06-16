package main

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/go-pg/pg/v10"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
)

var db *pg.DB

type User struct {
	ID   uint
	Name string
	Age  uint8
}

func TestRaw() {
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)`); err != nil {
		log.Printf("%v", err)
	}
}

func TestInsert() {
	user := User{Name: "opentelemetry", Age: 18}
	if _, err := db.Model(&user).Insert(); err != nil {
		log.Printf("%v", err)
	}
}

func TestQuery() {
	var users []User
	if err := db.Model(&users).Select(); err != nil {
		log.Printf("%v", err)
	}
}

func TestUpdate() {
	if _, err := db.Model(&User{ID: 1, Age: 10}).Update("name", "hello"); err != nil {
		log.Printf("%v", err)
	}
}

func TestDelete() {
	if _, err := db.Model(&User{ID: 1}).Delete(); err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	db = pg.Connect(&pg.Options{
		Addr:     "127.0.0.1:" + os.Getenv("POSTGRES_PORT"),
		User:     "user",
		Database: "database",
	})
	TestRaw()
	TestInsert()
	TestQuery()
	TestUpdate()
	TestDelete()
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "SELECT dual", "mysql", "127.0.0.1", "SELECT VERSION()", "SELECT", "dual", nil)
		verifier.VerifyDbAttributes(stubs[1][0], "ping", "mysql", "127.0.0.1", "ping", "ping", "", nil)
		verifier.VerifyDbAttributes(stubs[2][0], "raw", "mysql", "127.0.0.1", "", "raw", "", nil)
		verifier.VerifyDbAttributes(stubs[3][0], "START", "mysql", "127.0.0.1", "START TRANSACTION", "START", "", nil)
		verifier.VerifyDbAttributes(stubs[4][0], "create", "mysql", "127.0.0.1", "", "create", "", nil)
		verifier.VerifyDbAttributes(stubs[5][0], "query", "mysql", "127.0.0.1", "", "query", "", nil)
		verifier.VerifyDbAttributes(stubs[6][0], "row", "mysql", "127.0.0.1", "", "row", "", nil)
		verifier.VerifyDbAttributes(stubs[7][0], "START", "mysql", "127.0.0.1", "START TRANSACTION", "START", "", nil)
		verifier.VerifyDbAttributes(stubs[8][0], "update", "mysql", "127.0.0.1", "", "update", "", nil)
		verifier.VerifyDbAttributes(stubs[9][0], "START", "mysql", "127.0.0.1", "START TRANSACTION", "START", "", nil)
		verifier.VerifyDbAttributes(stubs[10][0], "delete", "mysql", "127.0.0.1", "", "delete", "", nil)
	}, 1)
}
