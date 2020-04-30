package client

import (
	"fmt"
	"github.com/nntaoli-project/goex"
	"testing"
	"time"
)

var client = NewClient()

func TestClient_GetSupportList(t *testing.T) {
	t.Log(client.GetSupportList())
}

func TestClient_SubscribeSpotDepth(t *testing.T) {
	t.Log(client.SubscribeSpotDepth(goex.BINANCE, "BTC_USDT", 200))
}

func TestClient_GetSpotDepth(t *testing.T) {
	client.SubscribeSpotDepth(goex.BINANCE, "BTC_USDT", 200)
	time.Sleep(1 * time.Second)
	t.Log(client.GetSpotDepth(goex.BINANCE, "BTC_USDT"))
}

func TestClient_GetSpotTicker(t *testing.T) {
	client.SubscribeSpotTicker(goex.BINANCE, "BTC_USDT", 200)
	time.Sleep(time.Second)
	t.Log(client.GetSpotTicker(goex.BINANCE, "BTC_USDT"))
}

func TestNewClient(t *testing.T) {
	c1 := NewClient()
	c2 := NewClient()
	c1.SubscribeSpotDepth("binance.com", "BTC_USDT", 500)
	c2.SubscribeSpotTicker("binance.com", "BTC_USDT", 600)
	c2.SubscribeSpotTicker("binance.com", "BTC_USDT", 300)
	//c1.EnableDebug()
	//c2.EnableDebug()
	for {
		fmt.Println(c1.GetSpotDepth("binance.com", "BTC_USDT"))
		fmt.Println(c2.GetSpotDepth("binance.com", "BTC_USDT"))
		fmt.Println(c2.GetSpotTicker("binance.com", "BTC_USDT"))
		time.Sleep(time.Second)
	}
}

func TestClient_GetFutureTicker(t *testing.T) {
	c1 := NewClient()
	c1.SubscribeFutureTicker("bitmex.com", "", "BTC_USDT", 1500)
	for {
		fmt.Println(c1.GetFutureTicker("bitmex.com", "", "BTC_USDT"))
		time.Sleep(time.Second)
	}
}
