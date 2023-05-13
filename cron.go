package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"go.mongodb.org/mongo-driver/bson"
)

func start_cron() {
	cron := gocron.NewScheduler(time.UTC)

	cron.Every(1).Day().At("03:00").Do(func() { database_dump() })
	cron.Every(1).Day().At("04:00").Do(func() { zip_images() })

	cron.Every(1).MonthLastDay().Do(func() {
		DB["users"].UpdateMany(ctx, bson.M{"last_time": bson.M{"$lt": time.Now().Unix() - 31536000}}, bson.D{{Key: "$set", Value: bson.D{{Key: "active", Value: false}}}})
	})

	cron.Every(1).Day().Do(func() {
		DB["users"].UpdateMany(ctx, bson.M{"premium": bson.M{"$lt": time.Now().Unix()}}, bson.D{{Key: "$set", Value: bson.D{{Key: "premium", Value: int64(0)}}}})
	})

	cron.StartAsync()
}
