package elkcollector

import (
	"context"
	"nhandd/bego/internal/elkclient"
	"nhandd/bego/internal/elkhandler"
	"nhandd/bego/internal/environment"
	"nhandd/bego/internal/postgresclient"
	"nhandd/bego/internal/redisclient"
	"nhandd/bego/internal/repository"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var initOnce sync.Once

var terminationSigChan chan os.Signal
var terminationCtx context.Context
var terminationFunc context.CancelFunc
var terminationWg *sync.WaitGroup

var elkRespChan chan *elkclient.ElkResponseWithCode
var elkRespHandler *elkhandler.ElkResponseHandler
var elkQueryBuilder *ElkQueryBuilder
var elkCollectors map[string]*ElkCollector

var scheduler *time.Ticker
var env *environment.EnvConfig
var pgDb *gorm.DB
var pgRepo *repository.DbRepository
var rdClient *redis.Client

func initialize() {
	var err error
	initOnce.Do(func() {
		terminationCtx, terminationFunc = context.WithCancel(context.Background())
		terminationSigChan = make(chan os.Signal, 1)
		signal.Notify(terminationSigChan, os.Interrupt, syscall.SIGTERM)
		terminationWg = &sync.WaitGroup{}

		// init environment
		env = environment.GetEnv()

		// init elk response channel with buffer size
		elkRespChan = make(chan *elkclient.ElkResponseWithCode, env.ELK_RESPONSE_BUFFER_SIZE)

		// init scheduler
		scheduler = time.NewTicker(time.Duration(env.SCHEDULER_TICKER) * time.Second)

		// init postgres and repository
		pgDb, err = postgresclient.GetPostgresClient(env.POSTGRES_DNS, &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			panic(err)
		}
		pgRepo = repository.NewDbRepository(pgDb)

		// init redis client
		rdClient, err = redisclient.GetRedisClient(env.REDIS_URL)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to create Redis client")
		}

		// init elk response handler
		elkRespHandler = elkhandler.NewElkRespHandler(rdClient, pgDb, pgRepo)

		// init elk query builder
		elkQueryBuilder = &ElkQueryBuilder{}

		// init elk collectors map
		elkCollectors = make(map[string]*ElkCollector)

		log.Info().Msgf("initialize complete, App running ...")
	})
}

func shutdown() {
	var err error
	if scheduler != nil {
		scheduler.Stop()
	}

	terminationFunc()

	log.Info().Msg("Waiting for other goroutines to finish...")
	terminationWg.Wait()

	log.Info().Msg("All goroutines finished, close other resources...")
	if elkRespChan != nil {
		close(elkRespChan)
		log.Info().Msg("Closed elkRespChan")
	}

	// close db connection
	// close redis connection
	if rdClient != nil {
		err = rdClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close Redis client")
		} else {
			log.Info().Msg("Redis client closed")
		}
	}

	// close collector
	for _, collector := range elkCollectors {
		if collector != nil {
			collector.Shutdown()
		}
	}
}
