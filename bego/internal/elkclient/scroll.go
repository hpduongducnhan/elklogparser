package elkclient

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

func ScrollLog[T any](
	elkClient *elasticsearch.Client,
	elkIndex, elkQuery string,
	resultChan chan *T,
	hitParser func(*ElkInnerHit) (*T, error),
) error {
	if elkClient == nil {
		log.Panic().Msg("Elasticsearch client is not initialized")
	}
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
	defer func() {
		if scrollID != "" {
			_, err := elkClient.ClearScroll(
				elkClient.ClearScroll.WithContext(context.Background()),
				elkClient.ClearScroll.WithScrollID(scrollID),
			)
			if err != nil {
				log.Error().Err(err).Msg("error clearing scroll")
			}
		}
	}()

	if scrollID == "" {
		log.Error().Msg("empty scroll ID")
		return fmt.Errorf("empty scroll ID")
	}

	if len(parsedResp.Hits.Hits) == 0 {
		// log.Info().Msg("No more hits")
		return nil
	}

	// push the parsed hits to the result channel
	for _, hit := range parsedResp.Hits.Hits {
		parsedHit, err := hitParser(&hit)
		if err != nil {
			log.Error().Err(err).Msg("error parsing hit")
		} else {
			resultChan <- parsedHit
		}
	}

	for {
		res, err := elkClient.Scroll(
			elkClient.Scroll.WithContext(ctx),
			elkClient.Scroll.WithScrollID(scrollID),
			elkClient.Scroll.WithScroll(elkScrollTime),
		)
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
			// log.Info().Msg("No more hits")
			return nil
		}
		// push the parsed hits to the result channel
		for _, hit := range parsedResp.Hits.Hits {
			parsedHit, err := hitParser(&hit)
			if err != nil {
				log.Error().Err(err).Msg("error parsing hit")
			} else {
				resultChan <- parsedHit
			}
		}
	}
}
