package config

import (
	"context"
	. "github.com/goex-top/market_center"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
	"time"

	"github.com/nntaoli-project/goex"
	"github.com/nntaoli-project/goex/builder"
)

type PairConfig struct {
	SpotApi    goex.API
	FutureApi  goex.FutureRestAPI
	Pair       goex.CurrencyPair
	Period     time.Duration
	Flag       DataFlag
	ctx        context.Context
	cancelFunc func()
	parent     *Config
}

type ExchangeConfig struct {
	ExchangeName string
	Pair         []PairConfig
}

type Config struct {
	ExchangesConfig []ExchangeConfig
	lock            sync.RWMutex
	logger          *log.Logger
}

func NewConfig() *Config {
	return &Config{ExchangesConfig: make([]ExchangeConfig, 0), logger: log.New()}
}

func (c *Config) AddConfig(parentCtx context.Context, exchange, pair string, period int64, flag DataFlag) *PairConfig {
	proxy := os.Getenv("HTTP_PROXY")
	if proxy != "" {
		c.logger.Printf("add config with proxy:%s", proxy)
	}
	ctx, cancelFunc := context.WithCancel(parentCtx)
	defer c.lock.Unlock()
	c.lock.Lock()
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
					cancelFunc: cancelFunc,
					ctx:        ctx,
					parent:     c,
				})
			} else {
				c.ExchangesConfig[k].Pair = append(c.ExchangesConfig[k].Pair, PairConfig{
					FutureApi:  builder.NewAPIBuilder().HttpProxy(proxy).BuildFuture(exchange),
					SpotApi:    nil,
					Pair:       goex.NewCurrencyPair2(pair),
					Period:     time.Duration(period * int64(time.Millisecond)),
					Flag:       flag,
					cancelFunc: cancelFunc,
					ctx:        ctx,
					parent:     c,
				})
			}
			return &c.ExchangesConfig[k].Pair[len(c.ExchangesConfig[k].Pair)-1]
		}
	}
	var pc PairConfig
	if !IsFutureExchange(exchange) {
		pc = PairConfig{
			SpotApi:    builder.NewAPIBuilder().HttpProxy(proxy).Build(exchange),
			FutureApi:  nil,
			Pair:       goex.NewCurrencyPair2(pair),
			Period:     time.Duration(period * int64(time.Millisecond)),
			Flag:       flag,
			cancelFunc: cancelFunc,
			ctx:        ctx,
			parent:     c,
		}
	} else {
		pc = PairConfig{
			FutureApi:  builder.NewAPIBuilder().HttpProxy(proxy).BuildFuture(exchange),
			SpotApi:    nil,
			Pair:       goex.NewCurrencyPair2(pair),
			Period:     time.Duration(period * int64(time.Millisecond)),
			Flag:       flag,
			cancelFunc: cancelFunc,
			ctx:        ctx,
			parent:     c,
		}
	}

	c.ExchangesConfig = append(c.ExchangesConfig, ExchangeConfig{
		ExchangeName: exchange,
		Pair:         []PairConfig{pc},
	})
	return &c.ExchangesConfig[len(c.ExchangesConfig)-1].Pair[0]
}

func (c *Config) FindConfig(exchange, pair string, flag DataFlag) *PairConfig {
	defer c.lock.RUnlock()
	c.lock.RLock()

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

func (c *Config) RemoveConfig(exchange, pair string, flag DataFlag) {
	defer c.lock.Unlock()
	c.lock.Lock()
	for k, ex := range c.ExchangesConfig {
		if ex.ExchangeName == exchange {
			for kk, p := range ex.Pair {
				if p.Pair.String() == pair && p.Flag == flag {
					c.ExchangesConfig[k].Pair = append(c.ExchangesConfig[k].Pair[:kk], c.ExchangesConfig[k].Pair[kk+1:]...)
					if len(c.ExchangesConfig[k].Pair) == 0 {
						c.ExchangesConfig = append(c.ExchangesConfig[:k], c.ExchangesConfig[k+1:]...)
					}
				}
			}
		}
	}
}

func (c *PairConfig) UpdatePeriod(period int64) bool {
	defer c.parent.lock.Unlock()
	c.parent.lock.Lock()
	if period*int64(time.Millisecond) < int64(c.Period) && period > 0 {
		c.parent.logger.Printf("update period from %dms to %dms", c.Period/time.Millisecond, period)
		c.Period = time.Duration(period * int64(time.Millisecond))
		return true
	}
	return false
}

func (c *PairConfig) Cancel() {
	defer c.parent.lock.Unlock()
	c.parent.lock.Lock()
	c.cancelFunc()
}

func (c *PairConfig) Context() context.Context {
	defer c.parent.lock.Unlock()
	c.parent.lock.Lock()
	return c.ctx
}

func (c *PairConfig) NewSubContext(parent context.Context) context.Context {
	defer c.parent.lock.Unlock()
	c.parent.lock.Lock()
	c.ctx, c.cancelFunc = context.WithCancel(parent)
	return c.ctx
}
