package main

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type signin struct {
	ID     string `form:"id"`
	Token  string `form:"token"`
	Email  string `form:"email"`
	Method string `form:"method"`
	Api    bool   `form:"api"`
}

func Signin(email string, c gin.Context, api bool) (int, error) {
	if !IsEmailValid(email) {
		return 0, errors.New("1")
	}

	code := MakeCode()

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "username": 1, "email": 1, "status": 1, "active": 1})
	notFound := DB["users"].FindOne(ctx, bson.M{"email": email}, opts).Decode(&user)

	if notFound == nil {
		if !user["status"].(bool) {
			return 0, errors.New("4")
		}

		if !user["active"].(bool) {
			DB["users"].UpdateOne(ctx, bson.M{"_id": user["_id"].(int32)}, bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: true}}}})
		}

		if api {
			MakeCookies(strconv.Itoa(int(user["_id"].(int32))), user["username"].(string), 86400*120, c)
			return int(user["_id"].(int32)), nil
		} else {
			DB["ensure"].DeleteOne(ctx, bson.M{"_id": user["_id"]})
			DB["ensure"].InsertOne(ctx, bson.D{
				{Key: "_id", Value: user["_id"]},
				{Key: "code", Value: code},
			})

			if err := SendCode(email, code, strconv.Itoa(int(user["_id"].(int32)))); err != nil {
				return 0, errors.New("2")
			}
		}
	} else {
		// Getting the last user for id
		var last bson.M
		opts = options.FindOne().SetProjection(bson.M{"_id": 1}).SetSort(bson.D{{Key: "_id", Value: -1}})
		DB["users"].FindOne(ctx, bson.M{}, opts).Decode(&last)

		id := 1
		if last["_id"] != nil {
			id = int(last["_id"].(int32)) + 1
		}
		date := time.Now().Unix()

		ObjectId, err := DB["users"].InsertOne(ctx, bson.D{
			{Key: "_id", Value: id},
			{Key: "username", Value: strconv.Itoa(id)},
			{Key: "email", Value: email},
			{Key: "title", Value: ""},
			{Key: "about", Value: ""},
			{Key: "is_about", Value: false},
			{Key: "sex", Value: 0},
			{Key: "age", Value: 0},
			{Key: "body", Value: 0},
			{Key: "height", Value: 0},
			{Key: "weight", Value: 0},
			{Key: "smokes", Value: 0},
			{Key: "drinks", Value: 0},
			{Key: "ethnicity", Value: 0},
			{Key: "search", Value: []int{}},
			{Key: "prefer", Value: 0},
			{Key: "income", Value: 0},
			{Key: "children", Value: 0},
			{Key: "industry", Value: 0},
			{Key: "country_id", Value: 0},
			{Key: "city_id", Value: 0},
			{Key: "premium", Value: int64(0)},
			{Key: "trial", Value: false},
			{Key: "status", Value: api},
			{Key: "active", Value: api},
			{Key: "created_at", Value: date},
			{Key: "last_time", Value: date},
			{Key: "online", Value: false},
			{Key: "avatar", Value: false},
			{Key: "public", Value: 0},
			{Key: "private", Value: 0},
			{Key: "images", Value: 0},
		})

		if err != nil {
			return 0, errors.New("3")
		}

		if api {
			id := strconv.Itoa(int(ObjectId.InsertedID.(int32)))
			MakeCookies(id, id, 86400*120, c)
			return int(ObjectId.InsertedID.(int32)), nil
		} else {
			if _, err := DB["ensure"].InsertOne(ctx, bson.D{
				{Key: "_id", Value: ObjectId.InsertedID},
				{Key: "code", Value: code},
			}); err != nil {
				// Delete new user, because code wasn't added
				DB["users"].DeleteOne(ctx, bson.M{"_id": int(ObjectId.InsertedID.(int32))})
				return 0, errors.New("3")
			}

			if err := SendCode(email, code, strconv.Itoa(int(ObjectId.InsertedID.(int32)))); err != nil {
				return 0, errors.New("2")
			}
		}
	}

	return 0, nil
}

func CheckApi(data signin, c gin.Context) (int, error) {
	var email string
	var err error = nil

	switch data.Method {
	case "Google":
		email, err = isGoogle(data.ID, data.Token)
	case "Facebook":
		email, err = isFacebook(data.ID, data.Token)
	}

	if err != nil {
		return 0, errors.New("0")
	}

	id, err := Signin(email, c, true)

	if err != nil {
		return id, errors.New(err.Error())
	}

	return id, nil
}
