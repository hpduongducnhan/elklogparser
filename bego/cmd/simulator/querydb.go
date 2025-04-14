package simulator

import (
	"nhandd/bego/internal/repository"

	"github.com/rs/zerolog/log"
)

func GetActiveCollectors() {
	if pgDb == nil {
		panic("Postgres DB is not initialized")
	}
	repo := repository.NewDbRepository(pgDb)
	activeConfigs, err := repo.GetActiveConfigsWithRelations()
	if err != nil {
		log.Error().Err(err).Msg("Error getting active configs")
		return
	}
	for _, config := range activeConfigs {
		log.Info().Msgf("Config ID: %d, Name: %s", config.ID, config.Name)
		log.Info().Msgf("Elk Config ID: %d, Name: %s:%d", config.ElkConfig.ID, config.ElkConfig.Host, config.ElkConfig.Port)
		for _, query := range config.ElkQueries {
			log.Info().Msgf("Query ID: %d, Name: %s", query.ID, query.Name)
		}
		for _, filter := range config.LogFilters {
			log.Info().Msgf("Filter ID: %d, Name: %s", filter.ID, filter.Name)
		}
	}

}
