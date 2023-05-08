package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// env
var MONGO_CONNECTION_STRING string = load_env_first_time("MONGO_CONNECTION_STRING")
var serverID string = os.Getenv("SERVER_ID")
var CLIENT_DOMAIN string = os.Getenv("CLIENT_DOMAIN")
var SELF_DOMAIN_NAME string = os.Getenv("SELF_DOMAIN_NAME")
var EMAIL_PASSWORD string = os.Getenv("EMAIL_PASSWORD")

// Connections
var ctx = context.TODO()
var r = gin.Default()        // https server (MAIN)
var router = gin.Default()   // http server
var mongo_client = connect() // database

var DB = map[string]*mongo.Collection{
	"users":     mongo_client.Database("db").Collection("users"),
	"ensure":    mongo_client.Database("db").Collection("ensure"),
	"countries": mongo_client.Database("db").Collection("countries"),
	"cities":    mongo_client.Database("db").Collection("cities"),
	"templates": mongo_client.Database("db").Collection("templates"),
	"views":     mongo_client.Database("db").Collection("views"),
	"likes":     mongo_client.Database("db").Collection("likes"),
	"private":   mongo_client.Database("db").Collection("private"),
	"access":    mongo_client.Database("db").Collection("access"),
	"messages":  mongo_client.Database("db").Collection("messages"),
	"visits":    mongo_client.Database("db").Collection("visits"),
	"payments":  mongo_client.Database("db").Collection("payments"),
}

func load_env_first_time(key string) string {
	godotenv.Load(".env")
	return os.Getenv(key)
}

func connect() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_CONNECTION_STRING))
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
// func redirect() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		if c.Request.Proto != "HTTP/2.0" {
// 			c.Redirect(302, "https://"+domainBack+c.Request.URL.String())
// 			return
// 		}

// 		c.Next()
// 	}
// }

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{CLIENT_DOMAIN},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "Access-Control-Max-Age"},
		AllowCredentials: true,
	})
}

func setHeaders() {
	r.SetTrustedProxies(nil)
	router.SetTrustedProxies(nil)
	// router.Use(redirect()) // bind endless redirect for NONE https requests
	r.Use(CORSMiddleware())
	router.Use(CORSMiddleware())
	r.MaxMultipartMemory = 8 << 20
	router.MaxMultipartMemory = 8 << 20
}

// Run both: http and https servers
func run() {
	// gin.SetMode(gin.ReleaseMode)
	// go r.RunTLS(serverID+":443", "/etc/ssl/wifer/__wifer-test_ru.full.crt", "/etc/ssl/wifer/__wifer-test_ru.key")
	r.Run(serverID + ":80")
}

func clearOnline() {
	DB["users"].UpdateMany(ctx, bson.M{"online": true},
		bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}},
	)
}
