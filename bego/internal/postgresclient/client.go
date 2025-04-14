package postgresclient

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var dbClient *gorm.DB

func MigrateModels(models ...any) {
	if dbClient == nil {
		panic("Database client is not initialized")
	}
	for _, model := range models {
		if err := dbClient.AutoMigrate(model); err != nil {
			panic(fmt.Sprintf("Failed to migrate model %v: %v", model, err))
		}
	}
}

func GetPostgresClient(dsn string, dbConfig *gorm.Config) (*gorm.DB, error) {
	// Chuỗi kết nối PostgreSQL
	// dsn := "host=localhost user=your_user password=your_password dbname=your_db port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"

	// Kết nối đến PostgreSQL
	var err error
	dbClient, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), dbConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to connect to PostgreSQL")
		return nil, err
	}
	return dbClient, err
}
