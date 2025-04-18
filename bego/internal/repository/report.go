package repository

import (
	"encoding/json"
	"nhandd/bego/internal/elkclient"
	"nhandd/bego/internal/models"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func (r *DbRepository) CreateOrUpdateElkCollectedLog(
	collectorCode, queryCode, logID, value string,
	logRaw elkclient.ElkInnerHit, logTimestamp time.Time,
) (*models.ReportElkcollectedlog, error) {
	// Prepare the log data
	logRawJSON, err := json.Marshal(logRaw)
	if err != nil {
		return nil, err
	}

	tempLog := &models.ReportElkcollectedlog{
		CollectorCode: collectorCode,
		LogType:       queryCode, // Giả sử QueryCode tương ứng với LogType
		LogID:         logID,
		LogRaw:        logRawJSON,
		Value:         value,
		LogTimestamp:  float64(logTimestamp.Unix()), // Chuyển time.Time thành Unix timestamp (seconds)
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// if logID == "hu9oMpYBDYFNSlkBg-qf" {
	// 	log.Info().Str("LogID", logID).Interface("logRawJSON", logRaw).Msg("check log raw")
	// }

	// Try to find existing record
	var existing models.ReportElkcollectedlog
	err = r.db.Where(&models.ReportElkcollectedlog{
		CollectorCode: collectorCode,
		LogType:       queryCode,
		LogID:         logID,
	}).First(&existing).Error

	if err == nil {
		// Record exists, update it
		updates := map[string]interface{}{
			"log_raw":       string(logRawJSON),
			"log_timestamp": tempLog.LogTimestamp,
			"value":         value,
			"updated_at":    time.Now(),
		}
		if err := r.db.Model(&existing).Updates(updates).Error; err != nil {
			log.Warn().Err(err).Msg("Failed to update existing log")
			return nil, err
		} else {
			log.Info().Str("LogID", logID).Msg("Updated existing log")
		}
		return &existing, nil
	}

	if err == gorm.ErrRecordNotFound {
		// No record found, create new one
		if err := r.db.Create(tempLog).Error; err != nil {
			return nil, err
		}
		return tempLog, nil
	}

	// Return any other error
	return nil, err
}
