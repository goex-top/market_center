package config

import (
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
	Api    goex.API
	Pair   goex.CurrencyPair
	Ticker *time.Ticker
	Period time.Duration
	Flag   DataFlag
}

type ExchangeConfig struct {
	ExchangName string
	Pair        []PairConfig
}

type Config struct {
	ExchangesConfig []ExchangeConfig
}

func NewConfig() *Config {
	return &Config{ExchangesConfig: make([]ExchangeConfig, 0)}
}

func (c *Config) AddConfig(exchange, pair string, period int64, flag DataFlag) *ExchangeConfig {
	proxy := os.Getenv("HTTP_PROXY")
	if proxy != "" {
		log.Printf("add config with proxy:%s", proxy)
	}
	c.ExchangesConfig = append(c.ExchangesConfig, ExchangeConfig{
		ExchangName: exchange,
		Pair: []PairConfig{
			{
				Api:    builder.NewAPIBuilder().HttpProxy(proxy).Build(exchange),
				Pair:   goex.NewCurrencyPair2(pair),
				Period: time.Duration(period * int64(time.Millisecond)),
				Ticker: time.NewTicker(time.Duration(period * int64(time.Millisecond))),
				Flag:   flag,
			},
		},
	})
	return &c.ExchangesConfig[len(c.ExchangesConfig)-1]
}

func (c *Config) FindConfig(exchange, pair string, period int64, flag DataFlag) *PairConfig {
	for k, ex := range c.ExchangesConfig {
		if ex.ExchangName == exchange {
			for kk, p := range ex.Pair {
				if p.Pair.String() == pair && p.Flag == flag {
					return &c.ExchangesConfig[k].Pair[kk]
				}
			}
		}
	}
	return nil
}

func (c *PairConfig) UpdatePeriod(period int64) {
	if period*int64(time.Millisecond) < int64(c.Period) && period > 0 {
		log.Printf("update period from %dms to %dms", c.Period, period)

		c.Ticker.Stop()
		c.Ticker = time.NewTicker(time.Duration(period * int64(time.Millisecond)))
		c.Period = time.Duration(period * int64(time.Millisecond))
	}
}
