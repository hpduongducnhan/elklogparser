package elkcollector

import (
	"strings"

	"github.com/rs/zerolog/log"
)

type ElkQueryBuilder struct {
}

func (b *ElkQueryBuilder) getBuilderByCode(code string) func(code, query string, elkCollector *ElkCollector) string {
	switch code {
	case "em-exc-info":
		return b.buildQueryCodeEmExcInfo
	case "em-ticket-code":
		return b.buildQueryCodeEmTicketCode
	default:
		return nil
	}
}

func (b *ElkQueryBuilder) buildWithTimeRange(query string, elkCollector *ElkCollector) string {
	fromDate, toDate := elkCollector.GetQueryTimeRange()

	strFromDate := fromDate.Format("2006-01-02T15:04:05.000000")
	strToDate := toDate.Format("2006-01-02T15:04:05.000000")

	builtQuery := strings.ReplaceAll(query, "<start-datetime>", strFromDate)
	builtQuery = strings.ReplaceAll(builtQuery, "<end-datetime>", strToDate)

	for _, char := range []string{"\t", "\n", ""} {
		builtQuery = strings.ReplaceAll(builtQuery, char, "")
	}
	log.Info().Str("fromDate", strFromDate).Str("toDate", strToDate).Str("query", query).Msg("Building query with time range")
	return builtQuery
}

func (b *ElkQueryBuilder) Build(code, query string, elkCollector *ElkCollector) string {
	builder := b.getBuilderByCode(code)
	if builder == nil {
		log.Error().Str("code", code).Msg("No builder found for this code")
		return ""
	}
	return builder(code, query, elkCollector)
}

func (b *ElkQueryBuilder) buildQueryCodeEmExcInfo(code, query string, elkCollector *ElkCollector) string {
	queryWithTimeRange := b.buildWithTimeRange(query, elkCollector)
	return queryWithTimeRange
}

func (b *ElkQueryBuilder) buildQueryCodeEmTicketCode(code, query string, elkCollector *ElkCollector) string {
	queryWithTimeRange := b.buildWithTimeRange(query, elkCollector)
	return queryWithTimeRange
}
