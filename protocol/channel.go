package protocol

type Channel string

const (
	ChannelMarketFeature Channel = "market_feature"
	ChannelPrediction    Channel = "prediction"
	ChannelTick          Channel = "tick"
	ChannelAIFeed        Channel = "ai_feed"
	ChannelNews          Channel = "news"
	ChannelAlphaLab      Channel = "alpha_lab"
)
