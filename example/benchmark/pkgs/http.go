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

package pkgs

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	redis "github.com/redis/go-redis/v9"
	"net/http"
	"strings"
	"time"
)

func service3(w http.ResponseWriter, r *http.Request) {
	header := r.Header
	for key, value := range header {
		values := strings.Join(value, ",")
		fmt.Printf("[Headers] Key is %s\tValue is %s\n", key, values)
	}
	m, err := RequestMySQL(context.Background())
	if err != nil {
		w.Write([]byte("error"))
		w.WriteHeader(500)
		return
	}
	d, err := redisService()
	if err != nil {
		w.Write([]byte("error"))
		w.WriteHeader(500)
		return
	}
	_, err = w.Write([]byte("Hello Http!" + "/" + m + "/" + d))
	if err != nil {
		w.Write([]byte("error"))
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)

}

var rdb *redis.Client
var err error

func init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "r-uf6x51fdb6tqmrmisq.redis.rds.aliyuncs.com:6379",
		Password: "Hello1234",
	})
	_ = redis.NewFailoverClient(&redis.FailoverOptions{
		MaxRetries: 3,
	})

	_ = redis.NewRing(&redis.RingOptions{})
	_ = redis.NewClusterClient(&redis.ClusterOptions{})
	mysql.ParseDSN("root:arms@123@tcp(rm-uf65qlp1f3q8752sd.mysql.rds.aliyuncs.com:3306)/mysqldb")
	db, err = sql.Open("mysql",
		"root:arms@123@tcp(rm-uf65qlp1f3q8752sd.mysql.rds.aliyuncs.com:3306)/mysqldb")
	if err != nil {
		return
	}
	// See "Important settings" section.
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if _, err := db.ExecContext(context.Background(), `CREATE TABLE IF NOT EXISTS usersx (id char(255), name VARCHAR(255), age INTEGER)`); err != nil {
		fmt.Printf("exec create error: %v", err)
	}

	// test insert
	if _, err := db.ExecContext(context.Background(), `INSERT INTO usersx (id, name, age) VALUE ( ?, ?, ?)`, "0", "foo", 10); err != nil {
		fmt.Printf("exec insert error: %v", err)
	}
}

func redisService() (string, error) {
	key := "key_TestSetAndGet"
	value := "value_TestSetAndGet"
	ctx := context.Background()
	res, err := rdb.Set(ctx, key, value, 10*time.Second).Result()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println("set result:", res)
	v, err := rdb.Get(ctx, key).Result()
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	fmt.Println("get result:", v)
	return "Finish Writing to Redis", nil
}

func SetupHttp() {
	http.Handle("/http-service1", http.HandlerFunc(service3))
	err := http.ListenAndServe("0.0.0.0:9000", nil)
	if err != nil {
		panic(err)
	}
}

var db *sql.DB

func RequestMySQL(ctx context.Context) (string, error) {
	var name string
	// test select
	if err := db.QueryRowContext(ctx, `SELECT name FROM usersx WHERE id = ?`, "0").Scan(&name); err != nil {
		fmt.Printf("query select error: %v", err)
		return "", err
	}

	return "Exec go-sql-driver command finished", nil
}
