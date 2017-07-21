package lazydog

import (
	"github.com/AutomaticCoinTrader/ACT/algorithm"
	"github.com/AutomaticCoinTrader/ACT/exchange"
	"github.com/AutomaticCoinTrader/ACT/notifier"
	"github.com/AutomaticCoinTrader/ACT/configurator"
	"github.com/pkg/errors"
	"log"
	"path"
)

// Stochastic RSI Oscillator と DMI を併用したアルゴリズム
// 買い -> 売りのみ
// https://www.moneypartners.co.jp/support/tech/sctrsi.html
// https://www.moneypartners.co.jp/support/tech/dmi-adx.html


const (
	algorithmName string = "lazydog"
)


type tradeConfig struct {
}

type arbitrageTradeConfig struct {
}

type config struct {
	Trade          *tradeConfig          `json:"trade"          yaml:"trade"          toml:"trade"`
	ArbitrageTrade *arbitrageTradeConfig `json:"arbitrageTrade" yaml:"arbitrageTrade" toml:"arbitrageTrade"`
}

type lazydog struct {
	name           string
	config         *tradeConfig
}

func (l *lazydog) GetName() (string) {
	return l.name
}

func (l *lazydog) Initialize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func (l *lazydog) Update(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	// trade
	log.Printf("trade")
	return nil
}

func (l *lazydog) Finalize(tradeContext exchange.TradeContext, notifier *notifier.Notifier) (error) {
	return nil
}

func newExample(configDir string) (algorithm.TradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	conf := new(config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &lazydog{
		name:           algorithmName,
		config:         conf.Trade,
	}, nil
}

type arbitrageLazydog struct {
	name           string
	config         *arbitrageTradeConfig
}

func (a *arbitrageLazydog) GetName() (string) {
	return a.name
}

func (a *arbitrageLazydog) Initialize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func (a *arbitrageLazydog) Update(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	// arbitrage trade
	log.Printf("arbitrage trade")
	return nil
}

func (a *arbitrageLazydog) Finalize(exchanges map[string]exchange.Exchange, notifier *notifier.Notifier) (error) {
	return nil
}

func newArbitrageExample(configDir string) (algorithm.ArbitrageTradeAlgorithm, error) {
	configFilePathPrefix := path.Join(configDir, algorithmName)
	cf, err := configurator.NewConfigurator(configFilePathPrefix)
	if err != nil {
		return nil, errors.Errorf("can not create configurator (config file path prefix = %v)", configFilePathPrefix)
	}
	conf := new(config)
	err = cf.Load(conf)
	if err != nil {
		return nil, errors.Errorf("can not load config (config file path prefix = %v)", configFilePathPrefix)
	}
	return &arbitrageLazydog{
		name:           algorithmName,
		config:         conf.ArbitrageTrade,
	}, nil
}

func init() {
	algorithm.RegisterAlgorithm(algorithmName, newExample, newArbitrageExample)
}


