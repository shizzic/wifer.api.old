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
			c.String(500, err.Error())
		} else {
			c.JSON(200, "inserted")
		}
	})

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "test")
	})

	r.GET("/token", func(c *gin.Context) {
		c.String(200, DecryptToken())
		// c.String(200, EncryptToken("kotcich"))
	})

	run()
}
