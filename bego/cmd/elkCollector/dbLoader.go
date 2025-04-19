package elkcollector

import (
	"sync"

	"github.com/rs/zerolog/log"
)

var reloadMutex sync.Mutex

func loadCollectorsFromDb() error {
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

func reloadCollectorConfigFromDb(collectorCodes []string) error {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	activeCollectorConf, err := pgRepo.GetActiveConfigByCodes(collectorCodes)
	if err != nil {
		log.Error().Err(err).Msg("Error getting active config")
		return err
	} else {
		log.Info().Msgf("Reloaded config for collector codes: %v", collectorCodes)
	}
	// update collector
	for _, conf := range activeCollectorConf {
		elkCollectors[conf.Code] = NewElkCollectorFromDbConf(conf)
	}
	return nil
}
