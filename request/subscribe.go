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
		// 定义一个订阅频道参数
		sub := NewChannelSubscription(channel, opts...)
		// 插入
		req.Channels = append(req.Channels, sub)
	}
}
func NewChannelSubscription(channel protocol.Channel, opts ...ChannelOption) ChannelSubscription {
	sub := ChannelSubscription{
		Channel: channel,
	}
	for _, opt := range opts {
		opt(&sub)
	}
	return sub
}

// WithChannels 批量订阅频道
func WithChannels(requests ...ChannelSubscription) SubscriptionOption {
	return func(req *SubscriptionRequest) {
		req.Channels = append(req.Channels, requests...)
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
	return WithChannels(
		NewChannelSubscription(protocol.ChannelOHLCV, opts...),
		NewChannelSubscription(protocol.ChannelTrades, opts...),
		NewChannelSubscription(protocol.ChannelOrderBook, opts...),
	)
}

// WithTrades 订阅逐笔交易数据
func WithTrades(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelTrades, opts...)
}

// WithOrderBook 订阅盘口数据
func WithOrderBook(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelOrderBook, opts...)
}

// WithOHLCV 订阅K线数据
func WithOHLCV(opts ...ChannelOption) SubscriptionOption {
	return WithChannel(protocol.ChannelOHLCV, opts...)
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

// WithSystemMetrics 订阅系统参数
func WithSystemMetrics() SubscriptionOption {
	return WithChannel(protocol.ChannelSystemMetrics)
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
