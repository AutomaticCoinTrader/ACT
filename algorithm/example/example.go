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

func newExample(configDir string) (algorithm.TradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		errors.Errorf("can not load config file (config file path prefix = %v)", configFilePathPrefix)
	}
	config := new(Config)
	cf.Load(config)
	return &Example{
		name:           algorithmName,
		config:         config.Trade,
	}, nil
}

type ArbitrageExample struct {
	name           string
	config         *ArbitrageTradeConfig
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

func newArbitrageExample(configDir string) (algorithm.ArbitrageTradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		errors.Errorf("can not load config file (config file path prefix = %v)", configFilePathPrefix)
	}
	config := new(Config)
	cf.Load(config)
	return &ArbitrageExample{
		name:           algorithmName,
		config:         config.ArbitrageTrade,
	}, nil
}

func init() {
	algorithm.RegisterAlgorithm(algorithmName, newExample, newArbitrageExample)
}



