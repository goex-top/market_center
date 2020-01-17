package market_center

import (
	"encoding/json"
	"testing"
)

func Test_ResponseUnmarshal(t *testing.T) {
	rspStr := `{"status":0, "error_message":"this is a test"}`
	rsp := &Response{}
	err := json.Unmarshal([]byte(rspStr), rsp)
	t.Log(err, rsp)
}

func Test_ResponseUnmarshal2(t *testing.T) {
	rspStr := `{"status":0, "data":"this is a test"}`
	rsp := &Response{}
	err := json.Unmarshal([]byte(rspStr), rsp)
	t.Log(err, rsp)
}
