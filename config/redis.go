package config

import (
	"context"
	"fmt"

	"boilerblade/helper"

	"github.com/redis/go-redis/v9"
)

func (e *Env) InitRedis() *redis.Client {
	addr := fmt.Sprintf("%s:%s", e.REDIS_HOST, e.REDIS_PORT)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: e.REDIS_PASSWORD,
		DB:       e.REDIS_DB,
	})

	// Test connection
	ctx := context.Background()
	_, err := client.Ping(ctx).Result()
	if err != nil {
		helper.LogError("Redis connection failed", err, addr, map[string]interface{}{
			"host":     e.REDIS_HOST,
			"port":     e.REDIS_PORT,
			"db":       e.REDIS_DB,
			"password": "***",
		})
		return client
	}

	helper.LogInfo("Redis connection initialized", map[string]interface{}{
		"host": e.REDIS_HOST,
		"port": e.REDIS_PORT,
		"db":   e.REDIS_DB,
		"addr": addr,
	})

	return client
}
