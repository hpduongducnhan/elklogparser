package repository

import (
	"nhandd/bego/internal/models"
	"time"
)

func (r *DbRepository) GetCollectorResults(collectorID int64, startTime, endTime time.Time) ([]models.DatasourceElkcollectresult, error) {
	var results []models.DatasourceElkcollectresult

	err := r.db.
		Where("collector_id = ? AND created_at BETWEEN ? AND ?",
			collectorID,
			startTime,
			endTime).
		Order("created_at DESC").
		Find(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}
