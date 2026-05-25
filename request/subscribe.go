package request

import (
	"time"

	"github.com/google/uuid"
	"github.com/simonks2016/hensparkIO-payload/protocol"
)

type SubscriptionRequest struct {
	Method         protocol.Method       `json:"method"`
	Channels       []ChannelSubscription `json:"channels"`
	UpdateInterval *int64                `json:"update_interval,omitempty"`
	ReturnType     *string               `json:"return_type,omitempty"`
	RequestId      *string               `json:"request_id,omitempty"`
}

type ChannelSubscription struct {
	Channel     protocol.Channel  `json:"channel"`
	Symbols     []string          `json:"symbols,omitempty"`
	TimeFrames  []string          `json:"time_frames,omitempty"`
	Exchange    *string           `json:"exchange,omitempty"`
	Preferences map[string]string `json:"preferences,omitempty"`
}

type SubscriptionOption func(*SubscriptionRequest)

func NewSubscriptionRequest(
	method protocol.Method,
	opts ...SubscriptionOption,
) SubscriptionRequest {

	requestId := uuid.New().String()

	req := SubscriptionRequest{
		Method:    method,
		Channels:  make([]ChannelSubscription, 0),
		RequestId: &requestId,
	}

	for _, opt := range opts {
		opt(&req)
	}

	return req
}

// WithUpdateInterval 设置更新频率
func WithUpdateInterval(duration time.Duration) SubscriptionOption {
	return func(req *SubscriptionRequest) {
		ms := duration.Milliseconds()
		req.UpdateInterval = &ms
	}
}

// WithChannel 设置订阅频道
func WithChannel(channel protocol.Channel, opts ...ChannelOption) SubscriptionOption {
	return func(req *SubscriptionRequest) {
		sub := ChannelSubscription{
			Channel: channel,
		}
		for _, opt := range opts {
			opt(&sub)
		}
		req.Channels = append(req.Channels, sub)
	}
}

// WithMarketFeatures 订阅市场特征频道
func WithMarketFeatures(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelMarketFeature, opts...)
}

// WithPredictions 订阅预测结果
func WithPredictions(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelPrediction, opts...)
}

// WithTicks 订阅Tick数据
func WithTicks(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelTick, opts...)
}

// WithAIFeeds 订阅信息流
func WithAIFeeds(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelAIFeed, opts...)
}

// WithNews 订阅新闻
func WithNews(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelNews, opts...)
}

// WithAlphaLab 订阅alpha实验室数据
func WithAlphaLab(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelAlphaLab, opts...)
}

type ChannelOption func(*ChannelSubscription)

func WithSymbols(symbols ...string) ChannelOption {
	return func(sub *ChannelSubscription) {
		sub.Symbols = append(sub.Symbols, symbols...)
	}
}

func WithTimeFrames(timeFrames ...string) ChannelOption {
	return func(sub *ChannelSubscription) {
		sub.TimeFrames = append(sub.TimeFrames, timeFrames...)
	}
}

func WithExchange(exchange string) ChannelOption {
	return func(sub *ChannelSubscription) {
		sub.Exchange = &exchange
	}
}

func WithPreferredLanguage(lang protocol.Lang) ChannelOption {

	return func(sub *ChannelSubscription) {
		if sub.Preferences == nil {
			sub.Preferences = make(map[string]string)
		}
		sub.Preferences["lang"] = string(lang)
	}
}
