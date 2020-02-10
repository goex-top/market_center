package api

import (
	. "github.com/goex-top/market_center"
	"strings"
	"time"
)

const defaultActiveTimerDuration = 10

func keyGen(exchangeName, pair string, flag DataFlag) string {
	return exchangeName + "/" + pair + "/" + flag.String()
}

func futureKeyGen(exchangeName, contactType, pair string, flag DataFlag) string {
	return exchangeName + "/" + contactType + "/" + pair + "/" + flag.String()
}

func (a *Api) addTimer(name string) {
	a.setTimer(name, defaultActiveTimerDuration)
}

func (a *Api) setTimer(name string, num int) {
	a.activeTimer.Store(name, num)
}

func (a *Api) resetTimer(name string) {
	_, ok := a.activeTimer.Load(name)
	if ok {
		a.setTimer(name, defaultActiveTimerDuration)
	}
}

func (a *Api) removeTimer(name string) {
	a.activeTimer.Delete(name)
}

func (a *Api) monitor() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	count := 0
	for {
		select {
		case <-a.ctx.Done():
			return
		case <-ticker.C:
			hasSubscribed := false
			a.activeTimer.Range(func(key, value interface{}) bool {
				hasSubscribed = true
				_timer := value.(int)
				if _timer > 0 {
					_timer--
					a.setTimer(key.(string), _timer)
				} else {
					a.removeTimer(key.(string))
					keys := strings.Split(key.(string), "/")

					flag := ParseDataFlag(keys[len(keys)-1])
					exc := a.cfg.FindConfig(keys[0], keys[1], flag)
					if exc != nil {
						a.logger.Warnf("[%s] %s %s timeout, cancel it", keys[0], keys[1], flag.String())
						exc.Cancel()
						time.Sleep(200 * time.Millisecond)
					}
					a.cfg.RemoveConfig(keys[0], keys[1], flag)
					if IsFutureExchange(keys[0]) {
						a.data.RemoveFuture(keys[0], keys[1], keys[2], flag)
					} else {
						a.data.RemoveSpot(keys[0], keys[1], flag)
					}
				}
				return true
			})

			if hasSubscribed {
				count++
				a.logger.Debugln("[api monitor] heartbeat:", count)
			} else {
				count = 0
			}
		}
	}
}
