package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

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
		return errors.New("something went wrong, try again later")
	}

	//inform user about success
	InfoAboutPasswordChange(fmt.Sprint(user["email"]), fmt.Sprint(user["username"]))

	return nil
}

// Delete account only from client side
func DeleteAccount(username, password string, c gin.Context) error {
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

	//inform user about success
	InfoAboutDelete(fmt.Sprint(user["email"]), fmt.Sprint(user["username"]))
	return nil
}
