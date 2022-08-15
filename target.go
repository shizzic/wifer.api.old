package main

import (
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type target struct {
	Target int    `form:"target"`
	Which  int    `form:"which"`
	Skip   int64  `form:"skip"`
	Limit  int64  `form:"limit"`
	Mode   bool   `form:"mode"`
	Count  bool   `form:"count"`
	Text   string `form:"text"`
}

type Target struct {
	Like    bson.M
	Private []bson.M
	Access  []bson.M
}

// _________________________GET_______________________________

// Compilation of all functions for target
func GetTarget(target int, c gin.Context) Target {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	var data Target

	if idInt > 0 && idInt != target && target > 0 {
		AddView(idInt, target, c)

		if err, like := GetLike(idInt, target, c); err == false {
			data.Like = like
		}

		if err, priv := GetPrivate(idInt, target, c); err == false {
			data.Private = priv
		}

		if err, access := GetAccess(idInt, target, c); err == false {
			data.Access = access
		}

		return data
	}

	return data
}

// Get like in profile
func GetLike(id, target int, c gin.Context) (bool, bson.M) {
	var like bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "text": 1})

	if err := DB["likes"].FindOne(ctx, bson.M{"user": id, "target": target}, opts).Decode(&like); err == nil {
		return false, like
	} else {
		return true, like
	}
}

func GetAccess(id, target int, c gin.Context) (bool, []bson.M) {
	arr := [2]int{}
	arr[0] = id
	arr[1] = target
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1})

	if cursor, err := DB["access"].Find(ctx, bson.M{"user": bson.M{"$in": arr}, "target": bson.M{"$in": arr}}, opts); err == nil {
		if e := cursor.All(ctx, &data); e == nil {
			return false, data
		} else {
			return true, data
		}
	}

	return true, data
}

// Get access for private images in profile
func GetPrivate(id, target int, c gin.Context) (bool, []bson.M) {
	arr := [2]int{}
	arr[0] = id
	arr[1] = target
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1})

	if cursor, err := DB["private"].Find(ctx, bson.M{"user": bson.M{"$in": arr}, "target": bson.M{"$in": arr}}, opts); err == nil {
		if e := cursor.All(ctx, &data); e == nil {
			return false, data
		} else {
			return true, data
		}
	}

	return true, data
}

// Get quantity of unseen notifications by user
func GetNotifications(c gin.Context) map[string]int64 {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	data := make(map[string]int64)

	iLikes, err := DB["likes"].CountDocuments(ctx, bson.M{"target": idInt, "viewed": false})
	if err == nil {
		data["likes"] = iLikes
	}

	iViews, err := DB["views"].CountDocuments(ctx, bson.M{"target": idInt, "viewed": false})
	if err == nil {
		data["views"] = iViews
	}

	iPrivates, err := DB["private"].CountDocuments(ctx, bson.M{"target": idInt, "viewed": false})
	if err == nil {
		data["privates"] = iPrivates
	}

	iAccesses, err := DB["access"].CountDocuments(ctx, bson.M{"target": idInt, "viewed": false})
	if err == nil {
		data["accesses"] = iAccesses
	}

	return data
}

// ___________________________________________________________

// _________________________ADD_______________________________

// Add view of another user's profile
func AddView(id, target int, c gin.Context) {
	DB["views"].DeleteOne(ctx, bson.M{"user": id, "target": target})

	date := time.Now().Unix()
	DB["views"].InsertOne(ctx, bson.D{
		{Key: "user", Value: id},
		{Key: "target", Value: target},
		{Key: "viewed", Value: false},
		{Key: "created_at", Value: date},
	})
}

// User likes another user
func AddPrivate(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt > 0 && idInt != target && target > 0 {
		date := time.Now().Unix()
		DB["private"].InsertOne(ctx, bson.D{
			{Key: "user", Value: idInt},
			{Key: "target", Value: target},
			{Key: "viewed", Value: false},
			{Key: "created_at", Value: date},
		})
	}
}

// User gives an access for chating
func AddAccess(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt > 0 && idInt != target && target > 0 {
		date := time.Now().Unix()
		DB["access"].InsertOne(ctx, bson.D{
			{Key: "user", Value: idInt},
			{Key: "target", Value: target},
			{Key: "viewed", Value: false},
			{Key: "created_at", Value: date},
		})
	}
}

// User likes another user
func AddLike(data target, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt > 0 && idInt != data.Target && data.Target > 0 {
		date := time.Now().Unix()
		viewed := false
		var like bson.M
		text := strings.TrimSpace(data.Text)

		opts := options.FindOne().SetProjection(bson.M{"_id": 0, "viewed": 1})
		if err := DB["likes"].FindOne(ctx, bson.M{"user": idInt, "target": data.Target}, opts).Decode(&like); err == nil {
			viewed = like["viewed"].(bool)

			DB["likes"].UpdateOne(ctx, bson.M{"user": idInt, "target": data.Target}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "text", Value: text}}},
				{Key: "$set", Value: bson.D{{Key: "viewed", Value: viewed}}},
				{Key: "$set", Value: bson.D{{Key: "created_at", Value: date}}},
			})
		} else {
			DB["likes"].InsertOne(ctx, bson.D{
				{Key: "user", Value: idInt},
				{Key: "target", Value: data.Target},
				{Key: "text", Value: text},
				{Key: "viewed", Value: viewed},
				{Key: "created_at", Value: date},
			})
		}
	}
}

// ___________________________________________________________

// _________________________DELETE_______________________________

func DeleteLike(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt > 0 && idInt != target && target > 0 {
		DB["likes"].DeleteOne(ctx, bson.M{"user": idInt, "target": target})
	}
}

func DeletePrivate(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt > 0 && idInt != target && target > 0 {
		DB["private"].DeleteOne(ctx, bson.M{"user": idInt, "target": target})
	}
}

func DeleteAccess(target int, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if idInt > 0 && idInt != target && target > 0 {
		DB["access"].DeleteOne(ctx, bson.M{"user": idInt, "target": target})
	}
}

// ___________________________________________________________

func GetTargets(data target, c gin.Context) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	if data.Which == 0 {
		q, res = GetViews(idInt, data)
	}

	if data.Which == 1 {
		q, res = GetLikes(idInt, data)
	}

	if data.Which == 2 {
		q, res = GetPrivates(idInt, data)
	}

	if data.Which == 3 {
		q, res = GetAccesses(idInt, data)
	}

	return q, res
}

func GetViews(id int, data target) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := DB["views"].CountDocuments(ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := DB["views"].Find(ctx, filter, opts1)
	cursor.All(ctx, &targets)
	ids = RetrieveTargets(targets, key)

	if data.Mode {
		DB["views"].UpdateMany(ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		DB["views"].UpdateMany(ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})
	cur, _ := DB["users"].Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

func GetLikes(id int, data target) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		projection["text"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := DB["likes"].CountDocuments(ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := DB["likes"].Find(ctx, filter, opts1)
	cursor.All(ctx, &targets)
	ids = RetrieveTargets(targets, key)

	if data.Mode {
		DB["likes"].UpdateMany(ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		DB["likes"].UpdateMany(ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})

	cur, _ := DB["users"].Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

func GetPrivates(id int, data target) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := DB["private"].CountDocuments(ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := DB["private"].Find(ctx, filter, opts1)
	cursor.All(ctx, &targets)
	ids = RetrieveTargets(targets, key)

	if data.Mode {
		DB["private"].UpdateMany(ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		DB["private"].UpdateMany(ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})

	cur, _ := DB["users"].Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

func GetAccesses(id int, data target) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := DB["access"].CountDocuments(ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := DB["access"].Find(ctx, filter, opts1)
	cursor.All(ctx, &targets)
	ids = RetrieveTargets(targets, key)

	if data.Mode {
		DB["access"].UpdateMany(ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		DB["access"].UpdateMany(ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})

	cur, _ := DB["users"].Find(ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

// Get clean array of ints
func RetrieveTargets(data []bson.M, key string) []int32 {
	res := []int32{}
	for _, value := range data {
		res = append(res, value[key].(int32))
	}
	return res
}
