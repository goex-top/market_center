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

func NewDepthWorker(ctx context.Context, depthData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new depth worker for [%s] %s, period is %dms", exchange, pair.String(), period/time.Millisecond)
	ticker := time.NewTicker(period)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dep, err := api.GetDepth(defaultDepthSize, pair)
			if err != nil {
				log.Printf("[%s] refresh depth error:%s", exchange, err.Error())
			}
			//log.Println("DEPTH:", dep)
			depthData.UpdateDepth(exchange, pair.String(), dep)
		}
	}
}

func NewTickerWorker(ctx context.Context, tickerData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	log.Printf("new ticker worker for [%s] %s, period is %dms ", exchange, pair.String(), period/time.Millisecond)
	ticker := time.NewTicker(period)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick, err := api.GetTicker(pair)
			if err != nil {
				log.Printf("[%s] refresh ticker error:%s", exchange, err.Error())
			}
			//log.Println("TICKER:", tick)
			tickerData.UpdateTicker(exchange, pair.String(), tick)
		}
	}
}
