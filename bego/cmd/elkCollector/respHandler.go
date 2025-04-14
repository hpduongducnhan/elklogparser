package elkcollector

import (
	"fmt"
	"nhandd/bego/internal/elkclient"
)

type ElkResponseHandler struct {
	// codeMap map[string]func(msg *elkclient.ElkResponseWithCode) error
}

func (h *ElkResponseHandler) getHandlerByCode(collectorCode, queryCode string) func(msg *elkclient.ElkResponseWithCode) error {
	code := fmt.Sprintf("%s__%s", collectorCode, queryCode)
	switch code {
	case "em-exc-info__em-exc-info":
		return elkRespHandleEMExcInfo
	default:
		return nil
	}
}

func (h *ElkResponseHandler) saveDBCollectedLog(msg *elkclient.ElkResponseWithCode) error {
	return elkRespHandleEMExcInfo(msg)
}

func (h *ElkResponseHandler) Handle(msg *elkclient.ElkResponseWithCode) error {
	handler := h.getHandlerByCode(msg.CollectorCode, msg.QueryCode)
	if handler == nil {
		return fmt.Errorf("no handler found for query code: %s", msg.QueryCode)
	}
	return handler(msg)
}
