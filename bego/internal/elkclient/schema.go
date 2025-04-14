package elkclient

type ElkInnerHit struct {
	ID     string                 `json:"_id"`
	Index  string                 `json:"_index"`
	Source map[string]interface{} `json:"_source"`
	Sort   []float64              `json:"sort"`
}

type ElkHits struct {
	Hits []ElkInnerHit `json:"hits"`
}

type ElkResponse struct {
	ScrollID string  `json:"_scroll_id"`
	Hits     ElkHits `json:"hits"`
}

type ElkResponseWithCode struct {
	CollectorCode string
	QueryCode     string
	*ElkInnerHit
}
