# A client for market center in FMZ
![HitCount](http://hits.dwyl.io/goex-top/market_center_client_fmz.svg)

Get market data from market center in [FMZ](https://fmz.com)

## Description
Ref to [FMZ](https://www.fmz.com/strategy/182185)

## API
* GetSupportList
* SubscribeSpotTicker
* SubscribeSpotDepth
* ~~SubscribeSpotTrade~~
* GetSpotTicker
* GetSpotDepth
* ~~GetSpotTrade~~
* SubscribeFutureTicker
* SubscribeFutureDepth
* ~~SubscribeFutureTrade~~
* GetFutureTicker
* GetFutureDepth
* ~~GetFutureTrade~~

## 行情中心客户端

此[客户端](https://github.com/goex-top/market_center_client_fmz)配合[行情中心](https://github.com/goex-top/market_center)使用

## 为什么有行情中心

对于行情的访问有REST与Websocket，而Websocket由于种种不可抗拒的因素，导致其连接被强行断开，这时策略需做好种种容错机制。REST往往还是最稳的方式。
当你有多个策略在相同平台上跑时，如OKEx，而这多策略又在同一台服务器上，对于REST API的访问是有频率限制的(详情请参考个平台API文档), 限制的方式有多种，如IP限制，账号限制，或两者组合等。
使用行情中心可降低REST请求频率

![传统方式](https://raw.githubusercontent.com/goex-top/market_center/master/trandition.png)

![行情中心](https://raw.githubusercontent.com/goex-top/market_center/master/market_center.png)

## 行情中心部署
1. [源码](https://github.com/goex-top/market_center)编译部署
2. [二进制部署](https://github.com/goex-top/market_center/releases)
3. 需要帮助请联系wx:btstarinfo, Q:6510575

## 请注意
* 行情中心只提供行情数据接口，不提供下单接口
* 行情中心只能运行在Linux, Unix上

### 观星者

![观星者](https://starchart.cc/goex-top/market_center_client_fmz.svg)