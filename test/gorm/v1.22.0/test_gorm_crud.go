package main

import (
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct {
	ID   uint
	Name string
	Age  uint8
}

func TestRaw() {
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)`).Error; err != nil {
		log.Printf("%v", err)
	}
}

func TestCreate() {
	user := User{Name: "opentelemetry", Age: 18}
	if err := db.Create(&user).Error; err != nil {
		log.Printf("%v", err)
	}
}

func TestQuery() {
	var user User
	if err := db.First(&user).Error; err != nil {
		log.Printf("%v", err)
	}
}

func TestRow() {
	var name string
	var age uint8
	row := db.Table("users").Where("name = ?", "opentelemetry").Select("name", "age").Row()
	row.Scan(&name, &age)
}

func TestUpdate() {
	tx := db.Model(&User{}).Where("name = ?", "opentelemetry").Update("name", "hello")
	if err := tx.Error; err != nil {
		log.Printf("%v", err)
	}
}

func TestDelete() {
	if err := db.Delete(&User{}, 1).Error; err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	tmpDB, err := gorm.Open(mysql.Open("test:test@tcp(127.0.0.1:"+os.Getenv("MYSQL_PORT")+")/test"), &gorm.Config{})
	if err != nil {
		log.Fatalf("open db error: %v \n", err)
	}
	db = tmpDB
	TestRaw()
	TestCreate()
	TestQuery()
	TestRow()
	TestUpdate()
	TestDelete()
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "SELECT", "", "mysql", "", "127.0.0.1", "SELECT VERSION()", "SELECT")
		verifier.VerifyDbAttributes(stubs[1][0], "ping", "", "mysql", "", "127.0.0.1", "ping", "ping")
		verifier.VerifyDbAttributes(stubs[2][0], "raw", "test", "mysql", "test", "127.0.0.1", "", "raw")
		verifier.VerifyDbAttributes(stubs[3][0], "START", "", "mysql", "", "127.0.0.1", "START TRANSACTION", "START")
		verifier.VerifyDbAttributes(stubs[4][0], "create", "test", "mysql", "test", "127.0.0.1", "", "create")
		verifier.VerifyDbAttributes(stubs[5][0], "query", "test", "mysql", "test", "127.0.0.1", "", "query")
		verifier.VerifyDbAttributes(stubs[6][0], "row", "test", "mysql", "test", "127.0.0.1", "", "row")
		verifier.VerifyDbAttributes(stubs[7][0], "START", "", "mysql", "", "127.0.0.1", "START TRANSACTION", "START")
		verifier.VerifyDbAttributes(stubs[8][0], "update", "test", "mysql", "test", "127.0.0.1", "", "update")
		verifier.VerifyDbAttributes(stubs[9][0], "START", "", "mysql", "", "127.0.0.1", "START TRANSACTION", "START")
		verifier.VerifyDbAttributes(stubs[10][0], "delete", "test", "mysql", "test", "127.0.0.1", "", "delete")

	})
}
