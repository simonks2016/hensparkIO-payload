package protocol

import "time"

type HeartBeat struct {
	Event string `json:"event"`
	Ts    int64  `json:"ts"`
}

func NewHeartBeat(event string) HeartBeat {
	return HeartBeat{
		Event: event,
		Ts:    time.Now().UnixMilli(),
	}
}
