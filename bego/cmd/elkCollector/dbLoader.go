package elkcollector

import (
	"sync"

	"github.com/rs/zerolog/log"
)

var reloadMutex sync.Mutex

func loadCollectorFromDatabases() error {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	if elkCollectors == nil {
		elkCollectors = make(map[string]*ElkCollector)
	}
	// load from databases
	activeConfigs, err := pgRepo.GetActiveConfigsWithRelations()
	if err != nil {
		log.Error().Err(err).Msg("Error getting active configs")
		return err
	} else {
		log.Info().Msgf("Loaded %d active configs", len(activeConfigs))
	}

	for _, config := range activeConfigs {
		elkCollectors[config.Code] = NewElkCollectorFromDbConf(config)
	}
	return nil
}
