package main

import (
	"context"
	"fmt"
	"github.com/alibaba/opentelemetry-go-auto-instrumentation/pkg/verifier"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"log"
	"os"
)

const (
	db = "otel_database"
)

var dsn = "mongodb://127.0.0.1:" + os.Getenv("MONGO_PORT")

var http_server_port int

// User model.
type User struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Name string             `bson:"name"`
	Age  int                `bson:"age"`
}

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(dsn))
	if err != nil {
		panic(fmt.Sprintf("connect mongodb error %v \n", err))
	}
	ctx := context.Background()
	err = TestCreateCollection(ctx, client)
	if err != nil {
		log.Printf("failed to create collection: %v", err)
	}
	err = TestCreate(ctx, client)
	if err != nil {
		log.Printf("failed to create: %v", err)
	}
	err = TestQuery(ctx, client)
	if err != nil {
		log.Printf("failed to query: %v", err)
	}
	err = TestUpdate(ctx, client)
	if err != nil {
		log.Printf("failed to update: %v", err)
	}
	err = TestDelete(ctx, client)
	if err != nil {
		log.Printf("failed to delete: %v", err)
	}

	verifier.WaitAndAssertTraces(func(stubs []tracetest.SpanStubs) {
		// TODO: add http server as root span
		verifier.VerifyDbAttributes(stubs[0][0], "create", "otel_database", "mongodb", "", "127.0.0.1", "create", "create")
		verifier.VerifyDbAttributes(stubs[1][0], "insert", "otel_database", "mongodb", "", "127.0.0.1", "insert", "insert")
		verifier.VerifyDbAttributes(stubs[2][0], "find", "otel_database", "mongodb", "", "127.0.0.1", "find", "find")
		verifier.VerifyDbAttributes(stubs[3][0], "find", "otel_database", "mongodb", "", "127.0.0.1", "find", "find")
		verifier.VerifyDbAttributes(stubs[4][0], "update", "otel_database", "mongodb", "", "127.0.0.1", "update", "update")
		verifier.VerifyDbAttributes(stubs[5][0], "delete", "otel_database", "mongodb", "", "127.0.0.1", "delete", "delete")
	})
}

func TestCreateCollection(ctx context.Context, client *mongo.Client) error {
	return client.Database(db).CreateCollection(ctx, "users")
}

func TestCreate(ctx context.Context, client *mongo.Client) error {
	collection := client.Database(db).Collection("users")
	objectID, err := primitive.ObjectIDFromHex("637334579a3d0cf34c31d08f")
	if err != nil {
		return err
	}
	_, err = collection.InsertOne(ctx, &User{
		ID:   objectID,
		Name: "Elza2",
		Age:  18,
	})
	return err
}

func TestQuery(ctx context.Context, client *mongo.Client) error {
	collection := client.Database(db).Collection("users")
	var user User
	err := collection.FindOne(ctx, bson.D{
		{Key: "name", Value: "Elza2"},
	}).Decode(&user)

	return err
}

func TestUpdate(ctx context.Context, client *mongo.Client) error {
	collection := client.Database(db).Collection("users")

	var user User
	err := collection.FindOne(ctx, bson.D{
		{Key: "name", Value: "Elza2"},
	}).Decode(&user)
	if err != nil {
		return err
	}

	_, err = collection.UpdateByID(ctx, user.ID, primitive.D{{
		Key: "$set", Value: primitive.D{
			{Key: "age", Value: 22},
		},
	}})
	return err
}

func TestDelete(ctx context.Context, client *mongo.Client) error {
	collection := client.Database(db).Collection("users")

	_, err := collection.DeleteOne(ctx, primitive.D{{Key: "name", Value: "Elza2"}})
	return err
}
