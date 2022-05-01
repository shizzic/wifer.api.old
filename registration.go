package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

func Registration(data user) error {
	if IsUsernameValid(data.Username) {
		return errors.New("incorrect username")
	}

	if !IsEmailValid(data.Email) {
		return errors.New("incorrect email")
	}

	if !IsPasswordValid(&data.Password) {
		return errors.New("incorrect password")
	}

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "active": 1})

	if err := users.FindOne(ctx, bson.M{"email": data.Email}, opts).Decode(&user); err == nil {
		if user["active"] != true {
			if r, err := users.DeleteOne(ctx, bson.M{"email": data.Email}); err != nil || r.DeletedCount == 0 {
				return errors.New("error with replace accounts")
			}

			ensure.DeleteOne(ctx, bson.M{"_id": user["_id"]})
		}
	}

	// last["_id"].(string)
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
		{Key: "username", Value: data.Username},
		{Key: "email", Value: data.Email},
		{Key: "password_hash", Value: data.Password},
		{Key: "title", Value: data.Title},
		{Key: "sex", Value: data.Sex},
		{Key: "age", Value: data.Age},
		{Key: "body", Value: data.Body},
		{Key: "height", Value: data.Height},
		{Key: "weight", Value: data.Weight},
		{Key: "smokes", Value: data.Smokes},
		{Key: "drinks", Value: data.Drinks},
		{Key: "ethnicity", Value: data.Ethnicity},
		{Key: "search", Value: data.Search},
		{Key: "income", Value: data.Income},
		{Key: "children", Value: data.Children},
		{Key: "industry", Value: data.Industry},
		{Key: "premium", Value: false},
		{Key: "status", Value: true},
		{Key: "active", Value: false},
		{Key: "avatar", Value: false},
		{Key: "public", Value: 0},
		{Key: "private", Value: 0},
		{Key: "created_at", Value: date},
		{Key: "about", Value: ""},
	})

	if err != nil {
		return errors.New("user not inserted")
	}

	token := MakeToken()

	if _, err := ensure.InsertOne(ctx, bson.D{
		{Key: "_id", Value: ObjectId.InsertedID},
		{Key: "token", Value: token},
	}); err != nil {
		return errors.New("ensure not inserted")
	}

	if err := SendVerifyEmail(data.Username, data.Email, token, fmt.Sprint(ObjectId.InsertedID)); err != nil {
		return errors.New("couldn't send message to your email")
	}

	return nil
}

// Check email on valid
func IsEmailValid(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	return regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$").MatchString(email)
}

// Check password on valid length and do hash
func IsPasswordValid(password *string) bool {
	if len(*password) < 8 || len(*password) > 128 {
		return false
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(*password), 8)
	*password = string(hashed)
	return true
}

// Check username on valid
func IsUsernameValid(username string) bool {
	if len(username) < 3 || len(username) > 20 {
		return true
	}

	return regexp.MustCompile(`\s`).MatchString(username)
}
