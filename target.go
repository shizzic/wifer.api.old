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

type Target struct {
	Like    bson.M
	Private []bson.M
}

// Compilation of all functions for target
func GetTarget(target int, c gin.Context) Target {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	var data Target

	if idInt != target && target != 0 {
		AddView(idInt, target, c)

		if err, like := GetLike(idInt, target, c); err == false {
			data.Like = like
		}

		if err, priv := GetPrivate(idInt, target, c); err == false {
			data.Private = priv
		}

		return data
	}

	return data
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

// Get access for private images in profile
func GetPrivate(id, target int, c gin.Context) (bool, []bson.M) {
	arr := [2]int{}
	arr[0] = id
	arr[1] = target
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1})

	if cursor, err := private.Find(ctx, bson.M{"user": bson.M{"$in": arr}, "target": bson.M{"$in": arr}}, opts); err == nil {
		if e := cursor.All(ctx, &data); e == nil {
			return false, data
		} else {
			return true, data
		}
	}

	return true, data
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
			{Key: "viewed", Value: false},
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
			{Key: "viewed", Value: false},
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

// User likes another user
func AddPrivate(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt != target && target != 0 {
		date := time.Now().Unix()
		private.InsertOne(ctx, bson.D{
			{Key: "user", Value: idInt},
			{Key: "target", Value: target},
			{Key: "viewed", Value: false},
			{Key: "created_at", Value: date},
		})
	}
}

// User deletes his like
func DeletePrivate(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt != target && target != 0 {
		private.DeleteOne(ctx, bson.M{"user": idInt, "target": target})
	}
}
