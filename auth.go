package main

import (
	"errors"
	"math/rand"
	h "net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const nums = "1234567890"
const letters = "1234567890_-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type auth struct {
	ID   int    `form:"id"`
	Code string `form:"code"`
}

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

// Check code's fit for ensure
func CheckCode(id int, code string, c gin.Context) error {
	if !isCode(code) {
		return errors.New("0")
	}

	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	if err := ensure.FindOne(ctx, bson.M{"_id": id, "code": code}, opts).Decode(&data); err != nil {
		return errors.New("1")
	}

	// Delete document in ensure collection, if given code was valid
	ensure.DeleteOne(ctx, bson.M{"_id": id, "code": code})

	var user bson.M
	opt := options.FindOne().SetProjection(bson.M{"username": 1})

	if err := users.FindOne(ctx, bson.M{"_id": id}, opt).Decode(&user); err != nil {
		return errors.New("2")
	}

	users.UpdateOne(ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: true}}}})

	MakeCookies(strconv.Itoa(id), user["username"].(string), c)

	return nil
}

// Cookies for auth
func MakeCookies(id, username string, c gin.Context) {
	c.SetSameSite(h.SameSiteNoneMode)
	c.SetCookie("token", EncryptToken(username), 86400*120, "/", domainBack, true, true)
	c.SetCookie("username", username, 86400*120, "/", domainBack, true, true)
	c.SetCookie("id", id, 86400*120, "/", domainBack, true, true)
}

// Make token for auth any email operations or something :)
func MakeCode() string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, 6)
	for i := range b {
		b[i] = nums[rand.Int63()%int64(len(nums))]
	}

	return string(b)
}
