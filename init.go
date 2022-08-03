package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
var countries = con.Database("wifer").Collection("countries")
var cities = con.Database("wifer").Collection("cities")
var templates = con.Database("wifer").Collection("templates")
var views = con.Database("wifer").Collection("views")
var likes = con.Database("wifer").Collection("likes")

const mongoConnect string = "mongodb://shizzic:WebDev77@wifer-test.ru:27017/test?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"

// mongodump "mongodb://shizzic:WebDev77@wifer-test.ru:27017/wifer?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false" -d wifer -o /var/www/default/site

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
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "Access-Control-Max-Age"},
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

func clearOnline() {
	users.UpdateMany(ctx, bson.M{"online": true},
		bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}},
	)
}
