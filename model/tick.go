package model

type Trade struct {
	Symbol      string  `json:"symbol"`
	Px          float64 `json:"px"`            // 成交价格
	Qty         float64 `json:"qty"`           // 成交量
	Side        string  `json:"side"`          // buy/sell
	EventTimeMs int64   `json:"event_time_ms"` // 成交发生时间
	Venue       string  `json:"venue"`         // OKX/Binance/Kraken
}

type OrderBook struct {
	Symbol      string  `json:"symbol"`
	BidLevels   []Level `json:"bid_levels"`
	AskLevels   []Level `json:"ask_levels"`
	EventTimeMs int64   `json:"event_time_ms"`
	Venue       string  `json:"venue"` // 交易所
}

type Level struct {
	Px  float64 `json:"px"`
	Qty float64 `json:"qty"`
}

type OHLCV struct {
	Symbol        string  `json:"symbol"`
	Open          float64 `json:"open"`
	High          float64 `json:"high"`
	Low           float64 `json:"low"`
	Close         float64 `json:"close"`
	Volume        float64 `json:"volume"`
	TradeCount    int64   `json:"trade_count"`
	VWAP          float64 `json:"vwap"`
	WindowStartMs int64   `json:"window_start_ms"`
	WindowEndMs   int64   `json:"window_end_ms"`
	EventTimeMs   int64   `json:"event_time_ms"`
	Venue         string  `json:"venue"` // 交易所
}
