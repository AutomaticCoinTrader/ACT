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

const (
	exchangeName = "zaif"
)

func updateFunds(exchangeName string, requester *Requester, funds *ExchageFunds) (error) {
	info2Response, _, _, err := requester.GetInfo2()
	if err != nil {
		return errors.Wrapf(err, "can not get info2 (ID = %v)", exchangeName)
	}
	if info2Response.Success != 1 {
		return errors.Errorf("can not buy (ID = %v, reason = %v)", exchangeName, info2Response.Error)
	}
	funds.update(map[string]float64{
		"btc":      info2Response.Return.Funds.Btc,
		"bch":      info2Response.Return.Funds.Bch,
		"eth":      info2Response.Return.Funds.Eth,
		"mona":     info2Response.Return.Funds.Mona,
		"xem":      info2Response.Return.Funds.Xem,
		"jpy":      info2Response.Return.Funds.Jpy,
		"zaif":     info2Response.Return.Funds.Zaif,
		"pepecash": info2Response.Return.Funds.Pepecash})
	return nil
}

type BoardCursor struct {
	index  int
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
	index  int
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

func (t *TradeHistoryCursor) GetTradeHistory {
	// XXXXX
	// XXXXX TODO
	// XXXXX
}



type OrderHistoryCursor struct {
	index  int
	keys   []string
	values map[string]TradeHistoryRecordResponse
	keysToken   []string
	valuesToken map[string]TradeHistoryRecordResponse
}

func (o *OrderHistoryCursor) Next() (int64, exchange.OrderAction, float64, float64, bool) {
	if o.index >= len(o.keys) + len(o.keysToken) {
		return 0, "", 0, 0, false
	}
	var key string
	var value TradeHistoryRecordResponse
	if o.index < len(o.keys) {
		key = o.keys[o.index]
		value = o.values[key]
	} else {
		key = o.keysToken[o.index - len(o.keys)]
		value = o.valuesToken[key]
	}
	o.index++
	var action exchange.OrderAction
	if value.Action == "ask" {
		action = exchange.OrderActSell
	} else if value.Action == "bid" {
		action = exchange.OrderActBuy
	} else {
		action = exchange.OrderActUnkown
	}
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Printf("can not parse id (reason = %v)", err)
	}
	return id, action, value.Price, value.Amount, true
}

func (o *OrderHistoryCursor) Reset() {
	o.index = 0
}

func (o *OrderHistoryCursor) Len() int {
	return len(o.keys) + len(o.keysToken)
}




type ActiveOrderCursor struct {
	index  int
	keys   []string
	values map[string]TradeActiveOrderRecordResponse
}

func (o *ActiveOrderCursor) Next() (int64, exchange.OrderAction, float64, float64, bool) {
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

func (o *ActiveOrderCursor) Reset() {
	o.index = 0
}

func (o *ActiveOrderCursor) Len() int {
	return len(o.keys)
}



type TradeContext struct {
	funds                  *ExchageFunds
	requester              *Requester
	exchangeName           string
	streamingCallback      exchange.StreamingCallback
	userCallbackData       interface{}
	currencyPairsBids      map[string][][]float64
	currencyPairsAsks      map[string][][]float64
	currencyPairsLastPrice map[string]float64
	currencyPairsTrades    map[string][]*StreamingTradesResponse
	currencyPairs          []string
	mutex                  *sync.Mutex
}

func (t *TradeContext) GetExchangeName() (string) {
	return t.exchangeName
}

func (t *TradeContext) Buy(currencyPair string, price float64, amount float64) (int64, error) {
	tradeParams := t.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = currencyPair
	tradeResponse, _, _, err := t.requester.TradeBuy(tradeParams)
	if err != nil {
		return 0, errors.Wrapf(err, "can not buy trade (currencyPair = %v)", currencyPair)
	}
	if tradeResponse.Success != 1 {
		return 0, errors.Errorf("can not buy trade (currencyPair = %v, reason = %v)", currencyPair, tradeResponse.Error)
	}
	err = updateFunds(t.exchangeName, t.requester, t.funds)
	if err != nil {
		return tradeResponse.Return.OrderID, errors.Wrapf(err, "can not update fund (currencyPair = %v)", currencyPair)
	}
	return tradeResponse.Return.OrderID, nil
}

func (t *TradeContext) Sell(currencyPair string, price float64, amount float64) (int64, error) {
	tradeParams := t.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = currencyPair
	tradeResponse, _, _, err := t.requester.TradeSell(tradeParams)
	if err != nil {
		return 0, errors.Wrapf(err, "can not sell trade (currencyPair = %v)", currencyPair)
	}
	if tradeResponse.Success != 1 {
		return 0, errors.Errorf("can not sell trade (currencyPair = %v, reason = %v)", currencyPair, tradeResponse.Error)
	}
	err = updateFunds(t.exchangeName, t.requester, t.funds)
	if err != nil {
		return tradeResponse.Return.OrderID, errors.Wrapf(err,"can not update fund (currencyPair = %v)", currencyPair)
	}
	return tradeResponse.Return.OrderID, nil
}

func (t *TradeContext) Cancel(orderID int64) (error) {
	tradeCancelOrderParams := t.requester.NewTradeCancelOrderParams()
	tradeCancelOrderParams.IsToken = false
	tradeCancelOrderParams.OrderId = orderID
	tradeCancelOrderResponse, _, _, err := t.requester.TradeCancelOrder(tradeCancelOrderParams)
	if err != nil {
		return errors.Wrapf(err, "can not cancel order (orderID = %v)", orderID)
	}
	if tradeCancelOrderResponse.Success != 1 {
		return errors.Errorf("can not cancel order (orderID = %v, reason = %v)", orderID, tradeCancelOrderResponse.Error)
	}
	err = updateFunds(t.exchangeName, t.requester, t.funds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not update fund (orderID = %v)", orderID))
	}
	return nil
}

func (t *TradeContext) GetFunds() (map[string]float64, error) {
	return t.funds.copyAll(), nil
}

func (t *TradeContext) GetLastPrice(currencyPair string) (float64, error) {
	return t.currencyPairsLastPrice[currencyPair], nil
}

func (t *TradeContext) GetBuyBoardCursor(currencyPair string) (exchange.BoardCursor, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return &BoardCursor{
		index:  0,
		values: t.currencyPairsBids[currencyPair],
	}, nil
}

func (t *TradeContext) GetSellBoardCursor(currencyPair string) (exchange.BoardCursor, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return &BoardCursor{
		index:  0,
		values: t.currencyPairsAsks[currencyPair],
	}, nil
}

func (t *TradeContext) GetTradesCursor(currencyPair string) (exchange.TradesCursor, error) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return &TradeHistoryCursor{
		index:  0,
		values: t.currencyPairsTrades[currencyPair],
	}, nil
}



func (t *TradeContext) GetMyTradeHistoryCursor(count int64) (exchange.OrderCursor, error) {
	tradeHistoryParams :=  t.requester.NewTradeHistoryParams()
	tradeHistoryParams.IsToken = false
	tradeHistoryParams.Count = count
	tradeHistoryResponse, _, _, err := t.requester.TradeHistory(tradeHistoryParams)
	if err != nil {
		return nil, err
	}
	if tradeHistoryResponse.Success != 1 {
		return nil, errors.Errorf("can not get trade history (reason = %v)", tradeHistoryResponse.Error)
	}
	tradeHistoryParams =  t.requester.NewTradeHistoryParams()
	tradeHistoryParams.IsToken = true
	tradeHistoryParams.Count = count
	tradeHistoryTokenResponse, _, _, err := t.requester.TradeHistory(tradeHistoryParams)
	if err != nil {
		return nil, err
	}
	if tradeHistoryTokenResponse.Success != 1 {
		return nil, errors.Errorf("can not get trade history (reason = %v)", tradeHistoryTokenResponse.Error)
	}
	newOrderCursor := &ActiveOrderCursor{
		index:  0,
		keys:   make([]string, 0, len(tradeHistoryResponse.Return) + len(tradeHistoryTokenResponse.Return)),
		values: make(map[string]TradeHistoryRecordResponse),
	}




}


func (t *TradeContext) GetMyActiveOrderCursor() (exchange.OrderCursor, error) {
	tradeActiveOrderParams := t.requester.NewTradeActiveOrderParams()

	tradeActiveOrderParams.IsTokenBoth = true
	tradeActiveOrderBothResponse, _, _, err := t.requester.TradeActiveOrderBoth(tradeActiveOrderParams)
	if err != nil {
		return nil, err
	}
	if tradeActiveOrderBothResponse.Success != 1 {
		return nil, errors.Errorf("can not get active order (reason = %v)", tradeActiveOrderBothResponse.Error)
	}
	newOrderCursor := &ActiveOrderCursor{
		index:  0,
		keys:   make([]string, 0, len(tradeActiveOrderBothResponse.Return)),
		values: tradeActiveOrderBothResponse.Return,
	}
	for key := range tradeActiveOrderBothResponse.Return {
		newOrderCursor.keys = append(newOrderCursor.keys, key)
	}
	return newOrderCursor, nil
}





func (t *TradeContext) GetMinPriceUnit() (float64) {
	switch t.currencyPair {
	case "btc_jpy":
		return 5
	case "xem_jpy":
		return 0.0001
	case "mona_jpy":
		return 0.1
	case "bch_jpy":
		return 5
	case "eth_jpy":
		return 5
	case "zaif_jpy":
		return 0.0001
	case "pepecash_jpy":
		return 0.0001
	case "xem_btc":
		return 0.00000001
	case "mona_btc":
		return 0.00000001
	case "bch_btc":
		return 0.0001
	case "eth_btc":
		return 0.0001
	case "zaif_btc":
		return 0.00000001
	case "pepecash_btc":
		return 0.00000001
	default:
		return 0
	}
}

func (t *TradeContext) GetMinAmountUnit() (float64) {
	switch t.currencyPair {
	case "btc_jpy":
		return 0.0001
	case "xem_jpy":
		return 0.1
	case "mona_jpy":
		return 1
	case "bch_jpy":
		return 0.0001
	case "eth_jpy":
		return 0.0001
	case "zaif_jpy":
		return 0.1
	case "pepecash_jpy":
		return 0.0001
	case "xem_btc":
		return 1
	case "mona_btc":
		return 1
	case "bch_btc":
		return 0.0001
	case "eth_btc":
		return 0.0001
	case "zaif_btc":
		return 1
	case "pepecash_btc":
		return 1
	default:
		return 0
	}
}

type TradeContextCursor struct {
	tradeContexts []*TradeContext
	index         int
}

func (t *TradeContextCursor) Next() (exchange.TradeContext, bool) {
	if t.index >= len(t.tradeContexts) {
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
	funds map[string]float64
	mutex *sync.Mutex
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

func (e *ExchageFunds) copyAll() (map[string]float64) {
	newFunds := make(map[string]float64)
	e.mutex.Lock()
	defer e.mutex.Unlock()
	for name, value, := range e.funds {
		newFunds[name] = value
	}
	return newFunds
}


type Exchange struct {
	config        *ExchangeConfig
	requester     *Requester
	name          string
	tradeContext  *TradeContext
	funds         *ExchageFunds
}

func (e *Exchange) exchangeStreamingCallback(currencyPair string, streamingResponse *StreamingResponse, StreamingCallbackData interface{}) (error) {
	myTradeContext := StreamingCallbackData.(*TradeContext)
	myTradeContext.mutex.Lock()
	myTradeContext.currencyPairsLastPrice[currencyPair] = streamingResponse.LastPrice.Price
	myTradeContext.currencyPairsBids[currencyPair] = streamingResponse.Bids
	myTradeContext.currencyPairsAsks[currencyPair] = streamingResponse.Asks
	myTradeContext.currencyPairsTrades[currencyPair] = streamingResponse.Trades
	myTradeContext.mutex.Unlock()
	err := myTradeContext.streamingCallback(currencyPair, myTradeContext)
	if err != nil {
		return errors.Wrap(err, "trade update callback error")
	}
	return nil
}

func (e *Exchange) GetName() string {
	return e.name
}

// Initialize is initalize exchange
func (e *Exchange) Initialize(streamingCallback exchange.StreamingCallback) (error) {
	e.tradeContext = &TradeContext{
		funds:             e.funds,
		requester:         e.requester,
		exchangeName:      e.name,
		streamingCallback: streamingCallback,
		currencyPairs:     e.config.CurrencyPairs,
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

func (e *Exchange) GetTradeContext() (exchange.TradeContext) {
	return e.tradeContext
}

// StreamingStart is start streaming
func (e *Exchange) StartStreaming(tradeContext exchange.TradeContext) (error) {
	// ストリーミングを開始する
	myTradeContext := tradeContext.(*TradeContext)
	for _, currencyPair := range myTradeContext.currencyPairs {
		currencyPair = strings.ToLower(currencyPair)
		err := e.requester.StreamingStart(currencyPair, e.exchangeStreamingCallback, tradeContext)
		if (err != nil) {
			return errors.Wrapf(err, "can not start streaming (currency_pair = %v)", currencyPair);
		}
	}
	return nil
}

// StopStreaming is stop streaming
func (e *Exchange) StopStreaming(tradeContext exchange.TradeContext) (error) {
	// ストリーミングを停止する
	myTradeContext := tradeContext.(*TradeContext)
	for _, currencyPair := range myTradeContext.currencyPairs {
		currencyPair = strings.ToLower(currencyPair)
		e.requester.StreamingStop(currencyPair)

	}
	return nil
}


type ExchangeConfig struct {
	Key           string                        `json:"key"          yaml:"key"          toml:"key"`
	Secret        string                        `json:"secret"       yaml:"secret"       toml:"secret"`
	Retry         int                           `json:"retry"        yaml:"retry"        toml:"retry"`
	RetryWait     int                           `json:"retryWait"    yaml:"retryWait"    toml:"retryWait"`
	Timeout       int                           `json:"timeout"      yaml:"timeout"      toml:"timeout"`
	ReadBufSize   int                           `json:"readBufSize"  yaml:"readBufSize"  toml:"readBufSize"`
	WriteBufSize  int                           `json:"writeBufSize" yaml:"writeBufSize" toml:"writeBufSize"`
	CurrencyPairs []string                      `json:"currencyPairs" yaml:"currencyPairs" toml:"currencyPairs"`
}

func newZaifExchange(config interface{}) (exchange.Exchange, error) {
	myConfig := config.(*ExchangeConfig)
	return &Exchange{
		config:        myConfig,
		name:          exchangeName,
		requester:     NewRequester(myConfig.Key, myConfig.Secret, myConfig.Retry, myConfig.RetryWait, myConfig.Timeout, myConfig.ReadBufSize, myConfig.WriteBufSize),
		tradeContext:  nil,
		funds: &ExchageFunds{
			funds: make(map[string]float64),
			mutex: new(sync.Mutex),
		},
	}, nil
}

func init() {
	exchange.RegisterExchange(exchangeName, newZaifExchange)
}
