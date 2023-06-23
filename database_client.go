package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbClient struct {
	mongoClient *mongo.Client
	database    *mongo.Database
	collection  *mongo.Collection
}

func (c *dbClient) Init(ctx context.Context, db, collection string) error {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))

	if err != nil {
		return err
	}

	database := client.Database(db)
	col := database.Collection(collection)

	c.mongoClient = client
	c.database = database
	c.collection = col
	return nil
}

func (c *dbClient) Close(ctx context.Context) error {
	err := c.mongoClient.Disconnect(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (c *dbClient) upsortItem(ctx context.Context, title, summary string) error {
	filter := bson.M{"title": title}
	update := bson.M{
		"$set": bson.M{
			"summary": summary,
		},
	}

	// Set the upsert option to true
	opt := options.Update().SetUpsert(true)

	// Perform an upsert operation
	_, err := c.collection.UpdateOne(ctx, filter, update, opt)

	if err != nil {
		return err
	}

	return nil
}
