package worker

import (
	"context"
	goex "github.com/nntaoli-project/GoEx"
	"log"
	"market_center/data"
	"time"
)

const (
	defaultDepthSize = 10
)

func NewDepthWorker(ctx context.Context, depthData *data.Data, api goex.API, pair goex.CurrencyPair, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dep, err := api.GetDepth(defaultDepthSize, pair)
			if err != nil {
				log.Printf("[%s] refresh depth error:%s", api.GetExchangeName(), err.Error())
			}
			log.Println(dep)
			depthData.UpdateDepth(api.GetExchangeName(), pair.String(), dep)
		}
	}
}

func NewTickerWorker(ctx context.Context, tickerData *data.Data, api goex.API, pair goex.CurrencyPair, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick, err := api.GetTicker(pair)
			if err != nil {
				log.Printf("[%s] refresh ticker error:%s", api.GetExchangeName(), err.Error())
			}
			log.Println(tick)
			tickerData.UpdateTicker(api.GetExchangeName(), pair.String(), tick)
		}
	}
}
