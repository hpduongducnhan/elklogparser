package environment

import (
	"github.com/hpduongducnhan/goscs"
	"github.com/rs/zerolog/log"
)

type EnvConfig struct {
	goscs.BaseSettings

	WEB_APP_HOST        string `mapstructure:"WEB_APP_HOST"`
	WEB_APP_PORT        int    `mapstructure:"WEB_APP_PORT"`
	WEB_APP_PREFORK     bool   `mapstructure:"WEB_APP_PREFORK"`
	WEB_APP_CONCURRENCY int    `mapstructure:"WEB_APP_CONCURRENCY"`
	WEB_APP_STATIC_PATH string `mapstructure:"WEB_APP_STATIC_PATH"`

	ELK_LOG_INDEX            string   `mapstructure:"ELK_LOG_INDEX"`
	ELK_LOG_ADDRS            []string `mapstructure:"ELK_LOG_ADDRS"`
	ELK_LOG_AUTH_USERNAME    string   `mapstructure:"ELK_LOG_AUTH_USERNAME"`
	ELK_LOG_AUTH_PASSWORD    string   `mapstructure:"ELK_LOG_AUTH_PASSWORD"`
	ELK_LOG_MAX_HANDLERS     int      `mapstructure:"ELK_LOG_MAX_HANDLERS"`
	ELK_LOG_SCHEDULE_TIMER   int      `mapstructure:"ELK_LOG_SCHEDULE_TIMER"` // seconds
	ELK_RESPONSE_BUFFER_SIZE int      `mapstructure:"ELK_RESPONSE_BUFFER_SIZE"`

	SCHEDULER_TICKER int `mapstructure:"SCHEDULER_TICKER"` // seconds

	MONGO_URL     string `mapstructure:"MONGO_URL"`
	MONGO_DB_NAME string `mapstructure:"MONGO_DB_NAME"`

	POSTGRES_DNS string `mapstructure:"POSTGRES_DNS"`

	REDIS_URL string `mapstructure:"REDIS_URL"`

	DISCORD_WEBHOOK_URL string `mapstructure:"DISCORD_WEBHOOK_URL"`

	PROXY_ENABLE bool   `mapstructure:"PROXY_ENABLE"`
	PROXY_ADDR   string `mapstructure:"PROXY_ADDR"`
	NO_PROXY     string `mapstructure:"NO_PROXY"`
}

func (e *EnvConfig) ConfigureProxy() {

}

func (e *EnvConfig) SetDefault() {
	if e.ELK_RESPONSE_BUFFER_SIZE == 0 {
		e.ELK_RESPONSE_BUFFER_SIZE = 1000
	}
	if e.SCHEDULER_TICKER == 0 {
		e.SCHEDULER_TICKER = 30
	}
	if e.PROXY_ADDR == "" {
		e.PROXY_ADDR = "http://proxy.hcm.fpt.vn:80"
	}
	if e.NO_PROXY == "" {
		e.NO_PROXY = "localhost,127.0.0.1,172.27.228.171,172.24.224.113,em-dev-redis-headless"
	}
	if e.ELK_LOG_INDEX == "" {
		e.ELK_LOG_INDEX = "filebeat-em-prod-8.17.0-2025*"
	}
	if e.ELK_LOG_MAX_HANDLERS == 0 {
		e.ELK_LOG_MAX_HANDLERS = 1
	}
	if e.ELK_LOG_SCHEDULE_TIMER == 0 {
		e.ELK_LOG_SCHEDULE_TIMER = 30
	}
	if e.WEB_APP_HOST == "" {
		e.WEB_APP_HOST = "0.0.0.0"
	}
	if e.WEB_APP_PORT == 0 {
		e.WEB_APP_PORT = 6001
	}
	if e.WEB_APP_STATIC_PATH == "" {
		e.WEB_APP_STATIC_PATH = "./static"
	}
	log.Info().Str("POSTGRES_DNS", e.POSTGRES_DNS).Msg("init env")
}

var envConfig *EnvConfig

func GetEnv() *EnvConfig {
	if envConfig == nil {
		envConfig = &EnvConfig{}
		goscs.LoadEnvVars[*EnvConfig](envConfig, ".env", "./")
		// log.Info().Interface("envConfig", envConfig).Msg("Loaded env config")

		// set default values after loading env vars
		envConfig.SetDefault()
	}
	return envConfig
}
