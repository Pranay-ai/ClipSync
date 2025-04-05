package db

import (
	"log"

	"clipsync.com/m/config"
	"clipsync.com/m/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := config.GetDBConnectionString()
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB = database

	// Auto-migrate the models
	database.AutoMigrate(&models.User{})
}

var RedisClient *redis.Client

func ConnectRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // or your Redis container address
	})
}
