package main

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type target struct {
	Target int `form:"target"`
}

// Compilation of all functions for target
func GetTarget(target int, c gin.Context) map[string]bson.M {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt != target && target != 0 {
		data := make(map[string]bson.M)

		AddView(idInt, target, c)

		if err, like := GetLike(idInt, target, c); err == false {
			data["like"] = like
		}

		return data
	}

	return nil
}

// Get like in profile
func GetLike(id, target int, c gin.Context) (bool, bson.M) {
	var like bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "text": 1})

	if err := likes.FindOne(ctx, bson.M{"user": id, "target": target}, opts).Decode(&like); err == nil {
		return false, like
	} else {
		return true, like
	}
}

// Add view of another user's profile
func AddView(id, target int, c gin.Context) {
	var view bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	if err := views.FindOne(ctx, bson.M{"user": id, "target": target}, opts).Decode(&view); err != nil {
		date := time.Now().Unix()
		views.InsertOne(ctx, bson.D{
			{Key: "user", Value: id},
			{Key: "target", Value: target},
			{Key: "created_at", Value: date},
		})
	}
}

// User likes another user
func AddLike(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt != target && target != 0 {
		date := time.Now().Unix()
		likes.InsertOne(ctx, bson.D{
			{Key: "user", Value: idInt},
			{Key: "target", Value: target},
			{Key: "text", Value: ""},
			{Key: "created_at", Value: date},
		})
	}
}

// User deletes his like
func DeleteLike(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt != target && target != 0 {
		likes.DeleteOne(ctx, bson.M{"user": idInt, "target": target})
	}
}
