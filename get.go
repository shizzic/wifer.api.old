package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var list []bson.M

func GetUsers() []bson.M {
	opts := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1}).SetLimit(2).SetSkip(1).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, _ := users.Find(ctx, bson.M{}, opts)
	cursor.All(ctx, &list)

	return list
}
