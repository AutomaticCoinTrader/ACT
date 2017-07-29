package coincheck

import (
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"log"
)

type CoincheckExchangeConfig struct {
	APIKey    string `json:"apikey" yaml:"apikey"`
	APISecret string `json:"apisecret" yaml:"apisecret"`
}

type CoincheckExchange struct {
	config *CoincheckExchangeConfig
	context []*CoincheckTradeContext
}

type CoincheckTradeContext struct {
	callback exchange.StreamingCallback
}

type CoincheckBoardCursor struct {
}

type CoincheckTradeHistoryCursor struct {
}

type CoincheckOrderCursor struct {
}

/// cursors

/// tradecontext

func (ct *CoincheckTradeContext) GetID() string {
	return "coincheck-trade-context-id"
}

func (ct *CoincheckTradeContext) GetExchangeName() string {
	return "coincheck-exchange-id"
}

func (ct *CoincheckTradeContext) Buy(price float64, amount float64) (error) {
	return nil
}

func (ct *CoincheckTradeContext) Sell(price float64, amount float64) (error) {
	return nil
}

func (ct *CoincheckTradeContext) Cancel(orderID int64) (error) {
	return nil
}

func (ct *CoincheckTradeContext) GetSrcCurrencyFund() (float64, error) {
	return 0.0, nil
}

func (ct *CoincheckTradeContext) GetDstCurrencyFund() (float64, error) {
	return 0.0, nil
}

func (ct *CoincheckTradeContext) GetSrcCurrencyName() (string) {
	return "coincheck-trade-context-src-currency"
}

func (ct *CoincheckTradeContext) GetDstCurrencyName() (string) {
	return "coincheck-trade-context-dst-currency"
}

func (ct *CoincheckTradeContext) GetPrice() (float64, error) {
	return 0.0, nil
}

func (ct *CoincheckTradeContext) GetBuyBoardCursor() (exchange.BoardCursor, error) {
	return nil, nil
}

func (ct *CoincheckTradeContext) GetSellBoardCursor() (exchange.BoardCursor, error) {
	return nil, nil
}

func (ct *CoincheckTradeContext) GetTradeHistoryCursor() (exchange.TradeHistoryCursor, error) {
	return nil, nil
}

func (ct *CoincheckTradeContext) GetActiveOrderCursor() (exchange.OrderCursor, error) {
	return nil, nil
}

func (ct *CoincheckTradeContext) GetMinPriceUnit() (float64) {
	return 1.0
}

func (ct *CoincheckTradeContext) GetMinAmountUnit() (float64) {
	return 1.0
}


func (ce *CoincheckExchange) GetName() string {
	return "Coincheck"
}

// ここで tradecontext を作る
func (ce *CoincheckExchange) Initialize(streamingCallback exchange.StreamingCallback, userCallbackData interface{}) (error) {
	ce.context = make([]*CoincheckTradeContext, 0)
	ce.context = append(ce.context, &CoincheckTradeContext{callback: streamingCallback})
	log.Println("coincheckexchange Initialize")
	return nil
}

func (ce *CoincheckExchange) Finalize() (error) {
	return nil
}

func (ce *CoincheckExchange) GetTradeContext(srcCurrency string, dstCurrency string) (exchange.TradeContext, bool) {
	return nil, true
}

type CoincheckTradeContextCursor struct {
	context []*CoincheckTradeContext
	index int
}

func (ctcc *CoincheckTradeContextCursor) Next() (tradeContext exchange.TradeContext, ok bool) {
	if ctcc.index >= ctcc.Len() {
		return nil, false
	}
	res := ctcc.context[ctcc.index]
	ctcc.index++
	return res, true
}

func (ctcc *CoincheckTradeContextCursor) Reset() {
}

func (ctcc *CoincheckTradeContextCursor) Len() int {
	return len(ctcc.context)
}

func (ce *CoincheckExchange) GetTradeContextCursor() (exchange.TradeContextCursor) {
	return &CoincheckTradeContextCursor{ context: ce.context, index: 0 }
}

func (ce *CoincheckExchange) StartStreaming(tradeContext exchange.TradeContext) (error) {
	log.Println("CoincheckExchange StartStreaming")
	return nil
}

func (ce *CoincheckExchange) StopStreaming(tradeContext exchange.TradeContext) (error) {
	log.Println("CoincheckExchange StopStreaming")
	return nil
}

func newCoincheckExchange(config interface {}) (exchange.Exchange, error) {
	coincheckConfig := config.(*CoincheckExchangeConfig)
	return &CoincheckExchange { config: coincheckConfig, context: nil }, nil
}

func init() {
	exchange.RegisterExchange("coincheck", newCoincheckExchange)
}