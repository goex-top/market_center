package config

import (
	"context"
	"testing"
)

var c = NewConfig()

func TestConfig_AddConfig(t *testing.T) {
	c1 := c.AddConfig(context.Background(), "binance.com", "BTC_USDT", 100, 1)
	c1.cancelFunc()
}

func TestConfig_FindConfig(t *testing.T) {
	c.AddConfig(context.Background(), "binance.com", "BTC_USDT", 100, 1)
	c2 := c.FindConfig("binance.com", "BTC_USDT", 1)
	c2.cancelFunc()
}

func TestPairConfig_UpdatePeriod(t *testing.T) {
	c1 := c.AddConfig(context.Background(), "binance.com", "BTC_USDT", 100, 1)
	c2 := c.FindConfig("binance.com", "BTC_USDT", 1)
	t.Log(c1)
	c2.UpdatePeriod(50)
	t.Log(c2)
}
