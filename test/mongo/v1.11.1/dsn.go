package main

import "os"

const (
	db = "otel_database"
)

var dsn = "mongodb://127.0.0.1:" + os.Getenv("MONGO_PORT")

type Restaurant struct {
	Name         string
	RestaurantId string        `bson:"restaurant_id,omitempty"`
	Cuisine      string        `bson:"cuisine,omitempty"`
	Address      interface{}   `bson:"address,omitempty"`
	Borough      string        `bson:"borough,omitempty"`
	Grades       []interface{} `bson:"grades,omitempty"`
}
