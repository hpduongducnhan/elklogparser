package tools

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"nhandd/bego/internal/environment"

	"github.com/rs/zerolog/log"
)

var envConfig *environment.EnvConfig

func newTransportWithProxy() (*http.Transport, error) {
	envConfig = environment.GetEnv()

	proxyURL, err := url.Parse(envConfig.PROXY_ADDR)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to parse proxy URL: %v", err)
		return nil, err
	}
	// Create transport with proxy
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	return transport, nil
}

func SendToDiscord(jobID, webhookUrl string, notifyMsg string) error {
	// Create the message payload
	message := map[string]string{
		"content": fmt.Sprintf("new message at %s\n```%s```", time.Now().UTC().Format(time.RFC1123), notifyMsg),
	}
	payload, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msgf("sendToDiscord task %s error marshalling JSON: %v", jobID, err)
		return err
	}
	// Create the POST request
	req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(payload))
	if err != nil {
		log.Error().Err(err).Msgf("sendToDiscord task %s error creating request: %v", jobID, err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	transport, err := newTransportWithProxy()
	if err != nil {
		log.Error().Err(err).Msgf("sendToDiscord task %s error creating transport: %v", jobID, err)
		return err
	}

	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msgf("sendToDiscord task %s error sending request: %v", jobID, err)
		return err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode == http.StatusNoContent {
		log.Info().Msgf("sendToDiscord task %s message sent successfully!", jobID)
	} else {
		log.Warn().Msgf("sendToDiscord task %s failed to send message. Status code: %d", jobID, resp.StatusCode)
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}
	return nil
}
