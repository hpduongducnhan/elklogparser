package elkcollector

import "encoding/json"

type QueryResult struct {
	Code  string `json:"code"`
	Error error  `json:"error"`
}

type JobResult struct {
	Jid         string                 `json:"jid"`
	Status      string                 `json:"status"`
	QueryResult map[string]QueryResult `json:"query_result"`
}

func (jr *JobResult) AddQueryResult(code string, err error) {
	jr.QueryResult[code] = QueryResult{
		Code:  code,
		Error: err,
	}
}

func (jr *JobResult) GetStatus() string {
	return jr.Status
}

func (jr *JobResult) IsSuccess() bool {
	for _, qr := range jr.QueryResult {
		if qr.Error != nil {
			return false
		}
	}
	return true
}

func (jr *JobResult) ExportDetail() string {
	jsonBytes, err := json.Marshal(jr)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func NewJobResult(jid string) *JobResult {
	return &JobResult{
		Jid:         jid,
		Status:      "running",
		QueryResult: make(map[string]QueryResult),
	}
}
