package main

import "github.com/gin-gonic/gin"

func main() {
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

	r.GET("/profile", func(c *gin.Context) {
		var data user
		c.Bind(&data)

		if user, err := GetProfile(data.ID); err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, user)
		}
	})

	r.PUT("/online", func(c *gin.Context) {
		var data user
		c.Bind(&data)
		ChangeOnline(data.ID, data.Online)
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
	})

	r.POST("/upload", Auth(), func(c *gin.Context) {
		UploadImage(c.PostForm("dir"), *c)
		c.String(200, "nice")
	})

	r.PUT("/changeImageDir", Auth(), func(c *gin.Context) {
		ChangeImageDir(c.Query("isAvatar"), c.Query("dir"), c.Query("new"), c.Query("number"), *c)
		c.String(200, "nice")
	})

	r.PUT("/replaceAvatar", Auth(), func(c *gin.Context) {
		ReplaceAvatar(c.Query("dir"), c.Query("number"), *c)
		c.String(200, "nice")
	})

	r.DELETE("/deleteImage", Auth(), func(c *gin.Context) {
		DeleteImage(c.Query("isAvatar"), c.Query("dir"), c.Query("number"), *c)
		c.String(200, "nice")
	})

	r.GET("/test", func(c *gin.Context) {
		// cursor, _ := cities.Find(ctx, bson.M{"country_id": 231})
		// var data []bson.M
		// cursor.All(ctx, &data)
		// c.JSON(200, data)

		// if username, e := c.Cookie("username"); e != nil {
		// 	c.String(500, "error")
		// } else {
		// 	c.String(200, username)
		// }

		// if err := SendCode("kotcich@mail.ru", "123456"); err != nil {
		// 	c.String(500, "error")
		// } else {
		// 	c.String(200, "good")
		// }
	})

	run()
}
