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

package main

import (
	"bytes"
	"encoding/json"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/test/verifier"
	"github.com/elastic/go-elasticsearch/v8"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
	"strings"
)

var (
	client *elasticsearch.Client
	url    = "http://127.0.0.1:" + os.Getenv("OTEL_ES_PORT")
)

func main() {
	var err error
	client, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
		Password:  "123456",
		Username:  "elastic",
	})
	if err != nil {
		panic(err)
	}
	// creating an index
	_, err = client.Indices.Create("my_index")
	if err != nil {
		log.Printf("failed to create index %v\n", err)
	}
	// indexing documents
	document := struct {
		Name string `json:"name"`
	}{
		"go-elasticsearch",
	}
	data, _ := json.Marshal(document)
	_, err = client.Index("my_index", bytes.NewReader(data))
	if err != nil {
		log.Printf("failed to index document %v\n", err)
	}
	// getting documents
	_, err = client.Get("my_index", "id")
	if err != nil {
		log.Printf("failed to get documents %v\n", err)
	}
	// searching documents
	query := `{ "query": { "match_all": {} } }`
	_, err = client.Search(
		client.Search.WithIndex("my_index"),
		client.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		log.Printf("failed to search documents %v\n", err)
	}
	// updating documents
	_, err = client.Update("my_index", "id", strings.NewReader(`{doc: { language: "Go" }}`))
	if err != nil {
		log.Printf("failed to update document %v\n", err)
	}
	// deleting documents
	_, err = client.Delete("my_index", "id")
	if err != nil {
		log.Printf("failed to delete document %v\n", err)
	}
	// deleting an index
	_, err = client.Indices.Delete([]string{"my_index"})
	if err != nil {
		log.Printf("failed to delete index %v\n", err)
	}
	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		for i, _ := range stubs {
			span := stubs[i][0]
			println(span.Name)
			for _, attr := range span.Attributes {
				log.Printf("%v %v", attr.Key, attr.Value)
			}
			println()
		}
	}, 1)
}
