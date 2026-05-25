# hensparkIO-payload

`hensparkIO-payload` 是一组 Go 结构体和构造函数，用来统一 WebSocket/消息通道里的请求参数、响应信封、ACK、心跳、行情特征、预测结果、Tick、订单簿、账户、组合、Feed、订单状态、系统指标和 AlphaLab 数据格式。

模块名：

```go
import "hensparkIO-payload/..."
```

当前主要 package：

| package | 用途 |
| --- | --- |
| `request` | 认证请求、订阅/取消订阅请求 |
| `protocol` | 通用响应信封、频道、方法、ACK、心跳 |
| `model` | 业务数据模型：行情、预测、Feed、账户、组合、订单、AlphaLab |
| `utils` | 时间窗口工具 |

> 注意：当前代码里的 `model/order.go` 引用了 `hensparkIO-payload/response`，但仓库中还没有 `response` package，所以 `go test ./...` 会因为缺少该 package 失败。`OrderState.Account` 的类型依赖这个缺失包。

## 安装和使用

在同一个 Go module 中直接导入：

```go
import (
	"hensparkIO-payload/model"
	"hensparkIO-payload/protocol"
	"hensparkIO-payload/request"
)
```

如果作为外部依赖使用，需要先把仓库发布到可访问的 Git 地址，然后：

```bash
go get <your-git-host>/hensparkIO-payload
```

## 通用约定

### 时间字段

多数时间字段使用毫秒时间戳：

| 字段 | 含义 |
| --- | --- |
| `event_time_ms` | 事件发生时间，Unix milliseconds |
| `window_start_ms` | 窗口开始时间，Unix milliseconds |
| `window_end_ms` | 窗口结束时间，Unix milliseconds |
| `inference_time_ms` | 推理时间，Unix milliseconds |
| `expires_at_ms` | 预测过期时间，Unix milliseconds |
| `update_time_ms` | 更新时间，Unix milliseconds |

### TimeFrame 字符串

项目中常见的 timeframe 示例：

```text
1S, 5S, 10S, 30S, 1MIN, 5MIN, 10MIN,30MIN,1H,6H,12H,18H,1D
```

`utils.GenTimeFrame(ms int64)` 会把毫秒差转换成字符串：

| 输入 | 输出 |
| --- | --- |
| `1000` | `1S` |
| `60000` | `1MIN` |
| `3600000` | `1H` |
| `86400000` | `24H` |
| `90000` | `90000MS` |
| `0` 或负数 | 空字符串 |

调用示例：

```go
tf := utils.GenTimeFrame(5 * 60 * 1000)
// tf == "5MIN"
```

> 注意：`WithMarketFeatureWindow`、`WithMarketFeatureWindowMs`、`WithPredictsWindow`、`WithPredictsWindowStartMs` 当前使用 `start - end` 生成 timeframe。如果 `start < end`，自动生成的 `time_frame` 会是空字符串。`utils.GenTimeFrame` 会先匹配小时再匹配天，所以 24 小时会输出 `24H`，不是 `1D`。需要时可以用 `WithTimeFrame("1H")` 手动覆盖预测数据的 timeframe。

## protocol package

### Method

订阅请求方法：

| 常量 | JSON 值 | 含义 |
| --- | --- | --- |
| `protocol.MethodSubscribe` | `subscribe` | 订阅 |
| `protocol.MethodUnsubscribe` | `unsubscribe` | 取消订阅 |

### Channel

频道常量：

| 常量 | JSON 值 | 含义 |
| --- | --- | --- |
| `protocol.ChannelMarketFeature` | `market_feature` | 市场特征 |
| `protocol.ChannelPrediction` | `prediction` | 预测结果 |
| `protocol.ChannelTick` | `tick` | Tick / OHLCV / 成交 / 订单簿 |
| `protocol.ChannelAIFeed` | `ai_feed` | AI 信息流 |
| `protocol.ChannelNews` | `news` | 新闻 |
| `protocol.ChannelAlphaLab` | `alpha_lab` | AlphaLab 数据 |

### Envelope

`Envelope[T]` 是统一响应外壳：

```go
type Envelope[T any] struct {
	Code        int               `json:"code"`
	Content     []T               `json:"data"`
	ContentType ContentType       `json:"content_type"`
	Message     map[string]string `json:"message,omitempty"`
	ErrorCode   *int              `json:"error_code,omitempty"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `code` | `int` | 是 | 状态码，例如 `200` / `400` |
| `data` | `[]T` | 否 | 业务数据数组 |
| `content_type` | `ContentType` | 是 | 数据类型枚举 |
| `message` | `map[string]string` | 否 | 附加消息 |
| `error_code` | `*int` | 否 | 业务错误码 |

状态码常量：

| 常量 | 值 |
| --- | --- |
| `protocol.OK200` | `200` |
| `protocol.Error400` | `400` |

`ContentType` 枚举值：

| 常量 | 值 | 对应数据 |
| --- | --- | --- |
| `protocol.MarketFeatureContentType` | `0` | `model.MarketFeature` |
| `protocol.PredictsContentType` | `1` | `model.Predicts` |
| `protocol.OHLCVContentType` | `2` | `model.OHLCV` |
| `protocol.TradesContentType` | `3` | `model.Trade` |
| `protocol.OrderBookContentType` | `4` | `model.OrderBook` |
| `protocol.AccountStatusContentType` | `5` | `model.AccountStatus` |
| `protocol.EventAckContentType` | `6` | `protocol.EventAck` |
| `protocol.PortfolioContentType` | `7` | `model.Portfolio` |
| `protocol.AlphaLab` | `8` | `model.AlphaLab` |
| `protocol.HeartBeatContentType` | `9` | `protocol.HeartBeat` |
| `protocol.FeedContentType` | `10` | `model.Feeds` |
| `protocol.OrderStateContentType` | `11` | `model.OrderState` |
| `protocol.SystemFeatureContentType` | `12` | `protocol.SysFeature` |

创建响应：

```go
env := protocol.NewEnvelope[model.MarketFeature](
	protocol.OK200,
	protocol.WithMarketFeatures(feature1, feature2),
	protocol.WithMessage[model.MarketFeature]("request_id", "req-001"),
)
```

可用 option：

| 函数 | 参数 | 说明 |
| --- | --- | --- |
| `protocol.WithContent[T](data ...T)` | 任意允许的内容类型 | 只追加数据，不设置 `content_type` |
| `protocol.WithMarketFeatures(data ...model.MarketFeature)` | 市场特征 | 追加数据并设置 `content_type=0` |
| `protocol.WithPredicts(data ...model.Predicts)` | 预测结果 | 追加数据并设置 `content_type=1` |
| `protocol.WithOHLCV(data ...model.OHLCV)` | K 线数据 | 追加数据并设置 `content_type=2` |
| `protocol.WithTrades(data ...model.Trade)` | 成交数据 | 追加数据并设置 `content_type=3` |
| `protocol.WithOrderBooks(data ...model.OrderBook)` | 订单簿 | 追加数据并设置 `content_type=4` |
| `protocol.WithAccountStatus(data ...model.AccountStatus)` | 账户状态 | 追加数据并设置 `content_type=5` |
| `protocol.WithEventAck(event ...protocol.EventAck)` | ACK | 追加数据并设置 `content_type=6` |
| `protocol.WithPortfolios(data ...model.Portfolio)` | 组合数据 | 追加数据并设置 `content_type=7` |
| `protocol.WithAlphaLab(data model.AlphaLab)` | AlphaLab | 追加单条数据并设置 `content_type=8` |
| `protocol.WithHeartBeat(event string)` | 事件名 | 创建心跳数据并设置 `content_type=9` |
| `protocol.WithFeeds(feed ...model.Feeds)` | Feed 数据 | 追加数据并设置 `content_type=10` |
| `protocol.WithOrderState(state ...model.OrderState)` | 订单状态 | 追加数据并设置 `content_type=11` |
| `protocol.WithSystemMetrics(f ...protocol.SysFeature)` | 系统指标 | 追加数据并设置 `content_type=12` |
| `protocol.WithMessage[T](key, value string)` | 消息键值 | 添加 `message` |
| `protocol.WithErrorMessage[T](errorCode int, errorMessage ...string)` | 错误码、错误消息 | 设置 `code=400`、`error_code`、`message.error_message` |

错误响应示例：

```go
env := protocol.NewEnvelope[protocol.EventAck](
	protocol.Error400,
	protocol.WithErrorMessage[protocol.EventAck](10001, "invalid request"),
)
```

输出类似：

```json
{
  "code": 400,
  "data": null,
  "content_type": 0,
  "message": {
    "error_message": "invalid request"
  },
  "error_code": 10001
}
```

### EventAck

订阅/取消订阅确认消息：

```go
type EventAck struct {
	Method    Method     `json:"method"`
	IsSuccess bool       `json:"is_success"`
	ErrorInfo *ErrorInfo `json:"error_info,omitempty"`
	Channel   Channel    `json:"channel,omitempty"`
	Symbols   []string   `json:"symbols,omitempty"`
	TimeFrame []string   `json:"time_frame,omitempty"`
	RequestID string     `json:"request_id"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `method` | `protocol.Method` | 是 | `subscribe` 或 `unsubscribe` |
| `is_success` | `bool` | 是 | 是否成功 |
| `error_info` | `*ErrorInfo` | 否 | 错误信息 |
| `channel` | `protocol.Channel` | 否 | 确认的频道 |
| `symbols` | `[]string` | 否 | 确认的交易对 |
| `time_frame` | `[]string` | 否 | 确认的时间窗口 |
| `request_id` | `string` | 是 | 请求 ID |

创建成功 ACK：

```go
ack := protocol.NewEventAck(
	protocol.MethodSubscribe,
	protocol.WithAckChannel(protocol.ChannelMarketFeature),
	protocol.WithAckRequestID("req-001"),
	protocol.WithAckSymbols("BTC-USD", "ETH-USD"),
	protocol.WithAckTimeFrame("1S", "1MIN"),
)

env := protocol.NewEnvelope[protocol.EventAck](
	protocol.OK200,
	protocol.WithEventAck(ack),
)
```

创建失败 ACK：

```go
ack := protocol.NewEventAck(
	protocol.MethodSubscribe,
	protocol.WithAckChannel(protocol.ChannelAIFeed),
	protocol.WithAckRequestID("req-002"),
	protocol.WithErrorAck(403, "permission denied"),
)
```

`EventAck` option：

| 函数 | 参数 | 说明 |
| --- | --- | --- |
| `protocol.WithAckSymbols(symbols ...string)` | 交易对列表 | 设置 `symbols` |
| `protocol.WithAckTimeFrame(timeFrame ...string)` | timeframe 列表 | 设置 `time_frame` |
| `protocol.WithAckRequestID(requestID string)` | 请求 ID | 设置 `request_id` |
| `protocol.WithAckChannel(channel protocol.Channel)` | 频道 | 设置 `channel` |
| `protocol.WithErrorAck(errCode int, errMessage string)` | 错误码、错误消息 | 设置 `is_success=false` 和 `error_info` |

### HeartBeat

心跳数据：

```go
type HeartBeat struct {
	Event string `json:"event"`
	Ts    int64  `json:"ts"`
}
```

调用：

```go
env := protocol.NewEnvelope[protocol.HeartBeat](
	protocol.OK200,
	protocol.WithHeartBeat("ping"),
)
```

输出类似：

```json
{
  "code": 200,
  "data": [
    {
      "event": "ping",
      "ts": 1710000000000
    }
  ],
  "content_type": 9
}
```

### SysFeature

系统延迟指标：

```go
type SysFeature struct {
	OverallLatency     int64 `json:"overall_latency"`
	CollectionLatency  int64 `json:"collection_latency"`
	ComputationLatency int64 `json:"computation_latency"`
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `overall_latency` | 总体延迟 |
| `collection_latency` | 采集延迟 |
| `computation_latency` | 计算延迟 |

调用：

```go
metrics := protocol.SysFeature{
	OverallLatency:     120,
	CollectionLatency:  40,
	ComputationLatency: 80,
}

env := protocol.NewEnvelope[protocol.SysFeature](
	protocol.OK200,
	protocol.WithSystemMetrics(metrics),
)
```

## request package

### Auth

认证请求体：

```go
type Auth struct {
	Account     *string `json:"account"`
	Password    *string `json:"password"`
	MobilePhone *string `json:"mobile_phone"`
	VerifyCode  *string `json:"verify_code"`
	VerifyBy    string  `json:"verify_by"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `account` | `*string` | 否 | 账号登录时使用 |
| `password` | `*string` | 否 | 密码 |
| `mobile_phone` | `*string` | 否 | 手机号登录时使用 |
| `verify_code` | `*string` | 否 | 验证码 |
| `verify_by` | `string` | 是 | 验证方式，例如 `password` / `sms`，具体取值由服务端约定 |

账号密码示例：

```go
account := "demo"
password := "secret"

auth := request.Auth{
	Account:  &account,
	Password: &password,
	VerifyBy: "password",
}
```

验证码示例：

```go
phone := "+8613800000000"
code := "123456"

auth := request.Auth{
	MobilePhone: &phone,
	VerifyCode:  &code,
	VerifyBy:    "sms",
}
```

### SubscriptionRequest

订阅/取消订阅请求：

```go
type SubscriptionRequest struct {
	Method         protocol.Method       `json:"method"`
	Channels       []ChannelSubscription `json:"channels"`
	UpdateInterval *int64                `json:"update_interval,omitempty"`
	ReturnType     *string               `json:"return_type,omitempty"`
	RequestId      *string               `json:"request_id,omitempty"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `method` | `protocol.Method` | 是 | `subscribe` 或 `unsubscribe` |
| `channels` | `[]ChannelSubscription` | 是 | 频道订阅列表 |
| `update_interval` | `*int64` | 否 | 更新间隔，当前 `WithUpdateInterval` 写入的是微秒值 |
| `return_type` | `*string` | 否 | 返回类型，当前没有内置 option，需要手动赋值 |
| `request_id` | `*string` | 否 | 请求 ID，`NewSubscriptionRequest` 会自动生成 UUID |

`ChannelSubscription`：

```go
type ChannelSubscription struct {
	Channel    protocol.Channel `json:"channel"`
	Symbols    []string         `json:"symbols,omitempty"`
	TimeFrames []string         `json:"time_frames,omitempty"`
	Exchange   *string          `json:"exchange,omitempty"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `channel` | `protocol.Channel` | 是 | 频道 |
| `symbols` | `[]string` | 否 | 交易对 |
| `time_frames` | `[]string` | 否 | 时间窗口 |
| `exchange` | `*string` | 否 | 交易所 |

创建订阅请求：

```go
req := request.NewSubscriptionRequest(
	protocol.MethodSubscribe,
	request.WithTicks(
		request.WithSymbols("BTC-USD", "ETH-USD"),
		request.WithTimeFrames("1S", "5S", "1MIN"),
		request.WithExchange("Kraken"),
	),
	request.WithUpdateInterval(100*time.Millisecond),
)
```

输出类似：

```json
{
  "method": "subscribe",
  "channels": [
    {
      "channel": "tick",
      "symbols": ["BTC-USD", "ETH-USD"],
      "time_frames": ["1S", "5S", "1MIN"],
      "exchange": "Kraken"
    }
  ],
  "update_interval": 100000,
  "request_id": "自动生成的 UUID"
}
```

> 注意：`WithUpdateInterval(100*time.Millisecond)` 当前使用 `duration.Microseconds()`，所以 JSON 里是 `100000`，单位是微秒，不是毫秒。

取消订阅：

```go
req := request.NewSubscriptionRequest(
	protocol.MethodUnsubscribe,
	request.WithMarketFeatures(
		request.WithSymbols("BTC-USD"),
		request.WithTimeFrames("1MIN"),
	),
)
```

支持的订阅 option：

| 函数 | 参数 | 对应频道 |
| --- | --- | --- |
| `request.WithChannel(channel, opts...)` | 任意 `protocol.Channel` | 自定义频道 |
| `request.WithMarketFeatures(opts...)` | `ChannelOption` | `market_feature` |
| `request.WithPredictions(opts...)` | `ChannelOption` | `prediction` |
| `request.WithTicks(opts...)` | `ChannelOption` | `tick` |
| `request.WithAIFeeds(opts...)` | `ChannelOption` | `ai_feed` |
| `request.WithNews(opts...)` | `ChannelOption` | `news` |
| `request.WithAlphaLab(opts...)` | `ChannelOption` | `alpha_lab` |
| `request.WithUpdateInterval(duration)` | `time.Duration` | 设置更新间隔 |

支持的频道参数 option：

| 函数 | 参数 | 说明 |
| --- | --- | --- |
| `request.WithSymbols(symbols ...string)` | 交易对 | 设置 `symbols` |
| `request.WithTimeFrames(timeFrames ...string)` | timeframe | 设置 `time_frames` |
| `request.WithExchange(exchange string)` | 交易所 | 设置 `exchange` |

多频道订阅：

```go
req := request.NewSubscriptionRequest(
	protocol.MethodSubscribe,
	request.WithMarketFeatures(
		request.WithSymbols("BTC-USD"),
		request.WithTimeFrames("1MIN"),
		request.WithExchange("Binance"),
	),
	request.WithPredictions(
		request.WithSymbols("BTC-USD"),
		request.WithTimeFrames("1H"),
	),
	request.WithNews(),
)
```

## model package

### MarketFeature

市场特征数据：

```go
type MarketFeature struct {
	Symbol        string             `json:"symbol"`
	WindowStartMs int64              `json:"window_start_ms"`
	WindowEndMs   int64              `json:"window_end_ms"`
	Data          map[string]float64 `json:"data,omitempty"`
	FeatureType   FeatureType        `json:"feature_type"`
	TimeFrame     string             `json:"time_frame"`
	Source        *FeatureSource     `json:"source,omitempty"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `symbol` | `string` | 是 | 交易对，例如 `BTC-USD` |
| `window_start_ms` | `int64` | 是 | 窗口开始时间 |
| `window_end_ms` | `int64` | 是 | 窗口结束时间 |
| `data` | `map[string]float64` | 否 | 特征指标，例如 `VWAP`、`volatility` |
| `feature_type` | `FeatureType` | 是 | 特征类型 |
| `time_frame` | `string` | 否 | 时间窗口 |
| `source` | `*FeatureSource` | 否 | 数据来源 |

`FeatureType`：

| 常量 | JSON 值 |
| --- | --- |
| `model.TradeFeature` | `trade` |
| `model.BookFeature` | `book` |
| `model.CombinedFeature` | `combined` |
| `model.Unknown` | `unknown` |

调用示例：

```go
feature := model.NewMarketFeature(
	"BTC-USD",
	model.CombinedFeature,
	model.WithMarketFeatureWindowMs(1710003600000, 1710000000000),
	model.WithFeatureSource(
		model.NewFeatureSource(
			model.MulticastSource,
			model.WithDominantSource("Kraken"),
			model.WithFeatureSources("Kraken", "Coinbase"),
		),
	),
	model.WithMarketFeatureMetric("VWAP", 68400.12),
	model.WithMarketFeatureMetrics(map[string]float64{
		"volatility": 0.12,
		"momentum":   0.34,
	}),
)

env := protocol.NewEnvelope[model.MarketFeature](
	protocol.OK200,
	protocol.WithMarketFeatures(feature),
)
```

`MarketFeature` option：

| 函数 | 参数 | 说明 |
| --- | --- | --- |
| `model.WithMarketFeatureWindow(start, end time.Time)` | 开始、结束时间 | 设置窗口毫秒时间戳和自动 timeframe |
| `model.WithMarketFeatureWindowMs(startMs, endMs int64)` | 开始、结束毫秒 | 设置窗口和自动 timeframe |
| `model.WithFeatureSource(source *FeatureSource)` | 来源对象 | 设置 `source` |
| `model.WithMarketFeatureMetric(key, value)` | 指标名、指标值 | 添加单个指标 |
| `model.WithMarketFeatureMetrics(metrics map[string]float64)` | 指标 map | 批量添加指标 |

### FeatureSource

市场特征来源：

```go
type FeatureSource struct {
	Source         []string          `json:"source,omitempty"`
	Type           FeatureSourceType `json:"type"`
	DominantSource *string           `json:"dominant_source"`
}
```

`FeatureSourceType`：

| 常量 | JSON 值 |
| --- | --- |
| `model.SingleExchange` | `single_exchange` |
| `model.CrossExchange` | `cross_exchange` |
| `model.Derived` | `derived` |
| `model.MulticastSource` | `multicast_source` |
| `model.Synthetic` | `synthetic` |
| `model.UnknowFs` | `unknow_fs` |

调用：

```go
source := model.NewFeatureSource(
	model.CrossExchange,
	model.WithFeatureSources("Binance", "OKX"),
	model.WithDominantSource("Binance"),
)
```

`FeatureSource` option：

| 函数 | 参数 | 说明 |
| --- | --- | --- |
| `model.WithFeatureSources(sources ...string)` | 来源列表 | 设置 `source` |
| `model.WithDominantSource(ds string)` | 主来源 | 设置 `dominant_source` |

### Predicts

预测结果：

```go
type Predicts struct {
	Symbol          string            `json:"symbol"`
	TimeFrame       string            `json:"time_frame"`
	WindowStartMs   int64             `json:"window_start_ms"`
	WindowEndMs     int64             `json:"window_end_ms"`
	Metrics         any               `json:"metrics"`
	InferenceTimeMs int64             `json:"inference_time_ms"`
	ExpiresAtMs     int64             `json:"expires_at_ms"`
	Meta            map[string]string `json:"meta,omitempty"`
}
```

字段说明：

| 字段 | 类型 | 必填 | 说明 |
| --- | --- | --- | --- |
| `symbol` | `string` | 是 | 交易对 |
| `time_frame` | `string` | 是 | 时间窗口 |
| `window_start_ms` | `int64` | 是 | 预测窗口开始 |
| `window_end_ms` | `int64` | 是 | 预测窗口结束 |
| `metrics` | `any` | 是 | 指标，可用 `PredictsOHLCV` 或 `PredictsTrend` |
| `inference_time_ms` | `int64` | 是 | 推理时间 |
| `expires_at_ms` | `int64` | 是 | 过期时间 |
| `meta` | `map[string]string` | 否 | 元信息 |

OHLCV 预测指标：

```go
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
```

趋势预测指标：

```go
type PredictsTrend struct {
	PUp   float64 `json:"p_up"`
	PDown float64 `json:"p_down"`
}
```

调用示例：

```go
prediction := model.NewPredicts(
	model.WithSymbol("BTC-USD"),
	model.WithPredictsWindowStartMs(1710003600000, 1710000000000),
	model.WithTimeFrame("1H"),
	model.WithPredictsMetrics(model.PredictsTrend{
		PUp:   0.62,
		PDown: 0.38,
	}),
	model.WithInferenceTime(time.Now(), 5*time.Minute),
	model.WithMeta("model", "trend-v1"),
)

env := protocol.NewEnvelope[model.Predicts](
	protocol.OK200,
	protocol.WithPredicts(prediction),
)
```

读取 metrics：

```go
trend, ok := model.GetMetrics[model.PredictsTrend](prediction)
if ok {
	_ = trend.PUp
}
```

`Predicts` option：

| 函数 | 参数 | 说明 |
| --- | --- | --- |
| `model.WithSymbol(symbol string)` | 交易对 | 设置 `symbol` |
| `model.WithPredictsWindow(start, end time.Time)` | 开始、结束时间 | 设置窗口毫秒时间戳和自动 timeframe |
| `model.WithPredictsWindowStartMs(startMs, endMs int64)` | 开始、结束毫秒 | 设置窗口毫秒时间戳和自动 timeframe |
| `model.WithTimeFrame(timeFrame string)` | timeframe | 手动设置 `time_frame` |
| `model.WithInferenceTime(inferenceTime, expiresDuration)` | 推理时间、有效期 | 设置 `inference_time_ms` 和 `expires_at_ms` |
| `model.WithExpiresAt(expiresAt time.Time)` | 过期时间 | 只设置 `expires_at_ms` |
| `model.WithMeta(key, value string)` | 元信息键值 | 添加 `meta` |
| `model.WithPredictsMetrics[T](metrics T)` | `PredictsOHLCV` 或 `PredictsTrend` | 设置 `metrics` |

### Trade

成交数据：

```go
type Trade struct {
	Symbol      string  `json:"symbol"`
	Px          float64 `json:"px"`
	Qty         float64 `json:"qty"`
	Side        string  `json:"side"`
	EventTimeMs int64   `json:"event_time_ms"`
	Venue       string  `json:"venue"`
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `symbol` | 交易对 |
| `px` | 成交价格 |
| `qty` | 成交数量 |
| `side` | 成交方向，例如 `buy` / `sell` |
| `event_time_ms` | 成交发生时间 |
| `venue` | 交易所，例如 `OKX` / `Binance` / `Kraken` |

调用：

```go
trade := model.Trade{
	Symbol:      "BTC-USD",
	Px:          68400.12,
	Qty:         0.5,
	Side:        "buy",
	EventTimeMs: time.Now().UnixMilli(),
	Venue:       "Binance",
}

env := protocol.NewEnvelope[model.Trade](
	protocol.OK200,
	protocol.WithTrades(trade),
)
```

### OrderBook

订单簿数据：

```go
type OrderBook struct {
	Symbol      string  `json:"symbol"`
	BidLevels   []Level `json:"bid_levels"`
	AskLevels   []Level `json:"ask_levels"`
	EventTimeMs int64   `json:"event_time_ms"`
	Venue       string  `json:"venue"`
}

type Level struct {
	Px  float64 `json:"px"`
	Qty float64 `json:"qty"`
}
```

调用：

```go
book := model.OrderBook{
	Symbol: "BTC-USD",
	BidLevels: []model.Level{
		{Px: 68399.9, Qty: 1.2},
	},
	AskLevels: []model.Level{
		{Px: 68400.1, Qty: 0.8},
	},
	EventTimeMs: time.Now().UnixMilli(),
	Venue:       "OKX",
}

env := protocol.NewEnvelope[model.OrderBook](
	protocol.OK200,
	protocol.WithOrderBooks(book),
)
```

### OHLCV

K 线数据：

```go
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
	Venue         string  `json:"venue"`
}
```

调用：

```go
ohlcv := model.OHLCV{
	Symbol:        "BTC-USD",
	Open:          68000,
	High:          68500,
	Low:           67900,
	Close:         68400,
	Volume:        123.45,
	TradeCount:    520,
	VWAP:          68250,
	WindowStartMs: 1710000000000,
	WindowEndMs:   1710000060000,
	EventTimeMs:   time.Now().UnixMilli(),
	Venue:         "Kraken",
}

env := protocol.NewEnvelope[model.OHLCV](
	protocol.OK200,
	protocol.WithOHLCV(ohlcv),
)
```

### Feeds

信息流/新闻/市场信号：

```go
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
```

`FeedType`：

| 常量 | JSON 值 |
| --- | --- |
| `model.FeedTypeAIAnalysis` | `ai_analysis` |
| `model.FeedTypeNews` | `news` |
| `model.FeedTypeMarketSignal` | `market_signal` |
| `model.FeedTypeMarketEvent` | `market_event` |

调用：

```go
title := "BTC momentum changed"

feed := model.Feeds{
	Id:          "feed-001",
	Type:        model.FeedTypeMarketSignal,
	Title:       &title,
	Text:        "BTC short-term momentum increased.",
	Source:      "alpha-engine",
	EventTimeMs: time.Now().UnixMilli(),
	Priority:    2,
	Tags:        []string{"BTC", "momentum"},
	Meta: map[string]string{
		"model": "signal-v1",
	},
}

env := protocol.NewEnvelope[model.Feeds](
	protocol.OK200,
	protocol.WithFeeds(feed),
)
```

### AccountStatus

账户状态：

```go
type AccountStatus struct {
	Username string `json:"user_name"`
	UserId   string `json:"user_id"`
	UserIcon string `json:"user_icon"`
}
```

调用：

```go
status := model.AccountStatus{
	Username: "demo",
	UserId:   "u-001",
	UserIcon: "https://example.com/avatar.png",
}

env := protocol.NewEnvelope[model.AccountStatus](
	protocol.OK200,
	protocol.WithAccountStatus(status),
)
```

### Portfolio

账户组合数据：

```go
type Portfolio struct {
	Holdings    []Holding      `json:"holdings,omitempty"`
	CcyBalance  []CcyBalance   `json:"ccy_balance,omitempty"`
	RealizedPnl *float64       `json:"realized_pnl,omitempty"`
	Account     *AccountStatus `json:"account,omitempty"`
}
```

持仓：

```go
type Holding struct {
	InstId        string  `json:"inst_id"`
	Size          float64 `json:"size"`
	AvgPrice      float64 `json:"avg_price"`
	UnrealizedPnl float64 `json:"unrealized_pnl"`
}
```

币种余额：

```go
type CcyBalance struct {
	CurrencyCode string  `json:"currency_code"`
	Balance      float64 `json:"balance"`
}
```

调用：

```go
realized := 123.45

portfolio := model.Portfolio{
	Holdings: []model.Holding{
		{
			InstId:        "BTC-USD",
			Size:          0.8,
			AvgPrice:      65000,
			UnrealizedPnl: 1200,
		},
	},
	CcyBalance: []model.CcyBalance{
		{CurrencyCode: "USD", Balance: 10000},
	},
	RealizedPnl: &realized,
	Account: &model.AccountStatus{
		Username: "demo",
		UserId:   "u-001",
	},
}

env := protocol.NewEnvelope[model.Portfolio](
	protocol.OK200,
	protocol.WithPortfolios(portfolio),
)
```

### OrderState

订单状态：

```go
type OrderState struct {
	OrderId     string            `json:"order_id"`
	Status      string            `json:"status"`
	FilledPrice float64           `json:"filled_price"`
	FilledSize  float64           `json:"filled_size"`
	Side        string            `json:"side"`
	EventTime   int64             `json:"event_time"`
	Account     *response.Account `json:"account"`
}
```

字段说明：

| 字段 | 说明 |
| --- | --- |
| `order_id` | 订单 ID |
| `status` | 订单状态，具体取值由业务系统约定 |
| `filled_price` | 已成交均价 |
| `filled_size` | 已成交数量 |
| `side` | 买卖方向，例如 `buy` / `sell` |
| `event_time` | 事件时间 |
| `account` | 账户信息，依赖缺失的 `response.Account` |

调用：

```go
state := model.OrderState{
	OrderId:     "ord-001",
	Status:      "filled",
	FilledPrice: 68400.12,
	FilledSize:  0.1,
	Side:        "buy",
	EventTime:   time.Now().UnixMilli(),
}

env := protocol.NewEnvelope[model.OrderState](
	protocol.OK200,
	protocol.WithOrderState(state),
)
```

> 注意：`model.OrderState` 的 `Account` 字段依赖当前仓库缺失的 `hensparkIO-payload/response` package。补齐该 package 前，包含 `model/order.go` 的项目编译会失败。

### AlphaLab

AlphaLab 策略数据：

```go
type AlphaLab struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`
	StrategyId    string            `json:"strategy_id"`
	InitialEquity float64           `json:"initial_equity"`
	TotalEquity   float64           `json:"total_equity"`
	PnLRatio      float64           `json:"pnl_ratio"`
	PnL           float64           `json:"pnl"`
	MaxDrawdown   float64           `json:"max_drawdown"`
	SharpeRatio   float64           `json:"sharpe_ratio"`
	Positions     []LabPosition     `json:"positions"`
	UpdateTimeMs  int64             `json:"update_time_ms"`
	Meta          map[string]string `json:"meta,omitempty"`
}
```

仓位：

```go
type LabPosition struct {
	Symbol        string  `json:"symbol"`
	PosSide       string  `json:"pos_side"`
	Pnl           float64 `json:"pnl"`
	PnlRatio      float64 `json:"pnl_ratio"`
	PositionValue float64 `json:"position_value"`
	PositionRatio float64 `json:"position_ratio"`
	EquityRatio   float64 `json:"equity_ratio"`
	Venue         string  `json:"venue"`
}
```

仓位方向：

| 常量 | JSON 值 |
| --- | --- |
| `model.Long` | `long` |
| `model.Short` | `short` |
| `model.Flat` | `flat` |

调用：

```go
lab := model.AlphaLab{
	Id:            "lab-001",
	Name:          "BTC Trend Strategy",
	StrategyId:    "strategy-001",
	InitialEquity: 100000,
	TotalEquity:   112000,
	PnLRatio:      0.12,
	PnL:           12000,
	MaxDrawdown:   0.08,
	SharpeRatio:   1.7,
	Positions: []model.LabPosition{
		{
			Symbol:        "BTC-USD",
			PosSide:       string(model.Long),
			Pnl:           3000,
			PnlRatio:      0.05,
			PositionValue: 50000,
			PositionRatio: 0.45,
			EquityRatio:   0.44,
			Venue:         "Binance",
		},
	},
	UpdateTimeMs: time.Now().UnixMilli(),
	Meta: map[string]string{
		"risk_level": "medium",
	},
}

env := protocol.NewEnvelope[model.AlphaLab](
	protocol.OK200,
	protocol.WithAlphaLab(lab),
)
```

## 完整示例

### 订阅 tick 数据

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"hensparkIO-payload/protocol"
	"hensparkIO-payload/request"
)

func main() {
	req := request.NewSubscriptionRequest(
		protocol.MethodSubscribe,
		request.WithTicks(
			request.WithSymbols("BTC-USD", "ETH-USD"),
			request.WithTimeFrames("1S", "1MIN"),
			request.WithExchange("Kraken"),
		),
		request.WithUpdateInterval(100*time.Millisecond),
	)

	b, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println(string(b))
}
```

### 返回市场特征响应

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"hensparkIO-payload/model"
	"hensparkIO-payload/protocol"
)

func main() {
	now := time.Now()

	feature := model.NewMarketFeature(
		"BTC-USD",
		model.CombinedFeature,
		model.WithMarketFeatureWindowMs(
			now.Add(time.Hour).UnixMilli(),
			now.UnixMilli(),
		),
		model.WithMarketFeatureMetric("VWAP", 68400.12),
		model.WithMarketFeatureMetric("volatility", 0.12),
	)

	env := protocol.NewEnvelope[model.MarketFeature](
		protocol.OK200,
		protocol.WithMarketFeatures(feature),
	)

	b, _ := json.MarshalIndent(env, "", "  ")
	fmt.Println(string(b))
}
```

### 返回预测响应

```go
package main

import (
	"encoding/json"
	"fmt"
	"time"

	"hensparkIO-payload/model"
	"hensparkIO-payload/protocol"
)

func main() {
	now := time.Now()

	p := model.NewPredicts(
		model.WithSymbol("BTC-USD"),
		model.WithPredictsWindowStartMs(
			now.Add(time.Hour).UnixMilli(),
			now.UnixMilli(),
		),
		model.WithTimeFrame("1H"),
		model.WithPredictsMetrics(model.PredictsOHLCV{
			OpenPrice:  68000,
			HighPrice:  69000,
			LowPrice:   67500,
			ClosePrice: 68800,
			Volume:     120.5,
			VWAP:       68420,
			Momentum:   0.35,
			Volatility: 0.12,
		}),
		model.WithInferenceTime(now, 5*time.Minute),
	)

	env := protocol.NewEnvelope[model.Predicts](
		protocol.OK200,
		protocol.WithPredicts(p),
	)

	b, _ := json.MarshalIndent(env, "", "  ")
	fmt.Println(string(b))
}
```

### 返回订阅 ACK

```go
ack := protocol.NewEventAck(
	protocol.MethodSubscribe,
	protocol.WithAckRequestID("req-001"),
	protocol.WithAckChannel(protocol.ChannelTick),
	protocol.WithAckSymbols("BTC-USD", "ETH-USD"),
	protocol.WithAckTimeFrame("1S", "1MIN"),
)

env := protocol.NewEnvelope[protocol.EventAck](
	protocol.OK200,
	protocol.WithEventAck(ack),
)
```

### 返回心跳

```go
env := protocol.NewEnvelope[protocol.HeartBeat](
	protocol.OK200,
	protocol.WithHeartBeat("ping"),
)
```

## 当前已知注意事项

1. `hensparkIO-payload/response` package 缺失，导致 `model.OrderState` 相关代码无法完整编译。
2. `WithUpdateInterval` 的注释写“更新频率”，传入 `time.Duration`，但实际写入 `duration.Microseconds()`，字段名又叫 `update_interval`。调用方需要统一单位约定。
3. timeframe 自动生成当前使用 `start - end`。如果业务语义是“开始时间早于结束时间”，自动生成会得到空字符串，建议调用时手动 `WithTimeFrame`，或在代码里改为 `end - start`。
4. `model.Feeds`、`model.OrderState`、`protocol.SysFeature` 已经加入 `protocol.ContentTypeInterface`，可以分别通过 `WithFeeds`、`WithOrderState`、`WithSystemMetrics` 包装进 `Envelope`。
