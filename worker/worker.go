package worker

import (
	"context"
	"github.com/goex-top/market_center/data"
	goex "github.com/nntaoli-project/GoEx"
	"log"
	"sort"
	"time"
)

var (
	defaultDepthSize = 20
)

func SetDefaultDepthSize(size int) {
	defaultDepthSize = size
}

func GetDefaultDepthSize() int {
	return defaultDepthSize
}

func prepare() {
	now := time.Now()
	diff := int(time.Second) - now.Nanosecond()
	time.Sleep(time.Duration(diff))
}

func updateSportDepth(depthData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair) {
	dep, err := api.GetDepth(defaultDepthSize, pair)
	if err != nil {
		log.Printf("[%s] refresh %s depth error:%s", exchange, pair.String(), err.Error())
	} else {
		if dep.AskList[0].Price > dep.AskList[1].Price {
			sort.Sort(dep.AskList)
		}
		if dep.BidList[0].Price < dep.BidList[1].Price {
			sort.Sort(sort.Reverse(dep.AskList))
		}
		depthData.UpdateSpotDepth(exchange, pair.String(), dep)
	}
}

func updateSportTicker(tickerData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair) {
	tick, err := api.GetTicker(pair)
	if err != nil {
		log.Printf("[%s] refresh %s ticker error:%s", exchange, pair.String(), err.Error())
	} else {
		tickerData.UpdateSpotTicker(exchange, pair.String(), tick)
	}
}

func updateFutureDepth(depthData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair) {
	dep, err := api.GetFutureDepth(pair, contactType, defaultDepthSize)
	if err != nil {
		log.Printf("[%s] %s refresh %s depth error:%s", exchange, contactType, pair.String(), err.Error())
	} else {
		if dep.AskList[0].Price > dep.AskList[1].Price {
			sort.Sort(dep.AskList)
		}
		if dep.BidList[0].Price < dep.BidList[1].Price {
			sort.Sort(sort.Reverse(dep.AskList))
		}
		depthData.UpdateFutureDepth(exchange, contactType, pair.String(), dep)
	}
}

func updateFutureTicker(tickerData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair) {
	tick, err := api.GetFutureTicker(pair, contactType)
	if err != nil {
		log.Printf("[%s] %s refresh %s ticker error:%s", exchange, contactType, pair.String(), err.Error())
	} else {
		tickerData.UpdateFutureTicker(exchange, contactType, pair.String(), tick)
	}
}

func NewSpotDepthWorker(ctx context.Context, depthData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new spot depth worker for [%s] %s, period is %dms", exchange, pair.String(), period/time.Millisecond)

	updateSportDepth(depthData, api, exchange, pair)
	prepare()
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateSportDepth(depthData, api, exchange, pair)
		}
	}
}

func NewSpotTickerWorker(ctx context.Context, tickerData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new spot ticker worker for [%s] %s, period is %dms ", exchange, pair.String(), period/time.Millisecond)

	updateSportTicker(tickerData, api, exchange, pair)
	prepare()
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateSportTicker(tickerData, api, exchange, pair)
		}
	}
}

func NewFutureDepthWorker(ctx context.Context, depthData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new future depth worker for [%s] %s %s, period is %dms", exchange, contactType, pair.String(), period/time.Millisecond)

	updateFutureDepth(depthData, api, exchange, contactType, pair)
	prepare()

	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateFutureDepth(depthData, api, exchange, contactType, pair)
		}
	}
}

func NewFutureTickerWorker(ctx context.Context, tickerData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new future ticker worker for [%s] %s %s, period is %dms ", exchange, contactType, pair.String(), period/time.Millisecond)

	updateFutureTicker(tickerData, api, exchange, contactType, pair)
	prepare()

	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			updateFutureTicker(tickerData, api, exchange, contactType, pair)
		}
	}
}
