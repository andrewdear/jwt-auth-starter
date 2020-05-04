package services

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

// This file tells everything about how to connect to ou mongo db instance and helper methods around it

var Database *mongo.Database

func CreateContext() (context.Context, context.CancelFunc) {
	// define a context - tell the connection how long to wait before throwing a timeout error, have to create a new one as it counts down from creation
	//TODO: _ is a cancel function this needs to be called so return this aswell so we can defer call it.
	ctx , cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return ctx, cancel
}

func ConnectToMongo(databaseName string) error {
	// Connect to the mongoDb client
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://andrew2:r3GRtOTGXK!5dM*wyp@cluster0-bqlsu.mongodb.net/test?retryWrites=true&w=majority"))

	if err != nil {
		return err
	}

	ctx, cancel := CreateContext()
	// if we end up connecting then cancel the context
	defer cancel()

	err = client.Connect(ctx)
	// check to see if there was an error connecting, otherwise we are connected
	if err != nil {
		return err
	}

	//Create a database or use the one that is already there
	quickstartDatabase := client.Database(databaseName)

	// after creating the database we save a reference to it in the global Database that all services have access to
	Database = quickstartDatabase

	return nil
}


func GetCollection(collectionName string) *mongo.Collection {

	return Database.Collection(collectionName)

}