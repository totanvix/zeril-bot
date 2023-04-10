package redis

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

var upstashUrl = fmt.Sprintf("rediss://default:%s@dear-hound-30219.upstash.io:30219", os.Getenv("UPSTASH_PASSWORD"))
var ctx = context.Background()

func client() *redis.Client {
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
