package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const letters = "1234567890_-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncryptToken(username string) (token string) {
	for i, char := range username {
		if char%2 == 0 {
			token += string(char - 1)
		} else {
			token += string(char + 1)
		}

		rand.Seed(time.Now().UnixNano())

		b := make([]byte, i)
		for i := range b {
			b[i] = letters[rand.Int63()%int64(len(letters))]
		}

		token += string(b)
	}

	return
}

// 30 ms speed average
func DecryptToken(token string) (username string) {
	key := 0
	minus := 0

	for i, char := range token {
		if key == i {
			if char%2 == 0 {
				username += string(char - 1)
			} else {
				username += string(char + 1)
			}

			key += minus + 1
			minus += 1
		}
	}

	return
}

// check if user loged in
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil {
			if username, e := c.Cookie("username"); e == nil {
				if DecryptToken(token) == username {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatus(401)
	}
}

// Login for Form from client side. Yea, i'm lil dick :D
func Login(email, password string, c gin.Context) error {
	var user bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1, "username": 1, "password_hash": 1, "status": 1, "active": 1})

	if err := users.FindOne(ctx, bson.M{"email": email}, opts).Decode(&user); err != nil {
		return errors.New("document not found")
	}

	// Verify password
	if err := ComparePassword(fmt.Sprint(user["password_hash"]), password); err != nil {
		return errors.New("wrong password")
	}

	// Check if user ever ensure his account or ever been deleted
	if user["active"] == false {
		return errors.New("email not ensure")
	} else if user["active"] == true && user["status"] == false {
		if r, err := users.UpdateOne(ctx, bson.M{"_id": user["_id"]}, bson.D{
			{Key: "$set", Value: bson.D{{Key: "status", Value: true}}},
		}); err != nil || r.ModifiedCount == 0 {
			return errors.New("account not activated")
		}
	}

	// create cookies
	var username string = fmt.Sprint(user["username"])
	c.SetCookie("token", EncryptToken(username), 86400*30, "/", "wifer-test.ru", true, true)
	c.SetCookie("username", username, 86400*30, "/", "wifer-test.ru", true, true)
	c.SetCookie("id", fmt.Sprint(user["_id"]), 86400*30, "/", "wifer-test.ru", true, true)

	return nil
}

// Compare password from client side and database side
func ComparePassword(hash, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New("wrong password")
	}

	return nil
}

// Make token for auth any email operations or something :)
func MakeToken() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 64)
	for i := range b {
		b[i] = letters[rand.Int63()%int64(len(letters))]
	}

	return string(b)
}
