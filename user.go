package main

import (
	"errors"
	net "net/http"
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
	Premium    int64  `json:"premium"`
}

func Change(data user, c *gin.Context) error {
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
	isAbout := true

	if cookUsername != username {
		if available := CheckUsernameAvailable(username); !available {
			return errors.New("1")
		}
	}

	if len(about) == 0 {
		isAbout = false
	}

	if _, err := DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
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

	MakeCookies(id, username, 86400*120, c)
	return nil
}

// Change when user open my website last time
func ChangeOnline(value bool, c *gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: value}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})
}

func CheckUsernameAvailable(username string) bool {
	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"username": 1})
	if err := DB["users"].FindOne(ctx, bson.M{"username": username}, opts).Decode(&data); err == nil {
		return false
	}

	return true
}

func CreateTemplates(text string, c *gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	DB["templates"].DeleteOne(ctx, bson.M{"_id": idInt})
	DB["templates"].InsertOne(ctx, bson.D{
		{Key: "_id", Value: idInt},
		{Key: "data", Value: text},
	})
}

func Logout(c *gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "online", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})

	MakeCookies("", "", -1, c)
	c.SetSameSite(net.SameSiteNoneMode)
	c.SetCookie("premium", "premium", -1, "/", "."+SELF_DOMAIN_NAME, true, true)
}

// Set active filed to false and handle user's content
func DeactivateAccount(c *gin.Context) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "active", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "online", Value: false}}},
		{Key: "$set", Value: bson.D{{Key: "last_time", Value: time.Now().Unix()}}},
	})

	MakeCookies("", "", -1, c)
	c.SetSameSite(net.SameSiteNoneMode)
	c.SetCookie("premium", "premium", -1, "/", "."+SELF_DOMAIN_NAME, true, true)
}

func GetParamsAfterLogin(c *gin.Context) (bson.M, []interface{}) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "username": 1, "avatar": 1, "trial": 1, "premium": 1})
	DB["users"].FindOne(ctx, bson.M{"_id": idInt}, opts).Decode(&user)

	if user["premium"].(int64) != 0 {
		if user["premium"].(int64) <= time.Now().Unix() {
			DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "premium", Value: 0}}},
			})

			user["premium"] = 0

			c.SetSameSite(net.SameSiteNoneMode)
			c.SetCookie("premium", "premium", -1, "/", "."+SELF_DOMAIN_NAME, true, true)
		} else {
			if _, err := c.Cookie("premium"); err != nil {
				c.SetSameSite(net.SameSiteNoneMode)
				c.SetCookie("premium", "premium", int(user["premium"].(int64)-time.Now().Unix()), "/", "."+SELF_DOMAIN_NAME, true, true)
			}
		}
	}

	newMessages, _ := DB["messages"].Distinct(ctx, "user", bson.M{"target": idInt, "viewed": false})

	return user, newMessages
}

func TriggerTrial(c *gin.Context) (int64, error) {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "trial": 1, "premium": 1})
	DB["users"].FindOne(ctx, bson.M{"_id": idInt}, opts).Decode(&user)

	if !user["trial"].(bool) {
		expires := time.Now().Unix() + int64(1*60*60*24*7)

		if user["premium"].(int64) == 0 {
			DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "trial", Value: true}}},
				{Key: "$set", Value: bson.D{{Key: "premium", Value: expires}}},
			})

			c.SetSameSite(net.SameSiteNoneMode)
			c.SetCookie("premium", "premium", 1*60*60*24*7, "/", "."+SELF_DOMAIN_NAME, true, true)

			return expires, nil
		} else {
			DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "trial", Value: true}}},
				{Key: "$set", Value: bson.D{{Key: "premium", Value: int64(user["premium"].(int64) + int64(1*60*60*24*7))}}},
			})

			c.SetSameSite(net.SameSiteNoneMode)
			c.SetCookie("premium", "premium", int(user["premium"].(int64)-time.Now().Unix()+(1*60*60*24*7)), "/", "."+SELF_DOMAIN_NAME, true, true)
		}
	}

	return 0, errors.New("0")
}

func CheckPremium(c *gin.Context) bool {
	id, _ := c.Cookie("id")
	idInt, _ := strconv.Atoi(id)

	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "premium": 1})
	DB["users"].FindOne(ctx, bson.M{"_id": idInt}, opts).Decode(&user)

	if user["premium"].(int64) < time.Now().Unix() {
		DB["users"].UpdateOne(ctx, bson.M{"_id": idInt}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "premium", Value: int64(0)}}},
		})

		c.SetSameSite(net.SameSiteNoneMode)
		c.SetCookie("premium", "premium", -1, "/", "."+SELF_DOMAIN_NAME, true, true)

		return false
	}

	return true
}
