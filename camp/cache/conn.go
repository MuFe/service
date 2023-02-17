package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"os"
	"strconv"
)

var client *redis.Client

func init() {
	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic(err)
	}
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	_, err = client.Ping().Result()
	if err != nil {
		panic(err)
	}
}
