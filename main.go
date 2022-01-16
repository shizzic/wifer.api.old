package main

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	setHeaders()

	r.POST("/registration", func(c *gin.Context) {
		var data registrat
		c.Bind(&data)

		if err := Registration(data); err != nil {
			c.String(400, err.Error())
		} else {
			c.JSON(200, "inserted")
		}
	})

	// Ensure if user set HIS email or not
	r.PUT("/ensure", func(c *gin.Context) {
		if err := EnsureEmail(c.PostForm("username"), c.PostForm("token"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account activated")
		}
	})

	// Ensure if user set HIS email or not
	r.DELETE("/ensureDelete", func(c *gin.Context) {
		if err := EnsureDelete(c.PostForm("username"), c.PostForm("token")); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account deleted")
		}
	})

	// Ensure if user set HIS email or not
	r.DELETE("/deleteAccount", func(c *gin.Context) {
		if err := DeleteAccount(c.PostForm("username"), c.PostForm("password"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account deleted")
		}
	})

	r.GET("/test", func(c *gin.Context) {
		rand.Seed(time.Now().UnixNano())

		b := make([]byte, 64)
		for i := range b {
			b[i] = letters[rand.Int63()%int64(len(letters))]
		}

		c.String(200, string(b))
	})

	r.GET("/token", func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil {
			c.String(400, "token not found")
		} else {
			c.String(200, token)
		}
	})

	run()
}
