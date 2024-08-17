package config

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var RDB *redis.Client

func InitRedis() {
	connStr := fmt.Sprintf("%s:%d", AppConfig.Redis.Host, AppConfig.Redis.Port)
	RDB = redis.NewClient(&redis.Options{
		Addr: connStr,
	})
	result := RDB.Ping(context.Background())

	if result.Val() != "PONG" {
		log.Fatalf("Failed to connect to redis")
	}

	log.Println("Redis connected successfully")
}
