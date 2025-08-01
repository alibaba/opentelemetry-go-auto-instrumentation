// Copyright (c) 2025 Alibaba Group Holding Ltd.
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

package main

import (
	"github.com/alibaba/loongsuite-go-agent/test/verifier"
	"github.com/gocql/gocql"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
)

var session *gocql.Session

type ShoppingCart struct {
	Userid    string
	ItemCount int
}

func TestCreateKeySpace() {
	if err := session.Query(`CREATE KEYSPACE IF NOT EXISTS cassandra WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : '1' };`).Exec(); err != nil {
		log.Printf("%v", err)
	}
}

func TestCreateTable() {
	if err := session.Query(`CREATE TABLE IF NOT EXISTS cassandra.shopping_cart (userid text PRIMARY KEY,item_count int,last_update_timestamp timestamp);`).Exec(); err != nil {
		log.Printf("%v", err)
	}
}

func TestInsert() {
	if err := session.Query("INSERT INTO cassandra.shopping_cart\n(userid, item_count, last_update_timestamp)\nVALUES ('9876', 2, toTimeStamp(now()));").Exec(); err != nil {
		log.Printf("%v", err)
	}
}

func TestQuery() {
	var shoppingCarts []ShoppingCart
	if err := session.Query("SELECT * FROM store.shopping_cart;").Scan(&shoppingCarts); err != nil {
		log.Printf("%v", err)
	}
}

func TestUpdate() {
	if err := session.Query("UPDATE cassandra.shopping_cart \nSET item_count = 6, last_update_timestamp = toTimeStamp(now()) \nWHERE userid = '1234';").Exec(); err != nil {
		log.Printf("%v", err)
	}
}

func TestDelete() {
	if err := session.Query("DELETE FROM cassandra.shopping_cart \nWHERE userid = '1234';").Exec(); err != nil {
		log.Printf("%v", err)
	}
}

func TestDropTable() {
	if err := session.Query("DROP TABLE IF EXISTS store.shopping_cart;").Exec(); err != nil {
		log.Printf("%v", err)
	}
}

func main() {
	s, err := gocql.NewSession(gocql.ClusterConfig{
		Hosts:         []string{"127.0.0.1:" + os.Getenv("CASSCANDRA_PORT")},
		Keyspace:      "cassandra",
		Authenticator: gocql.PasswordAuthenticator{Username: "cassandra"},
	})
	if err != nil {
		panic(err)
	}
	session = s
	TestCreateKeySpace()
	TestCreateTable()
	TestInsert()
	TestQuery()
	TestUpdate()
	TestDelete()
	TestDropTable()
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		verifier.VerifyDbAttributes(stubs[0][0], "postgresql", "postgresql", "127.0.0.1", "CREATE TABLE IF NOT EXISTS users (id char(255), name VARCHAR(255), age INTEGER)", "", "", nil)
		verifier.VerifyDbAttributes(stubs[1][0], "INSERT", "postgresql", "127.0.0.1", "INSERT INTO \"users\" (\"id\", \"name\", \"age\") VALUES (DEFAULT, 'opentelemetry', 18)", "INSERT", "", nil)
		verifier.VerifyDbAttributes(stubs[2][0], "SELECT", "postgresql", "127.0.0.1", "SELECT \"user\".\"id\", \"user\".\"name\", \"user\".\"age\" FROM \"users\" AS \"user\"", "SELECT", "", nil)
		verifier.VerifyDbAttributes(stubs[3][0], "UPDATE", "postgresql", "127.0.0.1", "UPDATE \"users\" AS \"user\" SET \"name\" = NULL, \"age\" = 10 WHERE \"user\".\"id\" = '1'", "UPDATE", "", nil)
		verifier.VerifyDbAttributes(stubs[4][0], "DELETE", "postgresql", "127.0.0.1", "DELETE FROM \"users\" AS \"user\" WHERE \"user\".\"id\" = '1'", "DELETE", "", nil)
		verifier.VerifyDbAttributes(stubs[5][0], "DROP TABLE", "postgresql", "127.0.0.1", "DROP TABLE \"users\"", "DROP TABLE", "", nil)
	}, 1)
}
