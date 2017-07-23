package coincheck

import (
	"github.com/AutomaticCoinTrader/ACT/exchange"
)

type CoincheckExchangeConfig struct {
	APIKey    string
	APISecret string
}

type CoincheckExchange struct {
}

func (ce *CoincheckExchange) GetName() string {
	return "Coincheck"
}

func (ce *CoincheckExchange) Initialize(streamingCallback exchange.StreamingCallback, userCallbackData interface{}) (error) {
	return nil
}

func (ce *CoincheckExchange) Finalize() (error) {
	return nil
}

func (ce *CoincheckExchange) GetTradeContext(srcCurrency string, dstCurrency string) (exchange.TradeContext, bool) {
	return nil, true
}

func (ce *CoincheckExchange) GetTradeContextCursor() (exchange.TradeContextCursor) {
	return nil
}

func (ce *CoincheckExchange) StartStreaming(tradeContext exchange.TradeContext) (error) {
	return nil
}

func (ce *CoincheckExchange) StopStreaming(tradeContext exchange.TradeContext) (error) {
	return nil
}

func newCoincheckExchange(config interface {}) (exchange.Exchange, error) {
	// coincheckConfig := config.(*CoincheckExchangeConfig)
	return &CoincheckExchange {}, nil
}
