package model

import (
	"encoding/json"
	"hensparkIO-payload/utils"
	"time"
)

type MarketFeature struct {
	Symbol        string             `json:"symbol"`
	WindowStartMs int64              `json:"window_start_ms"`
	WindowEndMs   int64              `json:"window_end_ms"`
	Data          map[string]float64 `json:"data,omitempty"`
	FeatureType   FeatureType        `json:"feature_type"`
	TimeFrame     string             `json:"time_frame"`
	Source        *FeatureSource     `json:"source,omitempty"`
}

type MarketFeatureOption func(*MarketFeature)

func WithMarketFeatureWindow(start, end time.Time) MarketFeatureOption {
	return func(feature *MarketFeature) {
		feature.WindowStartMs = start.UnixMilli()
		feature.WindowEndMs = end.UnixMilli()
		feature.TimeFrame = utils.GenTimeFrame(end.UnixMilli() - start.UnixMilli())
	}
}

func WithMarketFeatureWindowMs(windowStartMs, windowEndMs int64) MarketFeatureOption {
	return func(feature *MarketFeature) {
		feature.WindowStartMs = windowStartMs
		feature.WindowEndMs = windowEndMs
		feature.TimeFrame = utils.GenTimeFrame(windowEndMs - windowStartMs)
	}
}

func WithFeatureSource(source *FeatureSource) MarketFeatureOption {
	return func(feature *MarketFeature) {
		feature.Source = source
	}
}

func WithMarketFeatureMetric(
	key string,
	value float64,
) MarketFeatureOption {
	return func(feature *MarketFeature) {
		if feature.Data == nil {
			feature.Data = make(map[string]float64)
		}
		feature.Data[key] = value
	}
}

func WithMarketFeatureMetrics(
	metrics map[string]float64,
) MarketFeatureOption {
	return func(feature *MarketFeature) {
		if feature.Data == nil {
			feature.Data = make(map[string]float64)
		}

		for k, v := range metrics {
			feature.Data[k] = v
		}
	}
}

func (m MarketFeature) String() string {
	jso, _ := json.Marshal(m)
	return string(jso)
}

func NewMarketFeature(symbol string, featureType FeatureType, options ...MarketFeatureOption) MarketFeature {

	m := MarketFeature{
		Symbol:      symbol,
		Data:        make(map[string]float64),
		FeatureType: featureType,
	}
	for _, option := range options {
		option(&m)
	}
	return m
}

type FeatureSource struct {
	Source         []string          `json:"source,omitempty"`
	Type           FeatureSourceType `json:"type"`
	DominantSource *string           `json:"dominant_source"`
}

type FeatureSourceOption func(*FeatureSource)

func WithFeatureSources(sources ...string) FeatureSourceOption {
	return func(feature *FeatureSource) {
		feature.Source = sources
	}
}

func WithDominantSource(ds string) FeatureSourceOption {
	return func(feature *FeatureSource) {
		feature.DominantSource = &ds
	}
}

type FeatureType string

const (
	TradeFeature    FeatureType = "trade"
	BookFeature     FeatureType = "book"
	CombinedFeature FeatureType = "combined"
	Unknown         FeatureType = "unknown"
)

type FeatureSourceType string

const (
	SingleExchange  FeatureSourceType = "single_exchange"
	CrossExchange   FeatureSourceType = "cross_exchange"
	Derived         FeatureSourceType = "derived"
	MulticastSource FeatureSourceType = "multicast_source"
	Synthetic       FeatureSourceType = "synthetic"
	UnknowFs        FeatureSourceType = "unknow_fs"
)

func NewFeatureSource(t FeatureSourceType, opts ...FeatureSourceOption) *FeatureSource {
	fs := &FeatureSource{
		Source:         nil,
		Type:           t,
		DominantSource: nil,
	}

	for _, opt := range opts {
		opt(fs)
	}
	return fs
}
