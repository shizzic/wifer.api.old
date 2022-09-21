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
var domainClient string = "https://luckriza.com"

// var domainClient string = "http://localhost:8080"
var domainBack string = "wifer-test.ru"
var emailPass string = "jukdNRaVWf3Fvmg"

var ctx = context.TODO()
var r = gin.Default()      // https server (MAIN)
var router = gin.Default() // http server
var con = connect()        // database

var DB = map[string]*mongo.Collection{
	"users":     con.Database("wifer").Collection("users"),
	"ensure":    con.Database("wifer").Collection("ensure"),
	"countries": con.Database("wifer").Collection("countries"),
	"cities":    con.Database("wifer").Collection("cities"),
	"templates": con.Database("wifer").Collection("templates"),
	"views":     con.Database("wifer").Collection("views"),
	"likes":     con.Database("wifer").Collection("likes"),
	"private":   con.Database("wifer").Collection("private"),
	"access":    con.Database("wifer").Collection("access"),
	"messages":  con.Database("wifer").Collection("messages"),
}

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
			c.Redirect(302, "https://"+domainBack+":450"+c.Request.URL.String())
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{domainClient},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "Access-Control-Max-Age"},
		AllowCredentials: true,
	})
}

func setHeaders() {
	r.SetTrustedProxies(nil)
	router.SetTrustedProxies(nil)
	router.Use(redirect()) // bind endless redirect for NONE https requests
	r.Use(CORSMiddleware())
	router.Use(CORSMiddleware())
	r.MaxMultipartMemory = 8 << 20
	router.MaxMultipartMemory = 8 << 20
}

// Run both: http and https servers
func run() {
	go r.RunTLS(serverID+":450", "/etc/ssl/luckriza/luckriza_com.full.crt", "/etc/ssl/luckriza/luckriza_com.key")
	router.Run(serverID + ":449")
}

func clearOnline() {
	DB["users"].UpdateMany(ctx, bson.M{"online": true},
		bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}},
	)
}
