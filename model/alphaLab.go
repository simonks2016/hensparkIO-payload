package model

type AlphaLab struct {
	Id            string            `json:"id"`
	Name          string            `json:"name"`           // 显示名称
	StrategyId    string            `json:"strategy_id"`    // 策略ID
	InitialEquity float64           `json:"initial_equity"` // 初始资金
	TotalEquity   float64           `json:"total_equity"`   // 当前总资产
	PnLRatio      float64           `json:"pnl_ratio"`      // 当前收益率
	PnL           float64           `json:"pnl"`            // 收益金额
	MaxDrawdown   float64           `json:"max_drawdown"`   // 最大回撤
	SharpeRatio   float64           `json:"sharpe_ratio"`   // 夏普利率
	Positions     []LabPosition     `json:"positions"`      // 仓位
	UpdateTimeMs  int64             `json:"update_time_ms"` // 更新时间
	Meta          map[string]string `json:"meta,omitempty"` // meta信息
}

type LabPosition struct {
	Symbol        string  `json:"symbol"`
	PosSide       string  `json:"pos_side"`
	Pnl           float64 `json:"pnl"`            // 收益金额
	PnlRatio      float64 `json:"pnl_ratio"`      // 当前收益率
	PositionValue float64 `json:"position_value"` // 持仓市值
	PositionRatio float64 `json:"position_ratio"` // 持仓占总仓位比例
	EquityRatio   float64 `json:"equity_ratio"`   // 仓位市值占总市值多少
	Venue         string  `json:"venue"`          // 交易所
}

type PosSide string

const (
	Long  PosSide = "long"
	Short PosSide = "short"
	Flat  PosSide = "flat"
)
