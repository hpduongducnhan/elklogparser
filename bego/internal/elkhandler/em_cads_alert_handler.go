package elkhandler

import (
	"fmt"
	"nhandd/bego/internal/elkclient"
	"time"
)

type emCadsAlertLog struct {
	Agid      string                `json:"agid"`
	Message   string                `json:"message"`
	Alert     string                `json:"req_body"`
	LogID     string                `json:"log_id"`
	Loc       string                `json:"loc"`
	Timestamp time.Time             `json:"timestamp"`
	Raw       elkclient.ElkInnerHit `json:"raw"`
}

func (e *emCadsAlertLog) IsValid() bool {
	isValid := e.LogID != ""
	return isValid
}

func (e *emCadsAlertLog) Parse(hit *elkclient.ElkInnerHit) {
	e.LogID = hit.ID
	e.Raw = *hit
	if alert, alertOK := hit.Source["req_body"].(string); alertOK {
		e.Alert = alert
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

func (h *ElkResponseHandler) emCadsAlertHandler(msg *elkclient.ElkResponseWithCode) error {
	ticketInfo := &emTicketCodeLog{}
	ticketInfo.Parse(msg.ElkInnerHit)
	if ticketInfo.IsValid() {
		h.pgRepo.CreateOrUpdateElkCollectedLog(
			msg.CollectorCode, msg.QueryCode,
			ticketInfo.LogID, ticketInfo.Raw, ticketInfo.Timestamp,
		)
	} else {
		return fmt.Errorf("invalid log data")
	}
	return nil
}
