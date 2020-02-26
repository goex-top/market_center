package worker

import (
	"context"
	"github.com/goex-top/market_center/data"
	"github.com/nntaoli-project/goex"
	log "github.com/sirupsen/logrus"
	"sort"
	"time"
)

var (
	defaultDepthSize = 20
	heartbeatTimer   = time.Now()
	logger           = log.New()
	heartbeatCount   = 0
)

func init() {
	logger.SetLevel(log.DebugLevel)
}

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
		logger.Infof("[%s] refresh %s depth error:%s", exchange, pair.String(), err.Error())
	} else {
		if dep.AskList[0].Price > dep.AskList[1].Price {
			sort.Sort(dep.AskList)
		}
		if dep.BidList[0].Price < dep.BidList[1].Price {
			sort.Sort(sort.Reverse(dep.BidList))
		}
		if dep.UTime.IsZero() {
			dep.UTime = time.Now()
		}
		depthData.UpdateSpotDepth(exchange, pair.String(), dep)
	}
}

func updateSportTicker(tickerData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair) {
	tick, err := api.GetTicker(pair)
	if err != nil {
		logger.Infof("[%s] refresh %s ticker error:%s", exchange, pair.String(), err.Error())
	} else {
		tickerData.UpdateSpotTicker(exchange, pair.String(), tick)
	}
}

func updateFutureDepth(depthData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair) {
	dep, err := api.GetFutureDepth(pair, contactType, defaultDepthSize)
	if err != nil {
		logger.Infof("[%s] %s refresh %s depth error:%s", exchange, contactType, pair.String(), err.Error())
	} else {
		if dep.AskList[0].Price > dep.AskList[1].Price {
			sort.Sort(dep.AskList)
		}
		if dep.BidList[0].Price < dep.BidList[1].Price {
			sort.Sort(sort.Reverse(dep.AskList))
		}
		if dep.UTime.IsZero() {
			dep.UTime = time.Now()
		}
		depthData.UpdateFutureDepth(exchange, contactType, pair.String(), dep)
	}
}

func updateFutureTicker(tickerData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair) {
	tick, err := api.GetFutureTicker(pair, contactType)
	if err != nil {
		logger.Infof("[%s] %s refresh %s ticker error:%s", exchange, contactType, pair.String(), err.Error())
	} else {
		tickerData.UpdateFutureTicker(exchange, contactType, pair.String(), tick)
	}
}

func heartbeat(args ...string) {
	if heartbeatTimer.Add(time.Second).Before(time.Now()) {
		heartbeatTimer = time.Now()
		heartbeatCount++
		logger.Debugln("[work] heartbeat:", heartbeatCount, args)
	}
}

func NewSpotDepthWorker(ctx context.Context, depthData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	logger.Infof("new spot depth worker for [%s] %s, period is %dms", exchange, pair.String(), period/time.Millisecond)

	updateSportDepth(depthData, api, exchange, pair)
	prepare()
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("sport depth worker for [%s] %s, period is %dms exit!", exchange, pair.String(), period/time.Millisecond)
			depthData.RemoveSpotDepth(exchange, pair.String())
			return
		case <-ticker.C:
			heartbeat(exchange, pair.String(), "sport depth")
			updateSportDepth(depthData, api, exchange, pair)
		}
	}
}

func NewSpotTickerWorker(ctx context.Context, tickerData *data.Data, api goex.API, exchange string, pair goex.CurrencyPair, period time.Duration) {
	logger.Infof("new spot ticker worker for [%s] %s, period is %dms ", exchange, pair.String(), period/time.Millisecond)

	updateSportTicker(tickerData, api, exchange, pair)
	prepare()
	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("sport ticker worker for [%s] %s, period is %dms exit!", exchange, pair.String(), period/time.Millisecond)
			tickerData.RemoveSpotTicker(exchange, pair.String())
			return
		case <-ticker.C:
			heartbeat(exchange, pair.String(), "sport ticker")
			updateSportTicker(tickerData, api, exchange, pair)
		}
	}
}

func NewFutureDepthWorker(ctx context.Context, depthData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair, period time.Duration) {
	logger.Infof("new future depth worker for [%s] %s %s, period is %dms", exchange, contactType, pair.String(), period/time.Millisecond)

	updateFutureDepth(depthData, api, exchange, contactType, pair)
	prepare()

	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("future depth worker for [%s] %s %s, period is %dms exit!", exchange, contactType, pair.String(), period/time.Millisecond)
			depthData.RemoveFutureDepth(exchange, contactType, pair.String())
			return
		case <-ticker.C:
			heartbeat(exchange, contactType, pair.String(), "future depth")
			updateFutureDepth(depthData, api, exchange, contactType, pair)
		}
	}
}

func NewFutureTickerWorker(ctx context.Context, tickerData *data.Data, api goex.FutureRestAPI, exchange, contactType string, pair goex.CurrencyPair, period time.Duration) {
	logger.Infof("new future ticker worker for [%s] %s %s, period is %dms ", exchange, contactType, pair.String(), period/time.Millisecond)

	updateFutureTicker(tickerData, api, exchange, contactType, pair)
	prepare()

	ticker := time.NewTicker(period)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			logger.Infof("future ticker worker for [%s] %s %s, period is %dms exit!", exchange, contactType, pair.String(), period/time.Millisecond)
			tickerData.RemoveFutureTicker(exchange, contactType, pair.String())
			return
		case <-ticker.C:
			heartbeat(exchange, contactType, pair.String(), "future ticker")
			updateFutureTicker(tickerData, api, exchange, contactType, pair)
		}
	}
}
