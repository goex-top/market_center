package market_center

import "fmt"

type Type int

const (
	_ Type = iota
	Type_GetSupportList
	Type_SubscribeSpotDepth
	Type_SubscribeSpotTicker
	Type_SubscribeFutureDepth
	Type_SubscribeFutureTicker
	Type_GetSpotDepth
	Type_GetSpotTicker
	Type_GetFutureDepth
	Type_GetFutureTicker
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
	DataFlag_Trade
	DataFlag_Kline

	//DataFlag_All = DataFlag_Depth | DataFlag_Ticker | DataFlag_Trade | DataFlag_Kline
)

func (df DataFlag) String() string {
	switch df {
	case DataFlag_Depth:
		return "depth"
	case DataFlag_Ticker:
		return "ticker"
	case DataFlag_Trade:
		return "trade"
	case DataFlag_Kline:
		return "kline"
	default:
		return "unknown"
	}
}

func ParseDataFlag(str string) DataFlag {
	switch str {
	case "depth":
		return DataFlag_Depth
	case "ticker":
		return DataFlag_Ticker
	case "trade":
		return DataFlag_Trade
	case "kline":
		return DataFlag_Kline
	default:
		return 0
	}
}
