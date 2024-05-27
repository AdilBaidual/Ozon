package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(options *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(options)

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return client, nil
}
