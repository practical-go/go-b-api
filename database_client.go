package main

import (
	"context"
	"time"

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

func (c *dbClient) upsertNews(ctx context.Context, title, summary, uuid string) error {
	filter := bson.M{"UUID": uuid}
	update := bson.M{
		"$set": bson.M{
			"Title":   title,
			"Summary": summary,
		},
	}

	opt := options.Update().SetUpsert(true)

	_, err := c.collection.UpdateOne(ctx, filter, update, opt)

	if err != nil {
		return err
	}

	return nil
}

func (c *dbClient) fetchNews() ([]News, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := c.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	var news []News

	for cursor.Next(ctx) {
		var result bson.M
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		news = append(news, News{
			Title:   result["Title"].(string),
			Summary: result["Summary"].(string),
		})
	}

	return news, nil
}
