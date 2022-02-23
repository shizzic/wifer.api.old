package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	setHeaders()

	r.POST("/registration", func(c *gin.Context) {
		var data user
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

	// If user got this email and he understand that hasn't registrated on resurce, then he delete account with this email
	r.DELETE("/ensureDelete", func(c *gin.Context) {
		if err := EnsureDelete(c.PostForm("id"), c.PostForm("token")); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account deleted")
		}
	})

	// Delete all user's data forever
	r.DELETE("/deleteAccount", Auth(), func(c *gin.Context) {
		if err := DeleteAccount(c.PostForm("password"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account deleted")
		}
	})

	// Delete user's data smootly
	r.PUT("/deactivateAccount", Auth(), func(c *gin.Context) {
		if err := DiactivateAccount(c.PostForm("password"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "account frozen")
		}
	})

	// r.DELETE("/restoreAccount", Auth(), func(c *gin.Context) {
	// 	if err := DeleteAccount(c.PostForm("password"), *c); err != nil {
	// 		c.String(400, err.Error())
	// 	} else {
	// 		c.String(200, "account restored")
	// 	}
	// })

	r.POST("/login", func(c *gin.Context) {
		if err := Login(c.PostForm("email"), c.PostForm("password"), *c); err != nil {
			c.String(401, err.Error())
		} else {
			c.String(200, "loged in")
		}
	})

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

	r.PUT("/changeEmail", func(c *gin.Context) {
		if err := ChangeEmail(c.PostForm("id"), c.PostForm("token"), c.PostForm("new"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "email changed")
		}
	})

	r.PUT("/changeTitle", Auth(), func(c *gin.Context) {
		if err := ChangeTitle(c.PostForm("new"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "title changed")
		}
	})

	r.PUT("/changeUsername", Auth(), func(c *gin.Context) {
		if err := ChangeUsername(c.PostForm("new"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "username changed")
		}
	})

	r.PUT("/changeParams", Auth(), func(c *gin.Context) {
		var data user
		c.Bind(&data)

		if err := ChangeParams(data, *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "Params changed")
		}
	})

	r.PUT("/changeAbout", Auth(), func(c *gin.Context) {
		if err := ChangeAbout(c.PostForm("text"), *c); err != nil {
			c.String(400, err.Error())
		} else {
			c.String(200, "About changed")
		}
	})

	r.GET("/getUsers", Auth(), func(c *gin.Context) {
		var data List
		c.ShouldBindJSON(&data)
		c.JSON(200, GetUsers(data))
		// c.JSON(200, data)
	})

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test")

		// var user bson.M
		// opts := options.FindOne().SetProjection(bson.M{"username": 1})
		// users.FindOne(ctx, bson.M{"username": "kotcich"}, opts).Decode(&user)

		// c.JSON(200, user)
	})

	run()
}
