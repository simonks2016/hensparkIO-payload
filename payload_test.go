package hensparkIO_payload

import (
	"testing"
	"time"

	"github.com/simonks2016/hensparkIO-payload/model"
	"github.com/simonks2016/hensparkIO-payload/protocol"
	"github.com/simonks2016/hensparkIO-payload/request"
)

func TestPayload(t *testing.T) {

	protocol.NewEnvelope[model.MarketFeature](
		protocol.OK200,
		protocol.WithMarketFeatures(
			model.NewMarketFeature("BTC-USD", model.CombinedFeature,
				model.WithMarketFeatureWindow(time.Now(), time.Now().Add(time.Hour)),
				model.WithFeatureSource(
					model.NewFeatureSource(
						model.MulticastSource,
						model.WithDominantSource("KRAKEN"),
						model.WithFeatureSources("Kraken", "Coinbase"),
					),
				),
				model.WithMarketFeatureMetric("VWAP", 86400),
				model.WithMarketFeatureMetric("", 965),
			),
		),
	)

	protocol.NewEnvelope[model.Predicts](
		protocol.OK200,
		protocol.WithPredicts(
			model.NewPredicts(
				model.WithSymbol("BTC-USD"),
				model.WithPredictsWindow(time.Now(), time.Now().Add(time.Hour)),
				model.WithPredictsMetrics[model.PredictsOHLCV](
					model.PredictsOHLCV{
						OpenPrice:  0,
						HighPrice:  0,
						LowPrice:   0,
						ClosePrice: 0,
						Volume:     0,
						VWAP:       0,
						Momentum:   0,
						Volatility: 0,
					},
				),
				model.WithInferenceTime(time.Now(), time.Second),
			),
		),
	)

	protocol.NewEnvelope[protocol.EventAck](
		200,
		protocol.WithEventAck(
			protocol.NewEventAck(protocol.MethodSubscribe,
				protocol.WithAckChannel(protocol.ChannelMarketFeature),
				protocol.WithAckRequestID(""),
				protocol.WithAckSymbols("BTC-USD", "ETH-USD"),
				protocol.WithAckTimeFrame("1S", "5S", "10S", "30S", "1MIN"),
			)))

	protocol.NewEnvelope[protocol.EventAck](
		protocol.Error403,
		protocol.WithEventAck(
			protocol.NewEventAck(protocol.MethodSubscribe,
				protocol.WithErrorAck(403, ""),
				protocol.WithAckChannel(protocol.ChannelAIFeed),
				protocol.WithAckSymbols("BTC-USD", "ETH-USD"),
				protocol.WithAckRequestID(""),
			)),
	)

	protocol.NewEnvelope[protocol.HeartBeat](
		200,
		protocol.WithHeartBeat("ping"))

	request.NewSubscriptionRequest(
		protocol.MethodSubscribe,
		request.WithTicks(
			request.WithSymbols("BTC-USD", "ETH-USD"),
			request.WithTimeFrames("1S", "5S", "10S", "30S", "1MIN"),
			request.WithExchange("Kraken"),
		),
		request.WithUpdateInterval(time.Millisecond*time.Duration(100)),
	)

}
