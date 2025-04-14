package simulator

import (
	"fmt"
	"strings"
	"time"

	"nhandd/bego/internal/elkclient"
	"nhandd/bego/internal/environment"
	"nhandd/bego/internal/redisclient"

	"github.com/rs/zerolog/log"
)

func newElkQuery(startTime, endTime time.Time) (query string) {
	template := `{
		"size": 500,
		"query": {
			"bool": {
				"must": [
					{
						"range": {
							"@timestamp": {
								"gte": "%s",
								"lte": "%s"
							}
						}
					},
					{
						"exists": {
							"field": "exc_info"
						}
					}
				]
			}
		},
		"_source": {
			"excludes": ["agent", "log", "container", "orchestrator", "ecs", "stream", "region", "input", "host", "docker"]
		},
		"sort": [
			{"@timestamp": {"order": "asc"}}
		]
	}`
	strFromDate := startTime.Add(-10 * time.Second).UTC().Format("2006-01-02T15:04:05.000000")
	strToDate := endTime.Add(10 * time.Second).UTC().Format("2006-01-02T15:04:05.000000")
	builtQuery := fmt.Sprintf(template, strFromDate, strToDate)
	// builtQuery := fmt.Sprintf(template, alertGroupID)
	query = strings.ReplaceAll(builtQuery, "\n", "")
	for _, char := range []string{"\t", "\n", ""} {
		query = strings.ReplaceAll(query, char, "")
	}
	return
}

type EMException struct {
	Agid      string    `json:"agid"`
	Message   string    `json:"message"`
	ExcInfo   string    `json:"exc_info"`
	LogID     string    `json:"log_id"`
	Loc       string    `json:"loc"`
	Timestamp time.Time `json:"timestamp"`
}

func processInnerHit(hit *elkclient.ElkInnerHit) (*EMException, error) {
	// Process the inner hit
	// For example, you can log the ID and Index
	// log.Info().Interface("hit", hit).Msg("Processing Inner Hit")

	excInfo, excInfoOK := hit.Source["exc_info"].(string)
	message, messageOK := hit.Source["message"].(string)
	agid, agidOK := hit.Source["agid"].(string)
	timestamp, timestampOK := hit.Source["@timestamp"].(string)
	loc, locOK := hit.Source["loc"].(string)
	if excInfoOK {
		_msg := ""
		_agid := ""
		_timestamp := time.Time{}
		_loc := ""
		if timestampOK {
			_timestamp, _ = time.Parse(time.RFC3339, timestamp)
		}
		if messageOK {
			_msg = message
		}
		if agidOK {
			_agid = agid
		}
		if locOK {
			_loc = loc
		}
		foundExc := &EMException{
			Agid:      _agid,
			Message:   _msg,
			ExcInfo:   excInfo,
			LogID:     hit.ID,
			Loc:       _loc,
			Timestamp: _timestamp,
		}
		// log.Info().Interface("foundExc", foundExc).Msg("Found exception")
		return foundExc, nil
	} else {
		// log.Info().Interface("excInfoOK", excInfoOK).Msg("No exc_info found in hit")
		return nil, fmt.Errorf("exc_info not found in hit")
	}
}

func RunQueryElk() {
	env = environment.GetEnv()
	rClient, _ = redisclient.GetRedisClient(env.REDIS_URL)
	log.Info().Str("REDIS_URL", env.REDIS_URL).Interface("rClient", rClient).Msg("Redis URL")

	eClient, _ := elkclient.NewElkClient(
		env.ELK_LOG_ADDRS,
		env.ELK_LOG_AUTH_USERNAME,
		env.ELK_LOG_AUTH_PASSWORD,
		nil,
	)

	query := newElkQuery(time.Now().Add(-30*time.Minute), time.Now())
	log.Info().Str("query", query).Msg("Elk query")
	resultChan := make(chan *EMException, 1000)
	elkclient.ScrollLog(eClient, env.ELK_LOG_INDEX, query, resultChan, processInnerHit)
	for {
		select {
		case result := <-resultChan:
			if result != nil {
				// log.Info().Interface("result", result).Msg("Result from Elk")
				// Do something with the result
			} else {
				// log.Info().Msg("No result from Elk")
			}
		case <-time.After(30 * time.Second):
			log.Info().Msg("Timeout waiting for results")
			return
		}
	}
}
