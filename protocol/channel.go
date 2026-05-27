package protocol

type Channel string

const (
	ChannelMarketFeature Channel = "market_feature"
	ChannelPrediction    Channel = "prediction"
	ChannelAIFeed        Channel = "ai_feed"
	ChannelNews          Channel = "news"
	ChannelAlphaLab      Channel = "alpha_lab"
	ChannelTrades        Channel = "trades"
	ChannelOHLCV         Channel = "ohlcv"
	ChannelOrderBook     Channel = "order_book"
	ChannelSystemMetrics Channel = "system_metrics"
)

func (channel Channel) String() string {
	return string(channel)
}
