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
		if err := EnsureEmail(c.PostForm("username"), c.PostForm("token"), c.PostForm("id"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account activated")
		}
	})

	// Ensure if user set HIS email or not
	r.DELETE("/ensureDelete", func(c *gin.Context) {
		if err := EnsureDelete(c.PostForm("id"), c.PostForm("token")); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account deleted")
		}
	})

	// Ensure if user set HIS email or not
	r.DELETE("/deleteAccount", Auth(), func(c *gin.Context) {
		if err := DeleteAccount(c.PostForm("password"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account deleted")
		}
	})

	// Just login
	r.POST("/login", func(c *gin.Context) {
		if err := Login(c.PostForm("email"), c.PostForm("password"), *c); err != nil {
			c.String(401, err.Error())
		} else {
			c.String(200, "loged in")
		}
	})

	// Password change
	r.PUT("/changePassword", Auth(), func(c *gin.Context) {
		if err := ChangePassword(c.PostForm("old"), c.PostForm("new"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "password changed")
		}
	})

	// Send link to new user's email
	r.POST("/sendChangeEmail", Auth(), func(c *gin.Context) {
		SendChangeEmail(c.PostForm("old"), c.PostForm("new"), *c)
	})

	// Email change
	r.PUT("/changeEmail", func(c *gin.Context) {
		if err := ChangeEmail(c.PostForm("id"), c.PostForm("token"), c.PostForm("new"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "email changed")
		}
	})

	r.GET("/test", Auth(), func(c *gin.Context) {
		c.String(200, "some")
	})

	run()
}
