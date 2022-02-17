package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()
var r = gin.Default()    // https server (MAIN)
var http = gin.Default() // http server
var con = connect()      // database
var users = con.Database("wifer").Collection("users")
var ensure = con.Database("wifer").Collection("ensure")

const uri string = "mongodb://shizzic:WebDev77@wifer-test.ru:27017/test?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"

type user struct {
	Id        string `form:"id"`
	Username  string `form:"username"`
	Email     string `form:"email"`
	Password  string `form:"password"`
	Title     string `form:"title"`
	Sex       uint8  `form:"sex"`
	Age       uint8  `form:"age"`
	Height    uint8  `form:"height"`
	Weight    uint8  `form:"weight"`
	Body      uint8  `form:"body"`
	Smokes    uint8  `form:"smokes"`
	Drinks    uint8  `form:"drinks"`
	Ethnicity uint8  `form:"ethnicity"`
	Search    uint8  `form:"search"`
	Income    uint8  `form:"income"`
	Children  uint8  `form:"children"`
	Industry  uint8  `form:"industry"`
}

func connect() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

// Redirect every NOT https request to https
func redirect() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Proto != "HTTP/2.0" {
			c.Redirect(302, "https://wifer-test.ru"+c.Request.URL.String())
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func setHeaders() {
	r.SetTrustedProxies(nil)
	http.SetTrustedProxies(nil)
	http.Use(redirect()) // bind endless redirect for NONE https requests
	r.Use(CORSMiddleware())
	http.Use(CORSMiddleware())
}

// Run both: http and https servers
func run() {
	go r.RunTLS("213.189.217.231:443", "/etc/ssl/wifer-test/__wifer-test_ru.full.crt", "/etc/ssl/wifer-test/__wifer-test_ru.key")
	http.Run("213.189.217.231:80")
}
