package elkhandler

import (
	"fmt"
	"nhandd/bego/internal/elkclient"
	"time"
)

type emTicketCodeLog struct {
	Agid       string                `json:"agid"`
	Message    string                `json:"message"`
	TicketCode string                `json:"ticket_code"`
	LogID      string                `json:"log_id"`
	Loc        string                `json:"loc"`
	Timestamp  time.Time             `json:"timestamp"`
	Raw        elkclient.ElkInnerHit `json:"raw"`
}

func (e *emTicketCodeLog) IsValid() bool {
	// isValid := e.LogID != "" && e.Agid != ""
	isValid := e.LogID != ""
	// log.Info().Str("LogID", e.LogID).Str("Agid", e.Agid).Bool("isValid", isValid).Msg("emExcInfoLog validation")
	return isValid
}

func (e *emTicketCodeLog) Parse(hit *elkclient.ElkInnerHit) {
	e.LogID = hit.ID
	e.Raw = *hit
	if ticketCode, ticketCodeOK := hit.Source["ticket_code"].(string); ticketCodeOK {
		e.TicketCode = ticketCode
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

func (h *ElkResponseHandler) emTicketCodeHandler(msg *elkclient.ElkResponseWithCode) error {
	ticketInfo := &emTicketCodeLog{}
	ticketInfo.Parse(msg.ElkInnerHit)
	if ticketInfo.IsValid() {
		h.pgRepo.CreateOrUpdateElkCollectedLog(
			msg.CollectorCode, msg.QueryCode,
			ticketInfo.LogID, ticketInfo.TicketCode,
			ticketInfo.Raw, ticketInfo.Timestamp,
		)
	} else {
		return fmt.Errorf("invalid log data")
	}
	return nil
}
