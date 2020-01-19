package main

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/goex-top/market_center"
	"net"
)

func SendErrorMsg(c net.Conn, status int64, err string) {
	var rsp Response
	rsp.Status = status
	rsp.ErrorMessage = err
	SendRespMsg(c, &rsp)
}

func SendRespMsg(c net.Conn, rsp *Response) {
	r, _ := json.Marshal(rsp)
	c.Write([]byte(r))
}

func ProcessMessage(c net.Conn, msg []byte) error {
	var req Request
	err := json.Unmarshal(msg, &req)
	fmt.Println("ProcessMessage:", string(msg))
	if err != nil {
		SendErrorMsg(c, -1, fmt.Sprintf(ErrMsg_RequestFormatError, err.Error()))
		return err
	}

	rsp := &Response{}

	switch req.Type {
	case Type_GetSupportList:
		rsp = Api.GetSupportList()
	case Type_SubscribeDepth:
		rsp = Api.SubscribeDepth(req.ExchangeName, req.CurrencyPair, req.Period)
	case Type_SubscribeTicker:
		rsp = Api.SubscribeTicker(req.ExchangeName, req.CurrencyPair, req.Period)
	case Type_GetDepth:
		rsp = Api.GetDepth(req.ExchangeName, req.CurrencyPair)
	case Type_GetTicker:
		rsp = Api.GetTicker(req.ExchangeName, req.CurrencyPair)
	default:
		SendErrorMsg(c, -1, fmt.Sprintf(ErrMsg_RequestTypeNotSupport, req.Type))
		return errors.New(fmt.Sprintf(ErrMsg_RequestTypeNotSupport, req.Type))
	}

	SendRespMsg(c, rsp)
	return nil
}
