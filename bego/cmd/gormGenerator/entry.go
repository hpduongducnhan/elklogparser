package gormgenerator

import (
	"nhandd/bego/internal/environment"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var env *environment.EnvConfig

func RunDbModelsGenerators() {
	env = environment.GetEnv()
	log.Info().Interface("environment", env).Msg("env")
	generator := gen.NewGenerator(gen.Config{
		OutPath: "./postgres_models/",                                               // Output directory for generated code
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // Generate default query and interface
	})
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  env.POSTGRES_DNS,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database: " + err.Error())
	}
	generator.UseDB(db) // Use the database connection
	generator.GenerateAllTable()
	generator.Execute()
}
