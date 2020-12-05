package cache

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

//NewRedisClient return a new redis client
func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
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
