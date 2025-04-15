package elkcollector

import (
	"fmt"
	"nhandd/bego/internal/elkclient"
	"nhandd/bego/internal/models"
	"strings"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog/log"
)

type ElkQuery struct {
	Code  string
	Index string
	Query string
}

type ElkConfig struct {
	Code      string
	Host      string
	Port      int32
	Username  string
	Password  string
	ProxyAddr string
}

func (e *ElkConfig) buildUrl() string {
	var url string = ""
	if e.Port != 0 {
		url = fmt.Sprintf("%s:%d", e.Host, e.Port)
	} else {
		url = e.Host
	}
	if !(strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")) {
		url = fmt.Sprintf("http://%s", url)
	}
	return url
}

type ElkCollector struct {
	Code      string
	Queries   map[string]ElkQuery // map code with ElkQuery
	ElkConfig ElkConfig
	Interval  int32
	LastRunAt time.Time
	NextRunAt time.Time
	ElkClient *elasticsearch.Client

	mutex sync.Mutex
}

func NewElkCollectorFromDbConf(dbConf models.DatasourceElkcollectorconfig) (newCollector *ElkCollector) {
	var collectorQueries = make(map[string]ElkQuery)
	for _, query := range dbConf.ElkQueries {
		collectorQueries[query.Code] = ElkQuery{
			Code:  query.Code,
			Index: query.Index,
			Query: query.Query,
		}
	}
	newElkConfig := ElkConfig{
		Code:      dbConf.ElkConfig.Code,
		Host:      dbConf.ElkConfig.Host,
		Port:      dbConf.ElkConfig.Port,
		Username:  dbConf.ElkConfig.Username,
		Password:  dbConf.ElkConfig.Password,
		ProxyAddr: dbConf.ElkConfig.Proxy.BuildProxyUrl(),
	}
	newElkClient, _ := elkclient.NewElkClient(
		[]string{newElkConfig.buildUrl()},
		newElkConfig.Username,
		newElkConfig.Password,
		nil,
	)

	newCollector = &ElkCollector{
		ElkClient: newElkClient,
		Code:      dbConf.Code,
		Queries:   collectorQueries,
		ElkConfig: newElkConfig,
		Interval:  dbConf.GetInterval(),
		LastRunAt: dbConf.LastRunAt,
	}

	log.Info().Str("code", newCollector.Code).Msg("New ElkCollector created")
	return
}

func (e *ElkCollector) updateLastRunAt(atTime time.Time) error {
	if atTime.IsZero() {
		atTime = time.Now()
	}
	e.LastRunAt = atTime
	return pgRepo.UpdateLastRunAt(e.Code, e.LastRunAt)
}

func (e *ElkCollector) calculateQueryTimeRange() (time.Time, time.Time) {
	now := time.Now().UTC()
	if e.LastRunAt.IsZero() || e.LastRunAt.After(now) {
		e.LastRunAt = now.Add(-time.Duration(e.Interval) * time.Second)
		e.NextRunAt = now
	} else if e.LastRunAt.Add(time.Duration(e.Interval) * time.Second).After(now) {
		e.NextRunAt = now
		return e.LastRunAt, e.NextRunAt
	} else if e.LastRunAt.Add(300 * time.Second).Before(now) {
		e.NextRunAt = e.LastRunAt.Add(300 * time.Second)
		return e.LastRunAt, e.NextRunAt
	} else if e.NextRunAt.IsZero() || e.NextRunAt.Before(e.LastRunAt) {
		if e.LastRunAt.Add(time.Duration(e.Interval) * time.Second).After(now) {
			e.NextRunAt = now
		} else {
			e.NextRunAt = e.LastRunAt.Add(300 * time.Second) // run with log in 5 minutes
		}
	}
	e.LastRunAt = e.LastRunAt.Add(-5 * time.Second)
	return e.LastRunAt, e.NextRunAt
}

func (e *ElkCollector) GetQueryTimeRange() (time.Time, time.Time) {
	return e.LastRunAt, e.NextRunAt
}

func (e *ElkCollector) Info() string {
	return fmt.Sprintf("ElkCollector{Code: %s, Interval: %d, LastRunAt: %s}", e.Code, e.Interval, e.LastRunAt.Format(time.RFC3339))
}

func (e *ElkCollector) collectWithQuery(query ElkQuery) error {
	builtQuery := elkQueryBuilder.Build(query.Code, query.Query, e)
	if builtQuery == "" {
		log.Error().Str("code", e.Code).Str("query", query.Code).Msg("Failed to build query")
		return fmt.Errorf("failed to build query for %s", query.Code)
	}
	res := elkclient.ScrollLogWithCode(e.Code, query.Code, e.ElkClient, query.Index, builtQuery, elkRespChan)
	// log.Info().Str("code", e.Code).Str("query", query.Code).Msg("Collecting log with query")
	return res
}

func (e *ElkCollector) AllowToRun() bool {
	if !e.mutex.TryLock() {
		// Nếu không khóa được, nghĩa là có goroutine khác đang chạy
		return false
	}
	// Mở khóa ngay lập tức để không ảnh hưởng đến luồng chính
	e.mutex.Unlock()

	if e.LastRunAt.IsZero() {
		return true
	}
	if e.NextRunAt.IsZero() {
		return true
	}
	if e.NextRunAt.Before(time.Now()) {
		return true
	}
	return false
}

func (e *ElkCollector) CollectLog() error {
	if !e.mutex.TryLock() {
		// Nếu không khóa được, nghĩa là có goroutine khác đang chạy
		return fmt.Errorf("another log collection is already in progress")
	}
	// Đảm bảo mutex được mở khóa khi hàm kết thúc
	defer e.mutex.Unlock()

	e.calculateQueryTimeRange()
	for _, query := range e.Queries {
		e.collectWithQuery(query)
	}
	e.updateLastRunAt(e.NextRunAt)
	return nil
}

func (e *ElkCollector) ReloadData() error {
	return nil
}

func (e *ElkCollector) ReloadLastRunAt() error {
	return nil
}

func (e *ElkCollector) Shutdown() error {
	// log.Info().Str("code", e.Code).Msg("Shutting down ElkCollector")
	// elkclient.CloseClientWithScrolls(e.ElkClient)
	return nil
}
