package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var collection *mongo.Collection

func ConnectMongo() error {

	// Create client
	// localhost --port 27017
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return err
	}

	// Create connect
	err = client.Connect(context.TODO())
	if err != nil {
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	// Get a handle for your collection
	collection = client.Database("proxyAlexSirmais").Collection("requests")

	// delete old db
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	if err = collection.Drop(ctx); err != nil {
		log.Fatal(err)
	}

	return nil
}
