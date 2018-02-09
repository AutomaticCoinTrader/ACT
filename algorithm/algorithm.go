package algorithm

import (
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/notifier"
)

const (
	AlgorithmConfigDir = "algorithm"
)

type InternalTradeAlgorithm interface {
	GetName() (string)
	Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error)
	Update(currencyPair string, tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error)
	Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error)
}

type ExternalTradeAlgorithm interface {
	GetName() (string)
	Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error)
	Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error)
	Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error)
}
