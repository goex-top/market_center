package market_center

import (
	"github.com/nntaoli-project/goex"
	"strings"
)

const (
	POLONIEX = "Poloniex"
	BITSTAMP = "Bitstamp"
	HUOBI    = "Huobi"
	BITFINEX = "Bitfinex"
	OKEX     = "Okex"
	BINANCE  = "Binance"
	BITTREX  = "Bittrex"
	BITHUMB  = "Bithumb"
	GDAX     = "Gdax"
	GATEIO   = "Gateio"
	ZB       = "Zb"
	COINEX   = "Coinex"
	BIGONE   = "Bigone"
	HITBTC   = "Hitbtc"

	FUTURE_BITMEX   = "Future_Bitmex"
	FUTURE_OKEX     = "Future_Okex"
	SWAP_OKEX       = "Swap_Okex"
	FUTURE_HBDM     = "Future_Hbdm"
	SWAP_BINANCE    = "Swap_Binance"
	FUTURE_COINBENE = "Future_Coinbene"
)

var (
	SupportList = []string{
		POLONIEX, BITSTAMP, HUOBI, BITFINEX, OKEX, BINANCE, BITTREX, BITHUMB, GDAX, GATEIO, ZB, COINEX, BIGONE, HITBTC, FUTURE_BITMEX,
		FUTURE_OKEX, SWAP_OKEX, FUTURE_HBDM, SWAP_BINANCE,
	}
	SupportAdapter = map[string]string{
		POLONIEX: goex.POLONIEX,
		BITSTAMP: goex.BITSTAMP,
		HUOBI:    goex.HUOBI_PRO,
		BITFINEX: goex.BITFINEX,
		OKEX:     goex.OKEX_V3,
		BINANCE:  goex.BINANCE,
		BITTREX:  goex.BITTREX,
		BITHUMB:  goex.BITHUMB,
		GDAX:     goex.GDAX,
		GATEIO:   goex.GATEIO,
		ZB:       goex.ZB,
		COINEX:   goex.COINEX,
		BIGONE:   goex.BIGONE,
		HITBTC:   goex.HITBTC,

		FUTURE_BITMEX:   goex.BITMEX,
		FUTURE_OKEX:     goex.OKEX_FUTURE,
		SWAP_OKEX:       goex.OKEX_SWAP,
		FUTURE_HBDM:     goex.HBDM,
		SWAP_BINANCE:    goex.BINANCE_SWAP,
		FUTURE_COINBENE: goex.COINBENE,
	}
	//FutureList = []string{
	//	goex.BITMEX,
	//	goex.OKEX_FUTURE,
	//	goex.OKEX_SWAP,
	//	goex.HBDM,
	//	goex.COINBENE,
	//	goex.FMEX,
	//	goex.BINANCE_SWAP,
	//}
	//SpotList = []string{
	//	goex.POLONIEX,
	//	goex.BITSTAMP,
	//	goex.HUOBI_PRO,
	//	goex.BITFINEX,
	//	goex.OKEX_V3,
	//	goex.BINANCE,
	//	goex.BITTREX,
	//	goex.BITHUMB,
	//	goex.GDAX,
	//	goex.GATEIO,
	//	goex.ZB,
	//	goex.COINEX,
	//	goex.FCOIN,
	//	goex.FCOIN_MARGIN,
	//	goex.BIGONE,
	//	goex.HITBTC,
	//	//goex.BITMEX,
	//	//goex.OKEX_FUTURE,
	//	//goex.OKEX_SWAP,
	//	//goex.HBDM,
	//	//goex.COINBENE,
	//	//goex.FMEX,
	//}
)

func IsFutureExchange(exchangeName string) bool {
	if strings.Contains(exchangeName, "Future") {
		return true
	} else if strings.Contains(exchangeName, "Swap") {
		return true
	}
	return false
}
