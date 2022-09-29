package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var con = connect() // database

var DB = map[string]*mongo.Collection{
	"users": con.Database("wifer").Collection("users"),
}

const mongoConnect string = "mongodb://shizzic:WebDev77@wifer-test.ru:27017/test?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"

func main() {
	DB["users"].UpdateMany(ctx, bson.M{"premium": bson.M{"$lt": time.Now().Unix()}}, bson.D{{Key: "$set", Value: bson.D{{Key: "premium", Value: int64(0)}}}})
}

func connect() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConnect))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
