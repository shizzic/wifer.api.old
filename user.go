package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// Delete account only from client side
func DeleteAccount(username, password string, c gin.Context) error {
	var account bson.M
	if err := users.FindOne(ctx, bson.M{"username": username}).Decode(&account); err != nil {
		return errors.New("account not deleted")
	}

	// Verify password
	if err := ComparePassword(fmt.Sprint(account["password_hash"]), password); err != nil {
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
	InfoAboutDelete(fmt.Sprint(account["email"]), fmt.Sprint(account["username"]))
	return nil
}
