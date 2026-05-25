package model

type FeedType string

const (
	FeedTypeAIAnalysis   FeedType = "ai_analysis"
	FeedTypeNews         FeedType = "news"
	FeedTypeMarketSignal FeedType = "market_signal"
	FeedTypeMarketEvent  FeedType = "market_event"
)

type Feeds struct {
	Id          string            `json:"id"`
	Type        FeedType          `json:"type"`
	Title       *string           `json:"title,omitempty"`
	Text        string            `json:"text"`
	Source      string            `json:"source"`
	EventTimeMs int64             `json:"event_time_ms"`
	Priority    int               `json:"priority,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
}
