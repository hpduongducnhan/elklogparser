package elkcollector

import (
	"github.com/rs/zerolog/log"
)

func runHandlers(hid int) {
	log.Info().Msgf("Starting handler %d...", hid+1)
	for {
		select {
		case <-terminationCtx.Done():
			log.Info().Msgf("Handler %d terminated", hid+1)
			return
		case val, ok := <-elkRespChan:
			if !ok {
				return
			}
			if val == nil {
				log.Warn().Msg("handler get nil value, ignore!")
				continue
			} else {
				err := elkRespHandler.Handle(val)
				if err != nil {
					log.Error().Err(err).Interface("val", val).Msgf("handler get error")
				}
			}
		}
	}
}

func runCollectors() {
	for _, collector := range elkCollectors {
		if collector != nil && collector.AllowToRun() {
			terminationWg.Add(1)
			go func() {
				log.Info().Str("code", collector.Code).Msg("Running collectors...")
				collector.CollectLog()
				log.Info().Str("code", collector.Code).Msg("Done collectors...")
				terminationWg.Done()
			}()
		}
	}
}

func RunCollector() {
	initialize()
	loadCollectorFromDatabases()

	// init handlers
	for i := range env.ELK_LOG_MAX_HANDLERS {
		terminationWg.Add(1)
		go func() {
			defer terminationWg.Done()
			runHandlers(i)
		}()
	}

	for {
		select {
		case <-terminationSigChan:
			log.Info().Msg("Received shutdown signal, shutting down...")
			shutdown()
			log.Info().Msg("App shutdown complete, bye bye!")
			return
		case <-scheduler.C:
			// run collector
			runCollectors()
		}
	}
}
