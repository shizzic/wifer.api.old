package main

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type List struct {
	Limit       int64  `json:"limit"`
	Skip        int64  `json:"skip"`
	Sort        string `json:"sort"`
	AgeMin      int    `json:"ageMin"`
	AgeMax      int    `json:"ageMax"`
	ImagesMin   int    `json:"imagesMin"`
	ImagesMax   int    `json:"imagesMax"`
	HeightMin   int    `json:"heightMin"`
	HeightMax   int    `json:"heightMax"`
	WeightMin   int    `json:"weightMin"`
	WeightMax   int    `json:"weightMax"`
	ChildrenMin int    `json:"childrenMin"`
	ChildrenMax int    `json:"childrenMax"`
	Body        []int  `json:"body"`
	Sex         []int  `json:"sex"`
	Smokes      []int  `json:"smokes"`
	Drinks      []int  `json:"drinks"`
	Ethnicity   []int  `json:"ethnicity"`
	Search      []int  `json:"search"`
	Income      []int  `json:"income"`
	Industry    []int  `json:"industry"`
	Premium     []int  `json:"premium"`
	Prefer      []int  `json:"prefer"`
	Country     []int  `json:"country"`
	City        []int  `json:"city"`
	Text        string `json:"text"`
	IsAbout     bool   `json:"is_about"`
	Avatar      bool   `json:"avatar"`
	Count       bool   `json:"count"`
}

// Fewer 40ms :D
func GetUsers(data List, filter bson.M) []bson.M {
	var list []bson.M

	opts := options.Find().SetProjection(bson.M{
		"username":   1,
		"title":      1,
		"age":        1,
		"weight":     1,
		"height":     1,
		"body":       1,
		"ethnicity":  1,
		"public":     1,
		"private":    1,
		"avatar":     1,
		"premium":    1,
		"country_id": 1,
		"city_id":    1,
		"online":     1,
		"is_about":   1,
	}).
		SetSort(bson.D{
			{Key: "premium", Value: -1},
			{Key: data.Sort, Value: -1},
			{Key: "_id", Value: 1},
		}).
		SetLimit(data.Limit).
		SetSkip(data.Skip)

	cursor, _ := DB["users"].Find(ctx, filter, opts)
	cursor.All(ctx, &list)

	return list
}

func GetProfile(id int) (bson.M, error) {
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{
		"username":   1,
		"title":      1,
		"about":      1,
		"sex":        1,
		"age":        1,
		"body":       1,
		"height":     1,
		"weight":     1,
		"smokes":     1,
		"drinks":     1,
		"ethnicity":  1,
		"search":     1,
		"income":     1,
		"children":   1,
		"industry":   1,
		"premium":    1,
		"avatar":     1,
		"public":     1,
		"private":    1,
		"prefer":     1,
		"created_at": 1,
		"last_time":  1,
		"online":     1,
		"country_id": 1,
		"city_id":    1,
	})

	if err := DB["users"].FindOne(ctx, bson.M{"_id": id, "status": true}, opts).Decode(&user); err != nil {
		return user, errors.New("0")
	}

	return user, nil
}

func GetCountries() []bson.M {
	var data []bson.M
	cursor, _ := DB["countries"].Find(ctx, bson.M{})
	cursor.All(ctx, &data)

	return data
}

func GetCities(country_id int) []bson.M {
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 1, "title": 1})
	cursor, _ := DB["cities"].Find(ctx, bson.M{"country_id": country_id}, opts)
	cursor.All(ctx, &data)

	return data
}

func GetTemplates(c *gin.Context) bson.M {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	var text bson.M

	opts := options.FindOne().SetProjection(bson.M{"data": 1})
	DB["templates"].FindOne(ctx, bson.M{"_id": idInt}, opts).Decode(&text)

	return text
}

func CountUsers(filter bson.M) int64 {
	count, err := DB["users"].CountDocuments(ctx, filter)

	if err != nil {
		return 0
	} else {
		return count
	}
}

// Filter for users search
func PrepareFilter(data List) bson.M {
	filter := bson.M{
		"age":      bson.M{"$gte": data.AgeMin, "$lte": data.AgeMax},
		"height":   bson.M{"$gte": data.HeightMin, "$lte": data.HeightMax},
		"weight":   bson.M{"$gte": data.WeightMin, "$lte": data.WeightMax},
		"children": bson.M{"$gte": data.ChildrenMin, "$lte": data.ChildrenMax},
		"images":   bson.M{"$gte": data.ImagesMin, "$lte": data.ImagesMax},
	}

	if len(data.Body) > 0 {
		filter["body"] = bson.M{"$in": data.Body}
	}

	if len(data.Sex) > 0 {
		filter["sex"] = bson.M{"$in": data.Sex}
	}

	if len(data.Smokes) > 0 {
		filter["smokes"] = bson.M{"$in": data.Smokes}
	}

	if len(data.Drinks) > 0 {
		filter["drinks"] = bson.M{"$in": data.Drinks}
	}

	if len(data.Ethnicity) > 0 {
		filter["ethnicity"] = bson.M{"$in": data.Ethnicity}
	}

	if len(data.Search) > 0 {
		filter["search"] = bson.M{"$in": data.Search}
	}

	if len(data.Income) > 0 {
		filter["income"] = bson.M{"$in": data.Income}
	}

	if len(data.Industry) > 0 {
		filter["industry"] = bson.M{"$in": data.Industry}
	}

	if len(data.Premium) > 0 {
		filter["premium"] = bson.M{"$in": data.Premium}
	}

	if len(data.Country) > 0 {
		filter["country_id"] = bson.M{"$in": data.Country}
	}

	if len(data.City) > 0 {
		filter["city_id"] = bson.M{"$in": data.City}
	}

	if data.IsAbout {
		filter["is_about"] = true
	}

	if data.Avatar {
		filter["avatar"] = true
	}

	if data.Text != "" {
		filter["$text"] = bson.M{"$search": data.Text}
	}

	filter["status"] = true
	filter["active"] = true

	return filter
}

func CountAll() int64 {
	count, err := DB["users"].CountDocuments(ctx, bson.M{"status": true})

	if err != nil {
		return 0
	} else {
		return count
	}
}

func Visit() {
	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"count": 1})
	DB["visits"].FindOne(ctx, bson.M{"_id": 1}, opts).Decode(&data)

	DB["visits"].UpdateOne(ctx, bson.M{"_id": 1}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "count", Value: data["count"].(int32) + 1}}},
	})
}

// Getter for all kinda files in project
func Get_file(c *gin.Context) (string, error) {
	path := config.PATH + "/" + c.Query("what") + "/" + c.Query("target_id") + "/" + c.Query("dir") + "/" + c.Query("filename")

	if c.Query("dir") == "private" {
		if id, err := c.Cookie("id"); err == nil {
			if id == c.Query("target_id") {
				return path, nil
			}

			idInt, _ := strconv.Atoi(id)
			user_id, _ := strconv.Atoi(c.Query("target_id"))
			found_accesses, err := DB["private"].CountDocuments(ctx, bson.M{"user": user_id, "target": idInt})

			if err == nil && found_accesses != 0 {
				return path, nil
			}
		}
	} else {
		return path, nil
	}

	return "", errors.New("0")
}
