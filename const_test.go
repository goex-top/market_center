package market_center

import "testing"

func TestDataFlag_String(t *testing.T) {
	df := DataFlag_Depth | DataFlag_Ticker | DataFlag_Kline
	t.Log(df.String())
}
