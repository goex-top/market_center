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

// spot
func (d *Data) UpdateSpotDepth(exchangeName, pair string, depth *goex.Depth) {
	key := key(exchangeName, pair)
	d.depths.Store(key, depth)
}

func (d *Data) GetSpotDepth(exchangeName, pair string) (*goex.Depth, error) {
	key := key(exchangeName, pair)
	dep, isOk := d.depths.Load(key)
	if isOk {
		return dep.(*goex.Depth), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoDepth, exchangeName, pair)
	}
}

func (d *Data) UpdateSpotTicker(exchangeName, pair string, ticker *goex.Ticker) {
	key := key(exchangeName, pair)
	d.tickers.Store(key, ticker)
}

func (d *Data) GetSpotTicker(exchangeName, pair string) (*goex.Ticker, error) {
	key := key(exchangeName, pair)
	ticker, isOk := d.tickers.Load(key)
	if isOk {
		return ticker.(*goex.Ticker), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoTicker, exchangeName, pair)
	}
}

// future
func (d *Data) UpdateFutureDepth(exchangeName, contractType, pair string, depth *goex.Depth) {
	key := futureKey(exchangeName, contractType, pair)
	d.depths.Store(key, depth)
}

func (d *Data) GetFutureDepth(exchangeName, contractType, pair string) (*goex.Depth, error) {
	key := futureKey(exchangeName, contractType, pair)
	dep, isOk := d.depths.Load(key)
	if isOk {
		return dep.(*goex.Depth), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoDepth, exchangeName, pair)
	}
}

func (d *Data) UpdateFutureTicker(exchangeName, contractType, pair string, ticker *goex.Ticker) {
	key := futureKey(exchangeName, contractType, pair)
	d.tickers.Store(key, ticker)
}

func (d *Data) GetFutureTicker(exchangeName, contractType, pair string) (*goex.Ticker, error) {
	key := futureKey(exchangeName, contractType, pair)
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

func futureKey(exchangeName, contactType, pair string) string {
	return exchangeName + "/" + contactType + "/" + pair
}
