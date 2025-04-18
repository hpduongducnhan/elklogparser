package repository

import (
	"nhandd/bego/internal/models"
	"time"

	"github.com/rs/zerolog/log"
)

// Query để lấy tất cả config đang active
func (r *DbRepository) GetAllActiveConfigs() ([]models.DatasourceElkcollectorconfig, error) {
	var configs []models.DatasourceElkcollectorconfig

	err := r.db.
		Where("active = ?", true).
		Find(&configs).Error

	if err != nil {
		return nil, err
	}

	return configs, nil
}

func (r *DbRepository) GetActiveConfigsWithRelations() ([]models.DatasourceElkcollectorconfig, error) {
	var configs []models.DatasourceElkcollectorconfig

	err := r.db.
		Preload("ElkQueries").      // Tải các queries liên quan
		Preload("LogFilters").      // Tải các filters liên quan
		Preload("ElkConfig").       // Tải các config liên quan
		Preload("ElkConfig.Proxy"). // Tải ProxyConfig thông qua ElkConfig
		Where("active = ?", true).
		Find(&configs).Error

	if err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *DbRepository) UpdateLastRunAt(collectorCode string, lastRunAt time.Time) error {
	err := r.db.Model(&models.DatasourceElkcollectorconfig{}).
		Where("code = ?", collectorCode).
		Update("last_run_at", lastRunAt).Error
	if err != nil {
		log.Error().Err(err).Msg("Failed to update last_run_at")
		return err
	}
	return nil
}
