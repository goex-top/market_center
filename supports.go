package market_center

import "github.com/nntaoli-project/goex"

var (
	FutureList = []string{
		goex.BITMEX,
		goex.OKEX_FUTURE,
		goex.OKEX_SWAP,
		goex.HBDM,
		goex.COINBENE,
		goex.FMEX,
	}
	SpotList = []string{
		goex.POLONIEX,
		goex.BITSTAMP,
		goex.HUOBI_PRO,
		goex.BITFINEX,
		goex.OKEX_V3,
		goex.BINANCE,
		goex.BITTREX,
		goex.BITHUMB,
		goex.GDAX,
		goex.GATEIO,
		goex.ZB,
		goex.COINEX,
		goex.FCOIN,
		goex.FCOIN_MARGIN,
		goex.BIGONE,
		goex.HITBTC,
		//goex.BITMEX,
		//goex.OKEX_FUTURE,
		//goex.OKEX_SWAP,
		//goex.HBDM,
		//goex.COINBENE,
		//goex.FMEX,
	}
)

func IsFutureExchange(exchangeName string) bool {
	for _, v := range FutureList {
		if exchangeName == v {
			return true
		}
	}
	return false
}
