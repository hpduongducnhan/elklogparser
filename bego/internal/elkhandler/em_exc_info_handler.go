package elkhandler

import (
	"fmt"
	"nhandd/bego/internal/elkclient"
	"time"
)

type emExcInfoLog struct {
	Agid      string                `json:"agid"`
	Message   string                `json:"message"`
	ExcInfo   string                `json:"exc_info"`
	LogID     string                `json:"log_id"`
	Loc       string                `json:"loc"`
	Timestamp time.Time             `json:"timestamp"`
	Raw       elkclient.ElkInnerHit `json:"raw"`
}

func (e *emExcInfoLog) IsValid() bool {
	// isValid := e.LogID != "" && e.Agid != ""
	isValid := e.LogID != ""
	// log.Info().Str("LogID", e.LogID).Str("Agid", e.Agid).Bool("isValid", isValid).Msg("emExcInfoLog validation")
	return isValid
}

func (e *emExcInfoLog) Parse(hit *elkclient.ElkInnerHit) {
	e.LogID = hit.ID
	e.Raw = *hit
	if excInfo, excInfoOK := hit.Source["exc_info"].(string); excInfoOK {
		e.ExcInfo = excInfo
	}
	if message, messageOK := hit.Source["message"].(string); messageOK {
		e.Message = message
	}
	if agid, agidOK := hit.Source["agid"].(string); agidOK {
		e.Agid = agid
	}
	if timestamp, timestampOK := hit.Source["@timestamp"].(string); timestampOK {
		e.Timestamp, _ = time.Parse(time.RFC3339, timestamp)
	}
	if loc, locOK := hit.Source["loc"].(string); locOK {
		e.Loc = loc
	}
}

func (h *ElkResponseHandler) emExcInfoHandler(msg *elkclient.ElkResponseWithCode) error {
	// log.Info().Msgf("elkRespHandleEMExcInfo: %s", msg.Code)
	excInfo := &emExcInfoLog{}
	excInfo.Parse(msg.ElkInnerHit)
	if excInfo.IsValid() {
		// log.Info().Msgf("elkRespHandleEMExcInfo: %v", excInfo)
		h.pgRepo.CreateOrUpdateElkCollectedLog(
			msg.CollectorCode, msg.QueryCode,
			excInfo.LogID, excInfo.Raw, excInfo.Timestamp,
		)
	} else {
		// log.Error().Interface("msg", msg).Msgf("elkRespHandleEMExcInfo: Invalid log data")
		return fmt.Errorf("invalid log data")
	}
	return nil
}
