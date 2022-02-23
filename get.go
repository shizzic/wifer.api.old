package main

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type List struct {
	SortKey   string  `json:"sortKey"`
	SortValue int64   `json:"sortValue"`
	AgeMin    uint8   `json:"ageMin"`
	AgeMax    uint8   `json:"ageMax"`
	HeightMin uint8   `json:"heightMin"`
	HeightMax uint8   `json:"heightMax"`
	Limit     int64   `json:"limit"`
	Skip      int64   `json:"skip"`
	Body      []int64 `json:"body"`
}

var list []bson.M

func GetUsers(data List) []bson.M {
	opts := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1}).SetLimit(data.Limit).SetSkip(data.Skip).SetSort(bson.D{{Key: data.SortKey, Value: data.SortValue}})
	filter := bson.M{"age": bson.M{"$gte": data.AgeMin, "$lte": data.AgeMax}, "height": bson.M{"$gte": data.HeightMin, "$lte": data.HeightMax}, "body": bson.M{"$in": data.Body}}
	cursor, _ := users.Find(ctx, filter, opts)
	cursor.All(ctx, &list)

	return list
}
