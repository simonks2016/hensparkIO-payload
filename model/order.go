package model

import "hensparkIO-payload/response"

type OrderState struct {
	OrderId     string            `json:"order_id"`
	Status      string            `json:"status"`
	FilledPrice float64           `json:"filled_price"`
	FilledSize  float64           `json:"filled_size"`
	Side        string            `json:"side"`
	EventTime   int64             `json:"event_time"`
	Account     *response.Account `json:"account"`
}
