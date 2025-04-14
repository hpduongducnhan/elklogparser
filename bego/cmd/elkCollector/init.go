package elkcollector

import (
	"context"
	"nhandd/bego/internal/elkclient"
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

var env *environment.EnvConfig
var pgDb *gorm.DB
var pgRepo *repository.DbRepository
var rClient *redis.Client
var err error
var terminateSigChan chan os.Signal
var terminateCtx context.Context
var terminateFunc context.CancelFunc
var terminateWaitGroup *sync.WaitGroup

var scheduler time.Ticker

var elkCollectors map[string]*ElkCollector
var elkQueryBuilder *ElkQueryBuilder = &ElkQueryBuilder{}
var elkResponseHandler *ElkResponseHandler = &ElkResponseHandler{}
var elkInnerHitChan chan *elkclient.ElkResponseWithCode

func initialize() {
	terminateCtx, terminateFunc = context.WithCancel(context.Background())
	terminateSigChan = make(chan os.Signal, 1)
	signal.Notify(terminateSigChan, os.Interrupt, syscall.SIGTERM)
	terminateWaitGroup = &sync.WaitGroup{}

	// init environment
	env = environment.GetEnv()

	// init scheduler
	scheduler = *time.NewTicker(time.Duration(env.SCHEDULER_TICKER) * time.Second)

	// init postgres
	pgDb, err = postgresclient.GetPostgresClient(env.POSTGRES_DNS, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic(err)
	}
	pgRepo = repository.NewDbRepository(pgDb)

	// init redis client
	rClient, err = redisclient.GetRedisClient(env.REDIS_URL)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to create Redis client")
	}

	// init elk inner hit response chan
	elkInnerHitChan = make(chan *elkclient.ElkResponseWithCode, 9999)
	log.Info().Msg("Initialized done")
}
