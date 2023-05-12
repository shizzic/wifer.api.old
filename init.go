package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Email struct {
	HOST     string
	USERNAME string
	PASSWORD string
	PORT     int
}

// env
type Config struct {
	MONGO_CONNECTION_STRING string
	SERVER_IP               string
	CLIENT_DOMAIN           string
	SELF_DOMAIN_NAME        string
	ADMIN_EMAIL             string
	EMAIL                   Email
}

var config = get_config()

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

// init is invoked before main()
func init() {
	clearOnline()
	setHeaders()
}

func get_config() *Config {
	// gin.SetMode(gin.ReleaseMode)

	if gin.Mode() == "release" {
		godotenv.Load(".env.production")
	} else {
		godotenv.Load(".env.development")
	}
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	return &Config{
		MONGO_CONNECTION_STRING: os.Getenv("MONGO_CONNECTION_STRING"),
		SERVER_IP:               os.Getenv("SERVER_IP"),
		CLIENT_DOMAIN:           os.Getenv("CLIENT_DOMAIN"),
		SELF_DOMAIN_NAME:        os.Getenv("SELF_DOMAIN_NAME"),
		ADMIN_EMAIL:             os.Getenv("ADMIN_EMAIL"),
		EMAIL: Email{
			HOST:     os.Getenv("EMAIL_HOST"),
			USERNAME: os.Getenv("EMAIL_USERNAME"),
			PASSWORD: os.Getenv("EMAIL_PASSWORD"),
			PORT:     port,
		},
	}
}

func connect() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MONGO_CONNECTION_STRING))
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
			c.Redirect(302, "https://"+config.SELF_DOMAIN_NAME+c.Request.URL.String())
			return
		}

		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{config.CLIENT_DOMAIN},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "origin", "Cache-Control", "X-Requested-With", "Access-Control-Max-Age"},
		AllowCredentials: true,
	})
}

func setHeaders() {
	r.SetTrustedProxies(nil)
	router.SetTrustedProxies(nil)
	r.Use(CORSMiddleware())
	router.Use(CORSMiddleware())
	r.MaxMultipartMemory = 8 << 20
	router.MaxMultipartMemory = 8 << 20
}

// Run both: http and https servers
func run() {
	if gin.Mode() == "release" {
		router.Use(redirect()) // bind endless redirect for NONE https requests
		go r.RunTLS(config.SERVER_IP+":443", "/etc/ssl/wifer/__wifer-test_ru.full.crt", "/etc/ssl/wifer/__wifer-test_ru.key")
		router.Run(config.SERVER_IP + ":80")
	} else {
		r.Run(config.SERVER_IP + ":80")
	}
}

func clearOnline() {
	DB["users"].UpdateMany(ctx, bson.M{"online": true},
		bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}},
	)
}
