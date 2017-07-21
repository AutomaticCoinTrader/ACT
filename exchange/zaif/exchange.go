package zaif

import (
	"github.com/pkg/errors"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"strings"
	"fmt"
	"sync"
	"strconv"
	"log"
)

const(
	exchangeName = "zaif"
)

func updateFunds(exchangeName string, requester *Requester, funds *ExchageFunds) (error){
	info2Response, _, _ , err := requester.GetInfo2()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not get info2 (ID = %v)", exchangeName))
	}
	if info2Response.Success != 1 {
		return errors.Errorf("can not buy (ID = %v, reason = %v)", exchangeName, info2Response.Error)
	}
	funds.update(map[string]float64{
		"btc": info2Response.Return.Funds.Btc,
		"mona": info2Response.Return.Funds.Mona,
		"xem": info2Response.Return.Funds.Xem,
		"jpy": info2Response.Return.Funds.Jpy})
	return nil
}

type BoardCursor struct {
	index int
	values [][]float64
}

func (b *BoardCursor) Next() (float64, float64, bool) {
	if b.index >= len(b.values) {
		return 0, 0, false
	}
	value := b.values[b.index]
	b.index++
	return value[0], value[1], true
}

func (b *BoardCursor) Reset() {
	b.index = 0
}

func (b *BoardCursor) Len() int {
	return len(b.values)
}

type TradeHistoryCursor struct {
	index int
	values []*StreamingTradesResponse
}

func (t *TradeHistoryCursor) Next() (time int64, peice float64, amount float64, tradeType string, ok bool) {
	if t.index >= len(t.values) {
		return 0, 0, 0, "", false
	}
	value := t.values[t.index]
	t.index++
	return value.Date, value.Price, value.Amount, value.TradeType, true
}

func (t *TradeHistoryCursor) Reset() {
	t.index = 0
}

func (t *TradeHistoryCursor) Len() int {
	return len(t.values)
}

type OrderCursor struct {
	index int
	keys []string
	values map[string]TradeActiveOrderRecordResponse
}

func (o *OrderCursor) Next() (int64, exchange.OrderAction, float64, float64, bool) {
	if o.index >= len(o.keys) {
		return 0, "", 0, 0, false
	}
	value := o.values[o.keys[o.index]]
	o.index++
	var action exchange.OrderAction
	if value.Action == "ask" {
		action = exchange.OrderActSell
	} else if value.Action == "bid" {
		action = exchange.OrderActBuy
	}
	id, err := strconv.ParseInt(o.keys[o.index], 10, 64)
	if err != nil {
		log.Printf("can not parse id (reason = %v)", err)
	}
	return id, action, value.Price, value.Amount, true
}

func (o *OrderCursor) Reset() {
	o.index = 0
}

func (o *OrderCursor) Len() int {
	return len(o.keys)
}

type TradeContext struct {
	funds                *ExchageFunds
	requester   	     *Requester
	exchangeName		 string
	currencySrc          string
	currencyDst          string
	currencyPair         string
	streamingCallback    exchange.StreamingCallback
	userCallbackData     interface{}
	currencyDatPrice     float64
	bids			 	 [][]float64
	asks			 	 [][]float64
	trades               []*StreamingTradesResponse
}

func (t *TradeContext) GetID() (string) {
	return exchange.MakeTradeID(t.exchangeName, t.currencyPair)
}

func (t *TradeContext) GetExchangeName() (string) {
	return t.exchangeName
}

func (t *TradeContext) Buy(price float64, amount float64) (error) {
	tradeParams := t.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = t.currencyPair
	tradeResponse, _, _, err := t.requester.TradeBuy(tradeParams)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not buy trade (ID = %v)", t.GetID()))
	}
	if tradeResponse.Success != 1 {
		return errors.Errorf("can not buy trade (ID = %v, reason = %v)", t.GetID(), tradeResponse.Error)
	}
	err = updateFunds(t.exchangeName, t.requester, t.funds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not update fund (ID = %v)", t.GetID()))
	}
	return nil
}

func (t *TradeContext) Sell(price float64, amount float64) (error) {
	tradeParams := t.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = t.currencyPair
	tradeResponse, _, _, err := t.requester.TradeSell(tradeParams)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not sell trade (ID = %v)", t.GetID()))
	}
	if tradeResponse.Success != 1 {
		return errors.Errorf("can not sell trade (ID = %v, reason = %v)", t.GetID(), tradeResponse.Error)
	}
	err = updateFunds(t.exchangeName, t.requester, t.funds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not update fund (ID = %v)", t.GetID()))
	}
	return nil
}

func (t *TradeContext) Cancel(orderID int64) (error) {
	tradeCancelOrderParams := t.requester.NewTradeCancelOrderParams()
	tradeCancelOrderParams.IsToken = false
	tradeCancelOrderParams.OrderId = orderID
	tradeCancelOrderResponse, _, _, err := t.requester.TradeCancelOrder(tradeCancelOrderParams)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not cancel order (ID = %v)", t.GetID()))
	}
	if tradeCancelOrderResponse.Success != 1 {
		return errors.Errorf("can not cancel order (ID = %v, reason = %v)", t.GetID(), tradeCancelOrderResponse.Error)
	}
	err = updateFunds(t.exchangeName, t.requester, t.funds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not update fund (ID = %v)", t.GetID()))
	}
	return nil
}

func (t *TradeContext) GetSrcCurrencyFund() (float64, error) {
	return t.funds.get(t.currencySrc), nil
}

func (t *TradeContext) GetDstCurrencyFund() (float64, error) {
	return t.funds.get(t.currencyDst), nil
}

func (t *TradeContext) GetSrcCurrencyName() (string) {
	return t.currencySrc
}

func (t *TradeContext) GetDstCurrencyName() (string) {
	return t.currencyDst
}

func (t *TradeContext) GetPrice() (float64, error) {
	return t.currencyDatPrice, nil
}

func (t *TradeContext) GetBuyBoardCursor() (exchange.BoardCursor, error) {
	return &BoardCursor{
		index:  0,
		values: t.bids,
	}, nil
}

func (t *TradeContext) GetSellBoardCursor() (exchange.BoardCursor, error) {
	return &BoardCursor{
		index:  0,
		values: t.asks,
	}, nil
}

func (t *TradeContext) GetTradeHistoryCursor() (exchange.TradeHistoryCursor, error) {
	return &TradeHistoryCursor{
		index: 0,
		values: t.trades,
	}, nil
}

func (t *TradeContext) GetActiveOrderCursor() (exchange.OrderCursor, error) {
	tradeActiveOrderParams := t.requester.NewTradeActiveOrderParams()
	tradeActiveOrderParams.CurrencyPair = t.currencyPair
	tradeActiveOrderParams.IsToken = false
	tradeActiveOrderResponse, _, _, err := t.requester.TradeActiveOrder(tradeActiveOrderParams)
	if err != nil {
		return nil, err
	}
	if tradeActiveOrderResponse.Success != 1 {
		return nil, errors.Errorf("can not get active order (ID = %v, reason = %v)", t.GetID(), tradeActiveOrderResponse.Error)
	}
	newOrderCursor := &OrderCursor{
		index: 0,
		keys: make([]string, 0, len(tradeActiveOrderResponse.Return)),
		values: tradeActiveOrderResponse.Return,
	}
	for key := range tradeActiveOrderResponse.Return {
		newOrderCursor.keys = append(newOrderCursor.keys, key)
	}
	return newOrderCursor, nil
}

func  (t *TradeContext)  GetMinPriceUnit() (float64) {
	switch t.currencyPair {
	case "btc_jpy":
		return 5
	case "mona_jpy":
		return 0.1
	case "mona_btc":
		return 0.00000001
	default:
		return 0
	}
}

func  (t *TradeContext)  GetMinAmountUnit() (float64) {
	switch t.currencyPair {
	case "btc_jpy":
		return 0.0001
	case "mona_jpy":
		return 1
	case "mona_btc":
		return 1
	default:
		return 0
	}
}

type TradeContextCursor struct {
	tradeContexts []*TradeContext
	index int
}

func (t *TradeContextCursor) Next() (exchange.TradeContext, bool) {
	if t.index >= len(t.tradeContexts)  {
		return nil, false
	}
	tradeContext := t.tradeContexts[t.index]
	t.index++
	return tradeContext, true
}

func (t *TradeContextCursor) Reset() {
	t.index = 0
}

func (t *TradeContextCursor) Len() int {
	return len(t.tradeContexts)
}

type ExchageFunds struct {
	funds  map[string]float64
	mutex  *sync.Mutex
}

func (e *ExchageFunds) update(funds map[string]float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for name, fund := range funds {
		e.funds[name] = fund
	}
}

func (e *ExchageFunds) get(name string) (float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	fund, ok := e.funds[name]
	if !ok {
		return 0
	}
	return fund
}

type Exchange struct {
	config        *ExchangeConfig
	requester     *Requester
	name          string
	tradeContexts []*TradeContext
	funds         *ExchageFunds
}

func (e *Exchange) exchangeStreamingCallback(currencyPair string, streamingResponse *StreamingResponse, StreamingCallbackData interface{}) (error) {
	tradeContext := StreamingCallbackData.(*TradeContext)
	tradeContext.currencyDatPrice = streamingResponse.LastPrice.Price
	tradeContext.bids = streamingResponse.Bids
	tradeContext.asks = streamingResponse.Asks
	tradeContext.trades = streamingResponse.Trades
	err := tradeContext.streamingCallback(tradeContext, tradeContext.userCallbackData)
	if err != nil {
		return errors.Wrap(err,"trade update callback error")
	}
	return nil
}

func (e *Exchange) GetName() string {
	return e.name
}

// Initialize is initalize exchange
func (e *Exchange) Initialize(streamingCallback exchange.StreamingCallback, userCallbackData interface{}) (error) {
	// 設定のcurrencyPairsに応じてコンテキストを作る
	for _, exchangeConfigCurrencyPair := range e.config.CurrencyPairs {
		log.Printf("create zaif trade context (currency pair = %v)", exchangeConfigCurrencyPair)
		srcCurrency := strings.ToLower(exchangeConfigCurrencyPair.Src)
		dstCurrency := strings.ToLower(exchangeConfigCurrencyPair.Dst)
		currencyPair := srcCurrency + "_" + dstCurrency
		tradeContext := &TradeContext{
			funds:             e.funds,
			requester:         e.requester,
			exchangeName:              e.name,
			currencySrc:       srcCurrency,
			currencyDst:       dstCurrency,
			currencyPair:      currencyPair,
			streamingCallback: streamingCallback,
			userCallbackData:  userCallbackData,
		}
		e.tradeContexts = append(e.tradeContexts, tradeContext)
	}
	// fundsを初期化時に更新しておく
	err := updateFunds(e.name, e.requester, e.funds)
	if err != nil {
		return errors.Wrap(err, "can not update fund")
	}
	return nil
}

// Finalize is finalize exchage
func (e *Exchange) Finalize() (error) {
	return nil
}

func (e *Exchange) GetTradeContext(srcCurrency string, dstCcurrency string) (exchange.TradeContext, bool) {
	srcCurrency = strings.ToLower(srcCurrency)
	dstCcurrency = strings.ToLower(dstCcurrency)
	for _, tradeContext := range e.tradeContexts {
		if srcCurrency == tradeContext.currencySrc && dstCcurrency == tradeContext.currencyDst {
			return tradeContext, true
		}
	}
	return nil, false
}

func (e *Exchange) GetTradeContextCursor() (exchange.TradeContextCursor) {
	tradeContextCursor := &TradeContextCursor {
		tradeContexts: e.tradeContexts,
		index: 0,
	}
	return tradeContextCursor
}

// StreamingStart is start streaming
func (e *Exchange) StartStreaming(tradeContext exchange.TradeContext) (error){
	// ストリーミングを開始する
	myTradeContest := tradeContext.(*TradeContext)
	return e.requester.StreamingStart(myTradeContest.currencyPair, e.exchangeStreamingCallback, tradeContext)
}

// StopStreaming is stop streaming
func (e *Exchange) StopStreaming(tradeContext exchange.TradeContext) (error) {
	// ストリーミングを停止する
	myTradeContest := tradeContext.(*TradeContext)
	e.requester.StreamingStop(myTradeContest.currencyPair)
	return nil
}

type exchangeConfigCurrencyPair struct {
	Src string `json:"src" yaml:"src" toml:"src"`
	Dst string `json:"dst" yaml:"dst" toml:"dst"`
}

type ExchangeConfig struct {
	Key           string                         `json:"key"          yaml:"key"          toml:"key"`
	Secret        string                         `json:"secret"       yaml:"secret"       toml:"secret"`
	Retry         int                            `json:"retry"        yaml:"retry"        toml:"retry"`
	Timeout       int                            `json:"timeout"      yaml:"timeout"      toml:"timeout"`
	ReadBufSize   int                            `json:"readBufSize"  yaml:"readBufSize"  toml:"readBufSize"`
	WriteBufSize  int                            `json:"writeBufSize" yaml:"writeBufSize" toml:"writeBufSize"`
	CurrencyPairs []*exchangeConfigCurrencyPair  `json:"currencyPairs" yaml:"currencyPairs" toml:"currencyPairs"`
}

func newZaifExchange(config interface{}) (exchange.Exchange, error)  {
	myConfig := config.(*ExchangeConfig)
	return &Exchange{
		config:        myConfig,
		name :         exchangeName,
		requester:     NewRequester(myConfig.Key, myConfig.Secret, myConfig.Retry, myConfig.Timeout, myConfig.ReadBufSize, myConfig.WriteBufSize),
		tradeContexts: make([]*TradeContext, 0),
		funds : &ExchageFunds{
			funds: make(map[string]float64),
			mutex: new(sync.Mutex),
		},
	}, nil
}

func init() {
	exchange.RegisterExchange(exchangeName, newZaifExchange)
}
