package data

import (
	"fmt"
	. "github.com/goex-top/market_center"
	"github.com/nntaoli-project/goex"
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
func (d *Data) RemoveSpotDepth(exchangeName, pair string) {
	key := keyGen(exchangeName, pair)
	d.depths.Delete(key)
}

func (d *Data) UpdateSpotDepth(exchangeName, pair string, depth *goex.Depth) {
	key := keyGen(exchangeName, pair)
	d.depths.Store(key, depth)
}

func (d *Data) GetSpotDepth(exchangeName, pair string) (*goex.Depth, error) {
	key := keyGen(exchangeName, pair)
	dep, isOk := d.depths.Load(key)
	if isOk {
		return dep.(*goex.Depth), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoDepth, exchangeName, pair)
	}
}

func (d *Data) RemoveSpotTicker(exchangeName, pair string) {
	key := keyGen(exchangeName, pair)
	d.tickers.Delete(key)
}

func (d *Data) UpdateSpotTicker(exchangeName, pair string, ticker *goex.Ticker) {
	key := keyGen(exchangeName, pair)
	d.tickers.Store(key, ticker)
}

func (d *Data) GetSpotTicker(exchangeName, pair string) (*goex.Ticker, error) {
	key := keyGen(exchangeName, pair)
	ticker, isOk := d.tickers.Load(key)
	if isOk {
		return ticker.(*goex.Ticker), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoTicker, exchangeName, pair)
	}
}

// future

func (d *Data) RemoveFutureDepth(exchangeName, contractType, pair string) {
	key := futureKeyGen(exchangeName, contractType, pair)
	d.depths.Delete(key)
}

func (d *Data) UpdateFutureDepth(exchangeName, contractType, pair string, depth *goex.Depth) {
	key := futureKeyGen(exchangeName, contractType, pair)
	d.depths.Store(key, depth)
}

func (d *Data) GetFutureDepth(exchangeName, contractType, pair string) (*goex.Depth, error) {
	key := futureKeyGen(exchangeName, contractType, pair)
	dep, isOk := d.depths.Load(key)
	if isOk {
		return dep.(*goex.Depth), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoDepth, exchangeName, pair)
	}
}

func (d *Data) RemoveFutureTicker(exchangeName, contractType, pair string) {
	key := futureKeyGen(exchangeName, contractType, pair)
	d.tickers.Delete(key)
}

func (d *Data) UpdateFutureTicker(exchangeName, contractType, pair string, ticker *goex.Ticker) {
	key := futureKeyGen(exchangeName, contractType, pair)
	d.tickers.Store(key, ticker)
}

func (d *Data) GetFutureTicker(exchangeName, contractType, pair string) (*goex.Ticker, error) {
	key := futureKeyGen(exchangeName, contractType, pair)
	ticker, isOk := d.tickers.Load(key)
	if isOk {
		return ticker.(*goex.Ticker), nil
	} else {
		return nil, fmt.Errorf(ErrMsg_ExchangeNoTicker, exchangeName, pair)
	}
}

func (d *Data) RemoveSpot(exchangeName, pair string, flag DataFlag) {
	if flag == DataFlag_Ticker {
		d.RemoveSpotTicker(exchangeName, pair)
	} else if flag == DataFlag_Depth {
		d.RemoveSpotDepth(exchangeName, pair)
	}
}

func (d *Data) RemoveFuture(exchangeName, contractType, pair string, flag DataFlag) {
	if flag == DataFlag_Ticker {
		d.RemoveFutureTicker(exchangeName, contractType, pair)
	} else if flag == DataFlag_Depth {
		d.RemoveFutureDepth(exchangeName, contractType, pair)
	}
}

func keyGen(exchangeName, pair string) string {
	return exchangeName + "/" + pair
}

func futureKeyGen(exchangeName, contactType, pair string) string {
	return exchangeName + "/" + contactType + "/" + pair
}
