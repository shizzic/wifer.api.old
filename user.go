package main

import (
	"errors"
	h "net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	ID         int    `form:"id"`
	Username   string `form:"username"`
	Email      string `form:"email"`
	Title      string `form:"title"`
	About      string `form:"about"`
	Country_id int    `json:"country_id"`
	City_id    int    `json:"city_id"`
	Country    int    `form:"country_id"`
	City       int    `form:"city_id"`
	Sex        int    `form:"sex"`
	Age        int    `form:"age"`
	Height     int    `form:"height"`
	Weight     int    `form:"weight"`
	Body       int    `form:"body"`
	Smokes     int    `form:"smokes"`
	Drinks     int    `form:"drinks"`
	Ethnicity  int    `form:"ethnicity"`
	Search     []int  `form:"search"`
	Prefer     int    `form:"prefer"`
	Income     int    `form:"income"`
	Children   int    `form:"children"`
	Industry   int    `form:"industry"`
	Online     bool   `form:"online"`
	Premium    int    `json:"premium"`
}

func Change(data user, c gin.Context) error {
	if !IsUsernameValid(data.Username) {
		return errors.New("0")
	}
	if !IsAboutValid(data.About) {
		return errors.New("0")
	}
	if !IsTitleValid(data.Title) {
		return errors.New("0")
	}
	if !IsSexValid(data.Sex) {
		return errors.New("0")
	}
	if !IsAgeValid(data.Age) {
		return errors.New("0")
	}
	if !IsHeightValid(data.Height) {
		return errors.New("0")
	}
	if !IsWeightValid(data.Weight) {
		return errors.New("0")
	}
	if !IsSmokeValid(data.Smokes) {
		return errors.New("0")
	}
	if !IsDrinkValid(data.Drinks) {
		return errors.New("0")
	}
	if !IsEthnicityValid(data.Ethnicity) {
		return errors.New("0")
	}
	if !IsBodyValid(data.Body) {
		return errors.New("0")
	}
	if !IsIncomeValid(data.Income) {
		return errors.New("0")
	}
	if !IsIndustryValid(data.Industry) {
		return errors.New("0")
	}
	if !IsPreferValid(data.Prefer) {
		return errors.New("0")
	}
	if !IsChildrenValid(data.Children) {
		return errors.New("0")
	}

	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)
	username := strings.TrimSpace(data.Username)
	cookUsername, _ := c.Cookie("username")
	title := strings.TrimSpace(data.Title)
	about := strings.TrimSpace(data.About)
	var isAbout bool

	if cookUsername != username {
		if available := CheckUsernameAvailable(username); !available {
			return errors.New("1")
		}
	}

	if len(about) == 0 {
		isAbout = false
	} else {
		isAbout = true
	}

	if _, err := users.UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "username", Value: username}}},
		{Key: "$set", Value: bson.D{{Key: "title", Value: title}}},
		{Key: "$set", Value: bson.D{{Key: "about", Value: about}}},
		{Key: "$set", Value: bson.D{{Key: "is_about", Value: isAbout}}},
		{Key: "$set", Value: bson.D{{Key: "sex", Value: data.Sex}}},
		{Key: "$set", Value: bson.D{{Key: "age", Value: data.Age}}},
		{Key: "$set", Value: bson.D{{Key: "body", Value: data.Body}}},
		{Key: "$set", Value: bson.D{{Key: "weight", Value: data.Weight}}},
		{Key: "$set", Value: bson.D{{Key: "height", Value: data.Height}}},
		{Key: "$set", Value: bson.D{{Key: "smokes", Value: data.Smokes}}},
		{Key: "$set", Value: bson.D{{Key: "drinks", Value: data.Drinks}}},
		{Key: "$set", Value: bson.D{{Key: "ethnicity", Value: data.Ethnicity}}},
		{Key: "$set", Value: bson.D{{Key: "search", Value: data.Search}}},
		{Key: "$set", Value: bson.D{{Key: "prefer", Value: data.Prefer}}},
		{Key: "$set", Value: bson.D{{Key: "income", Value: data.Income}}},
		{Key: "$set", Value: bson.D{{Key: "children", Value: data.Children}}},
		{Key: "$set", Value: bson.D{{Key: "industry", Value: data.Industry}}},
		{Key: "$set", Value: bson.D{{Key: "country_id", Value: data.Country_id}}},
		{Key: "$set", Value: bson.D{{Key: "city_id", Value: data.City_id}}},
	}); err != nil {
		return errors.New("2")
	}

	c.SetSameSite(h.SameSiteNoneMode)
	c.SetCookie("username", username, 86400*120, "/", "*."+domainBack, true, true)
	c.SetCookie("token", EncryptToken(username), 86400*120, "/", "*."+domainBack, true, true)

	return nil
}

// Change when user open my website last time
func ChangeOnline(value bool, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	users.UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: value}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})
}

func CheckUsernameAvailable(username string) bool {
	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"username": 1})
	if err := users.FindOne(ctx, bson.M{"username": username}, opts).Decode(&data); err == nil {
		return false
	}

	return true
}

func CreateTemplates(text string, c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	templates.DeleteOne(ctx, bson.M{"_id": idInt})
	templates.InsertOne(ctx, bson.D{
		{Key: "_id", Value: idInt},
		{Key: "data", Value: text},
	})
}

func Logout(c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	users.UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})

	MakeCookies("", "", -1, c)
}

// Set active filed to false and handle user's content
func DeactivateAccount(c gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	users.UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "active", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "online", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})

	MakeCookies("", "", -1, c)
}

// Delete account forever
// func DeleteAccount(password string, c gin.Context) error {
// 	id, _ := c.Cookie("id")
// 	username, _ := c.Cookie("username")
// 	var user bson.M
// 	opts := options.FindOne().SetProjection(bson.M{"password_hash": 1})

// 	if err := users.FindOne(ctx, bson.M{"username": username}, opts).Decode(&user); err != nil {
// 		return errors.New("account not deleted")
// 	}

// 	if r, err := users.DeleteOne(ctx, bson.M{"_id": id, "username": username}); err != nil || r.DeletedCount == 0 {
// 		return errors.New("account not deleted")
// 	}

// 	// delete cookie
// 	c.SetCookie("token", "", -1, "/", "wifer-test.ru", true, true)
// 	c.SetCookie("username", "", -1, "/", "wifer-test.ru", true, true)
// 	c.SetCookie("id", "", -1, "/", "wifer-test.ru", true, true)

// 	return nil
// }
