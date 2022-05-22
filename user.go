package main

import (
	"errors"
	h "net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type user struct {
	Id        int    `form:"id"`
	Country   int    `form:"country"`
	City      int    `form:"city"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	Title     string `form:"title"`
	About     string `form:"about"`
	Sex       uint8  `form:"sex"`
	Age       uint8  `form:"age"`
	Height    uint8  `form:"height"`
	Weight    uint8  `form:"weight"`
	Body      uint8  `form:"body"`
	Smokes    uint8  `form:"smokes"`
	Drinks    uint8  `form:"drinks"`
	Ethnicity uint8  `form:"ethnicity"`
	Search    uint8  `form:"search"`
	Income    uint8  `form:"income"`
	Children  uint8  `form:"children"`
	Industry  uint8  `form:"industry"`
}

func ChangeAbout(text string, c gin.Context) error {
	id, _ := c.Cookie("id")
	len := len(text)

	if len > 19 && len < 1501 {
		if r, err := users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "about", Value: text}}},
		}); err != nil || r.ModifiedCount == 0 {
			return errors.New("about not updated")
		}
	} else {
		return errors.New("short text")
	}

	return nil
}

// Change small params for user like height, weight, annual income etc.
func ChangeParams(data user, c gin.Context) error {
	id, _ := c.Cookie("id")
	if r, err := users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "sex", Value: data.Sex}}},
		{Key: "$set", Value: bson.D{{Key: "age", Value: data.Age}}},
		{Key: "$set", Value: bson.D{{Key: "weight", Value: data.Weight}}},
		{Key: "$set", Value: bson.D{{Key: "height", Value: data.Height}}},
		{Key: "$set", Value: bson.D{{Key: "body", Value: data.Body}}},
		{Key: "$set", Value: bson.D{{Key: "smokes", Value: data.Smokes}}},
		{Key: "$set", Value: bson.D{{Key: "drinks", Value: data.Drinks}}},
		{Key: "$set", Value: bson.D{{Key: "search", Value: data.Search}}},
		{Key: "$set", Value: bson.D{{Key: "income", Value: data.Income}}},
		{Key: "$set", Value: bson.D{{Key: "children", Value: data.Children}}},
		{Key: "$set", Value: bson.D{{Key: "industry", Value: data.Industry}}},
		{Key: "$set", Value: bson.D{{Key: "ethnicity", Value: data.Ethnicity}}},
		{Key: "$set", Value: bson.D{{Key: "country_id", Value: data.Country}}},
		{Key: "$set", Value: bson.D{{Key: "city_id", Value: data.City}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("error")
	}

	return nil
}

func ChangeUsername(new string, c gin.Context) error {
	username, _ := c.Cookie("username")
	if IsUsernameValid(new) || new == username {
		return errors.New("invalid username")
	}

	if r, err := users.UpdateOne(ctx, bson.M{"username": username}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "username", Value: new}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("error")
	}

	c.SetSameSite(h.SameSiteNoneMode)
	c.SetCookie("username", new, 120, "/", "wifer-test.ru", true, true)

	return nil
}

func ChangeTitle(new string, c gin.Context) error {
	id, _ := c.Cookie("id")
	if r, err := users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "title", Value: new}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("error")
	}

	return nil
}

// Change user's email only from clien't side
func ChangeEmail(id, token, newEmail string, c gin.Context) error {
	if !IsEmailValid(newEmail) {
		return errors.New("invalid email")
	}

	if r, err := ensure.DeleteOne(ctx, bson.M{"_id": id, "token": token}); err != nil || r.DeletedCount == 0 {
		return errors.New("not found")
	}

	if r, err := users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "email", Value: newEmail}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("error")
	}

	return nil
}

// Delete account forever
func DeleteAccount(password string, c gin.Context) error {
	id, _ := c.Cookie("id")
	username, _ := c.Cookie("username")
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"password_hash": 1})

	if err := users.FindOne(ctx, bson.M{"username": username}, opts).Decode(&user); err != nil {
		return errors.New("account not deleted")
	}

	if r, err := users.DeleteOne(ctx, bson.M{"_id": id, "username": username}); err != nil || r.DeletedCount == 0 {
		return errors.New("account not deleted")
	}

	// delete cookie
	c.SetCookie("token", "", -1, "/", "wifer-test.ru", true, true)
	c.SetCookie("username", "", -1, "/", "wifer-test.ru", true, true)
	c.SetCookie("id", "", -1, "/", "wifer-test.ru", true, true)

	return nil
}

// Set status to false for user
func DiactivateAccount(password string, c gin.Context) error {
	id, _ := c.Cookie("id")
	username, _ := c.Cookie("username")
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"password_hash": 1})

	if err := users.FindOne(ctx, bson.M{"_id": id, "username": username}, opts).Decode(&user); err != nil {
		return errors.New("account not frozen")
	}

	if r, err := users.UpdateOne(ctx, bson.M{"_id": id, "username": username}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "status", Value: false}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("account not frozen")
	}

	// delete cookie
	c.SetCookie("token", "", -1, "/", "wifer-test.ru", true, true)
	c.SetCookie("username", "", -1, "/", "wifer-test.ru", true, true)
	c.SetCookie("id", "", -1, "/", "wifer-test.ru", true, true)

	return nil
}

func Google(id, secret string) {
	// "https://www.googleapis.com/auth/userinfo.email"

	// conf := &oauth2.Config{
	// 	ClientID:     id,
	// 	ClientSecret: secret,
	// 	RedirectURL:  "http://localhost:8080/search",
	// 	Scopes:       []string{"email"},
	// 	Endpoint:     google.Endpoint,
	// }

	// tok, _ := conf.Exchange(context.TODO(), "")
}
