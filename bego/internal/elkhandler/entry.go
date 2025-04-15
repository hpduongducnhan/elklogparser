package elkhandler

import (
	"fmt"
	"nhandd/bego/internal/elkclient"
	"nhandd/bego/internal/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ElkResponseHandler struct {
	// codeMap map[string]func(msg *elkclient.ElkResponseWithCode) error
	rdClient *redis.Client
	pgRepo   *repository.DbRepository
	pgDb     *gorm.DB
}

func NewElkRespHandler(r *redis.Client, db *gorm.DB, repo *repository.DbRepository) *ElkResponseHandler {
	return &ElkResponseHandler{rdClient: r, pgRepo: repo, pgDb: db}
}

func (h *ElkResponseHandler) getHandlerByCode(collectorCode, queryCode string) func(msg *elkclient.ElkResponseWithCode) error {
	code := fmt.Sprintf("%s__%s", collectorCode, queryCode)
	switch code {
	case "em-exc-info__em-exc-info":
		return h.emExcInfoHandler
	case "em-ticket-code__em-ticket-code":
		return h.emTicketCodeHandler
	default:
		return nil
	}
}

func (h *ElkResponseHandler) Handle(msg *elkclient.ElkResponseWithCode) error {
	handler := h.getHandlerByCode(msg.CollectorCode, msg.QueryCode)
	if handler == nil {
		return fmt.Errorf("no handler found for query code: %s", msg.QueryCode)
	}
	return handler(msg)
}
