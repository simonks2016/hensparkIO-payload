package model

import (
	"hensparkIO-payload/utils"
	"time"
)

type Predicts struct {
	Symbol          string            `json:"symbol"`
	TimeFrame       string            `json:"time_frame"`
	WindowStartMs   int64             `json:"window_start_ms"`
	WindowEndMs     int64             `json:"window_end_ms"`
	Metrics         any               `json:"metrics"`
	InferenceTimeMs int64             `json:"inference_time_ms"` // 推理时间
	ExpiresAtMs     int64             `json:"expires_at_ms"`     //有效时间
	Meta            map[string]string `json:"meta,omitempty"`    // 元信息
}

type PredictsOHLCV struct {
	OpenPrice  float64 `json:"open_price"`
	HighPrice  float64 `json:"high_price"`
	LowPrice   float64 `json:"low_price"`
	ClosePrice float64 `json:"close_price"`
	Volume     float64 `json:"volume"`
	VWAP       float64 `json:"vwap"`
	Momentum   float64 `json:"momentum"`
	Volatility float64 `json:"volatility"`
}

type PredictsMetricsType interface {
	PredictsOHLCV | PredictsTrend
}

type PredictsTrend struct {
	PUp   float64 `json:"p_up"`
	PDown float64 `json:"p_down"`
}

type PredictsOption func(*Predicts)

func NewPredicts(opts ...PredictsOption) Predicts {
	p := Predicts{}

	for _, opt := range opts {
		opt(&p)
	}

	return p
}

func WithSymbol(symbol string) PredictsOption {
	return func(predicts *Predicts) {
		predicts.Symbol = symbol
	}
}

func WithPredictsWindow(start, end time.Time) PredictsOption {
	return func(predicts *Predicts) {
		predicts.WindowStartMs = start.UnixMilli()
		predicts.WindowEndMs = end.UnixMilli()
		predicts.TimeFrame = utils.GenTimeFrame(end.UnixMilli() - start.UnixMilli())
	}
}

func WithPredictsWindowStartMs(startMs, endMs int64) PredictsOption {
	return func(predicts *Predicts) {
		predicts.WindowStartMs = startMs
		predicts.WindowEndMs = endMs
		predicts.TimeFrame = utils.GenTimeFrame(endMs - startMs)
	}
}

func WithTimeFrame(timeFrame string) PredictsOption {
	return func(predicts *Predicts) {
		predicts.TimeFrame = timeFrame
	}
}

func WithInferenceTime(inferenceTime time.Time, expiresDuration time.Duration) PredictsOption {
	return func(predicts *Predicts) {
		predicts.InferenceTimeMs = inferenceTime.UnixMilli()
		predicts.ExpiresAtMs = inferenceTime.Add(expiresDuration).UnixMilli()
	}
}

func WithExpiresAt(expiresAt time.Time) PredictsOption {
	return func(predicts *Predicts) {
		predicts.ExpiresAtMs = expiresAt.UnixMilli()
	}
}

func WithMeta(key, value string) PredictsOption {
	return func(predicts *Predicts) {
		if predicts.Meta == nil {
			predicts.Meta = make(map[string]string)
		}
		predicts.Meta[key] = value
	}
}

func WithPredictsMetrics[T PredictsMetricsType](metrics T) PredictsOption {
	return func(predicts *Predicts) {
		predicts.Metrics = metrics
	}
}

func GetMetrics[T PredictsMetricsType](p Predicts) (T, bool) {
	v, ok := p.Metrics.(T)
	return v, ok
}
