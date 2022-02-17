package main

import (
	"errors"
	"fmt"
	h "net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

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

// Change password only from client side
func ChangePassword(old, new string, c gin.Context) error {
	var user bson.M
	username, _ := c.Cookie("username")
	if err := users.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return errors.New("account not found")
	}

	if err := ComparePassword(fmt.Sprint(user["password_hash"]), old); err != nil {
		return errors.New("wrong password")
	}

	if len(new) < 8 || len(new) > 128 {
		return errors.New("invalid password length")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(new), 8)

	if r, err := users.UpdateOne(ctx, bson.M{"username": username}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "password_hash", Value: string(hashed)}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("error")
	}

	//inform user about success
	InfoAboutPasswordChange(fmt.Sprint(user["email"]), fmt.Sprint(user["username"]))

	return nil
}

// Delete account only from client side
func DeleteAccount(password string, c gin.Context) error {
	username, _ := c.Cookie("username")
	var user bson.M

	if err := users.FindOne(ctx, bson.M{"username": username}).Decode(&user); err != nil {
		return errors.New("account not deleted")
	}

	// Verify password
	if err := ComparePassword(fmt.Sprint(user["password_hash"]), password); err != nil {
		return errors.New("account not deleted")
	}

	// Smooth delete
	if r, err := users.UpdateOne(ctx, bson.M{"username": username}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "status", Value: false}}},
	}); err != nil || r.ModifiedCount == 0 {
		return errors.New("account not deleted")
	}

	// delete cookie
	c.SetCookie("token", "", -1, "/", "https://wifer-test.ru", true, true)
	c.SetCookie("username", "", -1, "/", "https://wifer-test.ru", true, true)
	c.SetCookie("id", "", -1, "/", "https://wifer-test.ru", true, true)

	//inform user about success
	InfoAboutDelete(fmt.Sprint(user["email"]), fmt.Sprint(user["username"]))
	return nil
}
