package data

import (
	"fmt"
	. "github.com/goex-top/market_center"
	goex "github.com/nntaoli-project/GoEx"
	"sync"
)

type Data struct {
	depths  sync.Map
	tickers sync.Map
}

func NewData() *Data {
	return &Data{
		depths:  sync.Map{},
		tickers: sync.Map{},
	}
}

func (d *Data) UpdateDepth(exchangeName, pair string, depth *goex.Depth) {
	key := key(exchangeName, pair)
	d.depths.Store(key, depth)
}

func (d *Data) GetDepth(exchangeName, pair string) (*goex.Depth, error) {
	key := key(exchangeName, pair)
	dep, isOk := d.depths.Load(key)
	if isOk {
		return dep.(*goex.Depth), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoDepth, exchangeName, pair)
	}
}

func (d *Data) UpdateTicker(exchangeName, pair string, ticker *goex.Ticker) {
	key := key(exchangeName, pair)
	d.tickers.Store(key, ticker)
}

func (d *Data) GetTicker(exchangeName, pair string) (*goex.Ticker, error) {
	key := key(exchangeName, pair)
	ticker, isOk := d.tickers.Load(key)
	if isOk {
		return ticker.(*goex.Ticker), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoTicker, exchangeName, pair)
	}
}

func key(exchangeName, pair string) string {
	return exchangeName + "/" + pair
}
