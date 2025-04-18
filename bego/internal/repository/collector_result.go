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

func (r *DbRepository) CreateNewCollectorResult(
	jid string,
	collectorID int64,
	queryFrom, queryTo, startAt time.Time,
) error {
	// find exist
	var existing models.DatasourceElkcollectresult
	err := r.db.Where(&models.DatasourceElkcollectresult{Jid: jid}).First(&existing).Error
	if err == nil {
		// Record exists, return
		// log.Info().Err(err).Str("jid", jid).Msg("CreateNewCollectorResult: Record existed")
		return nil
	}

	result := models.DatasourceElkcollectresult{
		Jid:               jid,
		CollectorID:       collectorID,
		QueryFromDatetime: queryFrom,
		QueryToDatetime:   queryTo,
		StartAt:           startAt,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Success:           false, // Mặc định là false khi tạo mới
		Status:            "running",
		Detail:            "{}", // Có thể để trống hoặc giá trị mặc định
	}
	return r.db.Create(&result).Error
}

// UpdateElkCollectResult updates an existing DatasourceElkcollectresult record by jid
func (r *DbRepository) UpdateElkCollectResult(
	jid string,
	success bool,
	status, detail string, finishAt time.Time,
) error {
	var existing models.DatasourceElkcollectresult
	err := r.db.Where(&models.DatasourceElkcollectresult{Jid: jid}).First(&existing).Error
	if err != nil {
		// Record not exists, return
		// log.Info().Err(err).Str("jid", jid).Msg("UpdateElkCollectResult: Record not found")
		return nil
	}
	return r.db.Model(&models.DatasourceElkcollectresult{}).
		Where("jid = ?", jid).
		Updates(map[string]any{
			"finish_at":  finishAt,
			"status":     status,
			"success":    success,
			"detail":     detail,
			"updated_at": time.Now(),
		}).Error
}
