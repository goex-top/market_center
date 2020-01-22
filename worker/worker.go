package worker

import (
	"context"
	"github.com/goex-top/market_center/data"
	goex "github.com/nntaoli-project/GoEx"
	"log"
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

func NewSpotDepthWorker(ctx context.Context, depthData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new spot depth worker for [%s] %s, period is %dms", exchange, pair.String(), period/time.Millisecond)
	ticker := time.NewTicker(period)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dep, err := api.GetDepth(defaultDepthSize, pair)
			if err != nil {
				log.Printf("[%s] refresh %s depth error:%s", exchange, pair.String(), err.Error())
			} else {
				depthData.UpdateSpotDepth(exchange, pair.String(), dep)
			}
		}
	}
}

func NewSpotTickerWorker(ctx context.Context, tickerData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new spot ticker worker for [%s] %s, period is %dms ", exchange, pair.String(), period/time.Millisecond)
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick, err := api.GetTicker(pair)
			if err != nil {
				log.Printf("[%s] refresh %s ticker error:%s", exchange, pair.String(), err.Error())
			} else {
				tickerData.UpdateSpotTicker(exchange, pair.String(), tick)
			}
		}
	}
}

func NewFutureDepthWorker(ctx context.Context, depthData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new future depth worker for [%s] %s %s, period is %dms", exchange, contactType, pair.String(), period/time.Millisecond)
	ticker := time.NewTicker(period)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dep, err := api.GetFutureDepth(pair, contactType, defaultDepthSize)
			if err != nil {
				log.Printf("[%s] %s refresh %s depth error:%s", exchange, contactType, pair.String(), err.Error())
			} else {
				depthData.UpdateFutureDepth(exchange, contactType, pair.String(), dep)
			}
		}
	}
}

func NewFutureTickerWorker(ctx context.Context, tickerData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new future ticker worker for [%s] %s %s, period is %dms ", exchange, contactType, pair.String(), period/time.Millisecond)
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick, err := api.GetFutureTicker(pair, contactType)
			if err != nil {
				log.Printf("[%s] %s refresh %s ticker error:%s", exchange, contactType, pair.String(), err.Error())
			} else {
				tickerData.UpdateFutureTicker(exchange, contactType, pair.String(), tick)
			}
		}
	}
}
