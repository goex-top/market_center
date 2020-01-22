package api

import (
	"context"
	"fmt"
	. "github.com/goex-top/market_center"
	"github.com/goex-top/market_center/config"
	"github.com/goex-top/market_center/data"
	"github.com/goex-top/market_center/worker"
	log "github.com/sirupsen/logrus"
)

type Api struct {
	ctx    context.Context
	cfg    *config.Config
	data   *data.Data
	logger *log.Logger
}

func NewApi(ctx context.Context, cfg *config.Config, data *data.Data) *Api {
	return &Api{ctx: ctx, cfg: cfg, data: data, logger: log.New()}
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

//func (a *Api) GetTrade(exchange, pair string) *Response {
//	a.logger.Debugf("GetTicker %s %s", exchange, pair)
//	if !validate(exchange) {
//		return &Response{
//			Status:       -1,
//			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
//		}
//	}
//	panic("not support yet")
//}

// spot

func (a *Api) GetSpotDepth(exchange, pair string) *Response {
	a.logger.Debugf("GetDepth %s %s", exchange, pair)

	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	depth, err := a.data.GetSpotDepth(exchange, pair)
	if err != nil {
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	return &Response{
		Status: 0,
		Data:   depth,
	}
}

func (a *Api) GetSpotTicker(exchange, pair string) *Response {
	a.logger.Debugf("GetTicker %s %s", exchange, pair)
	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	ticker, err := a.data.GetSpotTicker(exchange, pair)
	if err != nil {
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	return &Response{
		Status: 0,
		Data:   ticker,
	}
}

func (a *Api) SubscribeSpotDepth(exchange, pair string, period int64) *Response {
	a.logger.Infof("SubscribeSpotDepth %s %s %d", exchange, pair, period)
	if !validateSpot(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, period, DataFlag_Depth)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.CancelFunc()
			exc.Ctx, exc.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewSpotDepthWorker(exc.Ctx, a.data, exc.SpotApi, exchange, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(exchange, pair, period, DataFlag_Depth)
		if c != nil {
			go worker.NewSpotDepthWorker(a.ctx, a.data, c.SpotApi, exchange, c.Pair, c.Period)
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
	exc := a.cfg.FindConfig(exchange, pair, period, DataFlag_Ticker)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.CancelFunc()
			exc.Ctx, exc.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewSpotTickerWorker(exc.Ctx, a.data, exc.SpotApi, exchange, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(exchange, pair, period, DataFlag_Ticker)
		if c != nil {
			c.Ctx, c.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewSpotTickerWorker(c.Ctx, a.data, c.SpotApi, exchange, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

// future

func (a *Api) GetFutureDepth(exchange, contractType, pair string) *Response {
	a.logger.Debugf("GetDepth %s %s", exchange, pair)

	if !validateFuture(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	depth, err := a.data.GetFutureTicker(exchange, contractType, pair)
	if err != nil {
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
	return &Response{
		Status: 0,
		Data:   depth,
	}
}

func (a *Api) GetFutureTicker(exchange, contractType, pair string) *Response {
	a.logger.Debugf("GetTicker %s %s", exchange, pair)
	if !validateFuture(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	ticker, err := a.data.GetFutureTicker(exchange, contractType, pair)
	if err != nil {
		return &Response{
			Status:       -1,
			ErrorMessage: err.Error(),
		}
	}
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
	exc := a.cfg.FindConfig(exchange, pair, period, DataFlag_Depth)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.CancelFunc()
			exc.Ctx, exc.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewFutureDepthWorker(exc.Ctx, a.data, exc.FutureApi, exchange, contractType, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(exchange, pair, period, DataFlag_Depth)
		if c != nil {
			go worker.NewFutureDepthWorker(a.ctx, a.data, c.FutureApi, exchange, contractType, c.Pair, c.Period)
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
	exc := a.cfg.FindConfig(exchange, pair, period, DataFlag_Ticker)
	if exc != nil {
		if exc.UpdatePeriod(period) {
			exc.CancelFunc()
			exc.Ctx, exc.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewFutureTickerWorker(exc.Ctx, a.data, exc.FutureApi, exchange, contractType, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(exchange, pair, period, DataFlag_Ticker)
		if c != nil {
			c.Ctx, c.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewFutureTickerWorker(c.Ctx, a.data, c.FutureApi, exchange, contractType, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

//
//func (a *Api) SubscribeTrade(exchange, pair string, period int64) *Response {
//	a.logger.Infof("SubscribeTrade %s %s %d", exchange, pair, period)
//	if !validate(exchange) {
//		return &Response{
//			Status:       -1,
//			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
//		}
//	}
//	panic("not support yet")
//}

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
