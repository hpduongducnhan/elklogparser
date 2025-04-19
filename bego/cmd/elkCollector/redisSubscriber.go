package elkcollector

import (
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

type PubsubPayload struct {
	Codes []string `json:"codes"`
}

func requestReloadCollectorConfigHandler(rpubsub *redis.PubSub, psChannel string, psMessage *redis.Message) error {
	log.Info().Msgf("Received message on channel %s: %s", psChannel, psMessage.Payload)
	var loadedPayload PubsubPayload
	err := json.Unmarshal([]byte(psMessage.Payload), &loadedPayload)
	if err != nil || psMessage.Payload == "" || len(loadedPayload.Codes) == 0 {
		err = loadCollectorsFromDb()
		if err != nil {
			log.Error().Err(err).Msg("Error reloading collector config")
			return err
		}
	} else {
		log.Info().Msgf("Reloading collector config for codes: %v", loadedPayload.Codes)
		err = reloadCollectorConfigFromDb(loadedPayload.Codes)
		if err != nil {
			log.Error().Err(err).Msg("Error reloading collector config")
			return err
		}
	}
	return nil
}
