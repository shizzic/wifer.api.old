package main

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type List struct {
	Limit     int64  `json:"limit"`
	Skip      int64  `json:"skip"`
	SortKey   string `json:"sortKey"`
	SortValue int64  `json:"sortValue"`
	AgeMin    uint8  `json:"ageMin"`
	AgeMax    uint8  `json:"ageMax"`
	HeightMin uint8  `json:"heightMin"`
	HeightMax uint8  `json:"heightMax"`
	WeightMin uint8  `json:"weightMin"`
	WeightMax uint8  `json:"weightMax"`
	Body      []int8 `json:"body"`
	Sex       []int8 `json:"sex"`
	Smokes    []int8 `json:"smokes"`
	Drinks    []int8 `json:"drinks"`
	Ethnicity []int8 `json:"ethnicity"`
	Search    []int8 `json:"search"`
	Income    []int8 `json:"income"`
	Children  []int8 `json:"children"`
	Industry  []int8 `json:"industry"`
	Premium   []bool `json:"premium"`
	Text      string `json:"text"`
}

var list []bson.M

// Fewer 40ms :D
func GetUsers(data List) []bson.M {
	opts := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1}).SetLimit(data.Limit).SetSkip(data.Skip).SetSort(bson.D{{Key: data.SortKey, Value: data.SortValue}})
	filter := bson.M{"age": bson.M{"$gte": data.AgeMin, "$lte": data.AgeMax}, "height": bson.M{"$gte": data.HeightMin, "$lte": data.HeightMax}, "weight": bson.M{"$gte": data.WeightMin, "$lte": data.WeightMax}, "body": bson.M{"$in": data.Body}, "sex": bson.M{"$in": data.Sex}, "smokes": bson.M{"$in": data.Smokes}, "drinks": bson.M{"$in": data.Drinks}, "ethnicity": bson.M{"$in": data.Ethnicity}, "search": bson.M{"$in": data.Search}, "income": bson.M{"$in": data.Income}, "children": bson.M{"$in": data.Children}, "industry": bson.M{"$in": data.Industry}, "premium": bson.M{"$in": data.Premium}}
	if data.Text != "" {
		filter["$text"] = bson.M{"$search": data.Text}
	}
	cursor, _ := users.Find(ctx, filter, opts)
	cursor.All(ctx, &list)

	return list
}

// Get profile that has status 1
func GetProfile(id int) (bson.M, error) {
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"username": 1, "title": 1, "about": 1, "sex": 1, "age": 1, "body": 1, "height": 1, "weight": 1, "smokes": 1, "drinks": 1, "ethnicity": 1, "search": 1, "income": 1, "children": 1, "industry": 1, "premium": 1, "active": 1, "avatar": 1, "public": 1, "private": 1, "created_at": 1, "last_time": 1, "online": 1, "country_id": 1, "city_id": 1})

	if err := users.FindOne(ctx, bson.M{"_id": id, "status": true}, opts).Decode(&user); err != nil {
		return user, errors.New("0")
	}

	return user, nil
}

// func GetCity(id int) bson.M {
// 	var user bson.M
// 	opts := options.FindOne().SetProjection(bson.M{"username": 1, "title": 1, "about": 1, "sex": 1, "age": 1, "body": 1, "height": 1, "weight": 1, "smokes": 1, "drinks": 1, "ethnicity": 1, "search": 1, "income": 1, "children": 1, "industry": 1, "premium": 1, "active": 1, "avatar": 1, "public": 1, "private": 1, "created_at": 1, "last_time": 1, "online": 1, "country_id": 1, "city_id": 1})

// 	if err := cities.FindOne(ctx, bson.M{"_id": id}, opts).Decode(&user); err != nil {
// 		return user, errors.New("0")
// 	}

// 	return user, nil
// }
