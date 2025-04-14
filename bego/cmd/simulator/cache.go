package simulator

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

const _key_exc_sent_discord = "em-sent-discord"
const _key_exc_by_logid = "em-log-exceptions"
const _key_exc_by_agid = "em-agid-exceptions"
const _key_exc_sorted = "em-sorted-log-exceptions"

type LogEntry struct {
	Agid      string `json:"agid"`
	Message   string `json:"message"`
	ExcInfo   string `json:"exc_info"`
	LogId     string `json:"log_id"`
	Loc       string `json:"loc"`
	Timestamp string `json:"timestamp"`
}

func GetLogsByAgid(client *redis.Client, sortedSetKey string, hashKeyPrefix string, limit int64, offset int64) ([]LogEntry, error) {
	ctx := context.Background()

	// Lấy tất cả logId và score từ sorted set theo thứ tự giảm dần
	zItems, err := client.ZRevRangeWithScores(ctx, sortedSetKey, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}

	// Map để giữ log có score cao nhất cho mỗi agid
	agidMap := make(map[string]LogEntry)
	scoreMap := make(map[string]float64)

	// Pipeline để lấy dữ liệu từ HSET hiệu quả hơn
	pipe := client.Pipeline()
	cmdMap := make(map[string]*redis.StringCmd)

	// Tạo lệnh HGET cho từng logId
	for _, z := range zItems {
		logId := z.Member.(string)
		cmdMap[logId] = pipe.HGet(ctx, hashKeyPrefix, logId)
	}

	// Thực thi pipeline
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	// Xử lý kết quả
	for _, z := range zItems {
		logId := z.Member.(string)
		jsonStr, err := cmdMap[logId].Result()
		if err != nil {
			continue // Bỏ qua nếu không tìm thấy log
		}

		var logEntry LogEntry
		if err := json.Unmarshal([]byte(jsonStr), &logEntry); err != nil {
			continue // Bỏ qua nếu JSON không hợp lệ
		}

		// Nếu chưa có agid này hoặc score hiện tại cao hơn score đã lưu
		if currentScore, exists := scoreMap[logEntry.Agid]; !exists || z.Score > currentScore {
			agidMap[logEntry.Agid] = logEntry
			scoreMap[logEntry.Agid] = z.Score
		}
	}

	// Chuyển map thành slice để trả về
	result := make([]LogEntry, 0, len(agidMap))
	for _, log := range agidMap {
		result = append(result, log)
	}

	// Sắp xếp theo timestamp nếu muốn (tuỳ chọn)
	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp > result[j].Timestamp
	})

	return result, nil
}

func convertUnixFloatToLocalTime(unixFloat float64, loc *time.Location) time.Time {
	utcTime := time.Unix(int64(unixFloat), 0).UTC()
	localTime := utcTime.In(loc)
	return localTime
}
func CacheStatisticLastWeek(client *redis.Client, userTz string) (map[string]int, error) {
	if userTz == "" {
		userTz = "Asia/Ho_Chi_Minh"
	}
	userLoc, err := time.LoadLocation(userTz)
	if err != nil {
		log.Error().Err(err).Msg("Failed to load location")
		return nil, err
	}

	currentTime := time.Now().In(userLoc)
	lastWeek := currentTime.AddDate(0, 0, -7)
	startOfLastWeek := time.Date(lastWeek.Year(), lastWeek.Month(), lastWeek.Day(), 0, 0, 0, 0, userLoc)
	unixStartOfLastWeek := startOfLastWeek.Unix()

	min := fmt.Sprint(unixStartOfLastWeek)
	max := "+inf"
	// log.Info().Interface("min", min).Interface("max", max).Msg("min | max")
	results, err := client.ZRangeByScoreWithScores(
		context.Background(),
		_key_exc_sorted,
		&redis.ZRangeBy{Min: min, Max: max},
	).Result()
	if err != nil {
		return nil, err
	}

	counter := make(map[string]int)
	// initialize the counter for the last 7 days with default value 0
	for i := 0; i < 7; i++ {
		date := startOfLastWeek.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")
		counter[dateStr] = 0
	}

	for _, z := range results {
		logTime := convertUnixFloatToLocalTime(z.Score, userLoc)
		logDate := logTime.Format("2006-01-02")
		counter[logDate]++
	}
	log.Info().Interface("counter", counter).Msg("counter")
	return counter, nil
}
