package api

import (
	"context"
	"fmt"
	. "github.com/goex-top/market_center"
	"github.com/goex-top/market_center/config"
	"github.com/goex-top/market_center/data"
	"github.com/goex-top/market_center/worker"
)

type Api struct {
	ctx  context.Context
	cfg  *config.Config
	data *data.Data
}

func NewApi(ctx context.Context, cfg *config.Config, data *data.Data) *Api {
	return &Api{ctx: ctx, cfg: cfg, data: data}
}

func (a *Api) GetSupportList() *Response {
	return &Response{
		Status: 0,
		Data:   SupportList,
	}
}

func (a *Api) GetDepth(exchange, pair string) *Response {
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

func (a *Api) SubscribeDepth(exchange, pair string, period int64) *Response {
	if !validate(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, period)
	if exc != nil {
		exc.UpdatePeriod(period)
	} else {
		cfg := a.cfg.AddConfig(exchange, pair, period)
		go worker.NewDepthWorker(a.ctx, a.data, cfg.Pair[0].Api, cfg.Pair[0].Pair, cfg.Pair[0].Ticker)
	}
	return &Response{
		Status: 0,
	}
}

func (a *Api) SubscribeTicker(exchange, pair string, period int64) *Response {
	if !validate(exchange) {
		return &Response{
			Status:       -1,
			ErrorMessage: fmt.Sprintf(ErrMsg_ExchangeNotSupport, exchange),
		}
	}
	exc := a.cfg.FindConfig(exchange, pair, period)
	if exc != nil {
		exc.UpdatePeriod(period)
	} else {
		cfg := a.cfg.AddConfig(exchange, pair, period)
		go worker.NewTickerWorker(a.ctx, a.data, cfg.Pair[0].Api, cfg.Pair[0].Pair, cfg.Pair[0].Ticker)
	}
	return &Response{
		Status: 0,
	}
}

func validate(exchangeName string) bool {
	for _, ex := range SupportList {
		if ex == exchangeName {
			return true
		}
	}
	return false
}
