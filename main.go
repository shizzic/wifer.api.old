package main

import (
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

	r.GET("/test", Auth(), func(c *gin.Context) {
		c.String(200, "some")
	})

	r.GET("/token", func(c *gin.Context) {

	})

	run()
}
