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

func SubscribeChannel(
	terminatedCtx context.Context,
	channel string,
	handler func(pubsub *redis.PubSub, pubsubChan string, msg *redis.Message) error,
) error {
	// should call within goroutine
	if rclient == nil {
		log.Panic().Msg("Redis client is not initialized")
		return nil
	}
	pubsub := rclient.Subscribe(context.Background(), channel)
	if err := pubsub.Ping(context.Background()); err != nil {
		log.Panic().Err(err).Msg("Failed to ping Redis PubSub")
		return nil
	}
	defer func() {
		if err := pubsub.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Redis Pub/Sub")
		}
		log.Info().Str("channel", channel).Msg("Redis subscriber closed")
	}()

	for {
		select {
		case <-terminatedCtx.Done():
			log.Info().Str("channel", channel).Msg("Redis subscriber terminated")
			return nil
		case msg, ok := <-pubsub.Channel():
			if !ok {
				log.Warn().Str("channel", channel).Msg("Redis Pub/Sub channel closed")
				return nil
			}
			err := handler(pubsub, channel, msg)
			if err != nil {
				log.Error().Err(err).Str("channel", channel).Interface("msg", msg).Msg("Failed to handle Redis Pub/Sub message")
				continue
			}
			// log.Info().Str("channel", channel).Interface("msg", msg).Msg("Received message from Redis Pub/Sub")
		}
	}
}
