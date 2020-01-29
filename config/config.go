package config

import (
	"context"
	. "github.com/goex-top/market_center"
	"log"
	"os"
	"time"

	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
)

//type Worker struct {
//	API goex.API
//	Pair goex.CurrencyPair
//	Duration time.Duration
//}

type PairConfig struct {
	SpotApi    goex.API
	FutureApi  goex.FutureRestAPI
	Pair       goex.CurrencyPair
	Period     time.Duration
	Flag       DataFlag
	Ctx        context.Context
	CancelFunc func()
	//Ticker *time.Ticker
}

type ExchangeConfig struct {
	ExchangeName string
	Pair         []PairConfig
}

type Config struct {
	ExchangesConfig []ExchangeConfig
}

func NewConfig() *Config {
	return &Config{ExchangesConfig: make([]ExchangeConfig, 0)}
}

func (c *Config) AddConfig(parentCtx context.Context, exchange, pair string, period int64, flag DataFlag) *PairConfig {
	proxy := os.Getenv("HTTP_PROXY")
	if proxy != "" {
		log.Printf("add config with proxy:%s", proxy)
	}
	ctx, cancelFunc := context.WithCancel(parentCtx)

	for k, ex := range c.ExchangesConfig {
		if ex.ExchangeName == exchange {
			for _, p := range ex.Pair {
				if p.Pair.String() == pair && p.Flag == flag {
					return nil
				}
			}
			if !IsFutureExchange(exchange) {
				c.ExchangesConfig[k].Pair = append(c.ExchangesConfig[k].Pair, PairConfig{
					SpotApi:    builder.NewAPIBuilder().HttpProxy(proxy).Build(exchange),
					FutureApi:  nil,
					Pair:       goex.NewCurrencyPair2(pair),
					Period:     time.Duration(period * int64(time.Millisecond)),
					Flag:       flag,
					CancelFunc: cancelFunc,
					Ctx:        ctx,
					//Ticker: time.NewTicker(time.Duration(period * int64(time.Millisecond))),
				})
			} else {
				c.ExchangesConfig[k].Pair = append(c.ExchangesConfig[k].Pair, PairConfig{
					FutureApi:  builder.NewAPIBuilder().HttpProxy(proxy).BuildFuture(exchange),
					SpotApi:    nil,
					Pair:       goex.NewCurrencyPair2(pair),
					Period:     time.Duration(period * int64(time.Millisecond)),
					Flag:       flag,
					CancelFunc: cancelFunc,
					Ctx:        ctx,
					//Ticker: time.NewTicker(time.Duration(period * int64(time.Millisecond))),
				})
			}
			return &c.ExchangesConfig[k].Pair[len(c.ExchangesConfig[k].Pair)-1]
		}
	}
	if !IsFutureExchange(exchange) {
		c.ExchangesConfig = append(c.ExchangesConfig, ExchangeConfig{
			ExchangeName: exchange,
			Pair: []PairConfig{
				{
					SpotApi:    builder.NewAPIBuilder().HttpProxy(proxy).Build(exchange),
					FutureApi:  nil,
					Pair:       goex.NewCurrencyPair2(pair),
					Period:     time.Duration(period * int64(time.Millisecond)),
					Flag:       flag,
					CancelFunc: cancelFunc,
					Ctx:        ctx,
					//Ticker: time.NewTicker(time.Duration(period * int64(time.Millisecond))),
				},
			},
		})
	} else {
		c.ExchangesConfig = append(c.ExchangesConfig, ExchangeConfig{
			ExchangeName: exchange,
			Pair: []PairConfig{
				{
					FutureApi:  builder.NewAPIBuilder().HttpProxy(proxy).BuildFuture(exchange),
					SpotApi:    nil,
					Pair:       goex.NewCurrencyPair2(pair),
					Period:     time.Duration(period * int64(time.Millisecond)),
					Flag:       flag,
					CancelFunc: cancelFunc,
					Ctx:        ctx,
					//Ticker: time.NewTicker(time.Duration(period * int64(time.Millisecond))),
				},
			},
		})
	}
	return &c.ExchangesConfig[len(c.ExchangesConfig)-1].Pair[0]
}

func (c *Config) FindConfig(exchange, pair string, period int64, flag DataFlag) *PairConfig {
	for k, ex := range c.ExchangesConfig {
		if ex.ExchangeName == exchange {
			for kk, p := range ex.Pair {
				if p.Pair.String() == pair && p.Flag == flag {
					return &c.ExchangesConfig[k].Pair[kk]
				}
			}
		}
	}
	return nil
}

func (c *PairConfig) UpdatePeriod(period int64) bool {
	if period*int64(time.Millisecond) < int64(c.Period) && period > 0 {
		log.Printf("update period from %dms to %dms", c.Period/time.Millisecond, period)

		//c.Ticker.Stop()
		//c.Ticker = time.NewTicker(time.Duration(period * int64(time.Millisecond)))
		c.Period = time.Duration(period * int64(time.Millisecond))
		return true
	}
	return false
}
