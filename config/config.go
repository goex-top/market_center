package config

import (
	goex "github.com/nntaoli-project/GoEx"
	"github.com/nntaoli-project/GoEx/builder"
	"os"
	"time"
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

func (c *Config) AddConfig(exchange, pair string, period int64) *ExchangeConfig {
	proxy := os.Getenv("HTTP_PROXY")
	c.ExchangesConfig = append(c.ExchangesConfig, ExchangeConfig{
		ExchangName: exchange,
		Pair: []PairConfig{
			{
				Api:    builder.NewAPIBuilder().HttpProxy(proxy).Build(exchange),
				Pair:   goex.NewCurrencyPair2(pair),
				Period: time.Duration(period * int64(time.Millisecond)),
				Ticker: time.NewTicker(time.Duration(period * int64(time.Millisecond))),
			},
		},
	})
	return &c.ExchangesConfig[len(c.ExchangesConfig)-1]
}

func (c *Config) FindConfig(exchange, pair string, period int64) *PairConfig {
	for k, ex := range c.ExchangesConfig {
		if ex.ExchangName == exchange {
			for kk, p := range ex.Pair {
				if p.Pair.String() == pair {
					return &c.ExchangesConfig[k].Pair[kk]
				}
			}
		}
	}
	return nil
}

func (c *PairConfig) UpdatePeriod(period int64) {
	if period*int64(time.Millisecond) < int64(c.Period) && period > 0 {
		c.Ticker.Stop()
		c.Ticker = time.NewTicker(time.Duration(period * int64(time.Millisecond)))
		c.Period = time.Duration(period * int64(time.Millisecond))
	}
}
