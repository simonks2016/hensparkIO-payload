package model

type Portfolio struct {
	Holdings    []Holding      `json:"holdings,omitempty"`
	CcyBalance  []CcyBalance   `json:"ccy_balance,omitempty"`
	RealizedPnl *float64       `json:"realized_pnl,omitempty"`
	Account     *AccountStatus `json:"account,omitempty"`
}
type Holding struct {
	InstId        string  `json:"inst_id"`
	Size          float64 `json:"size"`
	AvgPrice      float64 `json:"avg_price"`
	UnrealizedPnl float64 `json:"unrealized_pnl"`
}
type CcyBalance struct {
	CurrencyCode string  `json:"currency_code"`
	Balance      float64 `json:"balance"`
}
