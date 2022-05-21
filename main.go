package main

import (
	"github.com/gin-gonic/gin"
	fb "github.com/huandu/facebook/v2"
)

func main() {
	setHeaders()

	r.POST("/signin", func(c *gin.Context) {
		var data user
		c.Bind(&data)

		if err := Signin(data); err != nil {
			c.JSON(400, err.Error())
		} else {
			c.JSON(200, "inserted")
		}

		// c.JSON(200, data)
	})

	r.POST("/checkCode", func(c *gin.Context) {
		var data user
		c.Bind(&data)

		if err := Signin(data); err != nil {
			c.JSON(400, err.Error())
		} else {
			c.JSON(200, "inserted")
		}
	})

	r.POST("/registration", func(c *gin.Context) {
		var data user
		c.Bind(&data)

		if err := Signin(data); err != nil {
			c.JSON(400, err.Error())
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
		if id, err := Login(c.PostForm("email"), c.PostForm("password"), *c); err != nil {
			c.JSON(401, gin.H{"error": err.Error()})
		} else {
			c.JSON(200, gin.H{"id": id})
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

		// payload, _ := idtoken.Validate(context.Background(), c.Query("secret"), c.Query("id"))
		// c.JSON(200, payload.Claims)

		// req := curl.NewRequest()
		// resp, _ := req.Get("https://graph.facebook.com/v9.0/me?access_token=EAAvouaj2GXcBANUuLL8kfhjJsLJgtZChiDFSlDSZBhkP9XIv4X2tM0NMlrZBZCgmUlmtnnpa92wRmNMZCf8XXTYRQdjdnJhtodsZBfP5Dtkhzfke0LVGvD8rk3lANOkLKYEP0lqOfaOlyQZAV0ZARzXsg2ZBicaugwoMvB9QZBafBZBZCzEcRcr13XdZCGEiCtJZCfZCg9YVJ7HafzxJPgZApBsSLOkJ&fields=email&method=get&pretty=0&sdk=joey&suppress_http_code=1")

		res, _ := fb.Get("/"+c.Query("id"), fb.Params{
			"fields":       "email",
			"access_token": c.Query("token"),
		})

		c.JSON(200, res)
	})

	run()
}
