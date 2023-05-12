package main

import (
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r.POST("/signin", func(c *gin.Context) {
		var data signin
		c.Bind(&data)

		var err error
		var id int

		if data.Api {
			id, err = CheckApi(data, c)
		} else {
			id, err = Signin(data.Email, c, false)
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

		if err := CheckCode(data.ID, data.Code, c); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"id": data.ID})
		}
	})

	r.PUT("/logout", Auth(), func(c *gin.Context) {
		Logout(c)
	})

	r.GET("/profile", func(c *gin.Context) {
		var data user
		c.Bind(&data)

		target := GetTarget(data.ID, c)
		if user, err := GetProfile(data.ID); err != nil {
			c.JSON(404, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"user": user, "target": target})
		}
	})

	r.GET("/online", Auth(), func(c *gin.Context) {
		var data user
		c.Bind(&data)

		ChangeOnline(data.Online, c)
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

		if err := Change(data, c); err != nil {
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
		CreateTemplates(c.PostForm("text"), c)
	})

	r.GET("/templates", Auth(), func(c *gin.Context) {
		text := GetTemplates(c)
		c.JSON(200, text)
	})

	r.PUT("/deactivate", Auth(), func(c *gin.Context) {
		DeactivateAccount(c)
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

		AddLike(data, c)
	})

	r.POST("/private", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		AddPrivate(data.Target, c)
	})

	r.POST("/access", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		AddAccess(data.Target, c)
	})

	r.DELETE("/like", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		DeleteLike(data.Target, c)
	})

	r.DELETE("/private", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		DeletePrivate(data.Target, c)
	})

	r.DELETE("/access", Auth(), func(c *gin.Context) {
		var data target
		c.Bind(&data)

		DeleteAccess(data.Target, c)
	})

	r.POST("/targets", Auth(), func(c *gin.Context) {
		var data target
		c.BindJSON(&data)

		count, res := GetTargets(data, c)
		c.JSON(200, gin.H{
			"data":  res,
			"count": count,
		})
	})

	r.GET("/notifications", Auth(), func(c *gin.Context) {
		res := GetNotifications(c)
		c.JSON(200, res)
	})

	r.POST("/upload-image", Auth(), func(c *gin.Context) {
		var data File
		c.Bind(&data)
		data.ID, _ = c.Cookie("id")
		data.EntryPath, _ = filepath.Abs("./images/" + data.ID)

		if err := UploadImage(data, c); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{})
		}
	})

	r.PUT("/changeImageDir", Auth(), func(c *gin.Context) {
		var data File
		c.Bind(&data)
		data.ID, _ = c.Cookie("id")
		data.EntryPath, _ = filepath.Abs("./images/" + data.ID)

		ChangeImageDir(data)
		c.JSON(200, "chanched")
	})

	r.PUT("/replaceAvatar", Auth(), func(c *gin.Context) {
		var data File
		c.Bind(&data)
		data.ID, _ = c.Cookie("id")
		data.EntryPath, _ = filepath.Abs("./images/" + data.ID)

		ReplaceAvatar(data)
		c.JSON(200, "replaced")
	})

	r.DELETE("/deleteImage", Auth(), func(c *gin.Context) {
		var data File
		c.Bind(&data)
		data.ID, _ = c.Cookie("id")
		data.EntryPath, _ = filepath.Abs("./images/" + data.ID)
		err := DeleteImage(data)
		c.JSON(200, gin.H{
			"error": err,
		})
	})

	r.PUT("/translate", func(c *gin.Context) {
		text, err := Translate(c.PostForm("text"), c.PostForm("lang"))

		if err != nil {
			c.JSON(400, gin.H{"error": "0"})
		} else {
			c.JSON(200, gin.H{"text": text})
		}
	})

	r.GET("/chat", Auth(), func(c *gin.Context) {
		Chat(c.Writer, c.Request, c)
	})

	r.POST("/getRooms", Auth(), func(c *gin.Context) {
		var data rooms
		c.BindJSON(&data)

		res, ids := GetRooms(data, c)
		c.JSON(200, gin.H{
			"rooms": res["rooms"],
			"users": res["users"],
			"ids":   ids,
		})
	})

	r.POST("/checkOnlineInChat", Auth(), func(c *gin.Context) {
		var data rooms
		c.BindJSON(&data)

		users := CheckOnlineInChat(data)
		c.JSON(200, users)
	})

	r.GET("/getMessages", Auth(), func(c *gin.Context) {
		var data messages
		c.Bind(&data)

		res := GetMessages(data, c)
		c.JSON(200, res)
	})

	r.GET("/getParamsAfterLogin", Auth(), func(c *gin.Context) {
		user, messages := GetParamsAfterLogin(c)
		c.JSON(200, gin.H{
			"user":     user,
			"messages": messages,
		})
	})

	r.GET("/count", func(c *gin.Context) {
		quantity := CountAll()
		c.JSON(200, quantity)
	})

	r.POST("/contact", func(c *gin.Context) {
		err := ContactMe(c.PostForm("name"), c.PostForm("email"), c.PostForm("subject"), c.PostForm("message"))
		c.JSON(200, err)
	})

	r.PUT("/trial", Auth(), func(c *gin.Context) {
		expires, err := TriggerTrial(c)
		c.JSON(200, gin.H{
			"err":     err,
			"expires": expires,
		})
	})

	r.GET("/checkPremium", Auth(), func(c *gin.Context) {
		left := CheckPremium(c)
		c.JSON(200, left)
	})

	r.POST("/visit", func(c *gin.Context) {
		Visit()
	})

	r.POST("/checkPayment", Auth(), func(c *gin.Context) {
		premium, err := CheckPayment(c.PostForm("orderID"), c)

		if err != nil {
			c.JSON(200, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"premium": premium})
		}
	})

	r.GET("/file", func(c *gin.Context) {
		path, err := Get_file(c)

		if err == nil {
			c.File(path)
		} else {
			c.JSON(401, "You don't have an access to this file")
		}
	})

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test")
	})

	run()
}
