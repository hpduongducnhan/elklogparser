package redisclient

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/redis/go-redis/v9"
)

var rclient *redis.Client

func pingRedis() error {
	_, err := rclient.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisClient(url string) (*redis.Client, error) {
	if rclient == nil {
		opts, err := redis.ParseURL(url)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to parse Redis URL")
			return nil, err
		}
		rclient = redis.NewClient(opts)
		if err := pingRedis(); err != nil {
			log.Panic().Err(err).Msg("Failed to ping Redis")
			return nil, err
		}
	}
	return rclient, nil
}
