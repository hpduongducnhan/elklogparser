package elkclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

func ScrollLogWithCode(
	collectorCode, queryCode string,
	elkClient *elasticsearch.Client,
	elkIndex, elkQuery string,
	resultChan chan *ElkResponseWithCode,
) error {
	if elkClient == nil {
		log.Panic().Msg("Elasticsearch client is not initialized")
	}
	counter := 0
	defer func() {
		log.Info().Str("collector", collectorCode).Str("query", queryCode).Int("counter", counter).Msg("scroll result")
	}()

	// log.Info().Str("code", code).Str("index", elkIndex).Str("query", elkQuery).Msg("scroll log with code")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	var scrollID string
	res, err := elkClient.Search(
		elkClient.Search.WithContext(ctx),
		elkClient.Search.WithIndex(elkIndex),
		elkClient.Search.WithBody(strings.NewReader(elkQuery)),
		elkClient.Search.WithScroll(elkScrollTime),
	)
	if err != nil {
		log.Error().Err(err).Msg("error searching")
		return err
	}

	parsedResp, err := parseElkResponse(res)
	if err != nil {
		log.Error().Err(err).Msg("error parsing response")
		return err
	}

	scrollID = parsedResp.ScrollID
	// log.Info().Str("scrollID", scrollID).Msg("scroll ID")
	defer func() {
		err := ClearScroll(elkClient, scrollID)
		if err != nil {
			log.Error().Err(err).Msg("error clearing scroll")
		}
	}()

	if len(parsedResp.Hits.Hits) == 0 {
		return nil
	}

	if scrollID == "" {
		log.Error().Msg("empty scroll ID")
		return fmt.Errorf("empty scroll ID")
	}

	// push the parsed hits to the result channel
	for _, hit := range parsedResp.Hits.Hits {
		// log.Info().Msg("pushing hit to channel")
		resultChan <- &ElkResponseWithCode{
			CollectorCode: collectorCode,
			QueryCode:     queryCode,
			ElkInnerHit:   &hit,
		}
		counter++
	}

	for {
		res, err := elkClient.Scroll(
			elkClient.Scroll.WithContext(ctx),
			elkClient.Scroll.WithScrollID(scrollID),
			elkClient.Scroll.WithScroll(elkScrollTime),
		)
		// log.Info().Interface("result", res).Msg("scrolling")
		if err != nil {
			log.Error().Err(err).Msg("error scrolling")
			return err
		}
		parsedResp, err = parseElkResponse(res)
		if err != nil {
			log.Error().Err(err).Msg("error parsing response when scrolling")
			return err
		}
		if len(parsedResp.Hits.Hits) == 0 {
			// log.Info().Int("counter", counter).Msg("scroll done")
			return nil
		}
		// push the parsed hits to the result channel
		for _, hit := range parsedResp.Hits.Hits {
			resultChan <- &ElkResponseWithCode{
				CollectorCode: collectorCode,
				QueryCode:     queryCode,
				ElkInnerHit:   &hit,
			}
			counter++
		}
		// Cập nhật scrollID cho lần lặp kế tiếp
		scrollID = parsedResp.ScrollID
	}
}
