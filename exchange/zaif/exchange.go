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

func updateFunds(requester *Requester, funds *CurrencyFunds) (error) {
	info2Response, _, _, err := requester.GetInfo2()
	if err != nil {
		return errors.Wrapf(err, "can not get info2 (exchange = %v, reason = %v)", exchangeName)
	}
	if info2Response.Success != 1 {
		return errors.Errorf("can not buy (exchange = %v, reason = %v)", exchangeName, info2Response.Error)
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
		log.Printf("can not parse id (exchange = %v, reason = %v)", exchangeName, err)
	}
	return id, action, value.Price, value.Amount, true
}

func (o *OrderHistoryCursor) Reset() {
	o.index = 0
}

func (o *OrderHistoryCursor) Len() int {
	return len(o.keys) + len(o.keysToken)
}

func newOrderHistoryCursor(values map[string]TradeHistoryRecordResponse, valuesToken map[string]TradeHistoryRecordResponse) (*OrderHistoryCursor) {
	newOrderHistoryCursor := &OrderHistoryCursor {
		index: 0,
		keys: make([]string, 0, len(values)),
		values: values,
		keysToken: make([]string, 0, len(valuesToken)),
		valuesToken: valuesToken,
	}
	for k, _ := range values {
		newOrderHistoryCursor.keys = append(newOrderHistoryCursor.keys, k)
	}
	for k, _ := range valuesToken {
		newOrderHistoryCursor.keysToken = append(newOrderHistoryCursor.keysToken, k)
	}
	return newOrderHistoryCursor
}

type ActiveOrderCursor struct {
	index  int
	keys   []string
	values map[string]TradeActiveOrderRecordResponse
	keysToken   []string
	valuesToken map[string]TradeActiveOrderRecordResponse
}

func (o *ActiveOrderCursor) Next() (int64, exchange.OrderAction, float64, float64, bool) {
	if o.index >= len(o.keys) + len(o.keys) {
		return 0, "", 0, 0, false
	}
	var key string
	var value TradeActiveOrderRecordResponse
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
	}
	id, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		log.Printf("can not parse id (exchange = %v, reason = %v)", exchangeName, err)
	}
	return id, action, value.Price, value.Amount, true
}

func (o *ActiveOrderCursor) Reset() {
	o.index = 0
}

func (o *ActiveOrderCursor) Len() int {
	return len(o.keys)
}

func newActiveOrderCursor(values map[string]TradeActiveOrderRecordResponse, valuesToken map[string]TradeActiveOrderRecordResponse) (*ActiveOrderCursor) {
	newActiveOrderCursor := &ActiveOrderCursor {
		index: 0,
		keys: make([]string, 0, len(values)),
		values: values,
		keysToken: make([]string, 0, len(valuesToken)),
		valuesToken: valuesToken,
	}
	for k, _ := range values {
		newActiveOrderCursor.keys = append(newActiveOrderCursor.keys, k)
	}
	for k, _ := range valuesToken {
		newActiveOrderCursor.keysToken = append(newActiveOrderCursor.keysToken, k)
	}
	return newActiveOrderCursor
}

type CurrencyFunds struct {
	funds map[string]float64
	mutex *sync.Mutex
}

func (e *CurrencyFunds) update(funds map[string]float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.funds = funds
}

func (e *CurrencyFunds) get(name string) (float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	fund, ok := e.funds[name]
	if !ok {
		return 0
	}
	return fund
}

func (e *CurrencyFunds) all() (map[string]float64) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.funds
}

type currencyPairsInfo struct {
	Bids      map[string][][]float64
	Asks      map[string][][]float64
	LastPrice map[string]float64
	Trades    map[string][]*StreamingTradesResponse
	mutex     *sync.Mutex
}

func (c *currencyPairsInfo) update(currencyPair string, currencyPairsBids [][]float64, currencyPairsAsks [][]float64, currencyPairsLastPrice float64, currencyPairsTrades []*StreamingTradesResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Bids[currencyPair] = currencyPairsBids
	c.Asks[currencyPair] = currencyPairsAsks
	c.LastPrice[currencyPair] = currencyPairsLastPrice
	c.Trades[currencyPair] = currencyPairsTrades
}

func (c *currencyPairsInfo) getBids(currencyPair string) ([][]float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Bids[currencyPair]
}

func (c *currencyPairsInfo) getAsks(currencyPair string) ([][]float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Asks[currencyPair]
}

func (c *currencyPairsInfo) getLastPrice(currencyPair string) (float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.LastPrice[currencyPair]
}

func (c *currencyPairsInfo) getTrades(currencyPair string) ([]*StreamingTradesResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.Trades[currencyPair]
}

type TradeContext struct {
	funds                  *CurrencyFunds
	requester              *Requester
	streamingCallback      exchange.StreamingCallback
	currencyPairs          []string
	currencyPairsInfo      *currencyPairsInfo
	userCallbackData       interface{}
}

func (t *TradeContext) GetExchangeName() (string) {
	return exchangeName
}

func (t *TradeContext) GetCurrencyPairs() ([]string) {
	return t.currencyPairs
}

func (t *TradeContext) Buy(currencyPair string, price float64, amount float64) (int64, error) {
	tradeParams := t.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = currencyPair
	tradeResponse, _, _, err := t.requester.TradeBuy(tradeParams)
	if err != nil {
		return 0, errors.Wrapf(err, "can not buy trade (exchange = %v, currencyPair = %v)", exchangeName, currencyPair)
	}
	if tradeResponse.Success != 1 {
		return 0, errors.Errorf("can not buy trade (exchange = %v, currencyPair = %v, reason = %v)", exchangeName, currencyPair, tradeResponse.Error)
	}
	err = updateFunds(t.requester, t.funds)
	if err != nil {
		return tradeResponse.Return.OrderID, errors.Wrapf(err, "can not update fund (exchange = %v, currencyPair = %v)", exchangeName, currencyPair)
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
	err = updateFunds(t.requester, t.funds)
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
	err = updateFunds(t.requester, t.funds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not update fund (orderID = %v)", orderID))
	}
	return nil
}

func (t *TradeContext) GetFunds() (map[string]float64, error) {
	return t.funds.all(), nil
}

func (t *TradeContext) GetLastPrice(currencyPair string) (float64, error) {
	return t.currencyPairsInfo.getLastPrice(currencyPair), nil
}

func (t *TradeContext) GetSellBoardCursor(currencyPair string) (exchange.BoardCursor, error) {
	return &BoardCursor{
		index:  0,
		values: t.currencyPairsInfo.getAsks(currencyPair),
	}, nil
}

func (t *TradeContext) GetBuyBoardCursor(currencyPair string) (exchange.BoardCursor, error) {

	return &BoardCursor{
		index:  0,
		values: t.currencyPairsInfo.getBids(currencyPair),
	}, nil
}

func (t *TradeContext) GetTradesCursor(currencyPair string) (exchange.TradesCursor, error) {
	return &TradeHistoryCursor{
		index:  0,
		values: t.currencyPairsInfo.getTrades(currencyPair),
	}, nil
}

func (t *TradeContext) GetOrderHistoryCursor(count int64) (exchange.OrderCursor, error) {
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
	return newOrderHistoryCursor(tradeHistoryResponse.Return, tradeHistoryTokenResponse.Return), nil
}

func (t *TradeContext) GetActiveOrderCursor() (exchange.OrderCursor, error) {
	tradeActiveOrderParams := t.requester.NewTradeActiveOrderParams()

	tradeActiveOrderParams.IsTokenBoth = true
	tradeActiveOrderBothResponse, _, _, err := t.requester.TradeActiveOrderBoth(tradeActiveOrderParams)
	if err != nil {
		return nil, err
	}
	if tradeActiveOrderBothResponse.Success != 1 {
		return nil, errors.Errorf("can not get active order (reason = %v)", tradeActiveOrderBothResponse.Error)
	}
	return newActiveOrderCursor(tradeActiveOrderBothResponse.Return.ActiveOrders, tradeActiveOrderBothResponse.Return.TokenActiveOrders), nil
}

func (t *TradeContext) GetMinPriceUnit(currencyPair string) (float64) {
	switch currencyPair {
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

func (t *TradeContext) GetMinAmountUnit(currencyPair string) (float64) {
	switch currencyPair {
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

type Exchange struct {
	config        *ExchangeConfig
	requester     *Requester
	tradeContext  *TradeContext
}

func (e *Exchange) exchangeStreamingCallback(currencyPair string, streamingResponse *StreamingResponse, StreamingCallbackData interface{}) (error) {
	myTradeContext := StreamingCallbackData.(*TradeContext)
	myTradeContext.currencyPairsInfo.update(currencyPair, streamingResponse.Bids, streamingResponse.Asks, streamingResponse.LastPrice.Price, streamingResponse.Trades)
	err := myTradeContext.streamingCallback(currencyPair, myTradeContext)
	if err != nil {
		return errors.Wrap(err, "trade update callback error")
	}
	return nil
}

func (e *Exchange) GetName() string {
	return exchangeName
}

// Initialize is initalize exchange
func (e *Exchange) Initialize(streamingCallback exchange.StreamingCallback) (error) {
	e.tradeContext = &TradeContext{
		funds:             &CurrencyFunds{
			funds: make(map[string]float64),
			mutex: new(sync.Mutex),
		},
		requester:         e.requester,
		streamingCallback: streamingCallback,
		currencyPairs:     e.config.CurrencyPairs,
		currencyPairsInfo: &currencyPairsInfo{
			mutex: new(sync.Mutex),
		},
	}
	// fundsを初期化時に更新しておく
	err := updateFunds(e.requester, e.tradeContext.funds)
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
func (e *Exchange) StartStreamings(tradeContext exchange.TradeContext) (error) {
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
func (e *Exchange) StopStreamings(tradeContext exchange.TradeContext) (error) {
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
		requester:     NewRequester(myConfig.Key, myConfig.Secret, myConfig.Retry, myConfig.RetryWait, myConfig.Timeout, myConfig.ReadBufSize, myConfig.WriteBufSize),
		tradeContext:  nil,
	}, nil
}

func init() {
	exchange.RegisterExchange(exchangeName, newZaifExchange)
}
