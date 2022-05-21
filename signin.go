package main

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Signin(data user) error {
	if !IsEmailValid(data.Email) {
		return errors.New("0")
	}

	code := MakeCode()

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "active": 1})

	if err := users.FindOne(ctx, bson.M{"email": data.Email}, opts).Decode(&user); err == nil {
		if user["active"] != true {
			if r, err := ensure.UpdateOne(ctx, bson.M{"_id": user["_id"]}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "code", Value: code}}},
			}); err != nil || r.ModifiedCount == 0 {
				return errors.New("1")
			}
		}
	} else {
		// Getting the last user for id
		var last bson.M
		opts = options.FindOne().SetProjection(bson.M{"_id": 1}).SetSort(bson.D{{Key: "_id", Value: -1}})
		users.FindOne(ctx, bson.M{}, opts).Decode(&last)
		id := 1
		if last["_id"] != nil {
			id = last["_id"].(int) + 1
		}
		date := time.Now().Unix()

		ObjectId, err := users.InsertOne(ctx, bson.D{
			{Key: "_id", Value: id},
			{Key: "username", Value: strconv.Itoa(id)},
			{Key: "email", Value: data.Email},
			{Key: "title", Value: ""},
			{Key: "about", Value: ""},
			{Key: "sex", Value: 0},
			{Key: "age", Value: 0},
			{Key: "body", Value: 0},
			{Key: "height", Value: 0},
			{Key: "weight", Value: 0},
			{Key: "smokes", Value: 0},
			{Key: "drinks", Value: 0},
			{Key: "ethnicity", Value: 0},
			{Key: "search", Value: 0},
			{Key: "income", Value: 0},
			{Key: "children", Value: 0},
			{Key: "industry", Value: 0},
			{Key: "premium", Value: 0},
			{Key: "status", Value: true},
			{Key: "active", Value: false},
			{Key: "avatar", Value: false},
			{Key: "public", Value: 0},
			{Key: "private", Value: 0},
			{Key: "created_at", Value: date},
			{Key: "country_id", Value: 0},
			{Key: "city_id", Value: 0},
		})

		if err != nil {
			return errors.New("2")
		}

		if _, err := ensure.InsertOne(ctx, bson.D{
			{Key: "_id", Value: ObjectId.InsertedID},
			{Key: "code", Value: code},
		}); err != nil {
			return errors.New("3")
		}

		if err := SendCode(data.Email, code); err != nil {
			return errors.New("4")
		}
	}

	return nil
}

func checkCode() error {
	return nil
}

// Check email on valid
func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 320 {
		return false
	}

	return regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(email)
}

// Check username on valid
func IsUsernameValid(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return true
	}

	return regexp.MustCompile(`\s`).MatchString(username)
}
