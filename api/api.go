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
		Data:   SupportList,
	}
}

func (a *Api) GetDepth(exchange, pair string) *Response {
	a.logger.Debugf("GetDepth %s %s", exchange, pair)

	if !validate(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	depth, err := a.data.GetDepth(exchange, pair)
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

func (a *Api) GetTicker(exchange, pair string) *Response {
	a.logger.Debugf("GetTicker %s %s", exchange, pair)
	if !validate(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	ticker, err := a.data.GetTicker(exchange, pair)
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

func (a *Api) GetTrade(exchange, pair string) *Response {
	a.logger.Debugf("GetTicker %s %s", exchange, pair)
	if !validate(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	panic("not support yet")
}

func (a *Api) SubscribeDepth(exchange, pair string, period int64) *Response {
	a.logger.Infof("SubscribeDepth %s %s %d", exchange, pair, period)
	if !validate(exchange) {
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
			go worker.NewTickerWorker(exc.Ctx, a.data, exc.Api, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(exchange, pair, period, DataFlag_Depth)
		if c != nil {
			go worker.NewDepthWorker(a.ctx, a.data, c.Api, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

func (a *Api) SubscribeTicker(exchange, pair string, period int64) *Response {
	a.logger.Infof("SubscribeTicker %s %s %d", exchange, pair, period)
	if !validate(exchange) {
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
			go worker.NewTickerWorker(exc.Ctx, a.data, exc.Api, exc.Pair, exc.Period)
		}
	} else {
		c := a.cfg.AddConfig(exchange, pair, period, DataFlag_Ticker)
		if c != nil {
			c.Ctx, c.CancelFunc = context.WithCancel(a.ctx)
			go worker.NewTickerWorker(c.Ctx, a.data, c.Api, c.Pair, c.Period)
		}
	}
	return &Response{
		Status: 0,
	}
}

func (a *Api) SubscribeTrade(exchange, pair string, period int64) *Response {
	a.logger.Infof("SubscribeTrade %s %s %d", exchange, pair, period)
	if !validate(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	panic("not support yet")
}

func validate(exchangeName string) bool {
	for _, ex := range SupportList {
		if ex == exchangeName {
			return true
		}
	}
	return false
}
