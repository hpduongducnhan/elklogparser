package elkcollector

import (
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/rs/zerolog/log"
)

func listRunningGoroutines() {
	// Lấy profile của tất cả goroutine
	profile := pprof.Lookup("goroutine")
	if profile == nil {
		fmt.Println("Cannot lookup goroutine profile")
		return
	}

	// In thông tin goroutine ra stdout
	fmt.Println("List of running goroutines:")
	profile.WriteTo(os.Stdout, 1)
}

func waitForTerminate(wait bool) {
	if wait {
		<-terminateSigChan
	}
	log.Info().Msg("Received termination signal, shutting down...")

	// terminate context
	terminateFunc()
	// waiit for all goroutines to finish
	terminateWaitGroup.Wait()
	log.Info().Msg("waitgroup finished")

	// Close the scheduler
	close(elkInnerHitChan)

	// close the collectors
	for _, collector := range elkCollectors {
		if collector != nil {
			collector.Shutdown()
		}
	}

	// Close the Redis client
	if rClient != nil {
		err = rClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close Redis client")
		}
	}

	for _, collector := range elkCollectors {
		if collector != nil {
			collector.Shutdown()
		}
	}

	log.Info().Msg("Closed everything, bye bye!")
}

func loadDatabases() {
	if elkCollectors == nil {
		elkCollectors = make(map[string]*ElkCollector)
	}
	// load from databases
	activeConfigs, err := pgRepo.GetActiveConfigsWithRelations()
	if err != nil {
		log.Error().Err(err).Msg("Error getting active configs")
		return
	}

	for _, config := range activeConfigs {
		elkCollectors[config.Code] = NewElkCollectorFromDbConf(config)
	}
	return
}

func handleLogResponses(wid int) {
	// log.Info().Int("wid", wid).Msgf("Starting log response handler %d", wid)
	for {
		select {
		case val, ok := <-elkInnerHitChan:
			if !ok {
				terminateWaitGroup.Done()
				log.Warn().Msg("elkInnerHitChan closed")
				return
			}
			elkResponseHandler.Handle(val)
		case <-terminateCtx.Done():
			terminateWaitGroup.Done()
			return
		}
	}
}

func RunCollector() {
	initialize()
	loadDatabases()

	// init handlers
	for wid := range env.ELK_LOG_MAX_HANDLERS {
		terminateWaitGroup.Add(1)
		go handleLogResponses(wid)
	}
	log.Info().Msg("Started log response handlers")

	// schedule collector
	defer scheduler.Stop()
	for {
		select {
		case <-scheduler.C:
			for _, collector := range elkCollectors {
				if collector.AllowToRun() {
					go collector.CollectLog()
				}
			}
		case <-terminateSigChan:
			waitForTerminate(false)
			listRunningGoroutines()
			return
		}
	}
}
