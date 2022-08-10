package main

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	clearOnline()
	setHeaders()

	r.POST("/signin", func(c *gin.Context) {
		var data signin
		c.Bind(&data)
		var err error
		var id int

		if data.Api == true {
			id, err = CheckApi(data, *c)
		} else {
			id, err = Signin(data.Email, *c, false)
		}

		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"id": id})
		}
	})

	r.POST("/checkCode", func(c *gin.Context) {
		var data auth
		c.Bind(&data)

		if err := CheckCode(data.ID, data.Code, *c); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"id": data.ID})
		}
	})

	r.PUT("/logout", Auth(), func(c *gin.Context) {
		Logout(*c)
	})

	r.GET("/profile", func(c *gin.Context) {
		var data user
		c.Bind(&data)

		target := GetTarget(data.ID, *c)
		if user, err := GetProfile(data.ID); err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"user": user, "target": target})
		}
	})

	r.POST("/online", Auth(), func(c *gin.Context) {
		var data user
		c.Bind(&data)

		ChangeOnline(data.Online, *c)
	})

	r.GET("/country", func(c *gin.Context) {
		countries := GetCountries()
		c.JSON(200, countries)
	})

	r.GET("/city", func(c *gin.Context) {
		var data user
		c.Bind(&data)
		cities := GetCities(data.Country)
		c.JSON(200, cities)
	})

	r.PUT("/change", Auth(), func(c *gin.Context) {
		var data user
		c.BindJSON(&data)

		if err := Change(data, *c); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{})
		}
	})

	r.GET("/checkUsername", Auth(), func(c *gin.Context) {
		username := strings.TrimSpace(c.Query("username"))

		if IsUsernameValid(username) {
			result := CheckUsernameAvailable(c.Query("username"))
			c.JSON(200, result)
		}
	})

	r.POST("/templates", Auth(), func(c *gin.Context) {
		CreateTemplates(c.PostForm("text"), *c)
	})

	r.GET("/templates", Auth(), func(c *gin.Context) {
		text := GetTemplates(*c)
		c.JSON(200, text)
	})

	r.PUT("/deactivate", Auth(), func(c *gin.Context) {
		DeactivateAccount(*c)
	})

	r.POST("/getUsers", func(c *gin.Context) {
		var data List
		c.BindJSON(&data)
		filter := PrepareFilter(data)

		if data.Count {
			c.JSON(200, gin.H{
				"users": GetUsers(data, filter),
				"count": CountUsers(filter),
			})
		} else {
			c.JSON(200, gin.H{"users": GetUsers(data, filter)})
		}
	})

	r.POST("/like", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		AddLike(data, *c)
	})

	r.DELETE("/like", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		DeleteLike(data.Target, *c)
	})

	r.POST("/private", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		AddPrivate(data.Target, *c)
	})

	r.DELETE("/private", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		DeletePrivate(data.Target, *c)
	})

	// r.POST("/note", Auth(), func(c *gin.Context) {
	// 	var data target
	// 	c.Bind(&data)

	// 	text := strings.TrimSpace(data.Text)
	// 	AddNote(data.Target, text, *c)
	// })

	r.POST("/targets", Auth(), func(c *gin.Context) {
		var data target
		c.BindJSON(&data)

		count, res := GetTargets(data, *c)
		c.JSON(200, gin.H{
			"data":  res,
			"count": count,
		})
	})

	r.GET("/notifications", Auth(), func(c *gin.Context) {
		res := GetNotifications(*c)
		c.JSON(200, res)
	})

	r.POST("/upload", Auth(), func(c *gin.Context) {
		if err := UploadImage(c.PostForm("dir"), *c); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{})
		}
	})

	r.PUT("/changeImageDir", Auth(), func(c *gin.Context) {
		ChangeImageDir(c.Query("isAvatar"), c.Query("dir"), c.Query("new"), c.Query("number"), *c)
		c.JSON(200, "chanched")
	})

	r.PUT("/replaceAvatar", Auth(), func(c *gin.Context) {
		ReplaceAvatar(c.Query("dir"), c.Query("number"), *c)
		c.JSON(200, "replaced")
	})

	r.DELETE("/deleteImage", Auth(), func(c *gin.Context) {
		DeleteImage(c.Query("isAvatar"), c.Query("dir"), c.Query("number"), *c)
		c.JSON(200, "deleted")
	})

	r.PUT("/translate", func(c *gin.Context) {
		text, err := Translate(c.PostForm("text"), c.PostForm("lang"))

		if err != nil {
			c.JSON(500, gin.H{"error": "0"})
		} else {
			c.JSON(200, gin.H{"text": text})
		}
	})

	run()
}
