package redis

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func client() *redis.Client {
	upstashUrl := os.Getenv("UPSTASH_URL")
	opt, _ := redis.ParseURL(upstashUrl)
	client := redis.NewClient(opt)

	return client
}

func Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return client().Set(ctx, key, value, expiration)
}

func Get(key string) *redis.StringCmd {
	return client().Get(ctx, key)
}
