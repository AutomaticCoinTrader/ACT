package example

import (
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"github.com/AutomaticCoinTrader/ACT/configurator"
	"github.com/pkg/errors"
	"log"
	"path"
)

const (
	algorithmName string = "example"
)


type TradeConfig struct {
}

type ArbitrageTradeConfig struct {
}

type Config struct {
	Trade          *TradeConfig          `json:"trade"          yaml:"trade"          toml:"trade"`
	ArbitrageTrade *ArbitrageTradeConfig `json:"arbitrageTrade" yaml:"arbitrageTrade" toml:"arbitrageTrade"`
}

type Example struct {
	name           string
	config         *TradeConfig
}

func (e *Example) GetName() (string) {
	return e.name
}

func (e *Example) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func (e *Example) Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	// trade
	log.Printf("trade")
	return nil
}

func (e *Example) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func newExample(configDir string) (algorithm.TradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	config := new(Config)
	err = cf.Load(config)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &Example{
		name:           algorithmName,
		config:         config.Trade,
	}, nil
}

type ArbitrageExample struct {
	name           string
	config         *ArbitrageTradeConfig
}

func (a *ArbitrageExample) GetName() (string) {
	return a.name
}

func (a *ArbitrageExample) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func (a *ArbitrageExample) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	// arbitrage trade
	log.Printf("arbitrage trade")
	return nil
}

func (a *ArbitrageExample) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func newArbitrageExample(configDir string) (algorithm.ArbitrageTradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	config := new(Config)
	err = cf.Load(config)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &ArbitrageExample{
		name:           algorithmName,
		config:         config.ArbitrageTrade,
	}, nil
}

func init() {
	algorithm.RegisterAlgorithm(algorithmName, newExample, newArbitrageExample)
}



