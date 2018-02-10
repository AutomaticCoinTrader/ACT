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

func (b *BoardCursor) PriceAll() []float64 {
	newPriceAll := make([]float64, 0, len(b.values))
	for _, v := range b.values {
		newPriceAll = append(newPriceAll, v[0])
	}
	return newPriceAll
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
	for k := range values {
		newOrderHistoryCursor.keys = append(newOrderHistoryCursor.keys, k)
	}
	for k := range valuesToken {
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
	for k := range values {
		newActiveOrderCursor.keys = append(newActiveOrderCursor.keys, k)
	}
	for k := range valuesToken {
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


type Exchange struct {
	config        *ExchangeConfig
	requester     *Requester
	funds                  *CurrencyFunds
	streamingCallback      exchange.StreamingCallback
	currencyPairs          []string
	currencyPairsInfo      *currencyPairsInfo
}



func (e *Exchange) GetName() (string) {
	return exchangeName
}

func (e *Exchange) GetCurrencyPairs() ([]string) {
	return e.currencyPairs
}

func (e *Exchange) Buy(currencyPair string, price float64, amount float64) (int64, error) {
	tradeParams := e.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = currencyPair
	tradeResponse, _, _, err := e.requester.TradeBuy(tradeParams)
	if err != nil {
		return 0, errors.Wrapf(err, "can not buy trade (exchange = %v, currencyPair = %v)", exchangeName, currencyPair)
	}
	if tradeResponse.Success != 1 {
		return 0, errors.Errorf("can not buy trade (exchange = %v, currencyPair = %v, reason = %v)", exchangeName, currencyPair, tradeResponse.Error)
	}
	err = updateFunds(e.requester, e.funds)
	if err != nil {
		return tradeResponse.Return.OrderID, errors.Wrapf(err, "can not update fund (exchange = %v, currencyPair = %v)", exchangeName, currencyPair)
	}
	return tradeResponse.Return.OrderID, nil
}

func (e *Exchange) Sell(currencyPair string, price float64, amount float64) (int64, error) {
	tradeParams := e.requester.NewTradeParams()
	tradeParams.Price = price
	tradeParams.Amount = amount
	tradeParams.CurrencyPair = currencyPair
	tradeResponse, _, _, err := e.requester.TradeSell(tradeParams)
	if err != nil {
		return 0, errors.Wrapf(err, "can not sell trade (currencyPair = %v)", currencyPair)
	}
	if tradeResponse.Success != 1 {
		return 0, errors.Errorf("can not sell trade (currencyPair = %v, reason = %v)", currencyPair, tradeResponse.Error)
	}
	err = updateFunds(e.requester, e.funds)
	if err != nil {
		return tradeResponse.Return.OrderID, errors.Wrapf(err,"can not update fund (currencyPair = %v)", currencyPair)
	}
	return tradeResponse.Return.OrderID, nil
}

func (e *Exchange) Cancel(orderID int64) (error) {
	tradeCancelOrderParams := e.requester.NewTradeCancelOrderParams()
	tradeCancelOrderParams.IsToken = false
	tradeCancelOrderParams.OrderId = orderID
	tradeCancelOrderResponse, _, _, err := e.requester.TradeCancelOrder(tradeCancelOrderParams)
	if err != nil {
		return errors.Wrapf(err, "can not cancel order (orderID = %v)", orderID)
	}
	if tradeCancelOrderResponse.Success != 1 {
		return errors.Errorf("can not cancel order (orderID = %v, reason = %v)", orderID, tradeCancelOrderResponse.Error)
	}
	err = updateFunds(e.requester, e.funds)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("can not update fund (orderID = %v)", orderID))
	}
	return nil
}

func (e *Exchange) GetFunds() (map[string]float64, error) {
	return e.funds.all(), nil
}

func (e *Exchange) GetLastPrice(currencyPair string) (float64, error) {
	return e.currencyPairsInfo.getLastPrice(currencyPair), nil
}

func (e *Exchange) GetSellBoardCursor(currencyPair string) (exchange.BoardCursor, error) {
	return &BoardCursor{
		index:  0,
		values: e.currencyPairsInfo.getAsks(currencyPair),
	}, nil
}

func (e *Exchange) GetBuyBoardCursor(currencyPair string) (exchange.BoardCursor, error) {

	return &BoardCursor{
		index:  0,
		values: e.currencyPairsInfo.getBids(currencyPair),
	}, nil
}

func (e *Exchange) GetTradesCursor(currencyPair string) (exchange.TradesCursor, error) {
	return &TradeHistoryCursor{
		index:  0,
		values: e.currencyPairsInfo.getTrades(currencyPair),
	}, nil
}

func (e *Exchange) GetOrderHistoryCursor(count int64) (exchange.OrderCursor, error) {
	tradeHistoryParams :=  e.requester.NewTradeHistoryParams()
	tradeHistoryParams.IsToken = false
	tradeHistoryParams.Count = count
	tradeHistoryResponse, _, _, err := e.requester.TradeHistory(tradeHistoryParams)
	if err != nil {
		return nil, err
	}
	if tradeHistoryResponse.Success != 1 {
		return nil, errors.Errorf("can not get trade history (reason = %v)", tradeHistoryResponse.Error)
	}
	tradeHistoryParams = e.requester.NewTradeHistoryParams()
	tradeHistoryParams.IsToken = true
	tradeHistoryParams.Count = count
	tradeHistoryTokenResponse, _, _, err := e.requester.TradeHistory(tradeHistoryParams)
	if err != nil {
		return nil, err
	}
	if tradeHistoryTokenResponse.Success != 1 {
		return nil, errors.Errorf("can not get trade history (reason = %v)", tradeHistoryTokenResponse.Error)
	}
	return newOrderHistoryCursor(tradeHistoryResponse.Return, tradeHistoryTokenResponse.Return), nil
}

func (e *Exchange) GetActiveOrderCursor() (exchange.OrderCursor, error) {
	tradeActiveOrderParams := e.requester.NewTradeActiveOrderParams()

	tradeActiveOrderParams.IsTokenBoth = true
	tradeActiveOrderBothResponse, _, _, err := e.requester.TradeActiveOrderBoth(tradeActiveOrderParams)
	if err != nil {
		return nil, err
	}
	if tradeActiveOrderBothResponse.Success != 1 {
		return nil, errors.Errorf("can not get active order (reason = %v)", tradeActiveOrderBothResponse.Error)
	}
	return newActiveOrderCursor(tradeActiveOrderBothResponse.Return.ActiveOrders, tradeActiveOrderBothResponse.Return.TokenActiveOrders), nil
}

func (e *Exchange) GetMinPriceUnit(currencyPair string) (float64) {
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

func (e *Exchange) GetMinAmountUnit(currencyPair string) (float64) {
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

func (e *Exchange) exchangeStreamingCallback(currencyPair string, streamingResponse *StreamingResponse, StreamingCallbackData interface{}) (error) {

	e.currencyPairsInfo.update(currencyPair, streamingResponse.Bids, streamingResponse.Asks, streamingResponse.LastPrice.Price, streamingResponse.Trades)
	err := e.streamingCallback(currencyPair, e)
	if err != nil {
		return errors.Wrap(err, "trade update callback error")
	}
	return nil
}

// Initialize is initalize exchange
func (e *Exchange) Initialize() (error) {
	// fundsを初期化時に更新しておく
	err := updateFunds(e.requester, e.funds)
	if err != nil {
		return errors.Wrap(err, "can not update fund")
	}
	return nil
}

// Finalize is finalize exchage
func (e *Exchange) Finalize() (error) {
	return nil
}

// StreamingStart is start streaming
func (e *Exchange) StartStreamings(streamingCallback exchange.StreamingCallback) (error) {
	// ストリーミングを開始する
	e.streamingCallback = streamingCallback
	for _, currencyPair := range e.currencyPairs {
		currencyPair = strings.ToLower(currencyPair)
		err := e.requester.StreamingStart(currencyPair, e.exchangeStreamingCallback, e)
		if (err != nil) {
			return errors.Wrapf(err, "can not start streaming (currency_pair = %v)", currencyPair);
		}
	}
	return nil
}

// StopStreaming is stop streaming
func (e *Exchange) StopStreamings() (error) {
	// ストリーミングを停止する

	for _, currencyPair := range e.currencyPairs {
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
		funds:             &CurrencyFunds{
			funds: make(map[string]float64),
			mutex: new(sync.Mutex),
		},
		currencyPairs:     myConfig.CurrencyPairs,
		currencyPairsInfo: &currencyPairsInfo{
			Bids: make(map[string][][]float64),
			Asks: make(map[string][][]float64),
			LastPrice: make(map[string]float64),
			Trades: make(map[string][]*StreamingTradesResponse),
			mutex: new(sync.Mutex),
		},
	}, nil
}

func init() {
	exchange.RegisterExchange(exchangeName, newZaifExchange)
}
