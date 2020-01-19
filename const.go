package market_center

import "fmt"

type Type int

const (
	_ Type = iota
	Type_GetSupportList
	Type_SubscribeDepth
	Type_SubscribeTicker
	Type_GetDepth
	Type_GetTicker
)

func (t Type) String() string {
	if t > 0 && int(t) <= len(typeSymbol) {
		return typeSymbol[t-1]
	}
	return fmt.Sprintf("UNKNOWN_METHOD_TYPE (%d)", t)
}

var typeSymbol = [...]string{"GetSupportList", "SubscribeDepth", "SubscribeTicker", "GetDepth", "GetTicker"}

const (
	ErrMsg_ExchangeNotSupport = "exchange(%s) not support"
	ErrMsg_ExchangeNoDepth    = "exchange(%s) no %s depth data"
	ErrMsg_ExchangeNoTicker   = "exchange(%s) no %s ticker data"

	ErrMsg_RequestFormatError    = "request format error: %s"
	ErrMsg_RequestTypeNotSupport = "request type(%v) not support"
)

const (
	UDS_PATH = "/tmp/goex.market.center"
)

type DataFlag int

const (
	DataFlag_Depth DataFlag = 1 << iota
	DataFlag_Ticker

	DataFlag_All = DataFlag_Depth | DataFlag_Ticker
)
