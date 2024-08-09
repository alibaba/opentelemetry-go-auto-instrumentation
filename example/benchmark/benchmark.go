package main

import (
	"bytes"
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/baggage"
	_ "go.opentelemetry.io/otel/sdk/trace"
	"log"
	"net/http"
	"time"
)

var db *sql.DB
var rdb *redis.Client

func setupHttp() {
	ctx := context.Background()
	var err error
	db, err = sql.Open("mysql",
		"test:test@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	http.HandleFunc("/http-service", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("this is otel http service"))
	})
	http.HandleFunc("/error-service", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("wrong otel http endpoint"))
	})
	http.HandleFunc("/mysql-service", func(w http.ResponseWriter, r *http.Request) {
		mysqlCrud(ctx)
		_, _ = w.Write([]byte("crud mysql finish"))
	})
	http.HandleFunc("/redis-service", func(w http.ResponseWriter, r *http.Request) {
		redisCrud(ctx)
		_, _ = w.Write([]byte("crud redis finish"))
	})
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("otel-service-type") == "http" {
			getWithWriteBack("http-service", writer, request)
		} else if request.Header.Get("otel-service-type") == "mysql" {
			postWithWriteBack("mysql-service", writer, request)
		} else if request.Header.Get("otel-service-type") == "redis" {
			getWithWriteBack("redis-service", writer, request)
		} else {
			getWithWriteBack("error-service", writer, request)
		}
	})

	http.HandleFunc("/request-all", func(writer http.ResponseWriter, request *http.Request) {
		getWithWriteBack("http-service", writer, request)
		time.Sleep(1 * time.Second)
		postWithWriteBack("mysql-service", writer, request)
		time.Sleep(1 * time.Second)
		getWithWriteBack("redis-service", writer, request)
		time.Sleep(1 * time.Second)
		getWithWriteBack("error-service", writer, request)
		_, _ = writer.Write([]byte("send request to all services finish"))
	})
	println("HTTP server started")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

func getWithWriteBack(path string, w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	_, err := client.Get("http://127.0.0.1:8080/" + path)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	} else {
		_, _ = w.Write([]byte("request to " + path + "\n"))
	}
}

func postWithWriteBack(path string, w http.ResponseWriter, r *http.Request) {
	client := http.Client{}
	jsonData := []byte(`{"key1": "value1", "key2": "value2"}`)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8080/"+path, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = client.Do(req)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	} else {
		_, _ = w.Write([]byte("request to " + path + "\n"))

	}
}

func mysqlCrud(ctx context.Context) {
	var err error
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
	_, err = stmt.Query(1)
	if err != nil {
		log.Fatal(err)
	}
}

func redisCrud(ctx context.Context) {
	_, err := rdb.Set(ctx, "a", "b", 5*time.Second).Result()
	if err != nil {
		panic(err)
	}
	_, err = rdb.Get(ctx, "a").Result()
	if err != nil {
		panic(err)
	}
}

func main() {
	setupHttp()
}
