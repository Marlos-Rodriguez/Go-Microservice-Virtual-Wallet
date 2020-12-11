package cache

import (
	"context"
	"log"
	"os"

	//Autoload the env
	"github.com/go-redis/redis/v8"
	_ "github.com/joho/godotenv/autoload"
)

//NewRedisClient return a new redis client
func NewRedisClient(service string) *redis.Client {
	addr := os.Getenv(service + "_REDIS_ADDR")

	if addr == "" {
		log.Fatalln("Error in Getting the ADDR from ENV")
		return nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := rdb.Ping(context.Background())

	if status.Err() != nil {
		return nil
	}

	log.Println("Redis Cache Connected")

	return rdb
}
