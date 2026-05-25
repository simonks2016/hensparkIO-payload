package protocol

import (
	"strings"

	"github.com/simonks2016/hensparkIO-payload/model"
)

const (
	OK200    int = 200
	Error400 int = 400
	Error403 int = 403
)

type ContentType int

const (
	MarketFeatureContentType ContentType = iota
	PredictsContentType
	OHLCVContentType
	TradesContentType
	OrderBookContentType
	AccountStatusContentType
	EventAckContentType
	PortfolioContentType
	AlphaLab
	HeartBeatContentType
	FeedContentType
	OrderStateContentType
	SystemFeatureContentType
)

type Envelope[T any] struct {
	Code        int               `json:"code"`
	Content     []T               `json:"data"`
	ContentType ContentType       `json:"content_type"`
	Message     map[string]string `json:"message,omitempty"`
	ErrorCode   *int              `json:"error_code,omitempty"`
}

type ContentTypeInterface interface {
	model.MarketFeature | model.Predicts | model.OHLCV | model.Trade | model.OrderBook | model.Portfolio |
		model.AccountStatus | model.AlphaLab | EventAck | HeartBeat | model.Feeds | model.OrderState | SysFeature
}

type EnvelopeOption[T ContentTypeInterface] func(envelope *Envelope[T])

func NewEnvelope[T ContentTypeInterface](Code int, opts ...EnvelopeOption[T]) *Envelope[T] {
	env := &Envelope[T]{
		Code:      Code,
		Content:   nil,
		Message:   nil,
		ErrorCode: nil,
	}

	for _, opt := range opts {
		opt(env)
	}
	return env
}

func WithContent[T ContentTypeInterface](data ...T) EnvelopeOption[T] {
	return func(envelope *Envelope[T]) {
		envelope.Content = append(envelope.Content, data...)
	}
}

func WithMarketFeatures(data ...model.MarketFeature) EnvelopeOption[model.MarketFeature] {
	return func(envelope *Envelope[model.MarketFeature]) {
		envelope.Content = append(envelope.Content, data...)
		envelope.ContentType = MarketFeatureContentType
	}
}

func WithPredicts(data ...model.Predicts) EnvelopeOption[model.Predicts] {
	return func(envelope *Envelope[model.Predicts]) {
		envelope.Content = append(envelope.Content, data...)
		envelope.ContentType = PredictsContentType
	}
}

func WithOHLCV(data ...model.OHLCV) EnvelopeOption[model.OHLCV] {
	return func(envelope *Envelope[model.OHLCV]) {
		envelope.Content = append(envelope.Content, data...)
		envelope.ContentType = OHLCVContentType
	}
}

func WithTrades(data ...model.Trade) EnvelopeOption[model.Trade] {
	return func(envelope *Envelope[model.Trade]) {
		envelope.Content = append(envelope.Content, data...)
		envelope.ContentType = TradesContentType
	}
}
func WithOrderBooks(data ...model.OrderBook) EnvelopeOption[model.OrderBook] {
	return func(envelope *Envelope[model.OrderBook]) {
		envelope.Content = append(envelope.Content, data...)
		envelope.ContentType = OrderBookContentType
	}
}

func WithAccountStatus(data ...model.AccountStatus) EnvelopeOption[model.AccountStatus] {
	return func(envelope *Envelope[model.AccountStatus]) {
		envelope.Content = append(envelope.Content, data...)
		envelope.ContentType = AccountStatusContentType
	}
}

func WithEventAck(event ...EventAck) EnvelopeOption[EventAck] {
	return func(envelope *Envelope[EventAck]) {
		envelope.Content = append(envelope.Content, event...)
		envelope.ContentType = EventAckContentType
	}
}

func WithPortfolios(Portfolio ...model.Portfolio) EnvelopeOption[model.Portfolio] {
	return func(envelope *Envelope[model.Portfolio]) {
		envelope.Content = append(envelope.Content, Portfolio...)
		envelope.ContentType = PortfolioContentType
	}
}

func WithAlphaLab(data model.AlphaLab) EnvelopeOption[model.AlphaLab] {
	return func(envelope *Envelope[model.AlphaLab]) {
		envelope.Content = append(envelope.Content, data)
		envelope.ContentType = AlphaLab
	}
}

func WithErrorMessage[T ContentTypeInterface](errorCode int, errorMessage ...string) EnvelopeOption[T] {
	return func(envelope *Envelope[T]) {
		envelope.Code = 400
		envelope.ErrorCode = &errorCode
		if len(errorMessage) > 0 {
			if envelope.Message == nil {
				envelope.Message = make(map[string]string)
			}
			envelope.Message = map[string]string{
				"error_message": strings.Join(errorMessage, ","),
			}
		}
	}
}

func WithMessage[T ContentTypeInterface](key, value string) EnvelopeOption[T] {
	return func(envelope *Envelope[T]) {
		if envelope.Message == nil {
			envelope.Message = make(map[string]string)
		}
		envelope.Message[key] = value
	}
}

func WithHeartBeat(event string) EnvelopeOption[HeartBeat] {
	return func(envelope *Envelope[HeartBeat]) {
		envelope.Content = append(envelope.Content, NewHeartBeat(event))
		envelope.ContentType = HeartBeatContentType
	}
}

func WithFeeds(feed ...model.Feeds) EnvelopeOption[model.Feeds] {
	return func(envelope *Envelope[model.Feeds]) {
		envelope.Content = append(envelope.Content, feed...)
		envelope.ContentType = FeedContentType
	}
}

func WithOrderState(state ...model.OrderState) EnvelopeOption[model.OrderState] {
	return func(envelope *Envelope[model.OrderState]) {
		envelope.Content = append(envelope.Content, state...)
		envelope.ContentType = OrderStateContentType
	}
}

func WithSystemMetrics(f ...SysFeature) EnvelopeOption[SysFeature] {
	return func(envelope *Envelope[SysFeature]) {
		envelope.Content = append(envelope.Content, f...)
		envelope.ContentType = SystemFeatureContentType
	}
}
