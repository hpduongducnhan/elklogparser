package elkcollector

import (
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func requestReloadCollectorConfigHandler(rpubsub *redis.PubSub, psChannel string, psMessage *redis.Message) error {
	log.Info().Msgf("Received message on channel %s: %s", psChannel, psMessage.Payload)
	err := loadCollectorFromDatabases()
	if err != nil {
		log.Error().Err(err).Msg("Error reloading collector config")
		return err
	}
	return nil
}
