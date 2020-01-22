package market_center

import (
	goex "github.com/nntaoli-project/GoEx"
)

type Response struct {
	Status       int64       `json:"status"`
	ErrorMessage string      `json:"error_message, omitempty"`
	Data         interface{} `json:"data",omitempty"`
}

type Request struct {
	Type         Type   `json:"type"`
	ExchangeName string `json:"exchange_name"`
	CurrencyPair string `json:"currency_pair"`
	ContractType string `json:"contract_type",omitempty`
	Period       int64  `json:"period",omitempty"` //unit: ms
}

//type Subscribe struct {
//	ExchangeName string `json:"exchange_name"`
//	CurrencyPair string `json:"currency_pair"`
//	Period       int64  `json:"period"` //unit: ms
//}

type ExchangeList struct {
	List []string
}

type Depth struct {
	ExchangeName string     `json:"exchange_name"`
	Depth        goex.Depth `json:"depth"`
}

type Ticker struct {
	ExchangeName string      `json:"exchange_name"`
	Ticker       goex.Ticker `json:"ticker"`
}
