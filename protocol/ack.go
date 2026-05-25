package protocol

type EventAck struct {
	Method    Method     `json:"method"`
	IsSuccess bool       `json:"is_success"`
	ErrorInfo *ErrorInfo `json:"error_info,omitempty"`
	Channel   Channel    `json:"channel,omitempty"`
	Symbols   []string   `json:"symbols,omitempty"`
	TimeFrame []string   `json:"time_frame,omitempty"`
	RequestID string     `json:"request_id"`
}

type ErrorInfo struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EventAckOption func(*EventAck)

func WithAckSymbols(symbols ...string) EventAckOption {
	return func(e *EventAck) {
		e.Symbols = append(e.Symbols, symbols...)
	}
}

func WithAckTimeFrame(timeFrame ...string) EventAckOption {
	return func(e *EventAck) {
		e.TimeFrame = append(e.TimeFrame, timeFrame...)
	}
}

func WithAckRequestID(requestID string) EventAckOption {
	return func(e *EventAck) {
		e.RequestID = requestID
	}
}

func WithAckChannel(channel Channel) EventAckOption {
	return func(e *EventAck) {
		e.Channel = channel
	}
}
func WithErrorAck(errCode int, errMessage string) EventAckOption {
	return func(e *EventAck) {
		e.ErrorInfo = &ErrorInfo{
			Code:    errCode,
			Message: errMessage,
		}
		e.IsSuccess = false
	}
}

func NewEventAck(method Method, opts ...EventAckOption) EventAck {
	e := EventAck{
		Method:    method,
		IsSuccess: true,
		ErrorInfo: nil,
		Channel:   "",
		Symbols:   nil,
		TimeFrame: nil,
		RequestID: "",
	}
	for _, opt := range opts {
		opt(&e)
	}
	return e
}
