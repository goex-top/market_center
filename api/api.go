package api

import (
	"context"
	"fmt"
	. "github.com/goex-top/market_center"
	"github.com/goex-top/market_center/config"
	"github.com/goex-top/market_center/data"
	"github.com/goex-top/market_center/worker"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Api struct {
	ctx         context.Context
	cfg         *config.Config
	data        *data.Data
	logger      *log.Logger
	activeTimer sync.Map
}

func NewApi(ctx context.Context, cfg *config.Config, data *data.Data) *Api {
	_api := &Api{ctx: ctx, cfg: cfg, data: data, logger: log.New(), activeTimer: sync.Map{}}
	go _api.monitor()
	return _api
}

func (a *Api) EnableDebug() {
	a.logger.SetLevel(log.DebugLevel)
}

func (a *Api) GetSupportList() *Response {
	a.logger.Debugln("GetSupportList")
	return &Response{
		Status: 0,
		Data:   append(SpotList, FutureList...),
	}
}

// spot

func (a *Api) GetSpotDepth(exchange, pair string) *Response {

	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	a.resetTimer(keyGen(exchange, pair, DataFlag_Depth))
	depth, err := a.data.GetSpotDepth(exchange, pair)
	if err != nil {
		a.logger.Debugf("GetSpotDepth %s %s err:%v", exchange, pair, err)
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	a.logger.Debugf("GetSpotDepth %s %s success", exchange, pair)
	return &Response{
		Status: 0,
		Data:   depth,
	}
}

func (a *Api) GetSpotTicker(exchange, pair string) *Response {
	a.logger.Debugf("GetSpotTicker %s %s", exchange, pair)
	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	a.resetTimer(keyGen(exchange, pair, DataFlag_Ticker))
	ticker, err := a.data.GetSpotTicker(exchange, pair)
	if err != nil {
		a.logger.Debugf("GetSpotTicker %s %s err:%v", exchange, pair, err)
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	a.logger.Debugf("GetSpotTicker %s %s success", exchange, pair)
	return &Response{
		Status: 0,
		Data:   ticker,
	}
}

func (a *Api) SubscribeSpotDepth(exchange, pair string, period int64) *Response {
	a.logger.Infof("SubscribeSpotDepth %s %s %dms", exchange, pair, period)
	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, DataFlag_Depth)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.Cancel()
			ctx := exc.NewSubContext(a.ctx)
			go worker.NewSpotDepthWorker(ctx, a.data, exc.SpotApi, exchange, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(a.ctx, exchange, pair, period, DataFlag_Depth)
		if c != nil {
			a.addTimer(keyGen(exchange, pair, DataFlag_Depth))
			go worker.NewSpotDepthWorker(c.Context(), a.data, c.SpotApi, exchange, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

func (a *Api) SubscribeSpotTicker(exchange, pair string, period int64) *Response {
	a.logger.Infof("SubscribeSpotTicker %s %s %d", exchange, pair, period)
	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, DataFlag_Ticker)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.Cancel()
			ctx := exc.NewSubContext(a.ctx)
			go worker.NewSpotTickerWorker(ctx, a.data, exc.SpotApi, exchange, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(a.ctx, exchange, pair, period, DataFlag_Ticker)
		if c != nil {
			a.addTimer(keyGen(exchange, pair, DataFlag_Ticker))
			go worker.NewSpotTickerWorker(c.Context(), a.data, c.SpotApi, exchange, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

// future

func (a *Api) GetFutureDepth(exchange, contractType, pair string) *Response {

	if !validateFuture(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	a.resetTimer(futureKeyGen(exchange, contractType, pair, DataFlag_Depth))
	depth, err := a.data.GetFutureTicker(exchange, contractType, pair)
	if err != nil {
		a.logger.Debugf("GetFutureDepth %s %s %s err:%v", exchange, contractType, pair, err)
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	a.logger.Debugf("GetFutureDepth %s %s success", exchange, pair)
	return &Response{
		Status: 0,
		Data:   depth,
	}
}

func (a *Api) GetFutureTicker(exchange, contractType, pair string) *Response {
	a.logger.Debugf("GetFutureTicker %s %s", exchange, pair)
	if !validateFuture(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	a.resetTimer(futureKeyGen(exchange, contractType, pair, DataFlag_Ticker))
	ticker, err := a.data.GetFutureTicker(exchange, contractType, pair)
	if err != nil {
		a.logger.Debugf("GetFutureTicker %s %s %s err:%v", exchange, contractType, pair, err)
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	a.logger.Debugf("GetFutureTicker %s %s success", exchange, pair)
	return &Response{
		Status: 0,
		Data:   ticker,
	}
}

func (a *Api) SubscribeFutureDepth(exchange, contractType, pair string, period int64) *Response {
	a.logger.Infof("SubscribeFutureDepth %s %s %s %d", exchange, contractType, pair, period)
	if !validateFuture(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, DataFlag_Depth)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.Cancel()
			ctx := exc.NewSubContext(a.ctx)
			go worker.NewFutureDepthWorker(ctx, a.data, exc.FutureApi, exchange, contractType, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(a.ctx, exchange, pair, period, DataFlag_Depth)
		if c != nil {
			a.addTimer(futureKeyGen(exchange, contractType, pair, DataFlag_Depth))
			go worker.NewFutureDepthWorker(c.Context(), a.data, c.FutureApi, exchange, contractType, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

func (a *Api) SubscribeFutureTicker(exchange, contractType, pair string, period int64) *Response {
	a.logger.Infof("SubscribeFutureTicker %s %s %s %d", exchange, contractType, pair, period)
	if !validateFuture(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, DataFlag_Ticker)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.Cancel()
			ctx := exc.NewSubContext(a.ctx)
			go worker.NewFutureTickerWorker(ctx, a.data, exc.FutureApi, exchange, contractType, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(a.ctx, exchange, pair, period, DataFlag_Ticker)
		if c != nil {
			a.addTimer(futureKeyGen(exchange, contractType, pair, DataFlag_Ticker))
			go worker.NewFutureTickerWorker(c.Context(), a.data, c.FutureApi, exchange, contractType, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

func validateSpot(exchangeName string) bool {
	for _, ex := range SpotList {
		if ex == exchangeName {
			return true
		}
	}
	return false
}
func validateFuture(exchangeName string) bool {
	for _, ex := range FutureList {
		if ex == exchangeName {
			return true
		}
	}
	return false
}
