package simulator

import (
	"nhandd/bego/internal/environment"
	"nhandd/bego/internal/postgresclient"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var env *environment.EnvConfig
var pgDb *gorm.DB
var rClient *redis.Client

func RunSimulation() {
	var err error
	env = environment.GetEnv()
	pgDb, err = postgresclient.GetPostgresClient(env.POSTGRES_DNS, &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// GetActiveCollectors()
	RunQueryElk()
}
