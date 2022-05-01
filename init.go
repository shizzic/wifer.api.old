package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var serverID string = "213.189.217.231"
var domainClient string = "http://localhost:8080"
var domainBack string = "wifer-test.ru"
var emailPass string = "jukdNRaVWf3Fvmg"

var ctx = context.TODO()
var r = gin.Default()    // https server (MAIN)
var http = gin.Default() // http server
var con = connect()      // database
var users = con.Database("wifer").Collection("users")
var ensure = con.Database("wifer").Collection("ensure")

const mongoConnect string = "mongodb://shizzic:WebDev77@wifer-test.ru:27017/test?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"

// mongodump "mongodb://shizzic:WebDev77@wifer-test.ru:27017/wifer?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false" -d wifer -o /var/www/default/site

type user struct {
	Id        int    `form:"id"`
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
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoConnect))
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
			c.Redirect(302, "https://"+domainBack+c.Request.URL.String())
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{domainClient, "http://192.168.1.106:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	})
}

func setHeaders() {
	r.SetTrustedProxies(nil)
	http.SetTrustedProxies(nil)
	http.Use(redirect()) // bind endless redirect for NONE https requests
	r.Use(CORSMiddleware())
	http.Use(CORSMiddleware())
	r.MaxMultipartMemory = 8 << 20
	http.MaxMultipartMemory = 8 << 20
}

// Run both: http and https servers
func run() {
	go r.RunTLS(serverID+":443", "/etc/ssl/wifer-test/__wifer-test_ru.full.crt", "/etc/ssl/wifer-test/__wifer-test_ru.key")
	http.Run(serverID + ":80")
}
