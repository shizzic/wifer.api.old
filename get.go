package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var list []bson.M

func GetUsers(data List) []bson.M {
	filter := bson.M{"age": bson.M{"$gte": data.AgeMin, "$lte": data.AgeMax}, "height": bson.M{"$gte": data.HeightMin, "$lte": data.HeightMax}}
	opts := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1}).SetLimit(data.Limit).SetSkip(data.Skip).SetSort(bson.D{{Key: data.SortKey, Value: data.SortValue}})
	cursor, _ := users.Find(ctx, filter, opts)
	cursor.All(ctx, &list)

	return list
}
