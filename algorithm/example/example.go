package example

import (
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"log"
)

const (
	algorithmName string = "example"
)

type Example struct {
	name           string
	config         *Config
}

func (l *Example) GetName() (string) {
	return l.name
}

func (l *Example) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func (l *Example) Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	// trade
	log.Printf("trade")
	return nil
}

func (l *Example) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

type Config struct {
}

func newExample(config interface{}) (algorithm.TradeAlgorithm, error) {
	myConfig := config.(*Config)
	return &Example{
		name:           algorithmName,
		config:         myConfig,
	}, nil
}

type ArbitrageExample struct {
	name           string
	config         *ArbitrageConfig
}

func (l *ArbitrageExample) GetName() (string) {
	return l.name
}

func (l *ArbitrageExample) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func (l *ArbitrageExample) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	// arbitrage trade
	log.Printf("arbitrage trade")
	return nil
}

func (l *ArbitrageExample) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

type ArbitrageConfig struct {
}

func newArbitrageExample(config interface{}) (algorithm.ArbitrageTradeAlgorithm, error) {
	myConfig := config.(*ArbitrageConfig)
	return &ArbitrageExample{
		name:           algorithmName,
		config:         myConfig,
	}, nil
}

func init() {
	algorithm.RegisterAlgorithm(algorithmName, newExample, newArbitrageExample)
}
