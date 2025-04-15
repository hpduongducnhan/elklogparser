package elkclient

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog/log"
)

var elkScrollTime time.Duration = 2 * time.Minute
var client *elasticsearch.Client

func GetDefaultElkConfig(serverAddrs []string, authUser, authPassword string) *elasticsearch.Config {
	return &elasticsearch.Config{
		Addresses:     serverAddrs,
		Username:      authUser,
		Password:      authPassword,
		RetryOnStatus: []int{500, 502, 503, 504},
		MaxRetries:    3,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     50,
			IdleConnTimeout:     30 * time.Second,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		},
		Logger: &elastictransport.TextLogger{
			// Output:             os.Stderr,
			Output:             io.Discard,
			EnableRequestBody:  true,
			EnableResponseBody: true,
		},
	}
}

func GetElkClient(serverAddrs []string, authUser, authPassword string, elkConfig *elasticsearch.Config) (*elasticsearch.Client, error) {
	if client == nil {
		if elkConfig == nil {
			elkConfig = GetDefaultElkConfig(serverAddrs, authUser, authPassword)
		}
		var err error
		client, err = elasticsearch.NewClient(*elkConfig)
		if err != nil {
			return nil, err
		}
		_, err = client.Ping()
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func NewElkClient(serverAddrs []string, authUser, authPassword string, elkConfig *elasticsearch.Config) (*elasticsearch.Client, error) {
	if elkConfig == nil {
		elkConfig = GetDefaultElkConfig(serverAddrs, authUser, authPassword)
	}
	client, err := elasticsearch.NewClient(*elkConfig)
	if err != nil {
		return nil, err
	}
	_, err = client.Ping()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func parseElkResponse(elkResp *esapi.Response) (*ElkResponse, error) {
	var response ElkResponse

	bodyBytes, err := io.ReadAll(elkResp.Body)
	elkResp.Body.Close() // always close the response body
	if err != nil {
		log.Error().Err(err).Str("elkClient", "parseResp").Msg("error reading response body")
		return nil, err
	}

	// Parse JSON into the struct
	err = json.Unmarshal(bodyBytes, &response)
	if err != nil {
		log.Error().Err(err).Str("elkClient", "parseResp").Msg("error unmarshaling JSON")
		return nil, err
	}
	return &response, nil
}

func CloseClientWithScrolls(elkClient *elasticsearch.Client) error {
	if elkClient == nil {
		return nil
	}
	// Gọi API ClearScroll với scroll_id là "_all" để xóa tất cả scroll context
	res, err := elkClient.ClearScroll(
		elkClient.ClearScroll.WithContext(context.Background()),
		elkClient.ClearScroll.WithScrollID("_all"),
	)
	if err != nil {
		log.Error().Err(err).Msg("error clearing all scrolls")
		return err
	}
	defer res.Body.Close()
	log.Info().Msg("cleared all scroll contexts")

	// close all connections
	// if transport, ok := elkClient.Transport.(*elastictransport.Client); ok {

	// }
	return nil
}

func ClearScroll(elkClient *elasticsearch.Client, scrollID string) error {
	if elkClient == nil {
		log.Warn().Msg("elkClient is nil")
		return nil
	}
	if scrollID == "" {
		log.Warn().Msg("scrollID is empty")
		return nil
	}
	_, err := elkClient.ClearScroll(
		elkClient.ClearScroll.WithContext(context.Background()),
		elkClient.ClearScroll.WithScrollID(scrollID),
	)
	if err != nil {
		return err
	}
	return nil
}
